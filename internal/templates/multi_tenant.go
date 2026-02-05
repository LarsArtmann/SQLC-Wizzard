package templates

import (
	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/validation"
	"github.com/LarsArtmann/SQLC-Wizzard/pkg/config"
	"github.com/samber/lo"
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
	// Multi-tenant-specific defaults
	if data.Package.Name == "" {
		data.Package.Name = "db"
	}
	if data.Package.Path == "" {
		data.Package.Path = "internal/db"
	}
	if data.Output.BaseDir == "" {
		data.Output.BaseDir = "internal/db"
	}
	if data.Output.QueriesDir == "" {
		data.Output.QueriesDir = "internal/db/queries"
	}
	if data.Output.SchemaDir == "" {
		data.Output.SchemaDir = "internal/db/schema"
	}
	if data.Database.URL == "" {
		data.Database.URL = "${DATABASE_URL}"
	}

	// Build config
	cfg := &config.SqlcConfig{
		Version: "2",
		SQL: []config.SQLConfig{
			{
				Name:                 lo.Ternary(data.ProjectName != "", data.ProjectName, "multi-tenant"),
				Engine:               string(data.Database.Engine),
				Queries:              config.NewPathOrPaths([]string{data.Output.QueriesDir}),
				Schema:               config.NewPathOrPaths([]string{data.Output.SchemaDir}),
				StrictFunctionChecks: lo.ToPtr(true),
				StrictOrderBy:        lo.ToPtr(true),
				Database: &config.DatabaseConfig{
					URI:     data.Database.URL,
					Managed: data.Database.UseManaged,
				},
				Gen: config.GenConfig{
					Go: t.buildGoGenConfig(data),
				},
				Rules: []config.RuleConfig{},
			},
		},
	}

	// Apply emit options using type-safe helper function
	config.ApplyEmitOptions(&data.Validation.EmitOptions, cfg.SQL[0].Gen.Go)

	// Convert rule types using the centralized transformer
	transformer := validation.NewRuleTransformer()
	rules := transformer.TransformSafetyRules(&data.Validation.SafetyRules)
	configRules := lo.Map(rules, func(r generated.RuleConfig, _ int) config.RuleConfig {
		return config.RuleConfig{
			Name:    r.Name,
			Rule:    r.Rule,
			Message: r.Message,
		}
	})
	cfg.SQL[0].Rules = configRules

	return cfg, nil
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
