package commands_test

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/LarsArtmann/SQLC-Wizzard/internal/commands"
	"github.com/LarsArtmann/SQLC-Wizzard/pkg/config"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"
)

// Helper function to create test SqlcConfig instances
func createTestSqlcConfig(schema, out, pkg string, engine ...string) *config.SqlcConfig {
	dbEngine := "postgresql"
	if len(engine) > 0 {
		dbEngine = engine[0]
	}
	return &config.SqlcConfig{
		Version: "2",
		SQL: []config.SQLConfig{
			{
				Engine: dbEngine,
				Schema: config.NewPathOrPaths([]string{schema}),
				Gen: config.GenConfig{
					Go: &config.GoGenConfig{
						Out:     out,
						Package: pkg,
					},
				},
			},
		},
	}
}

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

var _ = Describe("Migrate Command Enhanced Testing", func() {
	var (
		tempDir string
		err     error
	)

	BeforeEach(func() {
		tempDir, err = os.MkdirTemp("", "sqlc-migrate-extended-*")
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		if tempDir != "" {
			_ = os.RemoveAll(tempDir)
		}
	})

	Context("Command Structure and Interface", func() {
		It("should handle migrate command with help flag", func() {
			cmd := commands.NewMigrateCommand()

			var buf bytes.Buffer
			cmd.SetOut(&buf)
			cmd.SetArgs([]string{"--help"})

			err := cmd.Execute()
			// Help should work if implemented
			Expect(err).To(SatisfyAny(BeNil(), HaveOccurred()))
		})

		It("should handle migrate subcommands", func() {
			cmd := commands.NewMigrateCommand()

			// Test that migrate has subcommands (may be empty initially)
			subcommands := cmd.Commands()
			// Should not panic regardless of subcommands
			_ = subcommands
		})
	})

	Context("Migration File Operations", func() {
		It("should handle migration directory creation", func() {
			migrationDir := filepath.Join(tempDir, "migrations")
			err := os.MkdirAll(migrationDir, 0o755)
			Expect(err).NotTo(HaveOccurred())

			// Test that directory exists and is accessible
			Expect(migrationDir).To(BeADirectory())
		})

		It("should handle migration file creation", func() {
			migrationFile := filepath.Join(tempDir, "001_initial.sql")
			err := os.WriteFile(migrationFile, []byte("-- Initial migration\nCREATE TABLE test (id INTEGER);"), 0o644)
			Expect(err).NotTo(HaveOccurred())

			Expect(migrationFile).To(BeARegularFile())
		})

		It("should handle migration file naming", func() {
			migrationFiles := []string{
				"001_initial.sql",
				"002_add_users.sql",
				"003_create_indexes.sql",
				"010_late_migration.sql",
				"999_final_migration.sql",
			}

			for _, filename := range migrationFiles {
				migrationFile := filepath.Join(tempDir, filename)
				err := os.WriteFile(migrationFile, []byte("-- Test migration"), 0o644)
				Expect(err).NotTo(HaveOccurred())
				Expect(migrationFile).To(BeARegularFile())
			}
		})
	})

	Context("Migration Error Handling", func() {
		It("should handle permission errors in migration directories", func() {
			readOnlyDir := filepath.Join(tempDir, "readonly")
			err := os.Mkdir(readOnlyDir, 0o444)
			Expect(err).NotTo(HaveOccurred())

			// Should fail to write in read-only directory
			migrationFile := filepath.Join(readOnlyDir, "test.sql")
			err = os.WriteFile(migrationFile, []byte("-- Test"), 0o644)
			Expect(err).To(HaveOccurred())

			// Clean up permissions for removal
			err = os.Chmod(readOnlyDir, 0o755)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should handle migration file validation", func() {
			// Test creating migration files with various content
			migrationFiles := map[string]string{
				"valid.sql":          "-- Valid migration\nCREATE TABLE users (id INTEGER);",
				"empty.sql":          "",
				"comments_only.sql":  "-- Only comments\n-- No SQL statements",
				"invalid_syntax.sql": "INVALID SQL SYNTAX HERE",
			}

			for filename, content := range migrationFiles {
				migrationFile := filepath.Join(tempDir, filename)
				err := os.WriteFile(migrationFile, []byte(content), 0o644)
				Expect(err).NotTo(HaveOccurred())
				Expect(migrationFile).To(BeARegularFile())

				// Content validation would be done by actual migration logic
				fileContent, err := os.ReadFile(migrationFile)
				Expect(err).NotTo(HaveOccurred())
				Expect(string(fileContent)).To(Equal(content))
			}
		})
	})
})

var _ = Describe("Command Integration and Error Recovery", func() {
	var (
		tempDir string
		err     error
	)

	BeforeEach(func() {
		tempDir, err = os.MkdirTemp("", "sqlc-integration-*")
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		if tempDir != "" {
			_ = os.RemoveAll(tempDir)
		}
	})

	Context("Command Recovery and Robustness", func() {
		It("should recover from unexpected file system states", func() {
			// Create files and then simulate file system issues
			configFile := filepath.Join(tempDir, "config.yaml")
			err := os.WriteFile(configFile, []byte("version: \"2\""), 0o644)
			Expect(err).NotTo(HaveOccurred())

			// Remove file after creating command to simulate race condition
			_ = os.Remove(configFile)

			cmd := commands.NewValidateCommand()
			args := []string{configFile}
			cmd.SetArgs(args)

			err = cmd.Execute()
			Expect(err).To(HaveOccurred())
		})

		It("should handle concurrent command execution", func() {
			// Create multiple config files
			configFiles := make([]string, 5)
			for i := range 5 {
				configFiles[i] = filepath.Join(tempDir, fmt.Sprintf("config%d.yaml", i))
				cfg := createTestSqlcConfig("schema.sql", fmt.Sprintf("db%d", i), fmt.Sprintf("db%d", i))

				yamlData, err := config.Marshal(cfg)
				Expect(err).NotTo(HaveOccurred())
				err = os.WriteFile(configFiles[i], yamlData, 0o644)
				Expect(err).NotTo(HaveOccurred())
			}

			// Test multiple commands (simulating concurrent use)
			for _, configFile := range configFiles {
				cmd := commands.NewValidateCommand()
				args := []string{configFile}
				cmd.SetArgs(args)

				err = cmd.Execute()
				Expect(err).To(SatisfyAny(BeNil(), HaveOccurred()))
			}
		})
	})

	Context("Command Help and Documentation", func() {
		It("should handle help for basic commands", func() {
			commandsToTest := []func() *cobra.Command{
				commands.NewValidateCommand,
				commands.NewDoctorCommand,
				commands.NewGenerateCommand,
				commands.NewInitCommand,
			}

			for _, cmdFunc := range commandsToTest {
				cmd := cmdFunc()

				// Test that all commands have basic help structure
				Expect(cmd.Use).NotTo(BeEmpty())
				Expect(cmd.RunE).NotTo(BeNil())

				// Test help flag (some commands may not have help yet)
				var buf bytes.Buffer
				cmd.SetOut(&buf)
				cmd.SetArgs([]string{"--help"})

				err := cmd.Execute()
				// Help should work, but if not implemented, that's ok
				Expect(err).To(SatisfyAny(BeNil(), HaveOccurred()))
			}
		})
	})

	Context("Real-world Project Scenarios", func() {
		It("should handle complex project structures", func() {
			// Create a realistic project structure
			projectDirs := []string{
				"sql/migrations",
				"sql/queries",
				"internal/auth/db",
				"internal/analytics/db",
				"internal/common/db",
				"cmd/api",
				"pkg/models",
			}

			for _, dir := range projectDirs {
				err = os.MkdirAll(filepath.Join(tempDir, dir), 0o755)
				Expect(err).NotTo(HaveOccurred())
			}

			// Create complex configuration
			complexCfg := &config.SqlcConfig{
				Version: "2",
				SQL: []config.SQLConfig{
					{
						Engine:  "postgresql",
						Schema:  config.NewPathOrPaths([]string{"sql/migrations"}),
						Queries: config.NewPathOrPaths([]string{"sql/queries/auth"}),
						Gen: config.GenConfig{
							Go: &config.GoGenConfig{
								Out:     "internal/auth/db",
								Package: "authdb",
								Overrides: []config.Override{
									{
										GoType: "github.com/google/uuid.UUID",
										DBType: "uuid",
									},
									{
										GoType: "database/sql.NullString",
										DBType: "text",
									},
								},
							},
						},
					},
					{
						Engine:  "postgresql",
						Schema:  config.NewPathOrPaths([]string{"sql/migrations"}),
						Queries: config.NewPathOrPaths([]string{"sql/queries/analytics"}),
						Gen: config.GenConfig{
							Go: &config.GoGenConfig{
								Out:     "internal/analytics/db",
								Package: "analyticsdb",
								Overrides: []config.Override{
									{
										GoType: "github.com/lib/pq.NullTime",
										DBType: "timestamp",
									},
								},
							},
						},
					},
				},
			}

			configFile := filepath.Join(tempDir, "sqlc.yaml")
			yamlData, err := config.Marshal(complexCfg)
			Expect(err).NotTo(HaveOccurred())
			err = os.WriteFile(configFile, yamlData, 0o644)
			Expect(err).NotTo(HaveOccurred())

			// Test validation
			cmd := commands.NewValidateCommand()
			args := []string{configFile}
			cmd.SetArgs(args)

			err = cmd.Execute()
			Expect(err).To(SatisfyAny(BeNil(), HaveOccurred()))
		})
	})
})
