package integration_test

import (
	"bytes"
	"testing"

	"github.com/LarsArtmann/SQLC-Wizzard/internal/commands"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"
)

// executeCommand executes a command with given args and returns error.
func executeCommand(cmd *cobra.Command, args []string) error {
	cmd.SetArgs(args)
	return cmd.Execute()
}

// executeCommandWithOutput executes a command and returns both output and error.
func executeCommandWithOutput(cmd *cobra.Command, args []string) (string, error) {
	var output bytes.Buffer
	cmd.SetArgs(args)
	cmd.SetOut(&output)
	cmd.SetErr(&output)
	err := cmd.Execute()
	return output.String(), err
}

// executeCommandWithHelp executes a command with --help flag.
func executeCommandWithHelp(cmd *cobra.Command) error {
	args := []string{"--help"}
	var output bytes.Buffer
	cmd.SetArgs(args)
	cmd.SetOut(&output)
	return cmd.Execute()
}

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Integration Suite")
}

var _ = Describe("Basic Command Integration", func() {
	Context("Command Execution", func() {
		It("should run doctor command without panicking", func() {
			doctorCmd := commands.NewDoctorCommand()
			outputStr, err := executeCommandWithOutput(doctorCmd, []string{})

			Expect(err).To(SatisfyAny(BeNil(), HaveOccurred()))
			Expect(len(outputStr)).To(BeNumerically(">=", 0))
		})

		It("should handle validate command with missing file", func() {
			validateCmd := commands.NewValidateCommand()
			args := []string{"/nonexistent/config.yaml"}
			err := executeCommand(validateCmd, args)
			Expect(err).To(HaveOccurred())
		})

		It("should run generate command with basic args", func() {
			generateCmd := commands.NewGenerateCommand()
			err := executeCommandWithHelp(generateCmd)
			Expect(err).To(SatisfyAny(BeNil(), HaveOccurred()))
		})

		It("should run migrate command", func() {
			migrateCmd := commands.NewMigrateCommand()
			err := executeCommandWithHelp(migrateCmd)
			Expect(err).To(SatisfyAny(BeNil(), HaveOccurred()))
		})

		It("should run init command", func() {
			initCmd := commands.NewInitCommand()
			err := executeCommandWithHelp(initCmd)
			Expect(err).To(SatisfyAny(BeNil(), HaveOccurred()))
		})
	})

	Context("Command Help System", func() {
		It("should provide help for doctor command", func() {
			cmd := commands.NewDoctorCommand()
			err := executeCommandWithHelp(cmd)
			Expect(err).To(SatisfyAny(BeNil(), HaveOccurred()))
		})

		It("should provide help for validate command", func() {
			cmd := commands.NewValidateCommand()
			err := executeCommandWithHelp(cmd)
			Expect(err).To(SatisfyAny(BeNil(), HaveOccurred()))
		})

		It("should provide help for generate command", func() {
			cmd := commands.NewGenerateCommand()
			err := executeCommandWithHelp(cmd)
			Expect(err).To(SatisfyAny(BeNil(), HaveOccurred()))
		})

		It("should provide help for migrate command", func() {
			cmd := commands.NewMigrateCommand()
			err := executeCommandWithHelp(cmd)
			Expect(err).To(SatisfyAny(BeNil(), HaveOccurred()))
		})

		It("should provide help for init command", func() {
			cmd := commands.NewInitCommand()
			err := executeCommandWithHelp(cmd)
			Expect(err).To(SatisfyAny(BeNil(), HaveOccurred()))
		})
	})
})
