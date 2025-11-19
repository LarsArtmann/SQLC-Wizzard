package creators

import (
	"context"
	"fmt"

	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/adapters"
	"github.com/LarsArtmann/SQLC-Wizzard/pkg/config"
)

// CreateConfig contains configuration for project creation
type CreateConfig struct {
	ProjectName     string
	ProjectType     generated.ProjectType
	Database        generated.DatabaseType
	TemplateData    generated.TemplateData
	Config          *config.SqlcConfig
	IncludeAuth     bool
	IncludeFrontend bool
	Force           bool
}

// ProjectCreator handles creating complete project structures
type ProjectCreator struct {
	fs  adapters.FileSystemAdapter
	cli adapters.CLIAdapter
}

// NewProjectCreator creates a new project creator
func NewProjectCreator(fs adapters.FileSystemAdapter, cli adapters.CLIAdapter) *ProjectCreator {
	return &ProjectCreator{
		fs:  fs,
		cli: cli,
	}
}

// CreateProject creates a complete project structure
func (pc *ProjectCreator) CreateProject(ctx context.Context, config *CreateConfig) error {
	pc.cli.Println("🏗️  Creating project structure...")

	// Create directory structure
	if err := pc.createDirectoryStructure(ctx, config); err != nil {
		return fmt.Errorf("failed to create directory structure: %w", err)
	}

	// Generate sqlc.yaml
	if err := pc.generateSQLCConfig(ctx, config); err != nil {
		return fmt.Errorf("failed to generate sqlc.yaml: %w", err)
	}

	// TODO: Complete implementation of project scaffolding
	// The following methods are not yet implemented
	_ = ctx // Use ctx to avoid unused variable error

	// Generate database schemas
	// if err := pc.generateSchemas(ctx, config); err != nil {
	// 	return fmt.Errorf("failed to generate schemas: %w", err)
	// }

	// Generate query files
	// if err := pc.generateQueries(ctx, config); err != nil {
	// 	return fmt.Errorf("failed to generate queries: %w", err)
	// }

	// Generate migration files
	// if err := pc.generateMigrations(ctx, config); err != nil {
	// 	return fmt.Errorf("failed to generate migrations: %w", err)
	// }

	// Generate Go module and basic structure
	// if err := pc.generateGoStructure(ctx, config); err != nil {
	// 	return fmt.Errorf("failed to generate Go structure: %w", err)
	// }

	// Generate Docker configuration
	// if err := pc.generateDockerConfig(ctx, config); err != nil {
	// 	return fmt.Errorf("failed to generate Docker configuration: %w", err)
	// }

	// Generate Makefile
	// if err := pc.generateMakefile(ctx, config); err != nil {
	// 	return fmt.Errorf("failed to generate Makefile: %w", err)
	// }

	// Generate development scripts
	// if err := pc.generateDevScripts(ctx, config); err != nil {
	// 	return fmt.Errorf("failed to generate development scripts: %w", err)
	// }

	// Generate README
	// if err := pc.generateREADME(ctx, config); err != nil {
	// 	return fmt.Errorf("failed to generate README: %w", err)
	// }

	// Generate project-specific files
	// if err := pc.generateProjectSpecificFiles(ctx, config); err != nil {
	// 	return fmt.Errorf("failed to generate project-specific files: %w", err)
	// }

	return nil
}

// createDirectoryStructure creates the basic directory structure
func (pc *ProjectCreator) createDirectoryStructure(ctx context.Context, config *CreateConfig) error {
	pc.cli.Println("📁 Creating directory structure...")

	dirs := []string{
		"db/schema",
		"db/migrations",
		"internal/db",
		"internal/db/queries",
		"cmd/server",
		"pkg/config",
		"scripts",
		"test",
		"docs",
	}

	// Add project-specific directories based on project type
	// Note: Some project types may not be in generated types yet
	switch config.ProjectType {
	case generated.ProjectTypeMicroservice:
		dirs = append(dirs, "api", "internal/api", "internal/handlers")
	// TODO: Add other project types when generated types are complete
	// case generated.ProjectTypeFullstack:
	// 	dirs = append(dirs, "web", "web/src", "web/public", "internal/api")
	// case generated.ProjectTypeAPIFirst:
	// 	dirs = append(dirs, "api", "internal/api", "internal/handlers")
	// case generated.ProjectTypeLibrary:
	// 	dirs = append(dirs, "examples", "internal/testutil")
	}

	for _, dir := range dirs {
		if err := pc.fs.MkdirAll(ctx, dir, 0o755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}

// generateSQLCConfig generates the sqlc.yaml file
func (pc *ProjectCreator) generateSQLCConfig(ctx context.Context, cfg *CreateConfig) error {
	pc.cli.Println("⚙️  Generating sqlc.yaml...")

	// Convert config to YAML using the marshaller
	yamlContent, err := config.Marshal(cfg.Config)
	if err != nil {
		return fmt.Errorf("failed to convert config to YAML: %w", err)
	}

	return pc.fs.WriteFile(ctx, "sqlc.yaml", yamlContent, 0o644)
}

// TODO: The following methods are not yet implemented
// They are commented out to allow the project to build
// Uncomment and implement when ready

// generateSchemas generates database schema files
// func (pc *ProjectCreator) generateSchemas(ctx context.Context, config *CreateConfig) error {
// 	pc.cli.Println("🗄️  Generating database schemas...")
// 	// Implementation needed
// 	return nil
// }

// generateQueries generates SQL query files
// func (pc *ProjectCreator) generateQueries(ctx context.Context, config *CreateConfig) error {
// 	pc.cli.Println("🔍 Generating query files...")
// 	// Implementation needed
// 	return nil
// }

// generateMigrations generates migration files
// func (pc *ProjectCreator) generateMigrations(ctx context.Context, config *CreateConfig) error {
// 	pc.cli.Println("🔄 Generating migration files...")
// 	// Implementation needed
// 	return nil
// }

// generateGoStructure generates Go module and basic structure
// func (pc *ProjectCreator) generateGoStructure(ctx context.Context, config *CreateConfig) error {
// 	pc.cli.Println("🐹 Generating Go structure...")
// 	// Implementation needed
// 	return nil
// }

// generateDockerConfig generates Docker configuration
// func (pc *ProjectCreator) generateDockerConfig(ctx context.Context, config *CreateConfig) error {
// 	pc.cli.Println("🐳 Generating Docker configuration...")
// 	// Implementation needed
// 	return nil
// }

// generateMakefile generates Makefile with common tasks
// func (pc *ProjectCreator) generateMakefile(ctx context.Context, config *CreateConfig) error {
// 	pc.cli.Println("🔨 Generating Makefile...")
// 	// Implementation needed
// 	return nil
// }

// generateDevScripts generates development scripts
// func (pc *ProjectCreator) generateDevScripts(ctx context.Context, config *CreateConfig) error {
// 	pc.cli.Println("📜 Generating development scripts...")
// 	// Implementation needed
// 	return nil
// }

// generateREADME generates README.md
// func (pc *ProjectCreator) generateREADME(ctx context.Context, config *CreateConfig) error {
// 	pc.cli.Println("📖 Generating README...")
// 	// Implementation needed
// 	return nil
// }

// generateProjectSpecificFiles generates project-specific additional files
// func (pc *ProjectCreator) generateProjectSpecificFiles(ctx context.Context, config *CreateConfig) error {
// 	// Implementation needed
// 	return nil
// }
