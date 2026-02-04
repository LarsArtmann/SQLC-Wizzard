package templates

import (
	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/validation"
	"github.com/LarsArtmann/SQLC-Wizzard/pkg/config"
	"github.com/samber/lo"
)

// TestingTemplate generates sqlc config for test projects and fixtures.
type TestingTemplate struct{}

// NewTestingTemplate creates a new testing template.
func NewTestingTemplate() *TestingTemplate {
	return &TestingTemplate{}
}

// Name returns the template name.
func (t *TestingTemplate) Name() string {
	return "testing"
}

// Description returns a human-readable description.
func (t *TestingTemplate) Description() string {
	return "Lightweight configuration for test suites and database fixtures"
}

// Generate creates a SqlcConfig from template data.
func (t *TestingTemplate) Generate(data generated.TemplateData) (*config.SqlcConfig, error) {
	// Set defaults
	packageConfig := data.Package
	if packageConfig.Name == "" {
		packageConfig.Name = "testdata"
		data.Package = packageConfig
	}
	if packageConfig.Path == "" {
		packageConfig.Path = "testdata/db"
		data.Package = packageConfig
	}

	outputConfig := data.Output
	if outputConfig.BaseDir == "" {
		outputConfig.BaseDir = "testdata/db"
		data.Output = outputConfig
	}
	if outputConfig.QueriesDir == "" {
		outputConfig.QueriesDir = "testdata/queries"
		data.Output = outputConfig
	}
	if outputConfig.SchemaDir == "" {
		outputConfig.SchemaDir = "testdata/schema"
		data.Output = outputConfig
	}

	databaseConfig := data.Database
	if databaseConfig.URL == "" {
		databaseConfig.URL = "file:testdata/test.db"
	}

	// Determine SQL package based on database type
	sqlPackage := t.getSQLPackage(databaseConfig.Engine)

	// Build config
	cfg := &config.SqlcConfig{
		Version: "2",
		SQL: []config.SQLConfig{
			{
				Name:                 lo.Ternary(data.ProjectName != "", data.ProjectName, "test"),
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

// DefaultData returns default TemplateData for testing template.
func (t *TestingTemplate) DefaultData() TemplateData {
	return generated.TemplateData{
		ProjectName: "",
		ProjectType: MustNewProjectType("testing"),

		Package: generated.PackageConfig{
			Name: "testdata",
			Path: "testdata/db",
		},

		Database: generated.DatabaseConfig{
			Engine:      MustNewDatabaseType("sqlite"),
			URL:         "file:testdata/test.db",
			UseManaged:  false,
			UseUUIDs:    false,
			UseJSON:     false,
			UseArrays:   false,
			UseFullText: false,
		},

		Output: generated.OutputConfig{
			BaseDir:    "testdata/db",
			QueriesDir: "testdata/queries",
			SchemaDir:  "testdata/schema",
		},

		Validation: generated.ValidationConfig{
			StrictFunctions: false,
			StrictOrderBy:   false,
			EmitOptions: generated.EmitOptions{
				EmitJSONTags:             false,
				EmitPreparedQueries:      false,
				EmitInterface:            false,
				EmitEmptySlices:          true,
				EmitResultStructPointers: false,
				EmitParamsStructPointers: false,
				EmitEnumValidMethod:      false,
				EmitAllEnumValues:        false,
				JSONTagsCaseStyle:        "snake",
			},
			SafetyRules: generated.SafetyRules{
				NoSelectStar: false,
				RequireWhere: false,
				NoDropTable:  true,
				NoTruncate:   true,
				RequireLimit: false,
				Rules:        []generated.SafetyRule{},
			},
		},
	}
}

// RequiredFeatures returns which features this template requires.
func (t *TestingTemplate) RequiredFeatures() []string {
	return []string{"empty_slices"}
}

// buildGoGenConfig builds the GoGenConfig from template data.
func (t *TestingTemplate) buildGoGenConfig(data generated.TemplateData, sqlPackage string) *config.GoGenConfig {
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
func (t *TestingTemplate) getSQLPackage(db DatabaseType) string {
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

// getBuildTags returns appropriate build tags.
func (t *TestingTemplate) getBuildTags(data TemplateData) string {
	switch data.Database.Engine {
	case DatabaseTypePostgreSQL:
		return "test,postgres"
	case DatabaseTypeMySQL:
		return "test,mysql"
	case DatabaseTypeSQLite:
		return "test,sqlite"
	default:
		return "test"
	}
}

// getTypeOverrides returns database-specific type overrides.
func (t *TestingTemplate) getTypeOverrides(data TemplateData) []config.Override {
	var overrides []config.Override

	// Testing typically uses simple types
	switch data.Database.Engine {
	case DatabaseTypePostgreSQL:
		if data.Database.UseUUIDs {
			overrides = append(overrides, config.Override{
				DBType:       "uuid",
				GoType:       "UUID",
				GoImportPath: "github.com/google/uuid",
			})
		}
	}

	return overrides
}

// getRenameRules returns common rename rules for better Go naming.
func (t *TestingTemplate) getRenameRules() map[string]string {
	return map[string]string{
		"id": "ID",
	}
}
