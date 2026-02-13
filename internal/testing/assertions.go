// Package testing provides test helpers for both regular tests and ginkgo tests
package testing

import (
	"testing"

	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TemplateTestHelper contains parameters for generic template tests using testify.
type TemplateTestHelper struct {
	Template interface {
		DefaultData() generated.TemplateData
		Generate(data generated.TemplateData) (*config.SqlcConfig, error)
	}
	ExpectedProjectType       generated.ProjectType
	ExpectedProjectName       string
	ExpectedEngine            string
	ExpectedPackageName       string                 // defaults to "db" if empty
	ExpectedPackagePath       string                 // defaults to "internal/db" if empty
	ExpectedDatabaseType      generated.DatabaseType // defaults to PostgreSQL if empty
	ExpectUUID                bool
	ExpectJSON                bool
	ExpectArrays              bool
	ExpectFullText            bool // UseFullText for Enterprise and Analytics templates
	ExpectJSONTags            bool
	ExpectInterface           bool   // EmitPreparedQueries for APIFirst, EmitInterface for Library
	ExpectStrictChecks        bool   // StrictFunctionChecks and StrictOrderBy
	ExpectPreparedQueries     bool   // defaults to true if not specified
	ExpectedJSONTagsCaseStyle string // defaults to "camel" if empty
}

// AssertTemplateDefaultData verifies common template DefaultData() expectations.
func AssertTemplateDefaultData(t *testing.T, helper TemplateTestHelper) {
	t.Helper()

	data := helper.Template.DefaultData()

	// Set defaults for optional fields
	expectedPackageName := helper.ExpectedPackageName
	if expectedPackageName == "" {
		expectedPackageName = "db"
	}
	expectedPackagePath := helper.ExpectedPackagePath
	if expectedPackagePath == "" {
		expectedPackagePath = "internal/db"
	}
	expectedDatabaseType := helper.ExpectedDatabaseType
	if expectedDatabaseType == "" {
		expectedDatabaseType = generated.DatabaseTypePostgreSQL
	}
	expectedJSONTagsCaseStyle := helper.ExpectedJSONTagsCaseStyle
	if expectedJSONTagsCaseStyle == "" {
		expectedJSONTagsCaseStyle = "camel"
	}

	assert.Equal(t, helper.ExpectedProjectType, data.ProjectType)
	assert.Equal(t, expectedPackageName, data.Package.Name)
	assert.Equal(t, expectedPackagePath, data.Package.Path)
	assert.Equal(t, expectedDatabaseType, data.Database.Engine)
	assert.Equal(t, helper.ExpectUUID, data.Database.UseUUIDs)
	assert.Equal(t, helper.ExpectJSON, data.Database.UseJSON)
	assert.Equal(t, helper.ExpectArrays, data.Database.UseArrays)
	assert.Equal(t, helper.ExpectFullText, data.Database.UseFullText)
	assert.Equal(t, helper.ExpectJSONTags, data.Validation.EmitOptions.EmitJSONTags)
	if helper.ExpectInterface {
		assert.True(t, data.Validation.EmitOptions.EmitInterface || data.Validation.EmitOptions.EmitPreparedQueries)
	}
	assert.JSONEq(t, expectedJSONTagsCaseStyle, data.Validation.EmitOptions.JSONTagsCaseStyle)

	// Check prepared queries - defaults to true unless explicitly set
	expectedPreparedQueries := helper.ExpectPreparedQueries
	if !expectedPreparedQueries {
		// Only verify if explicitly set to false (like in TestingTemplate)
		assert.False(t, data.Validation.EmitOptions.EmitPreparedQueries)
	}
}

// AssertTemplateGenerateBasic verifies common template Generate() expectations.
func AssertTemplateGenerateBasic(t *testing.T, helper TemplateTestHelper) {
	t.Helper()

	data := helper.Template.DefaultData()
	data.ProjectName = helper.ExpectedProjectName

	result, err := helper.Template.Generate(data)

	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, "2", result.Version)
	assert.Len(t, result.SQL, 1)

	sqlConfig := result.SQL[0]
	assert.Equal(t, helper.ExpectedProjectName, sqlConfig.Name)
	assert.Equal(t, helper.ExpectedEngine, sqlConfig.Engine)
	assert.NotNil(t, sqlConfig.Database)
	if helper.ExpectStrictChecks {
		assert.True(t, *sqlConfig.StrictFunctionChecks)
		assert.True(t, *sqlConfig.StrictOrderBy)
	} else {
		assert.False(t, *sqlConfig.StrictFunctionChecks)
		assert.False(t, *sqlConfig.StrictOrderBy)
	}
}

// AssertTemplateGenerateBasicWithDefaults is a convenience function for testing templates
// that share identical Generate() expectations (UUID, JSON, Arrays, JSONTags, Interface, StrictChecks).
// This eliminates duplicate test code across templates with the same behavior profile.
//
// Parameters:
//   - t: testing.TB
//   - template: The template to test
//   - expectedProjectType: The expected project type enum
//   - expectedProjectName: The expected project name string
//
// Example usage:
//
//	func TestEnterpriseTemplate_Generate_Basic(t *testing.T) {
//	    AssertTemplateGenerateBasicWithDefaults(t,
//	        &templates.EnterpriseTemplate{},
//	        generated.ProjectTypeEnterprise,
//	        "enterprise-service",
//	    )
//	}
func AssertTemplateGenerateBasicWithDefaults(t *testing.T, template interface {
	DefaultData() generated.TemplateData
	Generate(data generated.TemplateData) (*config.SqlcConfig, error)
}, expectedProjectType generated.ProjectType, expectedProjectName string,
) {
	t.Helper()

	data := template.DefaultData()
	data.ProjectName = expectedProjectName

	result, err := template.Generate(data)

	require.NoError(t, err, "Generate should not return an error for %s template", expectedProjectType)
	require.NotNil(t, result, "Generate should return a non-nil config")

	assert.Equal(t, "2", result.Version, "Version should be '2'")
	assert.Len(t, result.SQL, 1, "Should have exactly one SQL configuration")

	sqlConfig := result.SQL[0]
	assert.Equal(t, expectedProjectName, sqlConfig.Name, "SQL config name should match project name")
	assert.Equal(t, "postgresql", sqlConfig.Engine, "Engine should be postgresql")
	assert.NotNil(t, sqlConfig.Database, "Database should be configured")

	// Verify all default safety features are enabled
	assert.True(t, *sqlConfig.StrictFunctionChecks, "StrictFunctionChecks should be enabled")
	assert.True(t, *sqlConfig.StrictOrderBy, "StrictOrderBy should be enabled")
}

// AssertTemplateGenerateBasicWithConfigs is a convenience function for testing templates
// with specific configuration patterns. This eliminates duplicate test code across
// templates that require custom configuration loops.
//
// Parameters:
//   - t: testing.TB
//   - template: The template to test
//   - expectedProjectType: The expected project type enum
//   - expectedProjectName: The expected project name string
//   - expectedEngine: The expected database engine string
//   - commonConfigs: A slice of TemplateTestHelperOption to apply to the helper
//
// Example usage:
//
//	func TestAnalyticsTemplate_Generate_Basic(t *testing.T) {
//	    AssertTemplateGenerateBasicWithConfigs(t,
//	        &templates.AnalyticsTemplate{},
//	        generated.ProjectType("analytics"),
//	        "analytics-service",
//	        "postgresql",
//	        CommonTemplateConfigs.PostgreSQLAnalytics,
//	    )
//	}
func AssertTemplateGenerateBasicWithConfigs(t *testing.T, template interface {
	DefaultData() generated.TemplateData
	Generate(data generated.TemplateData) (*config.SqlcConfig, error)
}, expectedProjectType generated.ProjectType, expectedProjectName, expectedEngine string, commonConfigs []TemplateTestHelperOption,
) {
	t.Helper()

	helper := NewTemplateTestHelper(
		template,
		WithProjectType(expectedProjectType),
		WithProjectName(expectedProjectName),
		WithEngine(expectedEngine),
	)

	for _, opt := range commonConfigs {
		opt(&helper)
	}

	AssertTemplateGenerateBasic(t, helper)
}

// TemplateTestHelperOption is a functional option for creating TemplateTestHelper.
type TemplateTestHelperOption func(*TemplateTestHelper)

// WithProjectType sets the expected project type.
func WithProjectType(pt generated.ProjectType) TemplateTestHelperOption {
	return func(h *TemplateTestHelper) {
		h.ExpectedProjectType = pt
	}
}

// WithProjectName sets the expected project name.
func WithProjectName(name string) TemplateTestHelperOption {
	return func(h *TemplateTestHelper) {
		h.ExpectedProjectName = name
	}
}

// WithEngine sets the expected database engine.
func WithEngine(engine string) TemplateTestHelperOption {
	return func(h *TemplateTestHelper) {
		h.ExpectedEngine = engine
	}
}

// WithPackagePath sets the expected package path.
func WithPackagePath(path string) TemplateTestHelperOption {
	return func(h *TemplateTestHelper) {
		h.ExpectedPackagePath = path
	}
}

// WithDatabaseType sets the expected database type.
func WithDatabaseType(dt generated.DatabaseType) TemplateTestHelperOption {
	return func(h *TemplateTestHelper) {
		h.ExpectedDatabaseType = dt
	}
}

// WithUUID configures UUID expectation.
func WithUUID(enabled bool) TemplateTestHelperOption {
	return func(h *TemplateTestHelper) {
		h.ExpectUUID = enabled
	}
}

// WithJSON configures JSON expectation.
func WithJSON(enabled bool) TemplateTestHelperOption {
	return func(h *TemplateTestHelper) {
		h.ExpectJSON = enabled
	}
}

// WithArrays configures arrays expectation.
func WithArrays(enabled bool) TemplateTestHelperOption {
	return func(h *TemplateTestHelper) {
		h.ExpectArrays = enabled
	}
}

// WithFullText configures full text expectation.
func WithFullText(enabled bool) TemplateTestHelperOption {
	return func(h *TemplateTestHelper) {
		h.ExpectFullText = enabled
	}
}

// WithJSONTags configures JSON tags expectation.
func WithJSONTags(enabled bool) TemplateTestHelperOption {
	return func(h *TemplateTestHelper) {
		h.ExpectJSONTags = enabled
	}
}

// WithInterface configures interface expectation (maps to EmitInterface or EmitPreparedQueries).
func WithInterface(enabled bool) TemplateTestHelperOption {
	return func(h *TemplateTestHelper) {
		h.ExpectInterface = enabled
	}
}

// WithStrictChecks configures strict checks expectation.
func WithStrictChecks(enabled bool) TemplateTestHelperOption {
	return func(h *TemplateTestHelper) {
		h.ExpectStrictChecks = enabled
	}
}

// WithPreparedQueries configures prepared queries expectation.
func WithPreparedQueries(enabled bool) TemplateTestHelperOption {
	return func(h *TemplateTestHelper) {
		h.ExpectPreparedQueries = enabled
	}
}

// WithJSONTagsCaseStyle configures JSON tags case style.
func WithJSONTagsCaseStyle(style string) TemplateTestHelperOption {
	return func(h *TemplateTestHelper) {
		h.ExpectedJSONTagsCaseStyle = style
	}
}

// NewTemplateTestHelper creates a TemplateTestHelper with the given template and options.
// This function reduces duplicate code by providing a consistent way to create test helpers
// with common configuration patterns.
//
// Example usage:
//
//	helper := testing.NewTemplateTestHelper(
//	    &templates.APIFirstTemplate{},
//	    testing.WithProjectType(generated.ProjectTypeAPIFirst),
//	    testing.WithProjectName("api-service"),
//	    testing.WithEngine("postgresql"),
//	    testing.WithUUID(true),
//	    testing.WithJSON(true),
//	    testing.WithArrays(true),
//	    testing.WithJSONTags(true),
//	    testing.WithInterface(true),
//	    testing.WithPreparedQueries(true),
//	)
//	testing.AssertTemplateDefaultData(t, helper)
func NewTemplateTestHelper(template interface {
	DefaultData() generated.TemplateData
	Generate(data generated.TemplateData) (*config.SqlcConfig, error)
}, opts ...TemplateTestHelperOption,
) TemplateTestHelper {
	helper := TemplateTestHelper{
		Template: template,
	}

	for _, opt := range opts {
		opt(&helper)
	}

	return helper
}

// CommonTemplateConfigs provides predefined configurations for common template patterns.
// This reduces duplication when multiple templates share similar configuration expectations.
var CommonTemplateConfigs = struct {
	// PostgreSQLFullFeatures is for templates with UUID, JSON, Arrays, JSONTags, Interface, PreparedQueries
	PostgreSQLFullFeatures []TemplateTestHelperOption
	// PostgreSQLWithStrict is for templates with UUID, JSON, Arrays, JSONTags, Interface, StrictChecks
	PostgreSQLWithStrict []TemplateTestHelperOption
	// PostgreSQLWithFullText is for templates with UUID, JSON, Arrays, FullText, JSONTags, Interface, PreparedQueries
	PostgreSQLWithFullText []TemplateTestHelperOption
	// PostgreSQLAnalytics is for templates with JSON, Arrays, JSONTags, Interface, StrictChecks (no UUID)
	PostgreSQLAnalytics []TemplateTestHelperOption
	// SQLiteMinimal is for templates with no UUID, no JSON, no Arrays, no Interface, no PreparedQueries
	SQLiteMinimal []TemplateTestHelperOption
}{
	PostgreSQLFullFeatures: []TemplateTestHelperOption{
		WithUUID(true),
		WithJSON(true),
		WithArrays(true),
		WithJSONTags(true),
		WithInterface(true),
		WithPreparedQueries(true),
	},
	PostgreSQLWithStrict: []TemplateTestHelperOption{
		WithUUID(true),
		WithJSON(true),
		WithArrays(true),
		WithJSONTags(true),
		WithInterface(true),
		WithStrictChecks(true),
	},
	PostgreSQLWithFullText: []TemplateTestHelperOption{
		WithUUID(true),
		WithJSON(true),
		WithArrays(true),
		WithFullText(true),
		WithJSONTags(true),
		WithInterface(true),
		WithPreparedQueries(true),
	},
	PostgreSQLAnalytics: []TemplateTestHelperOption{
		WithUUID(false),
		WithJSON(true),
		WithArrays(true),
		WithJSONTags(true),
		WithInterface(true),
		WithStrictChecks(true),
	},
	SQLiteMinimal: []TemplateTestHelperOption{
		WithUUID(false),
		WithJSON(false),
		WithArrays(false),
		WithJSONTags(false),
		WithInterface(false),
		WithPreparedQueries(false),
	},
}
