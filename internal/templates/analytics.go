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
		false, // emitResultPointers
		false, // emitParamsPointers
		false, // emitEnumValidMethod
		"camel",
		false, // noSelectStar
		false, // requireWhere
		true,  // requireLimit
	)
}

// RequiredFeatures returns which features this template requires.
func (t *AnalyticsTemplate) RequiredFeatures() []string {
	return []string{"emit_interface", "json_tags", "full_text_search", "strict_checks"}
}
