package templates

import (
	"github.com/LarsArtmann/SQLC-Wizzard/pkg/config"
	"github.com/samber/lo"
)

// MicroserviceTemplate generates sqlc config for microservice projects
type MicroserviceTemplate struct{}

// NewMicroserviceTemplate creates a new microservice template
func NewMicroserviceTemplate() *MicroserviceTemplate {
	return &MicroserviceTemplate{}
}

// Name returns the template name
func (t *MicroserviceTemplate) Name() string {
	return string(ProjectTypeMicroservice)
}

// Description returns a human-readable description
func (t *MicroserviceTemplate) Description() string {
	return "Single database, container-optimized configuration for API services and microservices"
}

// Generate creates a SqlcConfig from template data
func (t *MicroserviceTemplate) Generate(data TemplateData) (*config.SqlcConfig, error) {
	// Set defaults
	if data.PackageName == "" {
		data.PackageName = "db"
	}
	if data.OutputDir == "" {
		data.OutputDir = "internal/db"
	}
	if data.QueriesDir == "" {
		data.QueriesDir = "internal/db/queries"
	}
	if data.SchemaDir == "" {
		data.SchemaDir = "internal/db/schema"
	}
	if data.DatabaseURL == "" {
		data.DatabaseURL = "${DATABASE_URL}"
	}

	// Determine SQL package based on database type
	sqlPackage := t.getSQLPackage(data.Database)

	// Build config
	cfg := &config.SqlcConfig{
		Version: "2",
		SQL: []config.SQLConfig{
			{
				Name:   lo.Ternary(data.ProjectName != "", data.ProjectName, "service"),
				Engine: string(data.Database),
				Queries: []string{
					data.QueriesDir,
				},
				Schema: []string{
					data.SchemaDir,
				},
				StrictFunctionChecks: lo.ToPtr(data.SafetyRules.StrictFunctions),
				StrictOrderBy:        lo.ToPtr(data.SafetyRules.StrictOrderBy),
				Database: &config.DatabaseConfig{
					URI:     data.DatabaseURL,
					Managed: data.UseManagedDB,
				},
				Gen: config.GenConfig{
					Go: &config.GoGenConfig{
						Package:                   data.PackageName,
						Out:                       data.OutputDir,
						SQLPackage:                sqlPackage,
						BuildTags:                 t.getBuildTags(data),
						EmitInterface:             data.Features.EmitInterface,
						EmitJSONTags:              data.Features.JSONTags,
						EmitDBTags:                data.Features.DBTags,
						EmitPreparedQueries:       data.Features.PreparedQueries,
						EmitExactTableNames:       data.Features.ExactTableNames,
						EmitEmptySlices:           data.Features.EmptySlices,
						EmitExportedQueries:       true,
						EmitPointersForNullTypes:  false,
						JSONTagsCaseStyle:         "camel",
						OmitUnusedStructs:         data.Features.OmitUnusedStructs,
						OmitSQLCVersion:           false,
						OutputDBFileName:          "db.go",
						OutputModelsFileName:      "models.go",
						OutputQuerierFileName:     "querier.go",
						OutputCopyfromFileName:    "copyfrom.go",
						OutputBatchFileName:       "batch.go",
						Overrides:                 t.getTypeOverrides(data),
						Rename:                    t.getRenameRules(),
					},
				},
				Rules: t.getSafetyRules(data.SafetyRules),
			},
		},
	}

	return cfg, nil
}

// DefaultData returns default TemplateData for microservice template
func (t *MicroserviceTemplate) DefaultData() TemplateData {
	return TemplateData{
		ProjectType:    ProjectTypeMicroservice,
		Database:       DatabaseTypePostgreSQL,
		UseManagedDB:   true,
		PackageName:    "db",
		OutputDir:      "internal/db",
		QueriesDir:     "internal/db/queries",
		SchemaDir:      "internal/db/schema",
		DatabaseURL:    "${DATABASE_URL}",
		Features:       DefaultFeatures(),
		SafetyRules:    DefaultSafetyRules(),
	}
}

// RequiredFeatures returns which features this template requires
func (t *MicroserviceTemplate) RequiredFeatures() []string {
	return []string{"emit_interface", "prepared_queries", "json_tags"}
}

// getSQLPackage returns the appropriate SQL package for the database
func (t *MicroserviceTemplate) getSQLPackage(db DatabaseType) string {
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

// getBuildTags returns appropriate build tags
func (t *MicroserviceTemplate) getBuildTags(data TemplateData) string {
	switch data.Database {
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

// getTypeOverrides returns database-specific type overrides
func (t *MicroserviceTemplate) getTypeOverrides(data TemplateData) []config.Override {
	var overrides []config.Override

	switch data.Database {
	case DatabaseTypePostgreSQL:
		// UUID support
		if data.Features.UUIDs {
			overrides = append(overrides, config.Override{
				DBType:       "uuid",
				GoType:       "UUID",
				GoImportPath: "github.com/google/uuid",
			})
		}

		// JSONB support
		if data.Features.JSON {
			overrides = append(overrides, config.Override{
				DBType:       "jsonb",
				GoType:       "RawMessage",
				GoImportPath: "encoding/json",
			})
		}

	case DatabaseTypeMySQL:
		// JSON support for MySQL
		if data.Features.JSON {
			overrides = append(overrides, config.Override{
				DBType:       "json",
				GoType:       "RawMessage",
				GoImportPath: "encoding/json",
			})
		}

	case DatabaseTypeSQLite:
		// SQLite specific overrides if needed
	}

	return overrides
}

// getRenameRules returns common rename rules for better Go naming
func (t *MicroserviceTemplate) getRenameRules() map[string]string {
	return map[string]string{
		"id":   "ID",
		"uuid": "UUID",
		"url":  "URL",
		"uri":  "URI",
		"json": "JSON",
		"api":  "API",
		"http": "HTTP",
		"db":   "DB",
	}
}

// getSafetyRules converts SafetyRules to config.RuleConfig
func (t *MicroserviceTemplate) getSafetyRules(rules SafetyRules) []config.RuleConfig {
	var configRules []config.RuleConfig

	if rules.NoSelectStar {
		configRules = append(configRules, config.RuleConfig{
			Name:    "no-select-star",
			Rule:    "!query.sql.contains(\"SELECT *\")",
			Message: "SELECT * is not allowed for security and performance reasons",
		})
	}

	if rules.RequireWhere {
		configRules = append(configRules, config.RuleConfig{
			Name:    "require-where-delete",
			Rule:    "query.cmd != \"exec\" || !query.sql.contains(\"DELETE\") || query.sql.contains(\"WHERE\")",
			Message: "DELETE statements must include a WHERE clause",
		})
	}

	if rules.NoDropTable {
		configRules = append(configRules, config.RuleConfig{
			Name:    "no-drop-table",
			Rule:    "!query.sql.contains(\"DROP TABLE\")",
			Message: "DROP TABLE statements are not allowed",
		})
	}

	if rules.RequireLimit {
		configRules = append(configRules, config.RuleConfig{
			Name:    "require-limit-select",
			Rule:    "query.cmd == \"many\" implies query.sql.contains(\"LIMIT\")",
			Message: "SELECT queries that return multiple rows should include a LIMIT clause",
		})
	}

	return configRules
}
