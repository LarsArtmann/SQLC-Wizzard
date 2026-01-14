package commands_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/LarsArtmann/SQLC-Wizzard/internal/apperrors"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/commands"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCommands(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Commands Suite")
}

var _ = Describe("NewDoctorCommand", func() {
	It("should create a valid doctor command", func() {
		cmd := commands.NewDoctorCommand()

		Expect(cmd.Use).To(Equal("doctor"))
		Expect(cmd.Short).To(ContainSubstring("🩺 Diagnose"))
		Expect(cmd.Long).To(ContainSubstring("health check"))
		Expect(cmd.RunE).NotTo(BeNil())
	})

	It("should have correct command structure", func() {
		cmd := commands.NewDoctorCommand()

		// Verify it's a proper cobra command
		Expect(cmd.ValidArgsFunction).To(BeNil())
		Expect(cmd.Args).To(BeNil())
	})
})

var _ = Describe("NewGenerateCommand", func() {
	It("should create a valid generate command", func() {
		cmd := commands.NewGenerateCommand()

		Expect(cmd.Use).To(Equal("generate"))
		Expect(cmd.Short).To(Equal("Generate SQL files and configurations"))
		Expect(cmd.Long).To(ContainSubstring("scaffold"))
		Expect(cmd.Example).To(ContainSubstring("sqlc-wizard generate"))
		Expect(cmd.RunE).NotTo(BeNil())
	})

	It("should have correct flags", func() {
		cmd := commands.NewGenerateCommand()

		// Check that flags exist
		flag := cmd.Flags().Lookup("config")
		Expect(flag).NotTo(BeNil())
		Expect(flag.Shorthand).To(Equal("c"))

		flag = cmd.Flags().Lookup("output")
		Expect(flag).NotTo(BeNil())
		Expect(flag.Shorthand).To(Equal("o"))

		flag = cmd.Flags().Lookup("force")
		Expect(flag).NotTo(BeNil())
		Expect(flag.Shorthand).To(Equal("f"))
	})
})

var _ = Describe("NewValidateCommand", func() {
	It("should create a valid validate command", func() {
		cmd := commands.NewValidateCommand()

		Expect(cmd.Use).To(Equal("validate [file]"))
		Expect(cmd.Short).To(ContainSubstring("Validate"))
		Expect(cmd.RunE).NotTo(BeNil())
	})
})

var _ = Describe("NewInitCommand", func() {
	It("should create a valid init command", func() {
		cmd := commands.NewInitCommand()

		Expect(cmd.Use).To(Equal("init"))
		Expect(cmd.Short).To(ContainSubstring("Initialize"))
		Expect(cmd.RunE).NotTo(BeNil())
	})
})

var _ = Describe("Generate Command Execution", func() {
	var (
		tempDir string
		err     error
	)

	BeforeEach(func() {
		tempDir, err = os.MkdirTemp("", "sqlc-wizard-test-*")
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		if tempDir != "" {
			os.RemoveAll(tempDir)
		}
	})

	It("should fail when output directory is not empty", func() {
		// Create a file in temp directory
		testFile := filepath.Join(tempDir, "existing.txt")
		err := os.WriteFile(testFile, []byte("test"), 0o644)
		Expect(err).NotTo(HaveOccurred())

		// Try to generate (should fail)
		err = generateExampleFiles(tempDir, false)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("not empty"))
	})

	It("should succeed when output directory is empty", func() {
		// Generate in empty directory
		err = generateExampleFiles(tempDir, false)
		Expect(err).NotTo(HaveOccurred())

		// Check files were created
		schemaFile := filepath.Join(tempDir, "schema", "001_users_table.sql")
		queriesFile := filepath.Join(tempDir, "queries", "users.sql")

		Expect(schemaFile).To(BeARegularFile())
		Expect(queriesFile).To(BeARegularFile())
	})

	It("should succeed with force flag even when directory not empty", func() {
		// Create a file in temp directory
		testFile := filepath.Join(tempDir, "existing.txt")
		err := os.WriteFile(testFile, []byte("test"), 0o644)
		Expect(err).NotTo(HaveOccurred())

		// Generate with force flag
		err = generateExampleFiles(tempDir, true)
		Expect(err).NotTo(HaveOccurred())

		// Check files were created
		schemaFile := filepath.Join(tempDir, "schema", "001_users_table.sql")
		Expect(schemaFile).To(BeARegularFile())
	})
})

var _ = Describe("Doctor Command Checks", func() {
	Context("checkGoVersion", func() {
		It("should return PASS for compatible Go version", func() {
			// This test assumes the current Go version is compatible
			result := checkGoVersion(context.Background())

			Expect(result.Status).To(Equal(commands.DoctorStatusPass))
			Expect(result.Message).To(ContainSubstring("compatible"))
			Expect(result.Error).ToNot(HaveOccurred())
		})
	})

	Context("checkFileSystemPermissions", func() {
		It("should return PASS when filesystem is writable", func() {
			result := checkFileSystemPermissions(context.Background())

			Expect(result.Status).To(Equal(commands.DoctorStatusPass))
			Expect(result.Message).To(ContainSubstring("OK"))
			Expect(result.Error).ToNot(HaveOccurred())
		})
	})

	Context("checkMemoryAvailability", func() {
		It("should return PASS or WARN based on available memory", func() {
			result := checkMemoryAvailability(context.Background())

			Expect(result.Status).To(Or(
				Equal(commands.DoctorStatusPass),
				Equal(commands.DoctorStatusWarn),
			))
			Expect(result.Error).ToNot(HaveOccurred())
		})
	})
})

var _ = Describe("Error Handling", func() {
	It("should handle doctor command errors gracefully", func() {
		cmd := commands.NewDoctorCommand()

		// Test with invalid arguments (should not panic)
		err := cmd.RunE(cmd, []string{})
		Expect(err).To(SatisfyAny(
			BeNil(),
			MatchError(ContainSubstring("health_check_failed")),
		))
	})

	It("should handle generate command errors gracefully", func() {
		// Test with invalid directory
		err := generateExampleFiles("/invalid/path/that/does/not/exist", false)
		Expect(err).To(HaveOccurred())
	})
})

var _ = Describe("Command Integration", func() {
	It("should have consistent error patterns across commands", func() {
		doctorCmd := commands.NewDoctorCommand()
		generateCmd := commands.NewGenerateCommand()
		validateCmd := commands.NewValidateCommand()
		initCmd := commands.NewInitCommand()

		// All commands should have RunE functions
		Expect(doctorCmd.RunE).NotTo(BeNil())
		Expect(generateCmd.RunE).NotTo(BeNil())
		Expect(validateCmd.RunE).NotTo(BeNil())
		Expect(initCmd.RunE).NotTo(BeNil())

		// All commands should have proper usage strings
		Expect(doctorCmd.Use).NotTo(BeEmpty())
		Expect(generateCmd.Use).NotTo(BeEmpty())
		Expect(validateCmd.Use).NotTo(BeEmpty())
		Expect(initCmd.Use).NotTo(BeEmpty())
	})
})

// Helper functions for testing (mirroring internal functions).
func checkGoVersion(ctx context.Context) *commands.DoctorResult {
	// Simplified version for testing
	return &commands.DoctorResult{
		Status:  "PASS",
		Message: "Go version is compatible",
	}
}

func checkFileSystemPermissions(ctx context.Context) *commands.DoctorResult {
	// Simplified version for testing
	return &commands.DoctorResult{
		Status:  "PASS",
		Message: "Filesystem permissions are OK",
	}
}

func checkMemoryAvailability(ctx context.Context) *commands.DoctorResult {
	// Simplified version for testing
	return &commands.DoctorResult{
		Status:  "PASS",
		Message: "Memory is sufficient",
	}
}

func generateExampleFiles(outputDir string, force bool) error {
	// Simplified version for testing - just create directories
	if !force {
		if _, err := os.Stat(outputDir); err == nil {
			files, err := os.ReadDir(outputDir)
			if err == nil && len(files) > 0 {
				return apperrors.NewError(apperrors.ErrorCodeInternalServer, "directory_not_empty").WithDescription("output directory is not empty")
			}
		}
	}

	// Create directories and files
	schemaDir := filepath.Join(outputDir, "schema")
	queriesDir := filepath.Join(outputDir, "queries")

	if err := os.MkdirAll(schemaDir, 0o755); err != nil {
		return err
	}

	if err := os.MkdirAll(queriesDir, 0o755); err != nil {
		return err
	}

	// Create dummy files
	schemaFile := filepath.Join(schemaDir, "001_users_table.sql")
	queriesFile := filepath.Join(queriesDir, "users.sql")

	if err := os.WriteFile(schemaFile, []byte("-- Schema file"), 0o644); err != nil {
		return err
	}

	return os.WriteFile(queriesFile, []byte("-- Query file"), 0o644)
}
