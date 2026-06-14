// Package templates provides sqlc configuration templates for different project types.
// It implements the Template interface to support consistent configuration generation.
//
// Architecture Overview
// ====================
//
// This package uses two patterns for implementing templates:
//
// 1. BaseTemplate Direct Embedding (simple, minimal overhead):
//    Used by: MicroserviceTemplate, HobbyTemplate, LibraryTemplate, AnalyticsTemplate,
//              TestingTemplate, MultiTenantTemplate
//
//    Pros:
//    - Direct method implementation, full control over Generate()
//    - No inheritance overhead
//    - Best for unique Generate() implementations
//
//    Cons:
//    - Code duplication in Generate() and DefaultData()
//    - Manual zero-value handling required
//
// 2. ConfiguredTemplate Embedding (reusable, consistent defaults):
//    Used by: EnterpriseTemplate, APIFirstTemplate
//
//    Pros:
//    - Built-in zero-value handling (sensible defaults)
//    - Consistent Generate() behavior via ConfiguredTemplate
//    - Less boilerplate in template implementations
//
//    Cons:
//    - Inheritance from ConfiguredTemplate
//    - May be overkill for simple templates
//
// When to Use Which Pattern
// ==========================
//
// Use ConfiguredTemplate when:
// - Template's Generate() uses mostly default handling
// - You want zero-value safety (sensible defaults)
// - You want consistent behavior across templates
//
// Use BaseTemplate directly when:
// - Generate() needs custom logic not covered by ConfiguredTemplate
// - You need full control over the generation process
// - The template is fundamentally different from ConfiguredTemplate
//
// Interface Compliance
// ====================
//
// All template types are verified at compile time to implement the Template interface.
// This ensures no runtime surprises from missing methods:
//
//	var (
//		_ Template = (*EnterpriseTemplate)(nil)
//		_ Template = (*APIFirstTemplate)(nil)
//		_ Template = (*MicroserviceTemplate)(nil)
//		// ... etc
//	)

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
func (t *BaseTemplate) BuildGoGenConfig(
	data generated.TemplateData,
	sqlPackage string,
) *config.GoGenConfig {
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
		return SQLPackagePostgreSQL
	case DatabaseTypeMySQL:
		return SQLPackageStdlib
	case DatabaseTypeSQLite:
		return SQLPackageStdlib
	default:
		return SQLPackageStdlib
	}
}

// GetBuildTags returns appropriate build tags based on database type.
func (t *BaseTemplate) GetBuildTags(data generated.TemplateData) string {
	switch data.Database.Engine {
	case DatabaseTypePostgreSQL:
		return BuildTagPostgreSQL
	case DatabaseTypeMySQL:
		return BuildTagMySQL
	case DatabaseTypeSQLite:
		return BuildTagSQLite
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

		fallthrough
	case DatabaseTypeMySQL:
		if data.Database.UseJSON {
			overrides = append(overrides, config.Override{
				DBType:       "json",
				GoType:       "RawMessage",
				GoImportPath: "encoding/json",
			})
		}
	case DatabaseTypeSQLite:
		// No SQLite-specific overrides currently
	default:
		// No default overrides
	}

	return overrides
}

// GetRenameRules returns common rename rules for better Go naming.
func (t *BaseTemplate) GetRenameRules() map[string]string {
	return CommonRenameRules()
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
		data.Package.Path = DefaultPackagePath
	}

	if data.Output.BaseDir == "" {
		data.Output.BaseDir = DefaultPackagePath
	}

	if data.Output.QueriesDir == "" {
		data.Output.QueriesDir = DefaultPackagePath + "/queries"
	}

	if data.Output.SchemaDir == "" {
		data.Output.SchemaDir = DefaultPackagePath + "/schema"
	}

	if data.Database.URL == "" {
		data.Database.URL = DefaultDatabaseURL
	}
}

// ApplyValidationRules applies emit options and safety rules to a config.
// This eliminates the duplicated validation code across all templates.
func (t *BaseTemplate) ApplyValidationRules(
	cfg *config.SqlcConfig,
	data generated.TemplateData,
) (*config.SqlcConfig, error) {
	// Apply emit options using type-safe helper function
	if len(cfg.SQL) > 0 {
		config.ApplyEmitOptions(&data.Validation.EmitOptions, cfg.SQL[0].Gen.Go)

		// Convert rule types using the centralized transformer
		cfg.SQL[0].Rules = TransformSafetyRulesToConfig(&data.Validation.SafetyRules)
	}

	return cfg, nil
}
