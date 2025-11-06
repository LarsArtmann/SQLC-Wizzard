package templates

import (
	"github.com/LarsArtmann/SQLC-Wizzard/internal/domain"
	"github.com/LarsArtmann/SQLC-Wizzard/pkg/config"
)

// ProjectType represents the type of project template
type ProjectType string

const (
	ProjectTypeHobby        ProjectType = "hobby"
	ProjectTypeMicroservice ProjectType = "microservice"
	ProjectTypeEnterprise   ProjectType = "enterprise"
	ProjectTypeAPIFirst     ProjectType = "api-first"
	ProjectTypeAnalytics    ProjectType = "analytics"
	ProjectTypeTesting      ProjectType = "testing"
	ProjectTypeMultiTenant  ProjectType = "multi-tenant"
	ProjectTypeLibrary      ProjectType = "library"
)

// DatabaseType represents the supported database engines
type DatabaseType string

const (
	DatabaseTypePostgreSQL DatabaseType = "postgresql"
	DatabaseTypeMySQL      DatabaseType = "mysql"
	DatabaseTypeSQLite     DatabaseType = "sqlite"
)

// TemplateData contains all data needed to render a template
type TemplateData struct {
	// Project information
	ProjectName string
	ProjectType ProjectType
	PackagePath string

	// Database configuration
	Database       DatabaseType
	DatabaseURL    string
	UseManagedDB   bool

	// Output directories
	OutputDir  string
	QueriesDir string
	SchemaDir  string

	// Go package settings
	PackageName string
	SQLPackage  string
	BuildTags   string

	// Database feature flags (affect type overrides and template logic)
	UseUUIDs          bool
	UseJSON           bool
	UseArrays         bool
	UseFullTextSearch bool

	// Code generation options
	EmitOptions domain.EmitOptions

	// Safety rules (CEL-based validation)
	SafetyRules domain.SafetyRules

	// Strict validation flags (config-level, not CEL rules)
	StrictFunctions bool
	StrictOrderBy   bool
}

// NOTE: Features split brain eliminated!
// - Database features (UUIDs, JSON, Arrays, FTS) are now individual boolean fields in TemplateData
// - Code generation options moved to internal/domain/emit_options.go
// Use domain.EmitOptions and domain.DefaultEmitOptions() instead of Features.

// NOTE: SafetyRules moved to internal/domain/rule.go to eliminate split brain.
// Use domain.SafetyRules and domain.DefaultSafetyRules() instead.

// Template represents a project template
type Template interface {
	// Name returns the template name
	Name() string

	// Description returns a human-readable description
	Description() string

	// Generate creates a SqlcConfig from template data
	Generate(data TemplateData) (*config.SqlcConfig, error)

	// DefaultData returns default TemplateData for this template
	DefaultData() TemplateData

	// RequiredFeatures returns which features this template requires
	RequiredFeatures() []string
}

// GetAllProjectTypes returns all available project types
func GetAllProjectTypes() []ProjectType {
	return []ProjectType{
		ProjectTypeHobby,
		ProjectTypeMicroservice,
		ProjectTypeEnterprise,
		ProjectTypeAPIFirst,
		ProjectTypeAnalytics,
		ProjectTypeTesting,
		ProjectTypeMultiTenant,
		ProjectTypeLibrary,
	}
}

// GetAllDatabaseTypes returns all supported database types
func GetAllDatabaseTypes() []DatabaseType {
	return []DatabaseType{
		DatabaseTypePostgreSQL,
		DatabaseTypeMySQL,
		DatabaseTypeSQLite,
	}
}

// String returns the string representation of ProjectType
func (pt ProjectType) String() string {
	return string(pt)
}

// String returns the string representation of DatabaseType
func (dt DatabaseType) String() string {
	return string(dt)
}

// NOTE: DefaultFeatures moved to domain.DefaultEmitOptions()
// NOTE: DefaultSafetyRules moved to domain.DefaultSafetyRules()
