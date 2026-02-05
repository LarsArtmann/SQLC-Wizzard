package templates

import (
	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/validation"
	"github.com/LarsArtmann/SQLC-Wizzard/pkg/config"
	"github.com/samber/lo"
)

// TestingTemplate generates sqlc config for test projects and fixtures.
type TestingTemplate struct {
	BaseTemplate
}

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
	// Testing-specific defaults
	if data.Package.Name == "" {
		data.Package.Name = "testdata"
	}
	if data.Package.Path == "" {
		data.Package.Path = "testdata/db"
	}
	if data.Output.BaseDir == "" {
		data.Output.BaseDir = "testdata/db"
	}
	if data.Output.QueriesDir == "" {
		data.Output.QueriesDir = "testdata/queries"
	}
	if data.Output.SchemaDir == "" {
		data.Output.SchemaDir = "testdata/schema"
	}
	if data.Database.URL == "" {
		data.Database.URL = "file:testdata/test.db"
	}

	// Build config
	cfg := &config.SqlcConfig{
		Version: "2",
		SQL: []config.SQLConfig{
			{
				Name:                 lo.Ternary(data.ProjectName != "", data.ProjectName, "test"),
				Engine:               string(data.Database.Engine),
				Queries:              config.NewPathOrPaths([]string{data.Output.QueriesDir}),
				Schema:               config.NewPathOrPaths([]string{data.Output.SchemaDir}),
				StrictFunctionChecks: lo.ToPtr(false),
				StrictOrderBy:        lo.ToPtr(false),
				Database: &config.DatabaseConfig{
					URI:     data.Database.URL,
					Managed: data.Database.UseManaged,
				},
				Gen: config.GenConfig{
					Go: t.buildGoGenConfig(data),
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
func (t *TestingTemplate) buildGoGenConfig(data generated.TemplateData) *config.GoGenConfig {
	sqlPackage := t.GetSQLPackage(data.Database.Engine)
	cfg := t.BuildGoGenConfig(data, sqlPackage)

	// Testing uses minimal rename rules
	cfg.Rename = map[string]string{
		"id": "ID",
	}

	return cfg
}
