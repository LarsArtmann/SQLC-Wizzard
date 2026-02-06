package templates

import (
	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/pkg/config"
)

// APIFirstTemplate generates sqlc config for API-first projects.
type APIFirstTemplate struct {
	BaseTemplate
}

// NewAPIFirstTemplate creates a new API-first template.
func NewAPIFirstTemplate() *APIFirstTemplate {
	return &APIFirstTemplate{}
}

// Name returns the template name.
func (t *APIFirstTemplate) Name() string {
	return "api-first"
}

// Description returns a human-readable description.
func (t *APIFirstTemplate) Description() string {
	return "Optimized for REST/GraphQL API development with JSON support and camelCase naming"
}

// Generate creates a SqlcConfig from template data.
func (t *APIFirstTemplate) Generate(data generated.TemplateData) (*config.SqlcConfig, error) {
	// Set defaults
	cb := ConfigBuilder{
		Data:               data,
		DefaultName:        "api",
		DefaultDatabaseURL: "${DATABASE_URL}",
		Strict:             false,
	}

	// Build config using base builder
	cfg, err := cb.Build()
	if err != nil {
		return nil, err
	}

	// Apply validation rules using base template helper
	return t.ApplyValidationRules(cfg, data)
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
		true,  // emitPreparedQueries
		true,  // emitResultStructPointers
		true,  // emitParamsStructPointers
		true,  // noSelectStar
		true,  // requireWhere
		false, // requireLimit
	)
}

// RequiredFeatures returns which features this template requires.
func (t *APIFirstTemplate) RequiredFeatures() []string {
	return []string{"emit_interface", "prepared_queries", "json_tags", "camel_case"}
}
