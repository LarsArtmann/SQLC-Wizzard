package templates

import (
	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/pkg/config"
)

// EnterpriseTemplate generates sqlc config for enterprise-scale projects.
type EnterpriseTemplate struct {
	BaseTemplate
}

// NewEnterpriseTemplate creates a new enterprise template.
func NewEnterpriseTemplate() *EnterpriseTemplate {
	return &EnterpriseTemplate{}
}

// Name returns the template name.
func (t *EnterpriseTemplate) Name() string {
	return "enterprise"
}

// Description returns a human-readable description.
func (t *EnterpriseTemplate) Description() string {
	return "Production-ready configuration with strict safety rules for enterprise applications"
}

// Generate creates a SqlcConfig from template data.
func (t *EnterpriseTemplate) Generate(data generated.TemplateData) (*config.SqlcConfig, error) {
	return t.BaseTemplate.GenerateWithDefaults(
		data,
		"db",                          // packageName
		"internal/db",                 // packagePath
		"internal/db",                 // baseDir
		"internal/db/queries",         // queriesDir
		"internal/db/schema",          // schemaDir
		"${DATABASE_URL}",             // databaseURL
		"enterprise",                  // projectName
		true,                          // strict
	)
}

// DefaultData returns default TemplateData for enterprise template.
func (t *EnterpriseTemplate) DefaultData() generated.TemplateData {
	return t.BuildDefaultData(
		"enterprise",
		"postgresql",
		"${DATABASE_URL}",
		"internal/db",
		"internal/db",
		true, // useManaged
		true, // useUUIDs
		true, // useJSON
		true, // useArrays
		true, // useFullText
		true, // emitPreparedQueries
		true, // emitResultStructPointers
		true, // emitParamsStructPointers
		true, // noSelectStar
		true, // requireWhere
		true, // requireLimit
	)
}

// RequiredFeatures returns which features this template requires.
func (t *EnterpriseTemplate) RequiredFeatures() []string {
	return []string{"emit_interface", "prepared_queries", "json_tags", "strict_checks"}
}