package commands

import (
	"fmt"
	"os"

	"github.com/LarsArtmann/SQLC-Wizzard/internal/domain"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/generators"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/templates"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/wizard"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

// InitOptions contains options for the init command
type InitOptions struct {
	ProjectType    string
	Database       string
	PackagePath    string
	OutputDir      string
	NonInteractive bool
}

// NewInitCommand creates the init command
func NewInitCommand() *cobra.Command {
	opts := &InitOptions{}

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize a new sqlc configuration",
		Long: `Interactive wizard to create a production-ready sqlc configuration.

The wizard will guide you through:
  • Selecting your project type
  • Choosing your database
  • Configuring output directories
  • Enabling features and safety rules
  • Generating example queries and schema

Example:
  sqlc-wizard init
  sqlc-wizard init --non-interactive --project-type=microservice --database=postgresql`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runInit(opts)
		},
	}

	// Add flags
	cmd.Flags().StringVar(&opts.ProjectType, "project-type", "", "Project type (hobby, microservice, enterprise, api-first, library)")
	cmd.Flags().StringVar(&opts.Database, "database", "", "Database engine (postgresql, mysql, sqlite)")
	cmd.Flags().StringVar(&opts.PackagePath, "package", "", "Go package path (e.g., github.com/user/project)")
	cmd.Flags().StringVarP(&opts.OutputDir, "output-dir", "o", ".", "Output directory for generated files")
	cmd.Flags().BoolVar(&opts.NonInteractive, "non-interactive", false, "Run in non-interactive mode using flags")

	return cmd
}

func runInit(opts *InitOptions) error {
	// Check if sqlc.yaml already exists
	if _, err := os.Stat("sqlc.yaml"); err == nil {
		return fmt.Errorf("sqlc.yaml already exists in current directory. Remove it first or run in a different directory")
	}

	var result *wizard.WizardResult
	var err error

	if opts.NonInteractive {
		// Non-interactive mode: use flags
		result, err = runNonInteractive(opts)
	} else {
		// Interactive mode: run wizard
		w := wizard.NewWizard()
		result, err = w.Run()
	}

	if err != nil {
		return fmt.Errorf("wizard failed: %w", err)
	}

	// Generate files
	gen := generators.NewGenerator(opts.OutputDir)
	if err := gen.GenerateAll(result.Config, result.TemplateData, result.GenerateQueries, result.GenerateSchema); err != nil {
		return fmt.Errorf("generation failed: %w", err)
	}

	// Show success message
	showSuccess(gen, result)

	return nil
}

func runNonInteractive(opts *InitOptions) (*wizard.WizardResult, error) {
	// Validate required flags
	if opts.ProjectType == "" {
		return nil, fmt.Errorf("--project-type is required in non-interactive mode")
	}
	if opts.Database == "" {
		return nil, fmt.Errorf("--database is required in non-interactive mode")
	}
	if opts.PackagePath == "" {
		return nil, fmt.Errorf("--package is required in non-interactive mode")
	}

	// Create template data from flags
	data := templates.TemplateData{
		ProjectType:  templates.ProjectType(opts.ProjectType),
		Database:     templates.DatabaseType(opts.Database),
		PackagePath:  opts.PackagePath,
		Features:     templates.DefaultFeatures(),
		SafetyRules:  domain.DefaultSafetyRules(),
	}

	// Generate config from template
	tmpl, err := templates.GetTemplate(data.ProjectType)
	if err != nil {
		return nil, fmt.Errorf("invalid project type: %w", err)
	}

	cfg, err := tmpl.Generate(data)
	if err != nil {
		return nil, fmt.Errorf("failed to generate config: %w", err)
	}

	return &wizard.WizardResult{
		Config:          cfg,
		TemplateData:    data,
		GenerateQueries: true,
		GenerateSchema:  true,
	}, nil
}

func showSuccess(gen *generators.Generator, result *wizard.WizardResult) {
	successStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("10")).
		Padding(1, 0)

	summaryStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Padding(0, 0, 1, 0)

	fmt.Println(successStyle.Render("✓ Successfully generated sqlc configuration!"))
	fmt.Println(summaryStyle.Render(gen.GenerateSummary(result.Config, result.GenerateQueries, result.GenerateSchema)))
}
