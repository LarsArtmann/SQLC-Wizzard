package templates

import (
	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/pkg/config"
)

// MultiTenantTemplate generates sqlc config for multi-tenant SaaS applications.
type MultiTenantTemplate struct {
	BaseTemplate
}

// NewMultiTenantTemplate creates a new multi-tenant template.
func NewMultiTenantTemplate() *MultiTenantTemplate {
	return &MultiTenantTemplate{}
}

// Name returns the template name.
func (t *MultiTenantTemplate) Name() string {
	return "multi-tenant"
}

// Description returns a human-readable description.
func (t *MultiTenantTemplate) Description() string {
	return "Optimized for SaaS multi-tenant architecture with tenant isolation and strict safety rules"
}

// Generate creates a SqlcConfig from template data.
func (t *MultiTenantTemplate) Generate(data generated.TemplateData) (*config.SqlcConfig, error) {
	// Apply default values first, before using them in ConfigBuilder
	t.ApplyDefaultValues(&data)

	// Build base config using shared builder
	builder := &ConfigBuilder{
		Data:             data,
		DefaultName:      "multi-tenant",
		DefaultDatabaseURL: "${DATABASE_URL}",
		Strict:           true,
	}
	cfg, _ := builder.Build()

	// Generate Go config with template-specific settings
	cfg.SQL[0].Gen.Go = t.buildGoGenConfig(data)

	// Apply validation rules using base helper
	return t.ApplyValidationRules(cfg, data)
}

// DefaultData returns default TemplateData for multi-tenant template.
func (t *MultiTenantTemplate) DefaultData() generated.TemplateData {
	return t.BuildDefaultData(
		"multi-tenant",
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
		true,  // requireLimit
	)
}

// RequiredFeatures returns which features this template requires.
func (t *MultiTenantTemplate) RequiredFeatures() []string {
	return []string{"emit_interface", "prepared_queries", "json_tags", "tenant_isolation", "strict_checks"}
}

// buildGoGenConfig builds the GoGenConfig from template data.
func (t *MultiTenantTemplate) buildGoGenConfig(data generated.TemplateData) *config.GoGenConfig {
	sqlPackage := t.GetSQLPackage(data.Database.Engine)
	cfg := t.BuildGoGenConfig(data, sqlPackage)

	// Multi-tenant uses extended rename rules
	cfg.Rename = map[string]string{
		"id":     "ID",
		"uuid":   "UUID",
		"tenant": "Tenant",
		"url":    "URL",
		"uri":    "URI",
		"json":   "JSON",
		"api":    "API",
		"http":   "HTTP",
		"db":     "DB",
		"sql":    "SQL",
	}

	return cfg
}
