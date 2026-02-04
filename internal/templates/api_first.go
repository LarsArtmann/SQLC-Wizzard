package templates

import (
	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/validation"
	"github.com/LarsArtmann/SQLC-Wizzard/pkg/config"
	"github.com/samber/lo"
)

// APIFirstTemplate generates sqlc config for API-first projects.
type APIFirstTemplate struct{}

// NewAPIFirstTemplate creates a new API-first template.
func NewAPIFirstTemplate() *APIFirstTemplate {
	return &APIFirstTemplate{}
}

// Name returns the template name.
func (t *APIFirstTemplate) Name() string {
	return "api-first"
}

// Description returns a human-readable description.
func (t *APIFirstTemplate) Description() string {
	return "Optimized for REST/GraphQL API development with JSON support and camelCase naming"
}

// Generate creates a SqlcConfig from template data.
func (t *APIFirstTemplate) Generate(data generated.TemplateData) (*config.SqlcConfig, error) {
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
				Name:                 lo.Ternary(data.ProjectName != "", data.ProjectName, "api"),
				Engine:               string(databaseConfig.Engine),
				Queries:              config.NewPathOrPaths([]string{outputConfig.QueriesDir}),
				Schema:               config.NewPathOrPaths([]string{outputConfig.SchemaDir}),
				StrictFunctionChecks: lo.ToPtr(false),
				StrictOrderBy:        lo.ToPtr(false),
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

// DefaultData returns default TemplateData for API-first template.
func (t *APIFirstTemplate) DefaultData() TemplateData {
	return generated.TemplateData{
		ProjectName: "",
		ProjectType: MustNewProjectType("api-first"),

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
			StrictFunctions: false,
			StrictOrderBy:   false,
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
				RequireLimit: false,
				Rules:        []generated.SafetyRule{},
			},
		},
	}
}

// RequiredFeatures returns which features this template requires.
func (t *APIFirstTemplate) RequiredFeatures() []string {
	return []string{"emit_interface", "prepared_queries", "json_tags", "camel_case"}
}

// buildGoGenConfig builds the GoGenConfig from template data.
func (t *APIFirstTemplate) buildGoGenConfig(data generated.TemplateData, sqlPackage string) *config.GoGenConfig {
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
func (t *APIFirstTemplate) getSQLPackage(db DatabaseType) string {
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
func (t *APIFirstTemplate) getBuildTags(data TemplateData) string {
	switch data.Database.Engine {
	case DatabaseTypePostgreSQL:
		return "postgres,pgx"
	case DatabaseTypeMySQL:
		return "mysql"
	case DatabaseTypeSQLite:
		return "sqlite"
	default:
		return ""
	}
}

// getTypeOverrides returns database-specific type overrides.
func (t *APIFirstTemplate) getTypeOverrides(data TemplateData) []config.Override {
	var overrides []config.Override

	switch data.Database.Engine {
	case DatabaseTypePostgreSQL:
		// UUID support
		if data.Database.UseUUIDs {
			overrides = append(overrides, config.Override{
				DBType:       "uuid",
				GoType:       "UUID",
				GoImportPath: "github.com/google/uuid",
			})
		}

		// JSONB support for APIs
		if data.Database.UseJSON {
			overrides = append(overrides, config.Override{
				DBType:       "jsonb",
				GoType:       "RawMessage",
				GoImportPath: "encoding/json",
			})
		}

	case DatabaseTypeMySQL:
		// JSON support for APIs
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
func (t *APIFirstTemplate) getRenameRules() map[string]string {
	return map[string]string{
		"id":   "ID",
		"uuid": "UUID",
		"url":  "URL",
		"uri":  "URI",
		"json": "JSON",
		"api":  "API",
		"http": "HTTP",
		"db":   "DB",
		"created_at": "CreatedAt",
		"updated_at": "UpdatedAt",
	}
}
