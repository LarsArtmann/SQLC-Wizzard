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
	return t.GenerateWithDefaults(
		data,
		"testdata",
		"testdata/db",
		"testdata/db",
		"testdata/queries",
		"testdata/schema",
		"file:testdata/test.db",
		"test",
		false,
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
		false, // emitResultPointers
		false, // emitParamsPointers
		false, // emitEnumValidMethod
		"camel",
		false, // noSelectStar
		false, // requireWhere
		false, // requireLimit
	)
}

// RequiredFeatures returns which features this template requires.
func (t *TestingTemplate) RequiredFeatures() []string {
	return []string{"empty_slices"}
}
