package templates

import (
	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/pkg/config"
)

// MultiTenantTemplate generates sqlc config for multi-tenant SaaS applications.
type MultiTenantTemplate struct {
	ConfiguredTemplate
}

// NewMultiTenantTemplate creates a new multi-tenant template.
func NewMultiTenantTemplate() *MultiTenantTemplate {
	t := &MultiTenantTemplate{}
	t.ConfiguredTemplate = ConfiguredTemplate{
		TemplateName:        "multi-tenant",
		TemplateDescription: "Optimized for SaaS multi-tenant architecture with tenant isolation and strict safety rules",
		DefaultPackageName:  "multi-tenant",
		DefaultProjectName:  "multi-tenant-app",
		StrictMode:          true,
		ProjectType:         "multi-tenant",
		DbEngine:            "postgresql",
		PackagePath:         "internal/db",
		BaseOutput:          "internal/db",
		UseManaged:          true,
		UseUUIDs:            true,
		UseJSON:             true,
		UseArrays:           true,
		UseFullText:         false,
		EmitJSONTags:        true,
		EmitPreparedQueries: true,
		EmitInterface:       true,
		EmitEmptySlices:     false,
		EmitResultStructPointers: true,
		EmitParamsStructPointers: true,
		EmitEnumValidMethod:      false,
		EmitAllEnumValues:        false,
		JSONTagsCaseStyle:        "camel",
		StrictFunctions:         true,
		StrictOrderBy:           true,
		NoSelectStar:             true,
		RequireWhere:             true,
		NoDropTable:             true,
		NoTruncate:              false,
		RequireLimit:            true,
		Features:                []string{"emit_interface", "prepared_queries", "json_tags", "tenant_isolation", "strict_checks"},
		CustomRenameRules: map[string]string{
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
		},
	}
	return t
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
		Data:               data,
		DefaultName:        "multi-tenant",
		DefaultDatabaseURL: "${DATABASE_URL}",
		Strict:             true,
	}
	cfg, _ := builder.Build()

	// Generate Go config with template-specific settings (including custom rename rules)
	cfg.SQL[0].Gen.Go = t.BuildGoConfigWithOverrides(data)

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
		true,  // noDropTable
		false, // noTruncate
		true,  // requireLimit
	)
}

// RequiredFeatures returns which features this template requires.
func (t *MultiTenantTemplate) RequiredFeatures() []string {
	return []string{"emit_interface", "prepared_queries", "json_tags", "tenant_isolation", "strict_checks"}
}