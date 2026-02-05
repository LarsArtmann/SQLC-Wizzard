package commands

import (
	"fmt"

	"github.com/LarsArtmann/SQLC-Wizzard/pkg/config"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

// ValidateOptions contains options for the validate command.
type ValidateOptions struct {
	ConfigPath string
	Fix        bool
	Strict     bool
}

// NewValidateCommand creates the validate command.
func NewValidateCommand() *cobra.Command {
	opts := &ValidateOptions{}

	cmd := &cobra.Command{
		Use:   "validate [file]",
		Short: "Validate sqlc configuration files",
		Long: `Validate an sqlc.yaml configuration file for correctness and best practices.

The validator checks for:
  • Required fields and valid values
  • Database engine compatibility
  • Path configurations
  • Best practice recommendations

Example:
  sqlc-wizard validate
  sqlc-wizard validate sqlc.yaml
  sqlc-wizard validate --strict`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Determine config path
			if len(args) > 0 {
				opts.ConfigPath = args[0]
			} else {
				opts.ConfigPath = "sqlc.yaml"
			}

			return runValidate(opts)
		},
	}

	// Add flags
	cmd.Flags().BoolVar(&opts.Fix, "fix", false, "Attempt to auto-fix common issues")
	cmd.Flags().BoolVar(&opts.Strict, "strict", false, "Enable strict validation mode")

	return cmd
}

func runValidate(opts *ValidateOptions) error {
	// Parse config file
	cfg, err := config.ParseFile(opts.ConfigPath)
	if err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	// Validate
	result := config.Validate(cfg)

	// Display results
	displayValidationResults(result, opts)

	// Return error if validation failed
	if !result.IsValid() {
		return fmt.Errorf("validation failed with %d error(s)", len(result.Errors))
	}

	return nil
}

func displayValidationResults(result *config.ValidationResult, opts *ValidateOptions) {
	successStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("10"))

	errorStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("9"))

	warningStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("11"))

	// Show errors
	if len(result.Errors) > 0 {
		fmt.Println(errorStyle.Render(fmt.Sprintf("✗ Found %d error(s):", len(result.Errors))))
		for _, err := range result.Errors {
			fmt.Printf("  • %s: %s\n", err.Field, err.Message)
		}
		fmt.Println()
	}

	// Show warnings
	if len(result.Warnings) > 0 {
		fmt.Println(warningStyle.Render(fmt.Sprintf("⚠ Found %d warning(s):", len(result.Warnings))))
		for _, warn := range result.Warnings {
			fmt.Printf("  • %s: %s\n", warn.Field, warn.Message)
		}
		fmt.Println()
	}

	// Show success if no errors
	if result.IsValid() {
		if len(result.Warnings) == 0 {
			fmt.Println(successStyle.Render("✓ Configuration is valid!"))
		} else {
			fmt.Println(successStyle.Render("✓ Configuration is valid (with warnings)"))
		}
	}

	// Show fix suggestion if applicable
	if !result.IsValid() && opts.Fix {
		fmt.Println("Note: Auto-fix is not yet implemented. Please fix errors manually.")
	}
}
