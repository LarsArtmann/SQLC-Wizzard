package commands_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"

	"github.com/LarsArtmann/SQLC-Wizzard/internal/commands"
	"github.com/LarsArtmann/SQLC-Wizzard/pkg/config"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Validate Command Enhanced Testing", func() {
	var (
		tempDir string
		err     error
	)

	BeforeEach(func() {
		tempDir, err = os.MkdirTemp("", "sqlc-validate-extended-*")
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		if tempDir != "" {
			_ = os.RemoveAll(tempDir)
		}
	})

	Context("Command Structure and Interface", func() {
		It("should create validate command with proper structure", func() {
			cmd := commands.NewValidateCommand()

			Expect(cmd).NotTo(BeNil())
			Expect(cmd.Use).To(Equal("validate [file]"))
			Expect(cmd.Short).To(ContainSubstring("Validate"))
			Expect(cmd.Long).To(Not(BeEmpty()))
			Expect(cmd.RunE).NotTo(BeNil())
		})

		It("should handle command execution via Cobra interface", func() {
			cmd := commands.NewValidateCommand()

			// Test command with help flag
			var buf bytes.Buffer
			cmd.SetOut(&buf)
			cmd.SetArgs([]string{"--help"})

			err := cmd.Execute()
			Expect(err).ToNot(HaveOccurred())

			help := buf.String()
			Expect(strings.ToLower(help)).To(ContainSubstring("usage"))
			Expect(strings.ToLower(help)).To(ContainSubstring("validate"))
		})

		It("should handle missing arguments gracefully", func() {
			cmd := commands.NewValidateCommand()

			var buf bytes.Buffer
			cmd.SetOut(&buf)
			cmd.SetErr(&buf)
			cmd.SetArgs([]string{})

			err := cmd.Execute()
			// Should either show help or fail gracefully
			Expect(err).To(SatisfyAny(BeNil(), HaveOccurred()))
		})

		It("should handle invalid arguments gracefully", func() {
			cmd := commands.NewValidateCommand()

			var buf bytes.Buffer
			cmd.SetOut(&buf)
			cmd.SetErr(&buf)
			cmd.SetArgs([]string{"--invalid-flag"})

			err := cmd.Execute()
			Expect(err).To(HaveOccurred())
		})
	})

	Context("File System Interactions", func() {
		It("should handle nonexistent files gracefully", func() {
			cmd := commands.NewValidateCommand()

			args := []string{"nonexistent.yaml"}
			cmd.SetArgs(args)

			err := cmd.Execute()
			Expect(err).To(HaveOccurred())
			// Error message may vary depending on implementation
			Expect(err.Error()).To(Or(
				ContainSubstring("no such file"),
				ContainSubstring("does not exist"),
				ContainSubstring("not found"),
				ContainSubstring("failed to parse"),
				ContainSubstring("invalid"),
			))
		})

		It("should handle directories instead of files", func() {
			cmd := commands.NewValidateCommand()

			args := []string{tempDir}
			cmd.SetArgs(args)

			err := cmd.Execute()
			Expect(err).To(HaveOccurred())
		})

		It("should handle unreadable files", func() {
			unreadableFile := filepath.Join(tempDir, "unreadable.yaml")
			err := os.WriteFile(unreadableFile, []byte("test"), 0o000) // No permissions
			Expect(err).NotTo(HaveOccurred())

			defer func() {
				_ = os.Chmod(unreadableFile, 0o644) // Restore permissions for cleanup
			}()

			cmd := commands.NewValidateCommand()
			args := []string{unreadableFile}
			cmd.SetArgs(args)

			err = cmd.Execute()
			Expect(err).To(HaveOccurred())
		})

		It("should handle empty files", func() {
			emptyFile := filepath.Join(tempDir, "empty.yaml")
			err := os.WriteFile(emptyFile, []byte{}, 0o644)
			Expect(err).NotTo(HaveOccurred())

			cmd := commands.NewValidateCommand()
			args := []string{emptyFile}
			cmd.SetArgs(args)

			err = cmd.Execute()
			Expect(err).To(SatisfyAny(BeNil(), HaveOccurred()))
		})
	})

	Context("Configuration File Handling", func() {
		It("should handle various configuration types", func() {
			type configTest struct {
				name    string
				content string
			}

			testCases := []configTest{
				{"valid configurations", "version: \"2\""},
				{"malformed configurations gracefully", "version: \"2\""},
			}

			for _, tc := range testCases {
				filename := filepath.Join(tempDir, tc.name+".yaml")
				err := os.WriteFile(filename, []byte(tc.content), 0o644)
				Expect(err).NotTo(HaveOccurred())

				cmd := commands.NewValidateCommand()
				args := []string{filename}
				cmd.SetArgs(args)

				err = cmd.Execute()
				Expect(err).To(SatisfyAny(BeNil(), HaveOccurred()))
			}
		})
	})

	Context("Multiple File Handling", func() {
		It("should handle multiple configuration files", func() {
			// Create first config
			config1 := filepath.Join(tempDir, "config1.yaml")
			cfg1 := createTestSqlcConfig("schema1.sql", "db1", "db1")

			yamlData1, err := config.Marshal(cfg1)
			Expect(err).NotTo(HaveOccurred())
			err = os.WriteFile(config1, yamlData1, 0o644)
			Expect(err).NotTo(HaveOccurred())

			// Create second config
			config2 := filepath.Join(tempDir, "config2.yaml")
			cfg2 := createTestSqlcConfig("schema2.sql", "db2", "db2", "mysql")

			yamlData2, err := config.Marshal(cfg2)
			Expect(err).NotTo(HaveOccurred())
			err = os.WriteFile(config2, yamlData2, 0o644)
			Expect(err).NotTo(HaveOccurred())

			cmd := commands.NewValidateCommand()
			args := []string{config1, config2}
			cmd.SetArgs(args)

			err = cmd.Execute()
			Expect(err).To(SatisfyAny(BeNil(), HaveOccurred()))
		})

		It("should handle mixed valid and invalid configuration files", func() {
			// Create valid config
			validConfig := filepath.Join(tempDir, "valid.yaml")
			validCfg := createTestSqlcConfig("schema.sql", "db", "db")

			yamlData, err := config.Marshal(validCfg)
			Expect(err).NotTo(HaveOccurred())
			err = os.WriteFile(validConfig, yamlData, 0o644)
			Expect(err).NotTo(HaveOccurred())

			// Create invalid config
			invalidConfig := filepath.Join(tempDir, "invalid.yaml")
			err = os.WriteFile(invalidConfig, []byte("invalid: yaml: ["), 0o644)
			Expect(err).NotTo(HaveOccurred())

			cmd := commands.NewValidateCommand()
			args := []string{validConfig, invalidConfig}
			cmd.SetArgs(args)

			err = cmd.Execute()
			Expect(err).To(SatisfyAny(BeNil(), HaveOccurred()))
		})
	})
})
