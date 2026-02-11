package templates

import (
	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/pkg/config"
)

// APIFirstTemplate generates sqlc config for API-first projects.
type APIFirstTemplate struct {
	ConfiguredTemplate
}

// NewAPIFirstTemplate creates a new API-first template.
func NewAPIFirstTemplate() *APIFirstTemplate {
	base := NewConfiguredTemplate(
		"api-first",
		"Optimized for REST/GraphQL API development with JSON support and camelCase naming",
		"api",
		"api",
		false,
		"api-first",
		"postgresql",
	)

	// Override API-first specific settings
	base.JSONTagsCaseStyle = "camel"
	base.EmitJSONTags = true
	base.EmitInterface = true
	base.Features = []string{"emit_interface", "prepared_queries", "json_tags", "camel_case"}

	return &APIFirstTemplate{ConfiguredTemplate: base}
}

// Name returns the template name.
func (t *APIFirstTemplate) Name() string {
	return "api-first"
}

// Description returns the template description.
func (t *APIFirstTemplate) Description() string {
	return "Optimized for REST/GraphQL API development with JSON support and camelCase naming"
}

// DefaultData returns default TemplateData for API-first template.
func (t *APIFirstTemplate) DefaultData() generated.TemplateData {
	return t.BuildDefaultData(
		"api-first",
		"postgresql",
		"${DATABASE_URL}",
		"internal/db",
		"internal/db",
		true,  // useManaged
		true,  // useUUIDs
		true,  // useJSON
		true,  // useArrays
		false, // useFullText
		true,  // emitJSONTags
		true,  // emitPreparedQueries
		true,  // emitInterface
		false, // emitEmptySlices
		true,  // emitResultStructPointers
		true,  // emitParamsStructPointers
		false, // emitEnumValidMethod
		false, // emitAllEnumValues
		"camel", // jsonTagsCaseStyle
		false, // strictFunctions
		false, // strictOrderBy
		true,  // noSelectStar
		true,  // requireWhere
		false, // noDropTable
		false, // noTruncate
		false, // requireLimit
	)
}

// Generate creates a SqlcConfig from template data.
func (t *APIFirstTemplate) Generate(data generated.TemplateData) (*config.SqlcConfig, error) {
	return t.ConfiguredTemplate.Generate(data)
}
