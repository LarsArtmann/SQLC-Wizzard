package templates

import (
	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/pkg/config"
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
	// Apply testing-specific defaults
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

	// Build base config using shared builder
	builder := &ConfigBuilder{
		Data:             data,
		DefaultName:      "test",
		DefaultDatabaseURL: "file:testdata/test.db",
		Strict:           false,
	}
	cfg, _ := builder.Build()

	// Generate Go config with template-specific settings
	cfg.SQL[0].Gen.Go = t.BuildGoConfigWithOverrides(data)

	// Apply validation rules using base helper
	return t.ApplyValidationRules(cfg, data)
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
