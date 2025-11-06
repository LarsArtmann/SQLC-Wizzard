package commands

import (
	"fmt"
	"os"

	"github.com/LarsArtmann/SQLC-Wizzard/internal/generators"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/templates"
	"github.com/spf13/cobra"
)

// GenerateOptions contains options for generate command
type GenerateOptions struct {
	configPath string
	outputDir  string
	force     bool
}

// NewGenerateCommand creates the generate command
func NewGenerateCommand() *cobra.Command {
	opts := &GenerateOptions{}

	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate SQL files and configurations",
		Long: `Generate creates example SQL files and configurations.
Use this to quickly scaffold a working sqlc setup.`,
		Example: `  sqlc-wizard generate
  sqlc-wizard generate --output ./generated --force`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGenerate(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.configPath, "config", "c", "", "Path to sqlc.yaml configuration file")
	cmd.Flags().StringVarP(&opts.outputDir, "output", "o", ".", "Output directory for generated files")
	cmd.Flags().BoolVarP(&opts.force, "force", "f", false, "Overwrite existing files")

	return cmd
}

func runGenerate(opts *GenerateOptions) error {
	// For now, just generate example files
	return generateExampleFiles(opts.outputDir, opts.force)
}

func generateExampleFiles(outputDir string, force bool) error {
	// Create generator
	generator := generators.NewGenerator(outputDir)

	// Create template data with defaults
	templateData := templates.TemplateData{
		ProjectName: "generated-project",
		ProjectType: templates.MustNewProjectType("microservice"),
		
		Package: templates.PackageConfig{
			Name: "db",
			Path: "db",
		},
		
		Database: templates.DatabaseConfig{
			Engine:     templates.MustNewDatabaseType("postgresql"),
			UseUUIDs:   false,
			UseJSON:    true,
			UseArrays:  false,
		},
		
		Output: templates.OutputConfig{
			BaseDir:    outputDir,
			QueriesDir:  "queries",
			SchemaDir:   "schema",
		},
		
		Validation: templates.ValidationConfig{
			StrictFunctions: true,
			StrictOrderBy:   true,
		},
	}

	// Check if output directory exists and has files
	if !force {
		if _, err := os.Stat(outputDir); err == nil {
			files, err := os.ReadDir(outputDir)
			if err == nil && len(files) > 0 {
				return fmt.Errorf("output directory %s is not empty. Use --force to overwrite", outputDir)
			}
		}
	}

	// Generate example files
	if err := generator.GenerateExampleSchema(templateData); err != nil {
		return fmt.Errorf("failed to generate schema: %w", err)
	}

	if err := generator.GenerateExampleQueries(templateData); err != nil {
		return fmt.Errorf("failed to generate queries: %w", err)
	}

	fmt.Printf("✅ Successfully generated example SQL files to %s\n", outputDir)
	fmt.Printf("📄 Generated files:\n")
	fmt.Printf("   - %s/schema/001_users_table.sql\n", outputDir)
	fmt.Printf("   - %s/queries/users.sql\n", outputDir)

	return nil
}