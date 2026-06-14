// Package templates_test provides basic testing for template types
package templates_test

import (
	"testing"

	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/templates"
	internal_testing "github.com/LarsArtmann/SQLC-Wizzard/internal/testing"
	"github.com/stretchr/testify/assert"
)

func TestMicroserviceTemplate_Name(t *testing.T) {
	template := &templates.MicroserviceTemplate{}
	assert.Equal(t, "microservice", template.Name())
}

func TestMicroserviceTemplate_Description(t *testing.T) {
	template := &templates.MicroserviceTemplate{}
	assert.Contains(t, template.Description(), "microservice")
}

func TestMicroserviceTemplate_DefaultData(t *testing.T) {
	internal_testing.AssertTemplateDefaultData(t, internal_testing.TemplateTestHelper{
		Template:            &templates.MicroserviceTemplate{},
		ExpectedProjectType: generated.ProjectTypeMicroservice,
		ExpectUUID:          true,
		ExpectJSON:          true,
		ExpectArrays:        false,
		ExpectJSONTags:      true,
		ExpectInterface:     true,
	})
}

func TestMicroserviceTemplate_Generate_Basic(t *testing.T) {
	helper := internal_testing.NewTemplateTestHelper(
		&templates.MicroserviceTemplate{},
		internal_testing.WithProjectType(generated.ProjectTypeMicroservice),
		internal_testing.WithProjectName("test-service"),
		internal_testing.WithEngine("postgresql"),
		internal_testing.WithUUID(true),
		internal_testing.WithJSON(true),
		internal_testing.WithJSONTags(true),
		internal_testing.WithInterface(true),
	)
	internal_testing.AssertTemplateGenerateBasic(t, helper)
}
