package templates_test

import (
	"testing"

	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/templates"
	internal_testing "github.com/LarsArtmann/SQLC-Wizzard/internal/testing"
	"github.com/stretchr/testify/assert"
)

func TestEnterpriseTemplate_Name(t *testing.T) {
	template := &templates.EnterpriseTemplate{}
	assert.Equal(t, "enterprise", template.Name())
}

func TestEnterpriseTemplate_Description(t *testing.T) {
	template := &templates.EnterpriseTemplate{}
	assert.Contains(t, template.Description(), "enterprise")
}

func TestEnterpriseTemplate_DefaultData(t *testing.T) {
	internal_testing.AssertTemplateDefaultData(t, internal_testing.TemplateTestHelper{
		Template:              &templates.EnterpriseTemplate{},
		ExpectedProjectType:   generated.ProjectTypeEnterprise,
		ExpectUUID:            true,
		ExpectJSON:            true,
		ExpectArrays:          true,
		ExpectFullText:        true,
		ExpectJSONTags:        true,
		ExpectInterface:       true,
		ExpectPreparedQueries: true,
	})
}

func TestEnterpriseTemplate_Generate_Basic(t *testing.T) {
	internal_testing.AssertTemplateGenerateBasicWithDefaults(
		t,
		&templates.EnterpriseTemplate{},
		generated.ProjectTypeEnterprise,
		"enterprise-service",
	)
}
