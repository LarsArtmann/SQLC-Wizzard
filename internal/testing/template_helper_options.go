// Package testing provides test helpers for both regular tests and ginkgo tests
package testing

import (
	"github.com/LarsArtmann/SQLC-Wizzard/generated"
)

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
