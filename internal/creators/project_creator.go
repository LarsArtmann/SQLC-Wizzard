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

	// TODO: Full project scaffolding is not yet implemented
	// See GitHub issues for roadmap:
	// - Database schema generation
	// - Query file generation
	// - Migration file generation
	// - Go module structure
	// - Docker configuration
	// - Makefile generation
	// - Development scripts
	// - README generation
	// - Project-specific files
	//
	// For now, ProjectCreator only generates:
	// 1. Directory structure
	// 2. sqlc.yaml configuration file
	//
	// Additional scaffolding will be added based on user feedback and demand.

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

// NOTE: Additional scaffolding methods will be implemented based on demand
// See the TODO in CreateProject for planned features
