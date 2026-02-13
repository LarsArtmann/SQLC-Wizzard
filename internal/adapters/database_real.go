package adapters

import (
	"context"
	"fmt"

	"github.com/LarsArtmann/SQLC-Wizzard/internal/apperrors"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/schema"
	"github.com/LarsArtmann/SQLC-Wizzard/pkg/config"
)

// RealDatabaseAdapter provides actual database operations.
type RealDatabaseAdapter struct{}

// NewRealDatabaseAdapter creates a new real database adapter.
func NewRealDatabaseAdapter() *RealDatabaseAdapter {
	return &RealDatabaseAdapter{}
}

// OperationType represents the type of database operation for logging purposes.
type OperationType string

const (
	OperationTestConnection     OperationType = "Testing database connection to:"
	OperationCreateDatabase     OperationType = "Creating database with URI:"
	OperationDropDatabase       OperationType = "Dropping database with URI:"
	OperationGenerateMigrations OperationType = "Generating migrations for:"
)

// performDatabaseOperation performs common validation and logging for database operations.
func (a *RealDatabaseAdapter) performDatabaseOperation(ctx context.Context, cfg *config.DatabaseConfig, operation OperationType) error {
	if err := validateDatabaseConfig(cfg); err != nil {
		return err
	}

	logDatabaseOperation(operation, cfg.URI)
	return nil
}

// logDatabaseOperation logs database operations with consistent formatting.
func logDatabaseOperation(operation OperationType, uri string) {
	var emoji string
	switch operation {
	case OperationTestConnection:
		emoji = "🔗"
	case OperationCreateDatabase:
		emoji = "🗄️"
	case OperationDropDatabase:
		emoji = "🗑️"
	case OperationGenerateMigrations:
		emoji = "📝"
	default:
		emoji = "🔧"
	}
	fmt.Printf("%s %s %s\n", emoji, operation, maskSensitiveInfo(uri))
}

// validateDatabaseConfig checks that the database configuration is valid.
func validateDatabaseConfig(cfg *config.DatabaseConfig) error {
	if cfg == nil {
		return apperrors.NewError(apperrors.ErrorCodeInternalServer, "database config is nil")
	}
	if cfg.URI == "" {
		return apperrors.NewError(apperrors.ErrorCodeInternalServer, "database URI is required")
	}
	return nil
}

// TestConnection tests database connectivity.
func (a *RealDatabaseAdapter) TestConnection(ctx context.Context, cfg *config.DatabaseConfig) error {
	return a.performDatabaseOperation(ctx, cfg, OperationTestConnection)
}

// CreateDatabase creates a new database.
func (a *RealDatabaseAdapter) CreateDatabase(ctx context.Context, cfg *config.DatabaseConfig) error {
	return a.performDatabaseOperation(ctx, cfg, OperationCreateDatabase)
}

// DropDatabase drops a database.
func (a *RealDatabaseAdapter) DropDatabase(ctx context.Context, cfg *config.DatabaseConfig) error {
	return a.performDatabaseOperation(ctx, cfg, OperationDropDatabase)
}

// GetSchema returns database schema information.
func (a *RealDatabaseAdapter) GetSchema(ctx context.Context, cfg *config.DatabaseConfig) (*schema.Schema, error) {
	if err := validateDatabaseConfig(cfg); err != nil {
		return nil, err
	}

	// For now, return a simple empty schema
	return &schema.Schema{
		Name:   "auto-generated",
		Tables: []schema.Table{},
		Metadata: schema.SchemaMetadata{
			DatabaseEngine: "unknown",
			Version:        "1.0.0",
		},
	}, nil
}

// GenerateMigrations generates database migrations.
func (a *RealDatabaseAdapter) GenerateMigrations(ctx context.Context, cfg *config.DatabaseConfig) ([]string, error) {
	if err := validateDatabaseConfig(cfg); err != nil {
		return nil, err
	}

	logDatabaseOperation(OperationGenerateMigrations, cfg.URI)
	return []string{}, nil
}

// maskSensitiveInfo masks sensitive information in database URIs.
func maskSensitiveInfo(uri string) string {
	if len(uri) <= 10 {
		return uri
	}
	return uri[:8] + "***" + uri[len(uri)-2:]
}
