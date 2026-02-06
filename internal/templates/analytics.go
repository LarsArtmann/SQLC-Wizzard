package templates

import (
	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/pkg/config"
)

// AnalyticsTemplate generates sqlc config for analytics and data warehouse projects.
type AnalyticsTemplate struct {
	BaseTemplate
}

// NewAnalyticsTemplate creates a new analytics template.
func NewAnalyticsTemplate() *AnalyticsTemplate {
	return &AnalyticsTemplate{}
}

// Name returns the template name.
func (t *AnalyticsTemplate) Name() string {
	return "analytics"
}

// Description returns a human-readable description.
func (t *AnalyticsTemplate) Description() string {
	return "Optimized for data analytics and reporting with full-text search and array support"
}

// Generate creates a SqlcConfig from template data.
func (t *AnalyticsTemplate) Generate(data generated.TemplateData) (*config.SqlcConfig, error) {
	// Apply analytics-specific defaults
	if data.Package.Name == "" {
		data.Package.Name = "analytics"
	}
	if data.Package.Path == "" {
		data.Package.Path = "internal/analytics"
	}
	if data.Output.BaseDir == "" {
		data.Output.BaseDir = "internal/analytics"
	}
	if data.Output.QueriesDir == "" {
		data.Output.QueriesDir = "internal/analytics/queries"
	}
	if data.Output.SchemaDir == "" {
		data.Output.SchemaDir = "internal/analytics/schema"
	}
	if data.Database.URL == "" {
		data.Database.URL = "${ANALYTICS_DATABASE_URL}"
	}

	// Build base config using shared builder
	builder := &ConfigBuilder{
		Data:             data,
		DefaultName:      "analytics",
		DefaultDatabaseURL: "${ANALYTICS_DATABASE_URL}",
		Strict:           true,
	}
	cfg, _ := builder.Build()

	// Generate Go config with template-specific settings
	cfg.SQL[0].Gen.Go = t.BuildGoConfigWithOverrides(data)

	// Apply validation rules using base helper
	return t.ApplyValidationRules(cfg, data)
}

// DefaultData returns default TemplateData for analytics template.
func (t *AnalyticsTemplate) DefaultData() TemplateData {
	return createDefaultTemplateData(
		"analytics",
		"analytics",
		"internal/analytics",
		"${ANALYTICS_DATABASE_URL}",
		"postgresql",
		false,  // useManaged
		false,  // useUUIDs
		true,   // useJSON
		true,   // useArrays
		true,   // useFullText
		"internal/analytics",
		true,   // strictFunctions
		true,   // strictOrderBy
		true,   // emitJSONTags
		true,   // emitInterface
		true,   // emitAllEnumValues
		false,  // emitPreparedQueries
		false,  // emitResultPointers
		false,  // emitParamsPointers
		false,  // emitEnumValidMethod
		"snake",
		false,  // noSelectStar
		false,  // requireWhere
		true,   // requireLimit
	)
}

// RequiredFeatures returns which features this template requires.
func (t *AnalyticsTemplate) RequiredFeatures() []string {
	return []string{"emit_interface", "json_tags", "full_text_search", "strict_checks"}
}
