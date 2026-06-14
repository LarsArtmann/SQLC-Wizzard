package templates_test

import (
	"testing"

	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/templates"
	internal_testing "github.com/LarsArtmann/SQLC-Wizzard/internal/testing"
	"github.com/stretchr/testify/assert"
)

func TestLibraryTemplate_Name(t *testing.T) {
	template := &templates.LibraryTemplate{}
	assert.Equal(t, "library", template.Name())
}

func TestLibraryTemplate_Description(t *testing.T) {
	template := templates.NewLibraryTemplate()
	assert.Contains(t, template.Description(), "Library")
}

func TestLibraryTemplate_DefaultData(t *testing.T) {
	internal_testing.AssertTemplateDefaultData(t, internal_testing.TemplateTestHelper{
		Template:            templates.NewLibraryTemplate(),
		ExpectedProjectType: generated.ProjectType("library"),
		ExpectUUID:          false,
		ExpectJSON:          false,
		ExpectArrays:        false,
		ExpectJSONTags:      true,
		ExpectInterface:     true,
	})
}

func TestLibraryTemplate_Generate_Basic(t *testing.T) {
	internal_testing.AssertTemplateGenerateBasic(
		t,
		internal_testing.NewTemplateTestHelper(
			templates.NewLibraryTemplate(),
			internal_testing.WithProjectType(generated.ProjectTypeLibrary),
			internal_testing.WithProjectName("library-module"),
			internal_testing.WithEngine("postgresql"),
		),
	)
}
