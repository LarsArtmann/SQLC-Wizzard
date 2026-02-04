package templates

import (
	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/validation"
	"github.com/LarsArtmann/SQLC-Wizzard/pkg/config"
	"github.com/samber/lo"
)

// AnalyticsTemplate generates sqlc config for analytics and data warehouse projects.
type AnalyticsTemplate struct{}

// NewAnalyticsTemplate creates a new analytics template.
func NewAnalyticsTemplate() *AnalyticsTemplate {
	return &AnalyticsTemplate{}
}

// Name returns the template name.
func (t *AnalyticsTemplate) Name() string {
	return "analytics"
}

// Description returns a human-readable description.
func (t *AnalyticsTemplate) Description() string {
	return "Optimized for data analytics and reporting with full-text search and array support"
}

// Generate creates a SqlcConfig from template data.
func (t *AnalyticsTemplate) Generate(data generated.TemplateData) (*config.SqlcConfig, error) {
	// Set defaults
	packageConfig := data.Package
	if packageConfig.Name == "" {
		packageConfig.Name = "analytics"
		data.Package = packageConfig
	}
	if packageConfig.Path == "" {
		packageConfig.Path = "internal/analytics"
		data.Package = packageConfig
	}

	outputConfig := data.Output
	if outputConfig.BaseDir == "" {
		outputConfig.BaseDir = "internal/analytics"
		data.Output = outputConfig
	}
	if outputConfig.QueriesDir == "" {
		outputConfig.QueriesDir = "internal/analytics/queries"
		data.Output = outputConfig
	}
	if outputConfig.SchemaDir == "" {
		outputConfig.SchemaDir = "internal/analytics/schema"
		data.Output = outputConfig
	}

	databaseConfig := data.Database
	if databaseConfig.URL == "" {
		databaseConfig.URL = "${ANALYTICS_DATABASE_URL}"
	}

	// Determine SQL package based on database type
	sqlPackage := t.getSQLPackage(databaseConfig.Engine)

	// Build config
	cfg := &config.SqlcConfig{
		Version: "2",
		SQL: []config.SQLConfig{
			{
				Name:                 lo.Ternary(data.ProjectName != "", data.ProjectName, "analytics"),
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

// DefaultData returns default TemplateData for analytics template.
func (t *AnalyticsTemplate) DefaultData() TemplateData {
	return generated.TemplateData{
		ProjectName: "",
		ProjectType: MustNewProjectType("analytics"),

		Package: generated.PackageConfig{
			Name: "analytics",
			Path: "internal/analytics",
		},

		Database: generated.DatabaseConfig{
			Engine:      MustNewDatabaseType("postgresql"),
			URL:         "${ANALYTICS_DATABASE_URL}",
			UseManaged:  false,
			UseUUIDs:    false,
			UseJSON:     true,
			UseArrays:   true,
			UseFullText: true,
		},

		Output: generated.OutputConfig{
			BaseDir:    "internal/analytics",
			QueriesDir: "internal/analytics/queries",
			SchemaDir:  "internal/analytics/schema",
		},

		Validation: generated.ValidationConfig{
			StrictFunctions: true,
			StrictOrderBy:   true,
			EmitOptions: generated.EmitOptions{
				EmitJSONTags:             true,
				EmitPreparedQueries:      false,
				EmitInterface:            true,
				EmitEmptySlices:          true,
				EmitResultStructPointers: false,
				EmitParamsStructPointers: false,
				EmitEnumValidMethod:      false,
				EmitAllEnumValues:        true,
				JSONTagsCaseStyle:        "snake",
			},
			SafetyRules: generated.SafetyRules{
				NoSelectStar: false,
				RequireWhere: false,
				NoDropTable:  true,
				NoTruncate:   true,
				RequireLimit: true,
				Rules:        []generated.SafetyRule{},
			},
		},
	}
}

// RequiredFeatures returns which features this template requires.
func (t *AnalyticsTemplate) RequiredFeatures() []string {
	return []string{"emit_interface", "json_tags", "full_text_search"}
}

// buildGoGenConfig builds the GoGenConfig from template data.
func (t *AnalyticsTemplate) buildGoGenConfig(data generated.TemplateData, sqlPackage string) *config.GoGenConfig {
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
func (t *AnalyticsTemplate) getSQLPackage(db DatabaseType) string {
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
func (t *AnalyticsTemplate) getBuildTags(data TemplateData) string {
	switch data.Database.Engine {
	case DatabaseTypePostgreSQL:
		return "postgres,pgx,analytics"
	case DatabaseTypeMySQL:
		return "mysql,analytics"
	case DatabaseTypeSQLite:
		return "sqlite,analytics"
	default:
		return "analytics"
	}
}

// getTypeOverrides returns database-specific type overrides.
func (t *AnalyticsTemplate) getTypeOverrides(data TemplateData) []config.Override {
	var overrides []config.Override

	switch data.Database.Engine {
	case DatabaseTypePostgreSQL:
		// JSONB support for analytics data
		if data.Database.UseJSON {
			overrides = append(overrides, config.Override{
				DBType:       "jsonb",
				GoType:       "RawMessage",
				GoImportPath: "encoding/json",
			})
		}

		// Array support for analytics
		if data.Database.UseArrays {
			overrides = append(overrides, config.Override{
				DBType:       "_text",
				GoType:       "[]string",
				GoImportPath: "",
				Nullable:     true,
			})
			overrides = append(overrides, config.Override{
				DBType:       "_bigint",
				GoType:       "[]int64",
				GoImportPath: "",
				Nullable:     true,
			})
		}

		// Full text search support
		if data.Database.UseFullText {
			overrides = append(overrides, config.Override{
				DBType:       "tsvector",
				GoType:       "string",
				GoImportPath: "",
			})
			overrides = append(overrides, config.Override{
				DBType:       "tsquery",
				GoType:       "string",
				GoImportPath: "",
			})
		}

	case DatabaseTypeMySQL:
		// JSON support for analytics data
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
func (t *AnalyticsTemplate) getRenameRules() map[string]string {
	return map[string]string{
		"id":   "ID",
		"json": "JSON",
		"api":  "API",
		"http": "HTTP",
	}
}
