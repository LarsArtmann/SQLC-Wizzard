package adapters

import (
	"context"
)

// MigrationAdapter defines the interface for database migration operations
// This follows the Adapter pattern to external dependencies
type MigrationAdapter interface {
	// Migrate runs database migrations from a source
	Migrate(ctx context.Context, source string, databaseURL string) error
	
	// Rollback rolls back database migrations
	Rollback(ctx context.Context, source string, databaseURL string, steps int) error
	
	// Status checks migration status
	Status(ctx context.Context, source string, databaseURL string) (map[string]interface{}, error)
	
	// Validate validates migration files
	Validate(ctx context.Context, source string) error
	
	// CreateMigration creates a new migration file
	CreateMigration(ctx context.Context, name string, directory string) (string, error)
}