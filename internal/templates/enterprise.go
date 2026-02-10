package templates

// EnterpriseTemplate generates sqlc config for enterprise-scale projects.
type EnterpriseTemplate struct {
	ConfiguredTemplate
}

// NewEnterpriseTemplate creates a new enterprise template.
func NewEnterpriseTemplate() *EnterpriseTemplate {
	base := NewConfiguredTemplate(
		"enterprise",
		"Production-ready configuration with strict safety rules for enterprise applications",
		"db",
		"enterprise",
		true,
		"enterprise",
		"postgresql",
	)

	// Override enterprise-specific settings
	base.UseFullText = true
	base.RequireLimit = true
	base.Features = []string{"emit_interface", "prepared_queries", "json_tags", "strict_checks"}

	return &EnterpriseTemplate{ConfiguredTemplate: base}
}

// Name returns the template name.
func (t *EnterpriseTemplate) Name() string {
	return "enterprise"
}

// Description returns a human-readable description.
func (t *EnterpriseTemplate) Description() string {
	return "Production-ready configuration with strict safety rules for enterprise applications"
}
