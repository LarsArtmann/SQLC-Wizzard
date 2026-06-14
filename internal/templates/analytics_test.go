package templates_test

import (
	"testing"

	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/templates"
	internal_testing "github.com/LarsArtmann/SQLC-Wizzard/internal/testing"
	"github.com/stretchr/testify/assert"
)

func TestAnalyticsTemplate_Name(t *testing.T) {
	template := &templates.AnalyticsTemplate{}
	assert.Equal(t, "analytics", template.Name())
}

func TestAnalyticsTemplate_Description(t *testing.T) {
	template := &templates.AnalyticsTemplate{}
	assert.Contains(t, template.Description(), "analytics")
}

func TestAnalyticsTemplate_DefaultData(t *testing.T) {
	internal_testing.AssertTemplateDefaultData(t, internal_testing.TemplateTestHelper{
		Template:                  &templates.AnalyticsTemplate{},
		ExpectedProjectType:       generated.ProjectType("analytics"),
		ExpectedPackageName:       "analytics",
		ExpectedPackagePath:       "internal/analytics",
		ExpectUUID:                false,
		ExpectJSON:                true,
		ExpectArrays:              true,
		ExpectFullText:            true,
		ExpectJSONTags:            true,
		ExpectInterface:           true,
		ExpectStrictChecks:        true,
		ExpectedJSONTagsCaseStyle: "camel",
		ExpectPreparedQueries:     false,
	})
}

func TestAnalyticsTemplate_Generate_Basic(t *testing.T) {
	internal_testing.AssertTemplateGenerateBasicWithConfigs(
		t,
		&templates.AnalyticsTemplate{},
		generated.ProjectType("analytics"),
		"analytics-service",
		"postgresql",
		internal_testing.CommonTemplateConfigs.PostgreSQLAnalytics,
	)
}
