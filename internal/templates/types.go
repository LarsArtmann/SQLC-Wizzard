// Package templates provides sqlc configuration templates for different project types.
// It implements the Template interface to support consistent configuration generation.
//
// Architecture Overview
// ====================
//
// This package uses two patterns for implementing templates:
//
//  1. BaseTemplate Direct Embedding (simple, minimal overhead):
//     Used by: MicroserviceTemplate, HobbyTemplate, LibraryTemplate, AnalyticsTemplate,
//     TestingTemplate, MultiTenantTemplate
//
//     Pros:
//     - Direct method implementation, full control over Generate()
//     - No inheritance overhead
//     - Best for unique Generate() implementations
//
//     Cons:
//     - Code duplication in Generate() and DefaultData()
//     - Manual zero-value handling required
//
//  2. ConfiguredTemplate Embedding (reusable, consistent defaults):
//     Used by: EnterpriseTemplate, APIFirstTemplate
//
//     Pros:
//     - Built-in zero-value handling (sensible defaults)
//     - Consistent Generate() behavior via ConfiguredTemplate
//     - Less boilerplate in template implementations
//
//     Cons:
//     - Inheritance from ConfiguredTemplate
//     - May be overkill for simple templates
//
// When to Use Which Pattern
// ==========================
//
// Use ConfiguredTemplate when:
// - Template's Generate() uses mostly default handling
// - You want zero-value safety (sensible defaults)
// - You want consistent behavior across templates
//
// Use BaseTemplate directly when:
// - Generate() needs custom logic not covered by ConfiguredTemplate
// - You need full control over the generation process
// - The template is fundamentally different from ConfiguredTemplate
//
// Interface Compliance
// ====================
//
// All template types are verified at compile time to implement the Template interface.
// This ensures no runtime surprises from missing methods:
//
//	var (
//		_ Template = (*EnterpriseTemplate)(nil)
//		_ Template = (*APIFirstTemplate)(nil)
//		_ Template = (*MicroserviceTemplate)(nil)
//		// ... etc
//	)
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

// Constants - use generated types directly, NO MANUAL STRINGS!
const (
	ProjectTypeHobby        = generated.ProjectTypeHobby
	ProjectTypeMicroservice = generated.ProjectTypeMicroservice
	ProjectTypeEnterprise   = generated.ProjectTypeEnterprise
	ProjectTypeAPIFirst     = generated.ProjectTypeAPIFirst
	ProjectTypeAnalytics    = generated.ProjectTypeAnalytics
	ProjectTypeTesting      = generated.ProjectTypeTesting
	ProjectTypeMultiTenant  = generated.ProjectTypeMultiTenant
	ProjectTypeLibrary      = generated.ProjectTypeLibrary

	DatabaseTypePostgreSQL = generated.DatabaseTypePostgreSQL
	DatabaseTypeMySQL      = generated.DatabaseTypeMySQL
	DatabaseTypeSQLite     = generated.DatabaseTypeSQLite
)
