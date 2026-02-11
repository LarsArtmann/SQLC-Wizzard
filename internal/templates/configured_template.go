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
	EmitPreparedQueries, EmitResultStructPointers, EmitParamsStructPointers,
	EmitJSONTags, EmitInterface, EmitEmptySlices,
	EmitEnumValidMethod, EmitAllEnumValues bool
	JSONTagsCaseStyle string

	// Emit options - extended
	StrictFunctions, StrictOrderBy bool

	// Safety rules
	NoSelectStar, RequireWhere, NoDropTable, NoTruncate, RequireLimit bool

	// Required features
	Features []string

	// Custom rename rules for Go codegen (optional override)
	CustomRenameRules map[string]string
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
		EmitJSONTags:            false,
		EmitInterface:           false,
		EmitEmptySlices:          false,
		EmitEnumValidMethod:     false,
		EmitAllEnumValues:       false,
		JSONTagsCaseStyle:       "snake",

		// Emit options - extended
		StrictFunctions: false,
		StrictOrderBy:   false,

		// Safety rules - conservative defaults
		NoSelectStar: true,
		RequireWhere: true,
		NoDropTable:  false,
		NoTruncate:   false,
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

	// Use template's StrictMode or fall back to data's validation settings
	strictMode := t.StrictMode
	if !strictMode {
		// Use data's validation strict settings if template's StrictMode is not set
		strictMode = data.Validation.StrictFunctions || data.Validation.StrictOrderBy
	}

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

	// Use configured values or defaults for boolean fields
	useManaged := t.UseManaged
	if !useManaged {
		useManaged = true // sensible default for most templates
	}
	useUUIDs := t.UseUUIDs
	if !useUUIDs {
		useUUIDs = true // UUIDs are commonly needed
	}
	useJSON := t.UseJSON
	if !useJSON {
		useJSON = true // JSON support is often useful
	}
	useArrays := t.UseArrays
	if !useArrays {
		useArrays = true // Arrays are commonly used
	}
	useFullText := t.UseFullText
	// No default - false by default

	emitJSONTags := t.EmitJSONTags
	emitPreparedQueries := t.EmitPreparedQueries
	if !emitPreparedQueries {
		emitPreparedQueries = true // sensible default
	}
	emitInterface := t.EmitInterface
	emitEmptySlices := t.EmitEmptySlices
	emitResultStructPointers := t.EmitResultStructPointers
	if !emitResultStructPointers {
		emitResultStructPointers = true // sensible default
	}
	emitParamsStructPointers := t.EmitParamsStructPointers
	if !emitParamsStructPointers {
		emitParamsStructPointers = true // sensible default
	}
	emitEnumValidMethod := t.EmitEnumValidMethod
	emitAllEnumValues := t.EmitAllEnumValues
	jsonTagsCaseStyle := t.JSONTagsCaseStyle
	if jsonTagsCaseStyle == "" {
		jsonTagsCaseStyle = "snake" // default case style
	}
	strictFunctions := t.StrictFunctions
	strictOrderBy := t.StrictOrderBy

	noSelectStar := t.NoSelectStar
	if !noSelectStar {
		noSelectStar = true // conservative default
	}
	requireWhere := t.RequireWhere
	if !requireWhere {
		requireWhere = true // conservative default
	}
	noDropTable := t.NoDropTable
	noTruncate := t.NoTruncate
	requireLimit := t.RequireLimit

	return t.BuildDefaultData(
		projectType,
		dbEngine,
		"${DATABASE_URL}",
		packagePath,
		baseOutput,
		useManaged,
		useUUIDs,
		useJSON,
		useArrays,
		useFullText,
		emitJSONTags,
		emitPreparedQueries,
		emitInterface,
		emitEmptySlices,
		emitResultStructPointers,
		emitParamsStructPointers,
		emitEnumValidMethod,
		emitAllEnumValues,
		jsonTagsCaseStyle,
		strictFunctions,
		strictOrderBy,
		noSelectStar,
		requireWhere,
		noDropTable,
		noTruncate,
		requireLimit,
	)
}

// RequiredFeatures returns which features this template requires.
func (t *ConfiguredTemplate) RequiredFeatures() []string {
	return t.Features
}

// GetRenameRules returns custom rename rules if set, otherwise falls back to base rules.
func (t *ConfiguredTemplate) GetRenameRules() map[string]string {
	if t.CustomRenameRules != nil {
		return t.CustomRenameRules
	}
	return t.BaseTemplate.GetRenameRules()
}

// BuildGoConfigWithOverrides builds GoGenConfig with custom rename rules support.
func (t *ConfiguredTemplate) BuildGoConfigWithOverrides(data generated.TemplateData) *config.GoGenConfig {
	sqlPackage := t.GetSQLPackage(data.Database.Engine)
	cfg := t.BuildGoGenConfig(data, sqlPackage)

	// Use custom rename rules if set
	if t.CustomRenameRules != nil {
		cfg.Rename = t.CustomRenameRules
	}

	return cfg
}

// Compile-time interface compliance checks.
// These ensure all template types implement the Template interface.
// If a template type is missing an interface method, this will fail to compile.
var (
	_ Template = (*ConfiguredTemplate)(nil)
	_ Template = (*EnterpriseTemplate)(nil)
	_ Template = (*APIFirstTemplate)(nil)
	_ Template = (*MicroserviceTemplate)(nil)
	_ Template = (*HobbyTemplate)(nil)
	_ Template = (*LibraryTemplate)(nil)
	_ Template = (*AnalyticsTemplate)(nil)
	_ Template = (*TestingTemplate)(nil)
	_ Template = (*MultiTenantTemplate)(nil)
)
