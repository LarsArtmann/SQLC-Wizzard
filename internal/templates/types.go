package templates

import (
	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/apperrors"
	"github.com/LarsArtmann/SQLC-Wizzard/pkg/config"
)

// Type aliases for generated types - PROHIBIT DIRECT STRING USAGE!
type (
	ProjectType  = generated.ProjectType
	DatabaseType = generated.DatabaseType
	TemplateData = generated.TemplateData
	EmitOptions  = generated.EmitOptions
	SafetyRule   = generated.SafetyRule
	SafetyRules  = generated.SafetyRules
)

// Helper functions for validation.
func IsValidProjectType(projectType string) bool {
	return generated.ProjectType(projectType).IsValid()
}

func IsValidDatabaseType(database string) bool {
	return generated.DatabaseType(database).IsValid()
}

// Template interface defines behavior for all template implementations.
type Template interface {
	Generate(data generated.TemplateData) (*config.SqlcConfig, error)
	DefaultData() generated.TemplateData
	RequiredFeatures() []string
	Name() string
	Description() string
}

// Smart constructors with validation - PREVENT INVALID STATES!
func NewProjectType(projectType string) (ProjectType, error) {
	pt := ProjectType(projectType)
	if !pt.IsValid() {
		return "", apperrors.ValidationError("project_type", projectType)
	}
	return pt, nil
}

func NewDatabaseType(database string) (DatabaseType, error) {
	dt := DatabaseType(database)
	if !dt.IsValid() {
		return "", apperrors.ValidationError("database", database)
	}
	return dt, nil
}

// MustNew constructors - PANIC on invalid input (for constants).
func MustNewProjectType(projectType string) ProjectType {
	pt, err := NewProjectType(projectType)
	if err != nil {
		panic(err) // This is programmer error, not runtime error
	}
	return pt
}

func MustNewDatabaseType(database string) DatabaseType {
	dt, err := NewDatabaseType(database)
	if err != nil {
		panic(err) // This is programmer error, not runtime error
	}
	return dt
}

// Constants - use generated types, NO MANUAL STRINGS!
const (
	ProjectTypeHobby        = ProjectType(generated.ProjectTypeHobby)
	ProjectTypeMicroservice = ProjectType(generated.ProjectTypeMicroservice)
	ProjectTypeEnterprise   = ProjectType(generated.ProjectTypeEnterprise)
	ProjectTypeAPIFirst     = ProjectType(generated.ProjectTypeAPIFirst)
	ProjectTypeAnalytics    = ProjectType(generated.ProjectTypeAnalytics)
	ProjectTypeTesting      = ProjectType(generated.ProjectTypeTesting)
	ProjectTypeMultiTenant  = ProjectType(generated.ProjectTypeMultiTenant)
	ProjectTypeLibrary      = ProjectType(generated.ProjectTypeLibrary)

	DatabaseTypePostgreSQL = DatabaseType(generated.DatabaseTypePostgreSQL)
	DatabaseTypeMySQL      = DatabaseType(generated.DatabaseTypeMySQL)
	DatabaseTypeSQLite     = DatabaseType(generated.DatabaseTypeSQLite)
)
