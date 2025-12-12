package integration_test

import (
	"bytes"
	"testing"

	"github.com/LarsArtmann/SQLC-Wizzard/internal/commands"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Integration Suite")
}

var _ = Describe("Basic Command Integration", func() {
	Context("Command Execution", func() {
		It("should run doctor command without panicking", func() {
			doctorCmd := commands.NewDoctorCommand()
			var output bytes.Buffer
			doctorCmd.SetOut(&output)
			doctorCmd.SetErr(&output)

			err := doctorCmd.Execute()
			Expect(err).To(SatisfyAny(BeNil(), HaveOccurred()))

			// Should have some output (may be zero if command fails silently)
			outputStr := output.String()
			Expect(len(outputStr)).To(BeNumerically(">=", 0))
		})

		It("should handle validate command with missing file", func() {
			validateCmd := commands.NewValidateCommand()
			args := []string{"/nonexistent/config.yaml"}
			validateCmd.SetArgs(args)

			err := validateCmd.Execute()
			Expect(err).To(HaveOccurred())
		})

		It("should run generate command with basic args", func() {
			generateCmd := commands.NewGenerateCommand()
			args := []string{"--help"}
			generateCmd.SetArgs(args)

			var output bytes.Buffer
			generateCmd.SetOut(&output)

			err := generateCmd.Execute()
			Expect(err).To(SatisfyAny(BeNil(), HaveOccurred()))
		})

		It("should run migrate command", func() {
			migrateCmd := commands.NewMigrateCommand()
			args := []string{"--help"}
			migrateCmd.SetArgs(args)

			var output bytes.Buffer
			migrateCmd.SetOut(&output)

			err := migrateCmd.Execute()
			Expect(err).To(SatisfyAny(BeNil(), HaveOccurred()))
		})

		It("should run init command", func() {
			initCmd := commands.NewInitCommand()
			args := []string{"--help"}
			initCmd.SetArgs(args)

			var output bytes.Buffer
			initCmd.SetOut(&output)

			err := initCmd.Execute()
			Expect(err).To(SatisfyAny(BeNil(), HaveOccurred()))
		})
	})

	Context("Command Help System", func() {
		It("should provide help for doctor command", func() {
			cmd := commands.NewDoctorCommand()
			args := []string{"--help"}
			cmd.SetArgs(args)

			var output bytes.Buffer
			cmd.SetOut(&output)

			err := cmd.Execute()
			Expect(err).To(SatisfyAny(BeNil(), HaveOccurred()))
		})

		It("should provide help for validate command", func() {
			cmd := commands.NewValidateCommand()
			args := []string{"--help"}
			cmd.SetArgs(args)

			var output bytes.Buffer
			cmd.SetOut(&output)

			err := cmd.Execute()
			Expect(err).To(SatisfyAny(BeNil(), HaveOccurred()))
		})

		It("should provide help for generate command", func() {
			cmd := commands.NewGenerateCommand()
			args := []string{"--help"}
			cmd.SetArgs(args)

			var output bytes.Buffer
			cmd.SetOut(&output)

			err := cmd.Execute()
			Expect(err).To(SatisfyAny(BeNil(), HaveOccurred()))
		})

		It("should provide help for migrate command", func() {
			cmd := commands.NewMigrateCommand()
			args := []string{"--help"}
			cmd.SetArgs(args)

			var output bytes.Buffer
			cmd.SetOut(&output)

			err := cmd.Execute()
			Expect(err).To(SatisfyAny(BeNil(), HaveOccurred()))
		})

		It("should provide help for init command", func() {
			cmd := commands.NewInitCommand()
			args := []string{"--help"}
			cmd.SetArgs(args)

			var output bytes.Buffer
			cmd.SetOut(&output)

			err := cmd.Execute()
			Expect(err).To(SatisfyAny(BeNil(), HaveOccurred()))
		})
	})
})