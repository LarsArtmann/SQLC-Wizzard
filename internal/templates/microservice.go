package templates

import (
	"github.com/LarsArtmann/SQLC-Wizzard/internal/domain"
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
	packageConfig := data.Package
	if packageConfig.Name == "" {
		packageConfig.Name = "db"
	}
	if packageConfig.Path == "" {
		packageConfig.Path = "db"
	}
	
	outputConfig := data.Output
	if outputConfig.BaseDir == "" {
		outputConfig.BaseDir = "internal/db"
	}
	if outputConfig.QueriesDir == "" {
		outputConfig.QueriesDir = "internal/db/queries"
	}
	if outputConfig.SchemaDir == "" {
		outputConfig.SchemaDir = "internal/db/schema"
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
				Rules: data.Validation.SafetyRules.ToRuleConfigs(),
			},
		},
	}

	return cfg, nil
}

// DefaultData returns default TemplateData for microservice template
func (t *MicroserviceTemplate) DefaultData() TemplateData {
	return TemplateData{
		ProjectName: "",
		ProjectType: MustNewProjectType("microservice"),
		
		Package: PackageConfig{
			Name: "db",
			Path: "internal/db",
		},
		
		Database: DatabaseConfig{
			Engine:      MustNewDatabaseType("postgresql"),
			URL:         "${DATABASE_URL}",
			UseManaged:  true,
			UseUUIDs:    true,
			UseJSON:     true,
			UseArrays:   false,
			UseFullText: false,
		},
		
		Output: OutputConfig{
			BaseDir:    "internal/db",
			QueriesDir:  "internal/db/queries",
			SchemaDir:   "internal/db/schema",
		},
		
		Validation: ValidationConfig{
			StrictFunctions: false,
			StrictOrderBy:   false,
			EmitOptions:    domain.DefaultEmitOptions(),
			SafetyRules:     domain.DefaultSafetyRules(),
		},
	}
}

// RequiredFeatures returns which features this template requires
func (t *MicroserviceTemplate) RequiredFeatures() []string {
	return []string{"emit_interface", "prepared_queries", "json_tags"}
}

// buildGoGenConfig builds the GoGenConfig from template data.
// This eliminates split brain by using Validation.EmitOptions.ApplyToGoGenConfig().
func (t *MicroserviceTemplate) buildGoGenConfig(data TemplateData, sqlPackage string) *config.GoGenConfig {
	cfg := &config.GoGenConfig{
		Package:    data.Package.Name,
		Out:        data.Output.BaseDir,
		SQLPackage: sqlPackage,
		BuildTags:  t.getBuildTags(data),
		Overrides:  t.getTypeOverrides(data),
		Rename:     t.getRenameRules(),
	}

	// Apply emit options (eliminates field-by-field copying!)
	data.EmitOptions.ApplyToGoGenConfig(cfg)

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
		if data.UseJSON {
			overrides = append(overrides, config.Override{
				DBType:       "jsonb",
				GoType:       "RawMessage",
				GoImportPath: "encoding/json",
			})
		}

	case DatabaseTypeMySQL:
		// JSON support for MySQL
		if data.UseJSON {
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

