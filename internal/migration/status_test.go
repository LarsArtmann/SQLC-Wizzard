package migration

import (
	"testing"
	"time"
)

func TestNewMigrationStatus(t *testing.T) {
	tests := []struct {
		name        string
		source      string
		databaseURL string
		expectError bool
	}{
		{
			name:        "valid input",
			source:      "file://migrations",
			databaseURL: "sqlite://test.db",
			expectError: false,
		},
		{
			name:        "empty source",
			source:      "",
			databaseURL: "sqlite://test.db",
			expectError: true,
		},
		{
			name:        "both inputs valid",
			source:      "file://migrations",
			databaseURL: "postgres://localhost/test",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, err := NewMigrationStatus(tt.source, tt.databaseURL)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if status.Source != tt.source {
				t.Errorf("Expected source %s, got %s", tt.source, status.Source)
			}

			if status.DatabaseURL != tt.databaseURL {
				t.Errorf("Expected database URL %s, got %s", tt.databaseURL, status.DatabaseURL)
			}

			if status.CurrentVersion != nil {
				t.Error("Expected current version to be nil initially")
			}

			if status.Dirty {
				t.Error("Expected dirty to be false initially")
			}

			if len(status.Migrations) != 0 {
				t.Error("Expected migrations to be empty initially")
			}

			if status.CheckedAt.IsZero() {
				t.Error("Expected CheckedAt to be set")
			}
		})
	}
}

func TestMigrationStatus_WithVersion(t *testing.T) {
	status, err := NewMigrationStatus("test", "test")
	if err != nil {
		t.Fatalf("Failed to create MigrationStatus: %v", err)
	}

	version := uint(42)
	status.WithVersion(version)

	if status.GetCurrentVersion() == nil || *status.GetCurrentVersion() != version {
		t.Error("WithVersion should set current version")
	}

	if status.Dirty {
		t.Error("WithVersion should not modify other fields")
	}
}

func TestMigrationStatus_WithDirty(t *testing.T) {
	status, err := NewMigrationStatus("test", "test")
	if err != nil {
		t.Fatalf("Failed to create MigrationStatus: %v", err)
	}

	status.WithDirty(true)

	if !status.IsDirty() {
		t.Error("WithDirty should set dirty status")
	}

	if status.GetCurrentVersion() != nil {
		t.Error("WithDirty should not modify other fields")
	}
}

func TestMigrationStatus_WithMigrations(t *testing.T) {
	status, err := NewMigrationStatus("test", "test")
	if err != nil {
		t.Fatalf("Failed to create MigrationStatus: %v", err)
	}

	migrations := []Migration{
		{Version: 1, Applied: true, Name: "initial"},
		{Version: 2, Applied: false, Name: "add_users"},
	}

	status.WithMigrations(migrations)

	if status.GetMigrationCount() != 2 {
		t.Error("WithMigrations should set migrations")
	}

	if status.GetAppliedMigrations() != 1 {
		t.Error("GetAppliedMigrations should return correct count")
	}

	if status.GetPendingMigrations() != 1 {
		t.Error("GetPendingMigrations should return correct count")
	}
}

func TestMigrationStatus_HelperMethods(t *testing.T) {
	now := time.Now()
	appliedAt := &now

	migrations := []Migration{
		{Version: 1, Applied: true, AppliedAt: appliedAt, Name: "initial"},
		{Version: 2, Applied: false, Name: "add_users"},
		{Version: 3, Applied: true, AppliedAt: appliedAt, Name: "add_posts"},
	}

	status, err := NewMigrationStatus("test", "test")
	if err != nil {
		t.Fatalf("Failed to create MigrationStatus: %v", err)
	}

	status.WithVersion(3)
	status.WithDirty(false)
	status.WithMigrations(migrations)

	// Test IsDirty
	if status.IsDirty() {
		t.Error("Expected IsDirty to return false")
	}

	// Test GetCurrentVersion
	current := status.GetCurrentVersion()
	if current == nil || *current != 3 {
		t.Error("GetCurrentVersion should return 3")
	}

	// Test GetMigrationCount
	if status.GetMigrationCount() != 3 {
		t.Error("GetMigrationCount should return 3")
	}

	// Test GetAppliedMigrations
	if status.GetAppliedMigrations() != 2 {
		t.Error("GetAppliedMigrations should return 2")
	}

	// Test GetPendingMigrations
	if status.GetPendingMigrations() != 1 {
		t.Error("GetPendingMigrations should return 1")
	}
}

func TestValidationError(t *testing.T) {
	err := &ValidationError{
		Field:   "test_field",
		Message: "test message",
	}

	if err.Error() != "test message" {
		t.Errorf("Expected 'test message', got '%s'", err.Error())
	}

	if err.Field != "test_field" {
		t.Errorf("Expected field 'test_field', got '%s'", err.Field)
	}
}