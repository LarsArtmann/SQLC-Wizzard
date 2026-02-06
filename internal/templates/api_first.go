package templates

// APIFirstTemplate generates sqlc config for API-first projects.
type APIFirstTemplate struct {
	ConfiguredTemplate
}

// NewAPIFirstTemplate creates a new API-first template.
func NewAPIFirstTemplate() *APIFirstTemplate {
	return &APIFirstTemplate{
		ConfiguredTemplate: ConfiguredTemplate{
			// Defaults for Generate()
			DefaultPackageName: "api",
			DefaultProjectName: "api",
			StrictMode:         false,
			// Paths
			PackagePath: "internal/db",
			BaseOutput:  "internal/db",
			// Type and features
			ProjectType: "api-first",
			DbEngine:    "postgresql",
			// Database features
			UseManaged:  true,
			UseUUIDs:    true,
			UseJSON:     true,
			UseArrays:   true,
			UseFullText: false,
			// Emit options
			EmitPreparedQueries:      true,
			EmitResultStructPointers: true,
			EmitParamsStructPointers: true,
			// Safety rules
			NoSelectStar: true,
			RequireWhere: true,
			RequireLimit: false,
			// Required features
			Features: []string{"emit_interface", "prepared_queries", "json_tags", "camel_case"},
		},
	}
}

// Name returns the template name.
func (t *APIFirstTemplate) Name() string {
	return "api-first"
}

// Description returns a human-readable description.
func (t *APIFirstTemplate) Description() string {
	return "Optimized for REST/GraphQL API development with JSON support and camelCase naming"
}
