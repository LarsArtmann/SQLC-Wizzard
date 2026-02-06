package templates

import (
	"github.com/LarsArtmann/SQLC-Wizzard/generated"
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