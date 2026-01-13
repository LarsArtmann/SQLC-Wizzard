package commands

import (
	"github.com/LarsArtmann/SQLC-Wizzard/internal/adapters"
)

// MigrationError represents migration-specific errors.
type MigrationError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *MigrationError) Error() string {
	return e.Message
}

// MigrationResult represents the result of a migration operation.
type MigrationResult struct {
	Success     bool     `json:"success"`
	Message     string   `json:"message"`
	Source      string   `json:"source,omitempty"`
	Destination string   `json:"destination,omitempty"`
	Changes     []string `json:"changes,omitempty"`
}

// SQLCMigrator handles SQLC configuration migrations.
type SQLCMigrator struct {
	adapter adapters.MigrationAdapter
}

// NewSQLCMigrator creates a new SQLC migrator.
func NewSQLCMigrator(adapter adapters.MigrationAdapter) *SQLCMigrator {
	return &SQLCMigrator{
		adapter: adapter,
	}
}
