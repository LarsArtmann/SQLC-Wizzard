package templates_test

import (
	"testing"

	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/templates"
	internal_testing "github.com/LarsArtmann/SQLC-Wizzard/internal/testing"
	"github.com/stretchr/testify/assert"
)

func TestTestingTemplate_Name(t *testing.T) {
	template := &templates.TestingTemplate{}
	assert.Equal(t, "testing", template.Name())
}

func TestTestingTemplate_Description(t *testing.T) {
	template := &templates.TestingTemplate{}
	assert.Contains(t, template.Description(), "test")
}

func TestTestingTemplate_DefaultData(t *testing.T) {
	internal_testing.AssertTemplateDefaultData(t, internal_testing.TemplateTestHelper{
		Template:              &templates.TestingTemplate{},
		ExpectedProjectType:   generated.ProjectType("testing"),
		ExpectedPackageName:   "testdata",
		ExpectedPackagePath:   "testdata/db",
		ExpectedDatabaseType:  generated.DatabaseTypeSQLite,
		ExpectUUID:            false,
		ExpectJSON:            false,
		ExpectArrays:          false,
		ExpectJSONTags:        false,
		ExpectPreparedQueries: false,
	})
}

func TestTestingTemplate_Generate_Basic(t *testing.T) {
	internal_testing.AssertTemplateGenerateBasicWithConfigs(
		t,
		&templates.TestingTemplate{},
		generated.ProjectType("testing"),
		"test-project",
		"sqlite",
		internal_testing.CommonTemplateConfigs.SQLiteMinimal,
	)
}
