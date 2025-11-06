package templates

import "github.com/LarsArtmann/SQLC-Wizzard/pkg/config"

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

	// Feature flags
	Features Features

	// Safety settings
	SafetyRules SafetyRules
}

// Features represents optional database features
type Features struct {
	UUIDs              bool
	JSON               bool
	Arrays             bool
	FullTextSearch     bool
	EmitInterface      bool
	PreparedQueries    bool
	JSONTags           bool
	DBTags             bool
	ExactTableNames    bool
	EmptySlices        bool
	OmitUnusedStructs  bool
}

// SafetyRules represents validation and safety features
type SafetyRules struct {
	NoSelectStar      bool
	RequireWhere      bool
	RequireLimit      bool
	NoDropTable       bool
	StrictFunctions   bool
	StrictOrderBy     bool
}

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

// DefaultFeatures returns a sensible default feature set
func DefaultFeatures() Features {
	return Features{
		UUIDs:              true,
		JSON:               true,
		Arrays:             false,
		FullTextSearch:     false,
		EmitInterface:      true,
		PreparedQueries:    true,
		JSONTags:           true,
		DBTags:             false,
		ExactTableNames:    true,
		EmptySlices:        true,
		OmitUnusedStructs:  true,
	}
}

// DefaultSafetyRules returns recommended safety rules
func DefaultSafetyRules() SafetyRules {
	return SafetyRules{
		NoSelectStar:    true,
		RequireWhere:    true,
		RequireLimit:    false,
		NoDropTable:     true,
		StrictFunctions: true,
		StrictOrderBy:   true,
	}
}
