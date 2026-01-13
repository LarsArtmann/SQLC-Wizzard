package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/LarsArtmann/SQLC-Wizzard/pkg/config"
	"github.com/charmbracelet/log"
)

// writeMigrationResult writes migrated configuration to file.
func writeMigrationResult(cfg *config.SqlcConfig, destination string, force bool) error {
	// Check if file exists and force is not set
	if _, err := os.Stat(destination); err == nil && !force {
		return &MigrationError{
			Code:    "FILE_EXISTS",
			Message: "Destination file exists and force flag not set: " + destination,
		}
	}

	// Create destination directory if needed
	if err := os.MkdirAll(filepath.Dir(destination), 0o755); err != nil {
		return &MigrationError{
			Code:    "DIRECTORY_CREATION_FAILED",
			Message: fmt.Sprintf("Failed to create destination directory: %v", err),
		}
	}

	// Write configuration
	if err := config.WriteFileFormatted(cfg, destination); err != nil {
		return &MigrationError{
			Code:    "WRITE_FAILED",
			Message: fmt.Sprintf("Failed to write configuration: %v", err),
		}
	}

	log.Info("Migration configuration written", "destination", destination)
	return nil
}
