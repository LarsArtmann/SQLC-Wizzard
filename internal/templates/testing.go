package templates

import (
	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/pkg/config"
)

// TestingTemplate generates sqlc config for test projects and fixtures.
type TestingTemplate struct {
	ConfiguredTemplate
}

// NewTestingTemplate creates a new testing template.
func NewTestingTemplate() *TestingTemplate {
	base := NewConfiguredTemplate(
		"testing",
		"Lightweight configuration for test suites and database fixtures",
		"testdata",
		"test",
		false, // strictMode
		"testing",
		"sqlite",
	)

	// Override testing-specific settings
	base.UseManaged = false
	base.UseUUIDs = false
	base.UseJSON = false
	base.UseArrays = false
	base.UseFullText = false
	base.EmitJSONTags = false
	base.EmitInterface = false
	base.EmitEmptySlices = true
	base.EmitPreparedQueries = false
	base.JSONTagsCaseStyle = "camel"
	base.StrictFunctions = false
	base.StrictOrderBy = false
	base.NoSelectStar = false
	base.RequireWhere = false
	base.RequireLimit = false
	base.Features = []string{"empty_slices"}

	return &TestingTemplate{ConfiguredTemplate: base}
}

// Name returns the template name.
func (t *TestingTemplate) Name() string {
	return "testing"
}

// Description returns the template description.
func (t *TestingTemplate) Description() string {
	return "Lightweight configuration for test suites and database fixtures"
}

// Generate creates a SqlcConfig from template data.
func (t *TestingTemplate) Generate(data generated.TemplateData) (*config.SqlcConfig, error) {
	return t.GenerateWithDefaults(
		data,
		"testdata",              // packageName
		"testdata/db",           // packagePath
		"testdata/db",           // baseDir
		"testdata/queries",      // queriesDir
		"testdata/schema",       // schemaDir
		"file:testdata/test.db", // databaseURL
		"test",                  // projectName
		false,                   // strict
	)
}

// DefaultData returns default TemplateData for testing template.
func (t *TestingTemplate) DefaultData() TemplateData {
	return createDefaultTemplateData(
		"testing",
		"testdata",
		"testdata/db",
		"file:testdata/test.db",
		"sqlite",
		false, // useManaged
		false, // useUUIDs
		false, // useJSON
		false, // useArrays
		false, // useFullText
		"testdata/db",
		false, // strictFunctions
		false, // strictOrderBy
		false, // emitJSONTags
		false, // emitInterface
		false, // emitAllEnumValues
		false, // emitPreparedQueries
		true,  // emitEmptySlices
		false, // emitResultPointers
		false, // emitParamsPointers
		false, // emitEnumValidMethod
		"camel",
		false, // noSelectStar
		false, // requireWhere
		true,  // noDropTable
		true,  // noTruncate
		false, // requireLimit
	)
}

// RequiredFeatures returns which features this template requires.
func (t *TestingTemplate) RequiredFeatures() []string {
	return t.Features
}
