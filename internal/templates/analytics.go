package templates

import (
	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/pkg/config"
)

// AnalyticsTemplate generates sqlc config for analytics and data warehouse projects.
type AnalyticsTemplate struct {
	ConfiguredTemplate
}

// NewAnalyticsTemplate creates a new analytics template.
func NewAnalyticsTemplate() *AnalyticsTemplate {
	base := NewConfiguredTemplate(
		"analytics",
		"Optimized for data analytics and reporting with full-text search and array support",
		"analytics",
		"analytics",
		true, // strictMode - analytics needs strict checks
		"analytics",
		"postgresql",
	)

	// Override analytics-specific settings
	base.UseManaged = false
	base.UseUUIDs = false
	base.UseJSON = true
	base.UseArrays = true
	base.UseFullText = true
	base.EmitJSONTags = true
	base.EmitInterface = true
	base.EmitEmptySlices = false
	base.EmitPreparedQueries = false
	base.JSONTagsCaseStyle = "camel"
	base.StrictFunctions = true
	base.StrictOrderBy = true
	base.NoSelectStar = false
	base.RequireWhere = false
	base.RequireLimit = true
	base.Features = []string{"emit_interface", "json_tags", "full_text_search", "strict_checks"}

	return &AnalyticsTemplate{ConfiguredTemplate: base}
}

// Name returns the template name.
func (t *AnalyticsTemplate) Name() string {
	return "analytics"
}

// Description returns the template description.
func (t *AnalyticsTemplate) Description() string {
	return "Optimized for data analytics and reporting with full-text search and array support"
}

// Generate creates a SqlcConfig from template data.
func (t *AnalyticsTemplate) Generate(data generated.TemplateData) (*config.SqlcConfig, error) {
	return t.GenerateWithDefaults(
		data,
		"analytics",
		"internal/analytics",
		"internal/analytics",
		"internal/analytics/queries",
		"internal/analytics/schema",
		"${ANALYTICS_DATABASE_URL}",
		"analytics",
		true,
	)
}

// DefaultData returns default TemplateData for analytics template.
func (t *AnalyticsTemplate) DefaultData() TemplateData {
	return createDefaultTemplateData(
		"analytics",
		"analytics",
		"internal/analytics",
		"${ANALYTICS_DATABASE_URL}",
		"postgresql",
		false, // useManaged
		false, // useUUIDs
		true,  // useJSON
		true,  // useArrays
		true,  // useFullText
		"internal/analytics",
		true,  // strictFunctions
		true,  // strictOrderBy
		true,  // emitJSONTags
		true,  // emitInterface
		true,  // emitAllEnumValues
		false, // emitPreparedQueries
		true,  // emitEmptySlices
		false, // emitResultPointers
		false, // emitParamsPointers
		false, // emitEnumValidMethod
		"camel",
		false, // noSelectStar
		false, // requireWhere
		true,  // noDropTable
		true,  // noTruncate
		true,  // requireLimit
	)
}

// RequiredFeatures returns which features this template requires.
func (t *AnalyticsTemplate) RequiredFeatures() []string {
	return t.Features
}
