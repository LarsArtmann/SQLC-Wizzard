// Package adapters provides concrete implementations of adapter interfaces
package adapters

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/templates"
	"github.com/LarsArtmann/SQLC-Wizzard/pkg/config"
)

// RealSQLCAdapter provides actual sqlc operations
type RealSQLCAdapter struct{}

// NewRealSQLCAdapter creates a new real SQLC adapter
func NewRealSQLCAdapter() *RealSQLCAdapter {
	return &RealSQLCAdapter{}
}

// Generate generates Go code from SQL files
func (a *RealSQLCAdapter) Generate(ctx context.Context, cfg *config.SqlcConfig) error {
	cmd := exec.CommandContext(ctx, "sqlc", "generate")
	cmd.Dir = filepath.Dir(".")

	return cmd.Run()
}

// Validate validates sqlc configuration
func (a *RealSQLCAdapter) Validate(ctx context.Context, cfg *config.SqlcConfig) error {
	cmd := exec.CommandContext(ctx, "sqlc", "validate")
	cmd.Dir = filepath.Dir(".")

	return cmd.Run()
}

// Version returns sqlc version
func (a *RealSQLCAdapter) Version(ctx context.Context) (string, error) {
	cmd := exec.CommandContext(ctx, "sqlc", "version")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(output), nil
}

// CheckInstallation checks if sqlc is installed
func (a *RealSQLCAdapter) CheckInstallation(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, "which", "sqlc")
	return cmd.Run()
}

// RealDatabaseAdapter provides actual database operations
type RealDatabaseAdapter struct{}

// NewRealDatabaseAdapter creates a new real database adapter
func NewRealDatabaseAdapter() *RealDatabaseAdapter {
	return &RealDatabaseAdapter{}
}

// TestConnection tests database connectivity
func (a *RealDatabaseAdapter) TestConnection(ctx context.Context, cfg *config.DatabaseConfig) error {
	// Implementation would depend on database type
	// For now, just check if URI is provided
	if cfg.URI == "" {
		return fmt.Errorf("database URI is required")
	}
	return nil
}

// CreateDatabase creates a new database
func (a *RealDatabaseAdapter) CreateDatabase(ctx context.Context, cfg *config.DatabaseConfig) error {
	// Implementation would use appropriate database driver
	return fmt.Errorf("database creation not yet implemented")
}

// DropDatabase drops a database
func (a *RealDatabaseAdapter) DropDatabase(ctx context.Context, cfg *config.DatabaseConfig) error {
	// Implementation would use appropriate database driver
	return fmt.Errorf("database dropping not yet implemented")
}

// GetSchema returns database schema information
func (a *RealDatabaseAdapter) GetSchema(ctx context.Context, cfg *config.DatabaseConfig) (any, error) {
	// Implementation would query database schema
	return nil, fmt.Errorf("schema retrieval not yet implemented")
}

// GenerateMigrations generates database migrations
func (a *RealDatabaseAdapter) GenerateMigrations(ctx context.Context, cfg *config.DatabaseConfig) ([]string, error) {
	// Implementation would generate migration files
	return nil, fmt.Errorf("migration generation not yet implemented")
}

// RealCLIAdapter provides actual CLI operations
type RealCLIAdapter struct{}

// NewRealCLIAdapter creates a new real CLI adapter
func NewRealCLIAdapter() *RealCLIAdapter {
	return &RealCLIAdapter{}
}

// RunCommand executes a CLI command
func (a *RealCLIAdapter) RunCommand(ctx context.Context, cmd string, args ...string) (string, error) {
	c := exec.CommandContext(ctx, cmd, args...)
	output, err := c.Output()
	if err != nil {
		return "", err
	}

	return string(output), nil
}

// CheckCommand checks if a command is available
func (a *RealCLIAdapter) CheckCommand(ctx context.Context, cmd string) error {
	_, err := exec.LookPath(cmd)
	return err
}

// GetVersion returns version of a CLI tool
func (a *RealCLIAdapter) GetVersion(ctx context.Context, cmd string) (string, error) {
	c := exec.CommandContext(ctx, cmd, "version")
	output, err := c.Output()
	if err != nil {
		return "", err
	}

	return string(output), nil
}

// Install installs a CLI tool
func (a *RealCLIAdapter) Install(ctx context.Context, cmd string) error {
	// Implementation would use appropriate package manager
	return fmt.Errorf("CLI installation not yet implemented")
}

// RealTemplateAdapter provides actual template operations
type RealTemplateAdapter struct {
	registry *templates.Registry
}

// NewRealTemplateAdapter creates a new real template adapter
func NewRealTemplateAdapter() *RealTemplateAdapter {
	return &RealTemplateAdapter{
		registry: templates.NewRegistry(),
	}
}

// GetTemplate retrieves a template by type
func (a *RealTemplateAdapter) GetTemplate(projectType generated.ProjectType) (templates.Template, error) {
	return a.registry.Get(projectType)
}

// GenerateConfig generates configuration from template data
func (a *RealTemplateAdapter) GenerateConfig(ctx context.Context, data generated.TemplateData) (*config.SqlcConfig, error) {
	template, err := a.GetTemplate(data.ProjectType)
	if err != nil {
		return nil, err
	}

	return template.Generate(data)
}

// GenerateFiles generates files from template
func (a *RealTemplateAdapter) GenerateFiles(ctx context.Context, data generated.TemplateData, outputDir string) ([]string, error) {
	// Implementation would generate template files
	return nil, fmt.Errorf("template file generation not yet implemented")
}

// ValidateTemplateData validates template data
func (a *RealTemplateAdapter) ValidateTemplateData(ctx context.Context, data generated.TemplateData) error {
	// Validate project type
	if !data.ProjectType.IsValid() {
		return fmt.Errorf("invalid project type: %s", data.ProjectType)
	}

	// Validate database type
	if !data.Database.Engine.IsValid() {
		return fmt.Errorf("invalid database type: %s", data.Database.Engine)
	}

	// Validate required fields
	if data.Package.Name == "" {
		return fmt.Errorf("package name is required")
	}

	return nil
}

// ListTemplates returns all available templates
func (a *RealTemplateAdapter) ListTemplates(ctx context.Context) ([]templates.Template, error) {
	// Implementation would list all templates in registry
	return nil, fmt.Errorf("template listing not yet implemented")
}

// RealFileSystemAdapter provides actual file system operations
type RealFileSystemAdapter struct{}

// NewRealFileSystemAdapter creates a new real file system adapter
func NewRealFileSystemAdapter() *RealFileSystemAdapter {
	return &RealFileSystemAdapter{}
}

// ReadFile reads a file
func (a *RealFileSystemAdapter) ReadFile(ctx context.Context, path string) ([]byte, error) {
	return os.ReadFile(path)
}

// WriteFile writes a file
func (a *RealFileSystemAdapter) WriteFile(ctx context.Context, path string, data []byte, perm fs.FileMode) error {
	return os.WriteFile(path, data, perm)
}

// CreateDirectory creates a directory
func (a *RealFileSystemAdapter) CreateDirectory(ctx context.Context, path string, perm fs.FileMode) error {
	return os.MkdirAll(path, perm)
}

// Exists checks if a path exists
func (a *RealFileSystemAdapter) Exists(ctx context.Context, path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// ListFiles lists files in a directory
func (a *RealFileSystemAdapter) ListFiles(ctx context.Context, dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var files []string
	for _, entry := range entries {
		if !entry.IsDir() {
			files = append(files, entry.Name())
		}
	}

	return files, nil
}

// Remove removes a file or directory
func (a *RealFileSystemAdapter) Remove(ctx context.Context, path string) error {
	return os.RemoveAll(path)
}

// Copy copies a file or directory
func (a *RealFileSystemAdapter) Copy(ctx context.Context, src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	return os.WriteFile(dst, data, 0644)
}

// TempDir creates a temporary directory
func (a *RealFileSystemAdapter) TempDir(ctx context.Context, prefix string) (string, error) {
	return os.MkdirTemp("", prefix)
}
