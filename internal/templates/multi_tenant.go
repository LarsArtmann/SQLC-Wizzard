package templates

import (
	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/validation"
	"github.com/LarsArtmann/SQLC-Wizzard/pkg/config"
	"github.com/samber/lo"
)

// MultiTenantTemplate generates sqlc config for multi-tenant SaaS applications.
type MultiTenantTemplate struct{}

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
	// Set defaults
	packageConfig := data.Package
	if packageConfig.Name == "" {
		packageConfig.Name = "db"
		data.Package = packageConfig
	}
	if packageConfig.Path == "" {
		packageConfig.Path = "internal/db"
		data.Package = packageConfig
	}

	outputConfig := data.Output
	if outputConfig.BaseDir == "" {
		outputConfig.BaseDir = "internal/db"
		data.Output = outputConfig
	}
	if outputConfig.QueriesDir == "" {
		outputConfig.QueriesDir = "internal/db/queries"
		data.Output = outputConfig
	}
	if outputConfig.SchemaDir == "" {
		outputConfig.SchemaDir = "internal/db/schema"
		data.Output = outputConfig
	}

	databaseConfig := data.Database
	if databaseConfig.URL == "" {
		databaseConfig.URL = "${DATABASE_URL}"
	}

	// Determine SQL package based on database type
	sqlPackage := t.getSQLPackage(databaseConfig.Engine)

	// Build config
	cfg := &config.SqlcConfig{
		Version: "2",
		SQL: []config.SQLConfig{
			{
				Name:                 lo.Ternary(data.ProjectName != "", data.ProjectName, "multi-tenant"),
				Engine:               string(databaseConfig.Engine),
				Queries:              config.NewPathOrPaths([]string{outputConfig.QueriesDir}),
				Schema:               config.NewPathOrPaths([]string{outputConfig.SchemaDir}),
				StrictFunctionChecks: lo.ToPtr(true),
				StrictOrderBy:        lo.ToPtr(true),
				Database: &config.DatabaseConfig{
					URI:     databaseConfig.URL,
					Managed: databaseConfig.UseManaged,
				},
				Gen: config.GenConfig{
					Go: t.buildGoGenConfig(data, sqlPackage),
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
	return []string{"emit_interface", "prepared_queries", "json_tags", "tenant_isolation"}
}

// buildGoGenConfig builds the GoGenConfig from template data.
func (t *MultiTenantTemplate) buildGoGenConfig(data generated.TemplateData, sqlPackage string) *config.GoGenConfig {
	cfg := &config.GoGenConfig{
		Package:    data.Package.Name,
		Out:        data.Output.BaseDir,
		SQLPackage: sqlPackage,
		BuildTags:  t.getBuildTags(data),
		Overrides:  t.getTypeOverrides(data),
		Rename:     t.getRenameRules(),
	}

	return cfg
}

// getSQLPackage returns the appropriate SQL package for the database.
func (t *MultiTenantTemplate) getSQLPackage(db DatabaseType) string {
	switch db {
	case DatabaseTypePostgreSQL:
		return "pgx/v5"
	case DatabaseTypeMySQL:
		return "database/sql"
	case DatabaseTypeSQLite:
		return "database/sql"
	default:
		return "database/sql"
	}
}

// getBuildTags returns appropriate build tags.
func (t *MultiTenantTemplate) getBuildTags(data TemplateData) string {
	switch data.Database.Engine {
	case DatabaseTypePostgreSQL:
		return "postgres,pgx,multitenant"
	case DatabaseTypeMySQL:
		return "mysql,multitenant"
	case DatabaseTypeSQLite:
		return "sqlite,multitenant"
	default:
		return "multitenant"
	}
}

// getTypeOverrides returns database-specific type overrides.
func (t *MultiTenantTemplate) getTypeOverrides(data TemplateData) []config.Override {
	var overrides []config.Override

	switch data.Database.Engine {
	case DatabaseTypePostgreSQL:
		// UUID support for tenant IDs
		if data.Database.UseUUIDs {
			overrides = append(overrides, config.Override{
				DBType:       "uuid",
				GoType:       "UUID",
				GoImportPath: "github.com/google/uuid",
			})
		}

		// JSONB support for metadata
		if data.Database.UseJSON {
			overrides = append(overrides, config.Override{
				DBType:       "jsonb",
				GoType:       "RawMessage",
				GoImportPath: "encoding/json",
			})
		}

		// Array support for multi-tenant patterns
		if data.Database.UseArrays {
			overrides = append(overrides, config.Override{
				DBType:       "_text",
				GoType:       "[]string",
				GoImportPath: "",
				Nullable:     true,
			})
			overrides = append(overrides, config.Override{
				DBType:       "_uuid",
				GoType:       "[]UUID",
				GoImportPath: "github.com/google/uuid",
				Nullable:     true,
			})
		}

	case DatabaseTypeMySQL:
		// JSON support for metadata
		if data.Database.UseJSON {
			overrides = append(overrides, config.Override{
				DBType:       "json",
				GoType:       "RawMessage",
				GoImportPath: "encoding/json",
			})
		}
	}

	return overrides
}

// getRenameRules returns common rename rules for better Go naming.
func (t *MultiTenantTemplate) getRenameRules() map[string]string {
	return map[string]string{
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
}
