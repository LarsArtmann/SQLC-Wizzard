package templates

import (
	"github.com/LarsArtmann/SQLC-Wizzard/internal/errors"
)

// NewProjectType creates a validated ProjectType.
// This is a smart constructor that ensures only valid project types are created.
func NewProjectType(s string) (ProjectType, error) {
	pt := ProjectType(s)

	// Validate against known project types
	switch pt {
	case ProjectTypeHobby,
		ProjectTypeMicroservice,
		ProjectTypeEnterprise,
		ProjectTypeAPIFirst,
		ProjectTypeAnalytics,
		ProjectTypeTesting,
		ProjectTypeMultiTenant,
		ProjectTypeLibrary:
		return pt, nil
	default:
		return "", errors.NewValidationErrorf(
			"project_type",
			"invalid project type: %s (must be one of: hobby, microservice, enterprise, api-first, analytics, testing, multi-tenant, library)",
			s,
		).Error
	}
}

// MustNewProjectType creates a ProjectType or panics if invalid.
// Use only when the value is guaranteed to be valid (e.g., constants).
func MustNewProjectType(s string) ProjectType {
	pt, err := NewProjectType(s)
	if err != nil {
		panic(err)
	}
	return pt
}

// NewDatabaseType creates a validated DatabaseType.
// This is a smart constructor that ensures only supported databases are created.
func NewDatabaseType(s string) (DatabaseType, error) {
	dt := DatabaseType(s)

	// Validate against supported database types
	switch dt {
	case DatabaseTypePostgreSQL,
		DatabaseTypeMySQL,
		DatabaseTypeSQLite:
		return dt, nil
	default:
		return "", errors.NewValidationErrorf(
			"database_type",
			"invalid database type: %s (must be one of: postgresql, mysql, sqlite)",
			s,
		).Error
	}
}

// MustNewDatabaseType creates a DatabaseType or panics if invalid.
// Use only when the value is guaranteed to be valid (e.g., constants).
func MustNewDatabaseType(s string) DatabaseType {
	dt, err := NewDatabaseType(s)
	if err != nil {
		panic(err)
	}
	return dt
}

// IsValidProjectType checks if a string is a valid project type
func IsValidProjectType(s string) bool {
	_, err := NewProjectType(s)
	return err == nil
}

// IsValidDatabaseType checks if a string is a valid database type
func IsValidDatabaseType(s string) bool {
	_, err := NewDatabaseType(s)
	return err == nil
}
