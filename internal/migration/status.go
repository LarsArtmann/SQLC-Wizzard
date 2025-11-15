package migration

import (
	"time"
)

// MigrationStatus represents the current status of database migrations
// Type-safe replacement for map[string]interface{}
type MigrationStatus struct {
	CurrentVersion *string    `json:"current_version,omitempty"`
	Dirty          bool        `json:"dirty"`
	Migrations     []Migration `json:"migrations"`
	LastRun        *time.Time  `json:"last_run,omitempty"`
	NextRun        *time.Time  `json:"next_run,omitempty"`
}

// Migration represents an individual migration entry
type Migration struct {
	Version     string    `json:"version"`
	Description string    `json:"description"`
	AppliedAt   *time.Time `json:"applied_at,omitempty"`
	Dirty       bool      `json:"dirty"`
}

// NewMigrationStatus creates a new MigrationStatus with validation
// Prevents invalid states through constructor validation
func NewMigrationStatus(currentVersion *string, dirty bool, migrations []Migration) (*MigrationStatus, error) {
	if len(migrations) > 1000 {
		return nil, &MigrationError{
			Code:    "TOO_MANY_MIGRATIONS",
			Message: "Migration count exceeds reasonable limit of 1000",
		}
	}
	
	return &MigrationStatus{
		CurrentVersion: currentVersion,
		Dirty:          dirty,
		Migrations:     migrations,
	}, nil
}

// IsDirty returns true if migration state is dirty
func (ms *MigrationStatus) IsDirty() bool {
	return ms.Dirty
}

// GetCurrentVersion safely returns current version
func (ms *MigrationStatus) GetCurrentVersion() *string {
	return ms.CurrentVersion
}

// GetAppliedMigrations returns only applied migrations
func (ms *MigrationStatus) GetAppliedMigrations() []Migration {
	var applied []Migration
	for _, migration := range ms.Migrations {
		if migration.AppliedAt != nil {
			applied = append(applied, migration)
		}
	}
	return applied
}

// GetPendingMigrations returns only pending migrations
func (ms *MigrationStatus) GetPendingMigrations() []Migration {
	var pending []Migration
	for _, migration := range ms.Migrations {
		if migration.AppliedAt == nil {
			pending = append(pending, migration)
		}
	}
	return pending
}

// MigrationError represents migration-specific errors
type MigrationError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *MigrationError) Error() string {
	return e.Message
}