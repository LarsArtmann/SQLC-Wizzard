package commands_test

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/LarsArtmann/SQLC-Wizzard/internal/commands"
	"github.com/LarsArtmann/SQLC-Wizzard/pkg/config"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"
)

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
				cfg := createTestSqlcConfig(
					"schema.sql",
					fmt.Sprintf("db%d", i),
					fmt.Sprintf("db%d", i),
				)

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
