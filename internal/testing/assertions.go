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
	ExpectedProjectType generated.ProjectType
	ExpectedProjectName string
	ExpectedEngine      string
	ExpectUUID          bool
	ExpectJSON          bool
	ExpectArrays        bool
	ExpectJSONTags      bool
	ExpectInterface     bool // EmitPreparedQueries for APIFirst, EmitInterface for Library
	ExpectStrictChecks  bool // StrictFunctionChecks and StrictOrderBy
}

// AssertTemplateDefaultData verifies common template DefaultData() expectations.
func AssertTemplateDefaultData(t *testing.T, helper TemplateTestHelper) {
	t.Helper()

	data := helper.Template.DefaultData()

	assert.Equal(t, helper.ExpectedProjectType, data.ProjectType)
	assert.Equal(t, "db", data.Package.Name)
	assert.Equal(t, "internal/db", data.Package.Path)
	assert.Equal(t, generated.DatabaseTypePostgreSQL, data.Database.Engine)
	assert.Equal(t, helper.ExpectUUID, data.Database.UseUUIDs)
	assert.Equal(t, helper.ExpectJSON, data.Database.UseJSON)
	assert.Equal(t, helper.ExpectArrays, data.Database.UseArrays)
	assert.Equal(t, helper.ExpectJSONTags, data.Validation.EmitOptions.EmitJSONTags)
	if helper.ExpectInterface {
		assert.True(t, data.Validation.EmitOptions.EmitInterface || data.Validation.EmitOptions.EmitPreparedQueries)
	}
	assert.Equal(t, "camel", data.Validation.EmitOptions.JSONTagsCaseStyle)
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
