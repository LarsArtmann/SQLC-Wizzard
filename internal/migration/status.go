package migration

import "time"

// MigrationStatus represents a strongly-typed migration status
// This eliminates map[string]interface{} usage and provides type safety
type MigrationStatus struct {
	CurrentVersion *uint       `json:"current_version,omitempty"`
	Dirty          bool        `json:"dirty"`
	Migrations     []Migration `json:"migrations"`
	CheckedAt      time.Time   `json:"checked_at"`
	Source         string      `json:"source"`
	DatabaseURL    string      `json:"database_url,omitempty"`
}

// Migration represents a single migration with type safety
type Migration struct {
	Version   uint       `json:"version"`
	Applied   bool       `json:"applied"`
	AppliedAt *time.Time `json:"applied_at,omitempty"`
	Name      string     `json:"name"`
	Dirty     bool       `json:"dirty"`
	UpFile    string     `json:"up_file"`
	DownFile  string     `json:"down_file,omitempty"`
}

// NewMigrationStatus creates a new MigrationStatus with validation
func NewMigrationStatus(source, databaseURL string) (*MigrationStatus, error) {
	if source == "" {
		return nil, &ValidationError{
			Field:   "source",
			Message: "source cannot be empty",
		}
	}

	return &MigrationStatus{
		CurrentVersion: nil,
		Dirty:          false,
		Migrations:     []Migration{},
		CheckedAt:      time.Now(),
		Source:         source,
		DatabaseURL:    databaseURL,
	}, nil
}

// WithVersion sets the current version
func (ms *MigrationStatus) WithVersion(version uint) {
	ms.CurrentVersion = &version
}

// WithDirty sets the dirty status
func (ms *MigrationStatus) WithDirty(dirty bool) {
	ms.Dirty = dirty
}

// WithMigrations sets the migrations list
func (ms *MigrationStatus) WithMigrations(migrations []Migration) {
	ms.Migrations = migrations
}

// IsDirty returns whether the database is in a dirty state
// This checks BOTH the database-level dirty flag AND any per-migration dirty flags
// to prevent split brain where flags are inconsistent
func (ms *MigrationStatus) IsDirty() bool {
	// Check database-level dirty flag
	if ms.Dirty {
		return true
	}

	// Check if any individual migration is dirty
	for _, mig := range ms.Migrations {
		if mig.Dirty {
			return true
		}
	}

	return false
}

// GetCurrentVersion returns the current migration version
func (ms *MigrationStatus) GetCurrentVersion() *uint {
	return ms.CurrentVersion
}

// GetMigrationCount returns the total number of migrations
// Returns uint because migration counts cannot be negative
func (ms *MigrationStatus) GetMigrationCount() uint {
	return uint(len(ms.Migrations))
}

// GetAppliedMigrations returns count of applied migrations
// Returns uint because migration counts cannot be negative
func (ms *MigrationStatus) GetAppliedMigrations() uint {
	var count uint
	for _, mig := range ms.Migrations {
		if mig.Applied {
			count++
		}
	}
	return count
}

// GetPendingMigrations returns count of pending migrations
// Returns uint because migration counts cannot be negative
func (ms *MigrationStatus) GetPendingMigrations() uint {
	total := ms.GetMigrationCount()
	applied := ms.GetAppliedMigrations()

	// Defensive check: prevent uint underflow if applied > total (inconsistent state)
	if applied >= total {
		return 0 // Inconsistent state: no pending migrations possible
	}

	return total - applied
}

// ValidationError represents migration-specific validation errors
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e *ValidationError) Error() string {
	return e.Message
}
