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
	return createDefaultTemplateData(
		"testing",
		"testdata",
		"testdata/db",
		"file:testdata/test.db",
		"sqlite",
		false,  // useManaged
		false,  // useUUIDs
		false,  // useJSON
		false,  // useArrays
		false,  // useFullText
		"testdata/db",
		false,  // strictFunctions
		false,  // strictOrderBy
		false,  // emitJSONTags
		false,  // emitInterface
		false,  // emitAllEnumValues
		false,  // emitPreparedQueries
		false,  // emitResultPointers
		false,  // emitParamsPointers
		false,  // emitEnumValidMethod
		"snake",
		false,  // noSelectStar
		false,  // requireWhere
		false,  // requireLimit
	)
}

// RequiredFeatures returns which features this template requires.
func (t *TestingTemplate) RequiredFeatures() []string {
	return []string{"empty_slices"}
}
