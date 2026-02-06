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
func (t *MultiTenantTemplate) DefaultData() TemplateData {
	return generated.TemplateData{
		ProjectName: "",
		ProjectType: MustNewProjectType("multi-tenant"),

		Package: generated.PackageConfig{
			Name: "db",
			Path: "internal/db",
		},

		Database: generated.DatabaseConfig{
			Engine:      MustNewDatabaseType("postgresql"),
			URL:         "${DATABASE_URL}",
			UseManaged:  true,
			UseUUIDs:    true,
			UseJSON:     true,
			UseArrays:   true,
			UseFullText: false,
		},

		Output: generated.OutputConfig{
			BaseDir:    "internal/db",
			QueriesDir: "internal/db/queries",
			SchemaDir:  "internal/db/schema",
		},

		Validation: generated.ValidationConfig{
			StrictFunctions: true,
			StrictOrderBy:   true,
			EmitOptions: generated.EmitOptions{
				EmitJSONTags:             true,
				EmitPreparedQueries:      true,
				EmitInterface:            true,
				EmitEmptySlices:          true,
				EmitResultStructPointers: true,
				EmitParamsStructPointers: true,
				EmitEnumValidMethod:      true,
				EmitAllEnumValues:        true,
				JSONTagsCaseStyle:        "camel",
			},
			SafetyRules: generated.SafetyRules{
				NoSelectStar: true,
				RequireWhere: true,
				NoDropTable:  true,
				NoTruncate:   true,
				RequireLimit: true,
				Rules:        []generated.SafetyRule{},
			},
		},
	}
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
