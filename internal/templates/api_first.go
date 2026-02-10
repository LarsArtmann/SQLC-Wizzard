package templates

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
	base.Features = []string{"emit_interface", "prepared_queries", "json_tags", "camel_case"}

	return &APIFirstTemplate{ConfiguredTemplate: base}
}
