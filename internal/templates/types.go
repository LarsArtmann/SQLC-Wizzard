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

// PackageConfig represents complete package configuration
type PackageConfig struct {
	Name     string
	Path      string
	BuildTags string
}

// DatabaseConfig represents database-specific configuration
type DatabaseConfig struct {
	Engine      DatabaseType
	URL         string
	UseManaged  bool
	UseUUIDs    bool
	UseJSON     bool
	UseArrays   bool
	UseFullText bool
}

// OutputConfig represents output directory configuration
type OutputConfig struct {
	BaseDir    string
	QueriesDir  string
	SchemaDir   string
}

// ValidationConfig represents validation settings
type ValidationConfig struct {
	StrictFunctions bool
	StrictOrderBy   bool
	EmitOptions    domain.EmitOptions
	SafetyRules     domain.SafetyRules
}

// TemplateData contains all data needed to render a template
// SPLIT BRAIN FIXED: Unified configuration structure
type TemplateData struct {
	ProjectName string
	ProjectType ProjectType
	
	Package    PackageConfig
	Database   DatabaseConfig
	Output     OutputConfig
	Validation ValidationConfig
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
