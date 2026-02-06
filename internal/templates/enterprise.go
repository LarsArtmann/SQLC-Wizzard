package templates

// EnterpriseTemplate generates sqlc config for enterprise-scale projects.
type EnterpriseTemplate struct {
	ConfiguredTemplate
}

// NewEnterpriseTemplate creates a new enterprise template.
func NewEnterpriseTemplate() *EnterpriseTemplate {
	return &EnterpriseTemplate{
		ConfiguredTemplate: ConfiguredTemplate{
			// Template identification
			TemplateName:        "enterprise",
			TemplateDescription: "Production-ready configuration with strict safety rules for enterprise applications",

			// Defaults for Generate()
			DefaultPackageName: "db",
			DefaultProjectName: "enterprise",
			StrictMode:         true,

			// Paths
			PackagePath: "internal/db",
			BaseOutput:  "internal/db",

			// Type and features
			ProjectType: "enterprise",
			DbEngine:    "postgresql",

			// Database features
			UseManaged: true,
			UseUUIDs:   true,
			UseJSON:    true,
			UseArrays:  true,
			UseFullText: true,

			// Emit options
			EmitPreparedQueries:      true,
			EmitResultStructPointers: true,
			EmitParamsStructPointers: true,

			// Safety rules
			NoSelectStar: true,
			RequireWhere: true,
			RequireLimit: true,

			// Required features
			Features: []string{"emit_interface", "prepared_queries", "json_tags", "strict_checks"},
		},
	}
}
