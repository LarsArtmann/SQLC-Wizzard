// Package adapters provides adapter pattern interfaces for external dependencies
package adapters

import (
	"context"
	"io/fs"

	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/schema"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/templates"
	"github.com/LarsArtmann/SQLC-Wizzard/pkg/config"
)

// SQLCAdapter defines interface for sqlc operations
// This isolates direct sqlc usage and enables testing
type SQLCAdapter interface {
	// Generate generates Go code from SQL files
	Generate(ctx context.Context, cfg *config.SqlcConfig) error

	// Validate validates sqlc configuration
	Validate(ctx context.Context, cfg *config.SqlcConfig) error

	// Version returns sqlc version
	Version(ctx context.Context) (string, error)

	// CheckInstallation checks if sqlc is installed
	CheckInstallation(ctx context.Context) error
}

// DatabaseAdapter defines interface for database operations
// This isolates database-specific logic and enables testing
type DatabaseAdapter interface {
	// TestConnection tests database connectivity
	TestConnection(ctx context.Context, cfg *config.DatabaseConfig) error

	// CreateDatabase creates a new database
	CreateDatabase(ctx context.Context, cfg *config.DatabaseConfig) error

	// DropDatabase drops a database
	DropDatabase(ctx context.Context, cfg *config.DatabaseConfig) error

	// GetSchema returns database schema information
	GetSchema(ctx context.Context, cfg *config.DatabaseConfig) (*schema.Schema, error)

	// GenerateMigrations generates database migrations
	GenerateMigrations(ctx context.Context, cfg *config.DatabaseConfig) ([]string, error)
}

// CLIAdapter defines interface for CLI operations
// This isolates CLI-specific logic and enables testing
type CLIAdapter interface {
	// RunCommand executes a CLI command
	RunCommand(ctx context.Context, cmd string, args ...string) (string, error)

	// CheckCommand checks if a command is available
	CheckCommand(ctx context.Context, cmd string) error

	// GetVersion returns version of a CLI tool
	GetVersion(ctx context.Context, cmd string) (string, error)

	// Install installs a CLI tool
	Install(ctx context.Context, cmd string) error

	// Println prints a message to output
	Println(message string) error
}

// TemplateAdapter defines interface for template operations
// This isolates template logic and enables testing
type TemplateAdapter interface {
	// GetTemplate retrieves a template by type
	GetTemplate(projectType generated.ProjectType) (templates.Template, error)

	// GenerateConfig generates configuration from template data
	GenerateConfig(ctx context.Context, data generated.TemplateData) (*config.SqlcConfig, error)

	// GenerateFiles generates files from template
	GenerateFiles(ctx context.Context, data generated.TemplateData, outputDir string) ([]string, error)

	// ValidateTemplateData validates template data
	ValidateTemplateData(ctx context.Context, data generated.TemplateData) error

	// ListTemplates returns all available templates
	ListTemplates(ctx context.Context) ([]templates.Template, error)
}

// FileSystemAdapter defines interface for file system operations
// This isolates file system logic and enables testing
type FileSystemAdapter interface {
	// ReadFile reads a file
	ReadFile(ctx context.Context, path string) ([]byte, error)

	// WriteFile writes a file
	WriteFile(ctx context.Context, path string, data []byte, perm fs.FileMode) error

	// CreateDirectory creates a directory
	CreateDirectory(ctx context.Context, path string, perm fs.FileMode) error

	// MkdirAll creates a directory and all parent directories
	MkdirAll(ctx context.Context, path string, perm fs.FileMode) error

	// Exists checks if a path exists
	Exists(ctx context.Context, path string) (bool, error)

	// ListFiles lists files in a directory
	ListFiles(ctx context.Context, dir string) ([]string, error)

	// Remove removes a file or directory
	Remove(ctx context.Context, path string) error

	// Copy copies a file or directory
	Copy(ctx context.Context, src, dst string) error

	// TempDir creates a temporary directory
	TempDir(ctx context.Context, prefix string) (string, error)
}
