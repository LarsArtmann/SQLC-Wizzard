package commands

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/adapters"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/creators"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/templates"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

// CreateOptions contains options for the create command.
type CreateOptions struct {
	ProjectType     string
	Database        string
	OutputDir       string
	IncludeAuth     bool
	IncludeFrontend bool
	NonInteractive  bool
	Force           bool
}

// NewCreateCommand creates the create command for complete project generation.
func NewCreateCommand() *cobra.Command {
	opts := &CreateOptions{}

	cmd := &cobra.Command{
		Use:   "create [project-name]",
		Short: "🚀 Create a complete SQLC project instantly",
		Long: `Magical one-command project setup that creates a complete, production-ready SQLC project.

This command generates everything you need to start developing:
✅ Optimized sqlc.yaml configuration
✅ Database schema templates  
✅ Migration files structure
✅ Go package structure
✅ Development scripts
✅ Docker configuration
✅ Makefile with common tasks
✅ README with usage instructions

Project Types:
  microservice  - Single database, container-optimized API service
  library       - Reusable database library package
  fullstack     - Complete backend + frontend application
  api           - API-first project with OpenAPI specs

Databases:
  postgresql    - PostgreSQL with modern features (UUID, JSONB)
  mysql         - MySQL 8.0+ with JSON support
  sqlite        - SQLite for embedded applications

Examples:
  sqlc-wizard create my-service --type microservice --database postgresql
  sqlc-wizard create my-lib --type library --database sqlite --include-auth
  sqlc-wizard create my-webapp --type fullstack --database postgresql --include-frontend`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCreate(args[0], opts)
		},
	}

	// Add flags
	cmd.Flags().StringVar(&opts.ProjectType, "type", "microservice", "Project type (microservice, library, fullstack, api)")
	cmd.Flags().StringVar(&opts.Database, "database", "postgresql", "Database engine (postgresql, mysql, sqlite)")
	cmd.Flags().StringVarP(&opts.OutputDir, "output-dir", "o", ".", "Output directory for the project")
	cmd.Flags().BoolVar(&opts.IncludeAuth, "include-auth", false, "Include authentication setup")
	cmd.Flags().BoolVar(&opts.IncludeFrontend, "include-frontend", false, "Include frontend setup (for fullstack)")
	cmd.Flags().BoolVar(&opts.NonInteractive, "non-interactive", false, "Run with smart defaults, no prompts")
	cmd.Flags().BoolVarP(&opts.Force, "force", "f", false, "Overwrite existing files")

	return cmd
}

func runCreate(projectName string, opts *CreateOptions) error {
	// Validate project name
	if projectName == "" {
		return errors.New("project name cannot be empty")
	}

	// Create output directory if needed
	outputPath := filepath.Join(opts.OutputDir, projectName)
	if err := os.MkdirAll(outputPath, 0o755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Change to project directory for relative paths
	originalDir, _ := os.Getwd()
	if err := os.Chdir(outputPath); err != nil {
		return fmt.Errorf("failed to change to project directory: %w", err)
	}
	defer os.Chdir(originalDir)

	// Check if directory is empty (unless force is used)
	if !opts.Force {
		if entries, err := os.ReadDir("."); err == nil && len(entries) > 0 {
			return errors.New("directory is not empty. Use --force to overwrite")
		}
	}

	// Validate project type and database types
	projectType, err := templates.NewProjectType(opts.ProjectType)
	if err != nil {
		return fmt.Errorf("invalid project type: %w", err)
	}

	databaseType, err := templates.NewDatabaseType(opts.Database)
	if err != nil {
		return fmt.Errorf("invalid database type: %w", err)
	}

	// Create template data with smart defaults
	templateData := createTemplateData(projectName, projectType, databaseType, opts)

	// Get the appropriate template
	template, err := templates.GetTemplate(projectType)
	if err != nil {
		return fmt.Errorf("failed to get template: %w", err)
	}

	// Generate SQLC configuration
	config, err := template.Generate(templateData)
	if err != nil {
		return fmt.Errorf("failed to generate configuration: %w", err)
	}

	// Create project creator with real adapters
	fs := adapters.NewRealFileSystemAdapter()
	cli := adapters.NewRealCLIAdapter()

	creator := creators.NewProjectCreator(fs, cli)
	createConfig := &creators.CreateConfig{
		ProjectName:     projectName,
		ProjectType:     projectType,
		Database:        databaseType,
		TemplateData:    templateData,
		Config:          config,
		IncludeAuth:     opts.IncludeAuth,
		IncludeFrontend: opts.IncludeFrontend,
		Force:           opts.Force,
	}

	// Create the complete project
	if err := creator.CreateProject(context.Background(), createConfig); err != nil {
		return fmt.Errorf("failed to create project: %w", err)
	}

	// Show success message with next steps
	showCreateSuccess(projectName, projectType, databaseType, outputPath)

	return nil
}

func createTemplateData(projectName string, projectType generated.ProjectType, databaseType generated.DatabaseType, opts *CreateOptions) generated.TemplateData {
	return generated.TemplateData{
		ProjectName: projectName,
		ProjectType: projectType,

		Package: generated.PackageConfig{
			Name: "db",
			Path: "internal/db",
		},

		Database: generated.DatabaseConfig{
			Engine:      databaseType,
			URL:         "${DATABASE_URL}",
			UseManaged:  true,
			UseUUIDs:    databaseType == generated.DatabaseTypePostgreSQL,
			UseJSON:     databaseType != generated.DatabaseTypeSQLite,
			UseArrays:   databaseType == generated.DatabaseTypePostgreSQL,
			UseFullText: databaseType == generated.DatabaseTypePostgreSQL,
		},

		Output: generated.OutputConfig{
			BaseDir:    "internal/db",
			QueriesDir: "internal/db/queries",
			SchemaDir:  "internal/db/schema",
		},

		Validation: generated.ValidationConfig{
			StrictFunctions: false,
			StrictOrderBy:   false,
			EmitOptions:     generated.DefaultEmitOptions(),
			SafetyRules:     generated.DefaultSafetyRules(),
		},
	}
}

func showCreateSuccess(projectName string, projectType generated.ProjectType, databaseType generated.DatabaseType, outputPath string) {
	successStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("10")).
		Padding(1, 0)

	nextStepsStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("14")).
		Bold(true).
		MarginTop(1)

	commandStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("8")).
		PaddingLeft(2)

	fmt.Println(successStyle.Render("🚀 Successfully created SQLC project!"))
	fmt.Printf("Project: %s\n", projectName)
	fmt.Printf("Type: %s\n", projectType)
	fmt.Printf("Database: %s\n", databaseType)
	fmt.Printf("Location: %s\n", outputPath)

	fmt.Println(nextStepsStyle.Render("Next Steps:"))
	fmt.Println(commandStyle.Render("1. cd " + projectName))
	fmt.Println(commandStyle.Render("2. make setup          # Install dependencies"))
	fmt.Println(commandStyle.Render("3. make dev             # Start development environment"))
	fmt.Println(commandStyle.Render("4. sqlc generate        # Generate Go code"))
	fmt.Println(commandStyle.Render("5. make test            # Run tests"))

	fmt.Println()
	fmt.Println(successStyle.Render("✨ Your SQLC project is ready!"))
}
