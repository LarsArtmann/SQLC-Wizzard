package templates

import (
	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/pkg/config"
)

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

// GetSQLPackage returns the appropriate SQL package for the database.
func (t *BaseTemplate) GetSQLPackage(db DatabaseType) string {
	switch db {
	case DatabaseTypePostgreSQL:
		return "database/sql"
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
	}
}
