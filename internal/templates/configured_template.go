package templates

import (
	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/pkg/config"
)

// ConfiguredTemplate is a reusable template implementation that accepts
// configuration parameters instead of duplicating method implementations.
// This eliminates the duplicate code pattern across similar templates.
type ConfiguredTemplate struct {
	BaseTemplate

	// Template identification
	TemplateName        string
	TemplateDescription string

	// Template-specific values for Generate()
	DefaultPackageName string
	DefaultProjectName string
	StrictMode         bool

	// Template-specific values for DefaultData()
	ProjectType string
	DbEngine    string
	PackagePath string
	BaseOutput  string

	// Database features
	UseManaged, UseUUIDs, UseJSON, UseArrays, UseFullText bool

	// Emit options
	EmitPreparedQueries, EmitResultStructPointers, EmitParamsStructPointers bool

	// Safety rules
	NoSelectStar bool
	RequireWhere bool
	RequireLimit bool

	// Required features
	Features []string
}

// NewConfiguredTemplate creates a base ConfiguredTemplate with common defaults.
// Template-specific constructors should call this and override fields as needed.
// This eliminates the duplicate initialization code across template implementations.
func NewConfiguredTemplate(
	name, description string,
	defaultPackageName, defaultProjectName string,
	strictMode bool,
	projectType, dbEngine string,
) ConfiguredTemplate {
	return ConfiguredTemplate{
		// Template identification
		TemplateName:        name,
		TemplateDescription: description,

		// Defaults for Generate()
		DefaultPackageName: defaultPackageName,
		DefaultProjectName: defaultProjectName,
		StrictMode:         strictMode,

		// Paths - common default for most templates
		PackagePath: "internal/db",
		BaseOutput:  "internal/db",

		// Type and features
		ProjectType: projectType,
		DbEngine:    dbEngine,

		// Database features - common defaults for PostgreSQL-based templates
		UseManaged:  true,
		UseUUIDs:    true,
		UseJSON:     true,
		UseArrays:   true,
		UseFullText: false,

		// Emit options - common defaults for production templates
		EmitPreparedQueries:      true,
		EmitResultStructPointers: true,
		EmitParamsStructPointers: true,

		// Safety rules - conservative defaults
		NoSelectStar: true,
		RequireWhere: true,
		RequireLimit: false,

		// Required features - minimal default set
		Features: []string{"emit_interface", "prepared_queries", "json_tags"},
	}
}

// Name returns the template name.
func (t *ConfiguredTemplate) Name() string {
	return t.TemplateName
}

// Description returns a human-readable description.
func (t *ConfiguredTemplate) Description() string {
	return t.TemplateDescription
}

// Generate creates a SqlcConfig from template data using the configured defaults.
func (t *ConfiguredTemplate) Generate(data generated.TemplateData) (*config.SqlcConfig, error) {
	// Use either configured or default values to support empty struct initialization
	defaultPackageName := t.DefaultPackageName
	if defaultPackageName == "" {
		defaultPackageName = "db"
	}
	defaultProjectName := t.DefaultProjectName
	if defaultProjectName == "" {
		defaultProjectName = "project"
	}
	strictMode := t.StrictMode

	packagePath := t.PackagePath
	if packagePath == "" {
		packagePath = "internal/db"
	}
	baseOutput := t.BaseOutput
	if baseOutput == "" {
		baseOutput = "internal/db"
	}

	return t.GenerateWithDefaults(
		data,
		defaultPackageName,
		packagePath,
		baseOutput,
		baseOutput+"/queries",
		baseOutput+"/schema",
		"${DATABASE_URL}",
		defaultProjectName,
		strictMode,
	)
}

// DefaultData returns default TemplateData with the configured values.
func (t *ConfiguredTemplate) DefaultData() generated.TemplateData {
	// Use configured values or defaults to support empty struct initialization
	projectType := t.ProjectType
	if projectType == "" {
		projectType = "microservice" // sensible default
	}
	dbEngine := t.DbEngine
	if dbEngine == "" {
		dbEngine = "postgresql" // sensible default
	}
	packagePath := t.PackagePath
	if packagePath == "" {
		packagePath = "internal/db"
	}
	baseOutput := t.BaseOutput
	if baseOutput == "" {
		baseOutput = "internal/db"
	}

	return t.BuildDefaultData(
		projectType,
		dbEngine,
		"${DATABASE_URL}",
		packagePath,
		baseOutput,
		t.UseManaged,
		t.UseUUIDs,
		t.UseJSON,
		t.UseArrays,
		t.UseFullText,
		t.EmitPreparedQueries,
		t.EmitResultStructPointers,
		t.EmitParamsStructPointers,
		t.NoSelectStar,
		t.RequireWhere,
		t.RequireLimit,
	)
}

// RequiredFeatures returns which features this template requires.
func (t *ConfiguredTemplate) RequiredFeatures() []string {
	return t.Features
}
