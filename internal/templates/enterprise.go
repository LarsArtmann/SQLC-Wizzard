package templates

import (
	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/pkg/config"
)

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
	base.JSONTagsCaseStyle = "camel"
	base.EmitJSONTags = true
	base.EmitInterface = true
	base.Features = []string{"emit_interface", "prepared_queries", "json_tags", "strict_checks"}

	return &EnterpriseTemplate{ConfiguredTemplate: base}
}

// Name returns the template name.
func (t *EnterpriseTemplate) Name() string {
	return "enterprise"
}

// Description returns the template description.
func (t *EnterpriseTemplate) Description() string {
	return "Production-ready configuration with strict safety rules for enterprise applications"
}

// DefaultData returns default TemplateData for enterprise template.
func (t *EnterpriseTemplate) DefaultData() generated.TemplateData {
	return t.BuildDefaultData(
		"enterprise",
		"postgresql",
		"${DATABASE_URL}",
		"internal/db",
		"internal/db",
		true,  // useManaged
		true,  // useUUIDs
		true,  // useJSON
		true,  // useArrays
		true,  // useFullText
		true,  // emitJSONTags
		true,  // emitPreparedQueries
		true,  // emitInterface
		false, // emitEmptySlices
		true,  // emitResultStructPointers
		true,  // emitParamsStructPointers
		false, // emitEnumValidMethod
		false, // emitAllEnumValues
		"camel", // jsonTagsCaseStyle
		true,  // strictFunctions
		true,  // strictOrderBy
		true,  // noSelectStar
		true,  // requireWhere
		false, // noDropTable
		false, // noTruncate
		true,  // requireLimit
	)
}

// Generate creates a SqlcConfig from template data.
func (t *EnterpriseTemplate) Generate(data generated.TemplateData) (*config.SqlcConfig, error) {
	return t.ConfiguredTemplate.Generate(data)
}
