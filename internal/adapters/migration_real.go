package adapters

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/migration"
	"github.com/LarsArtmann/SQLC-Wizzard/pkg/config"
	"github.com/charmbracelet/log"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/samber/lo"
)

// RealMigrationAdapter implements MigrationAdapter interface using golang-migrate
// Database drivers are imported on-demand to reduce build time.
type RealMigrationAdapter struct{}

// NewRealMigrationAdapter creates a new real migration adapter.
func NewRealMigrationAdapter() *RealMigrationAdapter {
	return &RealMigrationAdapter{}
}

// Migrate runs database migrations from a source.
func (r *RealMigrationAdapter) Migrate(ctx context.Context, source, databaseURL string) error {
	log.Info("Starting database migration", "source", source, "database", databaseURL)

	m, err := migrate.New(source, databaseURL)
	if err != nil {
		log.Error("Failed to create migration instance", "error", err)
		return fmt.Errorf("failed to create migration instance: %w", err)
	}
	defer m.Close()

	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		log.Error("Failed to get migration version", "error", err)
		return fmt.Errorf("failed to get migration version: %w", err)
	}

	if dirty {
		log.Warn("Database is in dirty state", "version", version)
		return fmt.Errorf("database is dirty at version %d", version)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Error("Migration failed", "error", err)
		return fmt.Errorf("migration failed: %w", err)
	}

	if version, _, err := m.Version(); err == nil {
		log.Info("Migration completed successfully", "final_version", version)
	} else {
		log.Info("No migrations to run")
	}

	return nil
}

// Rollback rolls back database migrations.
func (r *RealMigrationAdapter) Rollback(ctx context.Context, source, databaseURL string, steps int) error {
	log.Info("Rolling back database migrations", "source", source, "database", databaseURL, "steps", steps)

	m, err := migrate.New(source, databaseURL)
	if err != nil {
		log.Error("Failed to create migration instance", "error", err)
		return fmt.Errorf("failed to create migration instance: %w", err)
	}
	defer m.Close()

	for i := range steps {
		if err := m.Steps(-1); err != nil && err != migrate.ErrNoChange {
			log.Error("Rollback step failed", "step", i+1, "error", err)
			return fmt.Errorf("rollback step %d failed: %w", i+1, err)
		}
	}

	version, _, err := m.Version()
	switch err {
	case migrate.ErrNilVersion:
		log.Info("All migrations rolled back")
	case nil:
		log.Info("Rollback completed", "current_version", version)
	default:
		log.Error("Failed to get version after rollback", "error", err)
		return fmt.Errorf("failed to get version after rollback: %w", err)
	}

	return nil
}

// Status checks migration status.
func (r *RealMigrationAdapter) Status(ctx context.Context, source, databaseURL string) (*migration.MigrationStatus, error) {
	log.Info("Checking migration status", "source", source, "database", databaseURL)

	m, err := migrate.New(source, databaseURL)
	if err != nil {
		log.Error("Failed to create migration instance", "error", err)
		return nil, fmt.Errorf("failed to create migration instance: %w", err)
	}
	defer m.Close()

	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		log.Error("Failed to get migration version", "error", err)
		return nil, fmt.Errorf("failed to get migration version: %w", err)
	}

	// Create typed migration status
	status, err := migration.NewMigrationStatus(source, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create migration status: %w", err)
	}

	if err == migrate.ErrNilVersion {
		// No migrations applied yet
		log.Info("No migrations applied yet")
	} else {
		status.WithVersion(version)
		status.WithDirty(dirty)
	}

	log.Info("Migration status retrieved", "version", status.GetCurrentVersion(), "dirty", status.IsDirty())
	return status, nil
}

// Validate validates migration files.
func (r *RealMigrationAdapter) Validate(ctx context.Context, source string) error {
	log.Info("Validating migration files", "source", source)

	m, err := migrate.New(source, "file://tmp")
	if err != nil {
		log.Error("Failed to create migration instance for validation", "error", err)
		return fmt.Errorf("failed to create migration instance for validation: %w", err)
	}
	defer m.Close()

	version, _, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		log.Error("Migration validation failed", "error", err)
		return fmt.Errorf("migration validation failed: %w", err)
	}

	log.Info("Migration files validated", "latest_version", version)
	return nil
}

// CreateMigration creates a new migration file.
func (r *RealMigrationAdapter) CreateMigration(ctx context.Context, name, directory string) (string, error) {
	log.Info("Creating migration", "name", name, "directory", directory)

	if err := os.MkdirAll(directory, 0o755); err != nil {
		log.Error("Failed to create migrations directory", "directory", directory, "error", err)
		return "", fmt.Errorf("failed to create migrations directory: %w", err)
	}

	// Generate timestamp for migration file
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	upFile := filepath.Join(directory, fmt.Sprintf("%s_%s.up.sql", timestamp, name))
	downFile := filepath.Join(directory, fmt.Sprintf("%s_%s.down.sql", timestamp, name))

	// Create up migration
	upContent := fmt.Sprintf(`-- Migration: %s
-- Generated at: %s

-- Add your SQL statements here

`, name, timestamp)

	if err := os.WriteFile(upFile, []byte(upContent), 0o644); err != nil {
		log.Error("Failed to create up migration file", "file", upFile, "error", err)
		return "", fmt.Errorf("failed to create up migration file: %w", err)
	}

	// Create down migration
	downContent := fmt.Sprintf(`-- Migration: %s (Rollback)
-- Generated at: %s

-- Add your rollback SQL statements here

`, name, timestamp)

	if err := os.WriteFile(downFile, []byte(downContent), 0o644); err != nil {
		log.Error("Failed to create down migration file", "file", downFile, "error", err)
		return "", fmt.Errorf("failed to create down migration file: %w", err)
	}

	log.Info("Migration files created successfully", "up", upFile, "down", downFile)
	return upFile, nil
}

// MigrateSQLCConfig migrates SQLC configuration from one version/database to another.
func (r *RealMigrationAdapter) MigrateSQLCConfig(ctx context.Context, sourceConfig *config.SqlcConfig, targetDatabase generated.DatabaseType, targetVersion string) (*config.SqlcConfig, error) {
	log.Info("Migrating SQLC configuration", "target_database", targetDatabase, "target_version", targetVersion)

	// Create a copy of the source config
	newConfig := *sourceConfig

	// Update SQLC version if needed
	if targetVersion != "" {
		newConfig.Version = targetVersion
		log.Info("Updated SQLC version", "from", sourceConfig.Version, "to", targetVersion)
	}

	// Update database engine if needed
	if targetDatabase != "" && len(newConfig.SQL) > 0 {
		for i := range newConfig.SQL {
			newConfig.SQL[i].Engine = string(targetDatabase)
		}
		log.Info("Updated database engine", "to", targetDatabase)
	}

	// Update database configuration based on target database type
	if err := r.updateDatabaseConfig(&newConfig, targetDatabase); err != nil {
		return nil, fmt.Errorf("failed to update database config: %w", err)
	}

	log.Info("SQLC configuration migration completed")
	return &newConfig, nil
}

// updateDatabaseConfig updates database-specific configuration.
func (r *RealMigrationAdapter) updateDatabaseConfig(config *config.SqlcConfig, targetDatabase generated.DatabaseType) error {
	for i := range config.SQL {
		sqlConfig := &config.SQL[i]

		switch targetDatabase {
		case generated.DatabaseTypePostgreSQL:
			// PostgreSQL-specific settings
			if sqlConfig.Database != nil {
				if !lo.ContainsBy([]string{"postgres", "postgresql"}, func(db string) bool {
					return sqlConfig.Database.URI == db
				}) {
					// Default PostgreSQL connection string template
					sqlConfig.Database.URI = "postgres://user:password@localhost:5432/dbname?sslmode=disable"
				}
			}

		case generated.DatabaseTypeMySQL:
			// MySQL-specific settings
			if sqlConfig.Database != nil {
				if !lo.ContainsBy([]string{"mysql"}, func(db string) bool {
					return sqlConfig.Database.URI == db
				}) {
					// Default MySQL connection string template
					sqlConfig.Database.URI = "user:password@tcp(localhost:3306)/dbname"
				}
			}

		case generated.DatabaseTypeSQLite:
			// SQLite-specific settings
			if sqlConfig.Database != nil {
				if !lo.ContainsBy([]string{"sqlite", "sqlite3"}, func(db string) bool {
					return sqlConfig.Database.URI == db
				}) {
					// Default SQLite file path
					sqlConfig.Database.URI = "./database.db"
				}
			}
		}
	}

	return nil
}
