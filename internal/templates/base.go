package templates

import (
	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/validation"
	"github.com/LarsArtmann/SQLC-Wizzard/pkg/config"
	"github.com/samber/lo"
)

// ConfigBuilder helps construct sqlc configurations with common patterns.
// This eliminates duplication across template implementations.
type ConfigBuilder struct {
	// Data holds the template configuration data.
	Data generated.TemplateData
	// DefaultName is used when ProjectName is empty.
	DefaultName string
	// DefaultDatabaseURL is used when Database.URL is empty.
	DefaultDatabaseURL string
	// Strict controls strict mode settings.
	Strict bool
}

// Build creates a SqlcConfig from the configured values.
func (cb *ConfigBuilder) Build() (*config.SqlcConfig, error) {
	base := &BaseTemplate{}

	config := &config.SqlcConfig{
		Version: "2",
		SQL: []config.SQLConfig{
			{
				Name:                 lo.Ternary(cb.Data.ProjectName != "", cb.Data.ProjectName, cb.DefaultName),
				Engine:               string(cb.Data.Database.Engine),
				Queries:              config.NewPathOrPaths([]string{cb.Data.Output.QueriesDir}),
				Schema:               config.NewPathOrPaths([]string{cb.Data.Output.SchemaDir}),
				StrictFunctionChecks: lo.ToPtr(cb.Strict),
				StrictOrderBy:        lo.ToPtr(cb.Strict),
				Database: &config.DatabaseConfig{
					URI:     lo.Ternary(cb.Data.Database.URL != "", cb.Data.Database.URL, cb.DefaultDatabaseURL),
					Managed: cb.Data.Database.UseManaged,
				},
				Gen: config.GenConfig{
					Go: base.BuildGoGenConfig(cb.Data, base.GetSQLPackage(cb.Data.Database.Engine)),
				},
				Rules: []config.RuleConfig{},
			},
		},
	}

	return config, nil
}

// BaseTemplate provides common functionality for all templates.
// Embed this struct in template implementations to inherit helper methods.
type BaseTemplate struct{}

// BuildGoGenConfig builds the base GoGenConfig from template data.
// This is the foundation method that templates can override or extend.
func (t *BaseTemplate) BuildGoGenConfig(data generated.TemplateData, sqlPackage string) *config.GoGenConfig {
	return &config.GoGenConfig{
		Package:    data.Package.Name,
		Out:        data.Output.BaseDir,
		SQLPackage: sqlPackage,
		BuildTags:  t.GetBuildTags(data),
		Overrides:  t.GetTypeOverrides(data),
		Rename:     t.GetRenameRules(),
	}
}

// GetSQLPackage returns appropriate SQL package for database.
// PostgreSQL uses pgx/v5 for better performance and feature support.
// MySQL and SQLite use database/sql for compatibility.
func (t *BaseTemplate) GetSQLPackage(db generated.DatabaseType) string {
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

// GetBuildTags returns appropriate build tags based on database type.
func (t *BaseTemplate) GetBuildTags(data generated.TemplateData) string {
	switch data.Database.Engine {
	case DatabaseTypePostgreSQL:
		return "postgres"
	case DatabaseTypeMySQL:
		return "mysql"
	case DatabaseTypeSQLite:
		return "sqlite"
	default:
		return ""
	}
}

// GetTypeOverrides returns database-specific type overrides.
func (t *BaseTemplate) GetTypeOverrides(data generated.TemplateData) []config.Override {
	var overrides []config.Override

	switch data.Database.Engine {
	case DatabaseTypePostgreSQL:
		if data.Database.UseUUIDs {
			overrides = append(overrides, config.Override{
				DBType:       "uuid",
				GoType:       "UUID",
				GoImportPath: "github.com/google/uuid",
			})
		}
		if data.Database.UseJSON {
			overrides = append(overrides, config.Override{
				DBType:       "json",
				GoType:       "RawMessage",
				GoImportPath: "encoding/json",
			})
		}
	case DatabaseTypeMySQL:
		if data.Database.UseJSON {
			overrides = append(overrides, config.Override{
				DBType:       "json",
				GoType:       "RawMessage",
				GoImportPath: "encoding/json",
			})
		}
	default:
		// No default overrides
	}

	return overrides
}

// GetRenameRules returns common rename rules for better Go naming.
func (t *BaseTemplate) GetRenameRules() map[string]string {
	return map[string]string{
		"id":   "ID",
		"uuid": "UUID",
		"url":  "URL",
		"uri":  "URI",
		"api":  "API",
		"http": "HTTP",
		"json": "JSON",
		"db":   "DB",
		"otp":  "OTP",
	}
}

// BuildGoConfigWithOverrides builds a GoGenConfig with template-specific overrides.
// Template implementations can override this to provide custom rename rules.
func (t *BaseTemplate) BuildGoConfigWithOverrides(data generated.TemplateData) *config.GoGenConfig {
	sqlPackage := t.GetSQLPackage(data.Database.Engine)
	return t.BuildGoGenConfig(data, sqlPackage)
}

// ApplyDefaultValues sets default values for empty fields in TemplateData.
// This is used by template implementations to ensure consistent default behavior.
func (t *BaseTemplate) ApplyDefaultValues(data *generated.TemplateData) {
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
}

// ApplyValidationRules applies emit options and safety rules to a config.
// This eliminates the duplicated validation code across all templates.
func (t *BaseTemplate) ApplyValidationRules(cfg *config.SqlcConfig, data generated.TemplateData) (*config.SqlcConfig, error) {
	// Apply emit options using type-safe helper function
	if len(cfg.SQL) > 0 {
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
	}

	return cfg, nil
}

// BuildDefaultData creates default TemplateData with the provided parameters.
// This eliminates duplication in template DefaultData() methods by providing
// a template method that accepts the variable configuration values.
func (t *BaseTemplate) BuildDefaultData(
	projectType string,
	useManaged, useUUIDs, useJSON, useArrays, useFullText bool,
	emitPreparedQueries, emitResultStructPointers, emitParamsStructPointers bool,
	noSelectStar, requireWhere, requireLimit bool,
) generated.TemplateData {
	return generated.TemplateData{
		ProjectName: "",
		ProjectType: MustNewProjectType(projectType),

		Package: generated.PackageConfig{
			Name: "db",
			Path: "internal/db",
		},

		Database: generated.DatabaseConfig{
			Engine:      MustNewDatabaseType("postgresql"),
			URL:         "${DATABASE_URL}",
			UseManaged:  useManaged,
			UseUUIDs:    useUUIDs,
			UseJSON:     useJSON,
			UseArrays:   useArrays,
			UseFullText: useFullText,
		},

		Output: generated.OutputConfig{
			BaseDir:    "internal/db",
			QueriesDir: "internal/db/queries",
			SchemaDir:  "internal/db/schema",
		},

		Validation: generated.ValidationConfig{
			StrictFunctions: useManaged, // Match strict mode to UseManaged for consistency
			StrictOrderBy:   useManaged,
			EmitOptions: generated.EmitOptions{
				EmitJSONTags:             true,
				EmitPreparedQueries:      emitPreparedQueries,
				EmitInterface:            true,
				EmitEmptySlices:          true,
				EmitResultStructPointers: emitResultStructPointers,
				EmitParamsStructPointers: emitParamsStructPointers,
				EmitEnumValidMethod:      true,
				EmitAllEnumValues:        true,
				JSONTagsCaseStyle:        "camel",
			},
			SafetyRules: generated.SafetyRules{
				NoSelectStar: noSelectStar,
				RequireWhere: requireWhere,
				NoDropTable:  true,
				NoTruncate:   true,
				RequireLimit: requireLimit,
				Rules:        []generated.SafetyRule{},
			},
		},
	}
}