package templates_test

import (
	"testing"

	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/templates"
	internal_testing "github.com/LarsArtmann/SQLC-Wizzard/internal/testing"
	"github.com/stretchr/testify/assert"
)

func TestHobbyTemplate_Name(t *testing.T) {
	template := &templates.HobbyTemplate{}
	assert.Equal(t, "hobby", template.Name())
}

func TestHobbyTemplate_Description(t *testing.T) {
	template := &templates.HobbyTemplate{}
	assert.Contains(t, template.Description(), "hobby")
}

func TestHobbyTemplate_DefaultData(t *testing.T) {
	internal_testing.AssertTemplateDefaultData(t, internal_testing.TemplateTestHelper{
		Template:             &templates.HobbyTemplate{},
		ExpectedProjectType:  generated.ProjectTypeHobby,
		ExpectedPackagePath:  "db",
		ExpectedDatabaseType: generated.DatabaseTypeSQLite,
		ExpectUUID:           false,
		ExpectJSON:           false,
		ExpectArrays:         false,
		ExpectJSONTags:       false,
		ExpectInterface:      false,
	})
}

func TestHobbyTemplate_Generate_Basic(t *testing.T) {
	helper := internal_testing.NewTemplateTestHelper(
		&templates.HobbyTemplate{},
		internal_testing.WithProjectType(generated.ProjectTypeHobby),
		internal_testing.WithProjectName("my-hobby-project"),
		internal_testing.WithEngine("sqlite"),
	)
	helper.ExpectedPackagePath = ""
	helper.ExpectedDatabaseType = generated.DatabaseTypeSQLite
	internal_testing.AssertTemplateGenerateBasic(t, helper)
}
