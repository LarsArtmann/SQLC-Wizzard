package templates_test

import (
	"testing"

	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/templates"
	internal_testing "github.com/LarsArtmann/SQLC-Wizzard/internal/testing"
	"github.com/stretchr/testify/assert"
)

func TestAPIFirstTemplate_Name(t *testing.T) {
	template := &templates.APIFirstTemplate{}
	assert.Equal(t, "api-first", template.Name())
}

func TestAPIFirstTemplate_Description(t *testing.T) {
	template := &templates.APIFirstTemplate{}
	assert.Contains(t, template.Description(), "API")
}

func TestAPIFirstTemplate_DefaultData(t *testing.T) {
	internal_testing.AssertTemplateDefaultData(t, internal_testing.NewTemplateTestHelper(
		&templates.APIFirstTemplate{},
		internal_testing.WithProjectType(generated.ProjectType("api-first")),
		internal_testing.WithProjectName("api-service"),
		internal_testing.WithEngine("postgresql"),
		internal_testing.WithUUID(true),
		internal_testing.WithJSON(true),
		internal_testing.WithArrays(true),
		internal_testing.WithJSONTags(true),
		internal_testing.WithInterface(true),
		internal_testing.WithPreparedQueries(true),
	))
}

func TestAPIFirstTemplate_Generate_Basic(t *testing.T) {
	helper := internal_testing.NewTemplateTestHelper(
		&templates.APIFirstTemplate{},
		internal_testing.WithProjectType(generated.ProjectTypeAPIFirst),
		internal_testing.WithProjectName("api-service"),
		internal_testing.WithEngine("postgresql"),
	)
	for _, opt := range internal_testing.CommonTemplateConfigs.PostgreSQLFullFeatures {
		opt(&helper)
	}

	internal_testing.AssertTemplateGenerateBasic(t, helper)
}
