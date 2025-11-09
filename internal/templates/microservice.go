package templates

import (
	"github.com/LarsArtmann/SQLC-Wizzard/generated"
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
	return "microservice"
}

// Description returns a human-readable description
func (t *MicroserviceTemplate) Description() string {
	return "Single database, container-optimized configuration for API services and microservices"
}

// Generate creates a SqlcConfig from template data
func (t *MicroserviceTemplate) Generate(data generated.TemplateData) (*config.SqlcConfig, error) {
	// Set defaults
	packageConfig := data.Package
	if packageConfig.Name == "" {
		packageConfig.Name = "db"
		data.Package = packageConfig
	}
	if packageConfig.Path == "" {
		packageConfig.Path = "db"
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
				Name:                 lo.Ternary(data.ProjectName != "", data.ProjectName, "service"),
				Engine:               string(databaseConfig.Engine),
				Queries:              config.NewPathOrPaths([]string{outputConfig.QueriesDir}),
				Schema:               config.NewPathOrPaths([]string{outputConfig.SchemaDir}),
				StrictFunctionChecks: lo.ToPtr(data.Validation.StrictFunctions),
				StrictOrderBy:        lo.ToPtr(data.Validation.StrictOrderBy),
				Database: &config.DatabaseConfig{
					URI:     databaseConfig.URL,
					Managed: databaseConfig.UseManaged,
				},
				Gen: config.GenConfig{
					Go: t.buildGoGenConfig(data, sqlPackage),
				},
				Rules: []config.RuleConfig{}, // Will be set after conversion
			},
		},
	}

	// Apply emit options (eliminates field-by-field copying!)
	data.Validation.EmitOptions.ApplyToGoGenConfig(&cfg.SQL[0].Gen.Go)

	// Convert rule types
	rules := data.Validation.SafetyRules.ToRuleConfigs()
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

// DefaultData returns default TemplateData for microservice template
func (t *MicroserviceTemplate) DefaultData() TemplateData {
	return generated.TemplateData{
		ProjectName: "",
		ProjectType: MustNewProjectType("microservice"),

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
			UseArrays:   false,
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
			EmitOptions:     generated.DefaultEmitOptions(),
			SafetyRules:     generated.DefaultSafetyRules(),
		},
	}
}

// RequiredFeatures returns which features this template requires
func (t *MicroserviceTemplate) RequiredFeatures() []string {
	return []string{"emit_interface", "prepared_queries", "json_tags"}
}

// buildGoGenConfig builds the GoGenConfig from template data.
// This eliminates split brain by using Validation.EmitOptions.ApplyToGoGenConfig().
func (t *MicroserviceTemplate) buildGoGenConfig(data generated.TemplateData, sqlPackage string) *config.GoGenConfig {
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

// getTypeOverrides returns database-specific type overrides
func (t *MicroserviceTemplate) getTypeOverrides(data TemplateData) []config.Override {
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

		// JSONB support
		if data.Database.UseJSON {
			overrides = append(overrides, config.Override{
				DBType:       "jsonb",
				GoType:       "RawMessage",
				GoImportPath: "encoding/json",
			})
		}

	case DatabaseTypeMySQL:
		// JSON support for MySQL
		if data.Database.UseJSON {
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
