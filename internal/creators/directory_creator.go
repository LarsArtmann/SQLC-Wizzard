package creators

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/adapters"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/apperrors"
)

// DirectoryCreator handles directory structure creation
// Implements FileSystemCreator interface with clean separation of concerns.
type DirectoryCreator struct {
	fs  adapters.FileSystemAdapter
	cli adapters.CLIAdapter
}

// NewDirectoryCreator creates a new directory creator.
func NewDirectoryCreator(fs adapters.FileSystemAdapter, cli adapters.CLIAdapter) *DirectoryCreator {
	return &DirectoryCreator{
		fs:  fs,
		cli: cli,
	}
}

// Create creates directory structure for the project
// Implements Creator[CreatorConfig] interface.
func (dc *DirectoryCreator) Create(ctx context.Context, config CreatorConfig) error {
	dc.cli.Println("📁 Creating directory structure...")

	dirs := dc.getProjectDirectories(config)

	for _, dir := range dirs {
		fullPath := filepath.Join(config.OutputPath, dir)
		if err := dc.fs.MkdirAll(ctx, fullPath, 0o755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
		dc.cli.Println("Created directory: " + fullPath)
	}

	return nil
}

// CanHandle returns true if this creator can handle the config.
func (dc *DirectoryCreator) CanHandle(config CreatorConfig) bool {
	return config.OutputPath != "" && config.ProjectName != ""
}

// Dependencies returns any creators that must run first.
func (dc *DirectoryCreator) Dependencies() []string {
	return []string{} // No dependencies - runs first
}

// FileSystem returns the file system adapter.
func (dc *DirectoryCreator) FileSystem() adapters.FileSystemAdapter {
	return dc.fs
}

// SetFileSystem sets the file system adapter (for testing).
func (dc *DirectoryCreator) SetFileSystem(fs adapters.FileSystemAdapter) {
	dc.fs = fs
}

// getProjectDirectories returns directory list based on project type
// Uses switch to ensure all cases are handled.
func (dc *DirectoryCreator) getProjectDirectories(config CreatorConfig) []string {
	// Base directories required for all project types
	baseDirs := []string{
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

	// Add project-type-specific directories
	projectDirs := dc.getProjectTypeDirectories(config.ProjectType)

	return append(baseDirs, projectDirs...)
}

// getProjectTypeDirectories returns directories specific to project type
// Ensures all project types have proper directory structure.
func (dc *DirectoryCreator) getProjectTypeDirectories(projectType generated.ProjectType) []string {
	switch projectType {
	case generated.ProjectTypeMicroservice:
		return []string{
			"api",
			"internal/api",
			"internal/handlers",
			"internal/middleware",
			"pkg/auth",
			"pkg/logger",
		}

	case generated.ProjectTypeEnterprise:
		return []string{
			"api",
			"internal/api",
			"internal/handlers",
			"internal/middleware",
			"internal/audit",
			"pkg/auth",
			"pkg/logger",
			"pkg/monitoring",
			"pkg/security",
		}

	case generated.ProjectTypeAPIFirst:
		return []string{
			"api",
			"internal/api",
			"internal/handlers",
			"internal/middleware",
			"pkg/auth",
			"pkg/ratelimiter",
		}

	case generated.ProjectTypeHobby:
		return []string{
			"internal/handlers",
			"pkg/utils",
		}

	case generated.ProjectTypeAnalytics:
		return []string{
			"internal/analytics",
			"internal/processors",
			"internal/aggregators",
			"pkg/metrics",
		}

	case generated.ProjectTypeLibrary:
		return []string{
			"internal",
			"internal/users",
			"internal/orders",
			"internal/products",
			"internal/payments",
			"pkg/shared",
		}

	// TODO: Add remaining project types as they are implemented
	// case generated.ProjectTypeFullstack:
	// case generated.ProjectTypeCLI:
	// case generated.ProjectTypePlugin:

	default:
		// Fallback to microservice pattern for unknown types
		return []string{
			"api",
			"internal/api",
			"internal/handlers",
		}
	}
}

// Validate ensures config is valid for directory creation.
func (dc *DirectoryCreator) Validate(config CreatorConfig) error {
	if config.ProjectName == "" {
		return apperrors.NewError(apperrors.ErrorCodeInternalServer, "project name cannot be empty")
	}

	if config.OutputPath == "" {
		return apperrors.NewError(apperrors.ErrorCodeInternalServer, "output path cannot be empty")
	}

	if !config.ProjectType.IsValid() {
		return fmt.Errorf("invalid project type: %s", string(config.ProjectType))
	}

	return nil
}
