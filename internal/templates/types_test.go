// Package templates_test provides basic testing for template types
package templates_test

import (
	"testing"

	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/apperrors"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/templates"
	internal_testing "github.com/LarsArtmann/SQLC-Wizzard/internal/testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewProjectType_Valid(t *testing.T) {
	result, err := templates.NewProjectType("microservice")

	require.NoError(t, err)
	assert.Equal(t, generated.ProjectTypeMicroservice, result)
}

func TestNewProjectType_Invalid(t *testing.T) {
	result, err := templates.NewProjectType("invalid")

	require.Error(t, err)
	assert.Empty(t, result)
	assert.True(t, apperrors.Is(err, apperrors.ErrInvalidType))
}

func TestNewDatabaseType_Valid(t *testing.T) {
	result, err := templates.NewDatabaseType("postgresql")

	require.NoError(t, err)
	assert.Equal(t, generated.DatabaseTypePostgreSQL, result)
}

func TestNewDatabaseType_Invalid(t *testing.T) {
	result, err := templates.NewDatabaseType("invalid")

	require.Error(t, err)
	assert.Empty(t, result)
	assert.True(t, apperrors.Is(err, apperrors.ErrInvalidType))
}

func TestMustNewProjectType_Valid(t *testing.T) {
	result := templates.MustNewProjectType("microservice")
	assert.Equal(t, generated.ProjectTypeMicroservice, result)
}

func TestMustNewProjectType_Invalid(t *testing.T) {
	assert.Panics(t, func() {
		templates.MustNewProjectType("invalid")
	})
}

func TestMustNewDatabaseType_Valid(t *testing.T) {
	result := templates.MustNewDatabaseType("postgresql")
	assert.Equal(t, generated.DatabaseTypePostgreSQL, result)
}

func TestMustNewDatabaseType_Invalid(t *testing.T) {
	assert.Panics(t, func() {
		templates.MustNewDatabaseType("invalid")
	})
}

func TestIsValidProjectType(t *testing.T) {
	assert.True(t, templates.IsValidProjectType("microservice"))
	assert.True(t, templates.IsValidProjectType("hobby"))
	assert.False(t, templates.IsValidProjectType("invalid"))
	assert.False(t, templates.IsValidProjectType(""))
}

func TestIsValidDatabaseType(t *testing.T) {
	assert.True(t, templates.IsValidDatabaseType("postgresql"))
	assert.True(t, templates.IsValidDatabaseType("mysql"))
	assert.False(t, templates.IsValidDatabaseType("invalid"))
	assert.False(t, templates.IsValidDatabaseType(""))
}

func TestProjectTypeConstants(t *testing.T) {
	assert.Equal(t, "hobby", string(templates.ProjectTypeHobby))
	assert.Equal(t, "microservice", string(templates.ProjectTypeMicroservice))
	assert.Equal(t, "enterprise", string(templates.ProjectTypeEnterprise))
}

func TestDatabaseTypeConstants(t *testing.T) {
	assert.Equal(t, "postgresql", string(templates.DatabaseTypePostgreSQL))
	assert.Equal(t, "mysql", string(templates.DatabaseTypeMySQL))
	assert.Equal(t, "sqlite", string(templates.DatabaseTypeSQLite))
}

func TestMicroserviceTemplate_Name(t *testing.T) {
	template := &templates.MicroserviceTemplate{}
	assert.Equal(t, "microservice", template.Name())
}

func TestMicroserviceTemplate_Description(t *testing.T) {
	template := &templates.MicroserviceTemplate{}
	assert.Contains(t, template.Description(), "microservice")
}

func TestMicroserviceTemplate_DefaultData(t *testing.T) {
	template := &templates.MicroserviceTemplate{}
	data := template.DefaultData()

	// Verify defaults
	assert.Equal(t, generated.ProjectTypeMicroservice, data.ProjectType)
	assert.Equal(t, "db", data.Package.Name)
	assert.Equal(t, "internal/db", data.Package.Path)
	assert.Equal(t, generated.DatabaseTypePostgreSQL, data.Database.Engine)
	assert.True(t, data.Database.UseUUIDs)
	assert.True(t, data.Database.UseJSON)
	assert.False(t, data.Database.UseArrays)
}

func TestMicroserviceTemplate_Generate_Basic(t *testing.T) {
	template := &templates.MicroserviceTemplate{}
	data := template.DefaultData()
	data.ProjectName = "test-service"

	result, err := template.Generate(data)

	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, "2", result.Version)
	assert.Len(t, result.SQL, 1)

	sqlConfig := result.SQL[0]
	assert.Equal(t, "test-service", sqlConfig.Name)
	assert.Equal(t, "postgresql", sqlConfig.Engine)
	assert.NotNil(t, sqlConfig.Database)
}

// HobbyTemplate Tests
func TestHobbyTemplate_Name(t *testing.T) {
	template := &templates.HobbyTemplate{}
	assert.Equal(t, "hobby", template.Name())
}

func TestHobbyTemplate_Description(t *testing.T) {
	template := &templates.HobbyTemplate{}
	assert.Contains(t, template.Description(), "hobby")
}

func TestHobbyTemplate_DefaultData(t *testing.T) {
	template := &templates.HobbyTemplate{}
	data := template.DefaultData()

	// Verify defaults
	assert.Equal(t, generated.ProjectTypeHobby, data.ProjectType)
	assert.Equal(t, "db", data.Package.Name)
	assert.Equal(t, "db", data.Package.Path)
	assert.Equal(t, generated.DatabaseTypeSQLite, data.Database.Engine)
	assert.False(t, data.Database.UseUUIDs)
	assert.False(t, data.Database.UseJSON)
	assert.False(t, data.Database.UseArrays)
}

func TestHobbyTemplate_Generate_Basic(t *testing.T) {
	template := &templates.HobbyTemplate{}
	data := template.DefaultData()
	data.ProjectName = "my-hobby-project"

	result, err := template.Generate(data)

	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, "2", result.Version)
	assert.Len(t, result.SQL, 1)

	sqlConfig := result.SQL[0]
	assert.Equal(t, "my-hobby-project", sqlConfig.Name)
	assert.Equal(t, "sqlite", sqlConfig.Engine)
	assert.NotNil(t, sqlConfig.Database)
}

// EnterpriseTemplate Tests
func TestEnterpriseTemplate_Name(t *testing.T) {
	template := &templates.EnterpriseTemplate{}
	assert.Equal(t, "enterprise", template.Name())
}

func TestEnterpriseTemplate_Description(t *testing.T) {
	template := &templates.EnterpriseTemplate{}
	assert.Contains(t, template.Description(), "enterprise")
}

func TestEnterpriseTemplate_DefaultData(t *testing.T) {
	template := &templates.EnterpriseTemplate{}
	data := template.DefaultData()

	// Verify defaults
	assert.Equal(t, generated.ProjectTypeEnterprise, data.ProjectType)
	assert.Equal(t, "db", data.Package.Name)
	assert.Equal(t, "internal/db", data.Package.Path)
	assert.Equal(t, generated.DatabaseTypePostgreSQL, data.Database.Engine)
	assert.True(t, data.Database.UseUUIDs)
	assert.True(t, data.Database.UseJSON)
	assert.True(t, data.Database.UseArrays)
	assert.True(t, data.Database.UseFullText)
}

func TestEnterpriseTemplate_Generate_Basic(t *testing.T) {
	template := &templates.EnterpriseTemplate{}
	data := template.DefaultData()
	data.ProjectName = "enterprise-service"

	result, err := template.Generate(data)

	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, "2", result.Version)
	assert.Len(t, result.SQL, 1)

	sqlConfig := result.SQL[0]
	assert.Equal(t, "enterprise-service", sqlConfig.Name)
	assert.Equal(t, "postgresql", sqlConfig.Engine)
	assert.NotNil(t, sqlConfig.Database)
	assert.True(t, *sqlConfig.StrictFunctionChecks)
	assert.True(t, *sqlConfig.StrictOrderBy)
}

// APIFirstTemplate Tests
func TestAPIFirstTemplate_Name(t *testing.T) {
	template := &templates.APIFirstTemplate{}
	assert.Equal(t, "api-first", template.Name())
}

func TestAPIFirstTemplate_Description(t *testing.T) {
	template := &templates.APIFirstTemplate{}
	assert.Contains(t, template.Description(), "API")
}

func TestAPIFirstTemplate_DefaultData(t *testing.T) {
	internal_testing.AssertTemplateDefaultData(t, internal_testing.TemplateTestHelper{
		Template:            &templates.APIFirstTemplate{},
		ExpectedProjectType: generated.ProjectType("api-first"),
		ExpectedProjectName: "api-service",
		ExpectedEngine:      "postgresql",
		ExpectUUID:          true,
		ExpectJSON:          true,
		ExpectArrays:        true,
		ExpectJSONTags:      true,
		ExpectInterface:     true,
	})
}

func TestAPIFirstTemplate_Generate_Basic(t *testing.T) {
	internal_testing.AssertTemplateGenerateBasic(t, internal_testing.TemplateTestHelper{
		Template:            &templates.APIFirstTemplate{},
		ExpectedProjectType: generated.ProjectType("api-first"),
		ExpectedProjectName: "api-service",
		ExpectedEngine:      "postgresql",
		ExpectUUID:          true,
		ExpectJSON:          true,
		ExpectArrays:        true,
		ExpectJSONTags:      true,
		ExpectInterface:     true,
	})

// AnalyticsTemplate Tests
func TestAnalyticsTemplate_Name(t *testing.T) {
	template := &templates.AnalyticsTemplate{}
	assert.Equal(t, "analytics", template.Name())
}

func TestAnalyticsTemplate_Description(t *testing.T) {
	template := &templates.AnalyticsTemplate{}
	assert.Contains(t, template.Description(), "analytics")
}

func TestAnalyticsTemplate_DefaultData(t *testing.T) {
	template := &templates.AnalyticsTemplate{}
	data := template.DefaultData()

	assert.Equal(t, generated.ProjectType("analytics"), data.ProjectType)
	assert.Equal(t, "analytics", data.Package.Name)
	assert.Equal(t, "internal/analytics", data.Package.Path)
	assert.Equal(t, generated.DatabaseTypePostgreSQL, data.Database.Engine)
	assert.False(t, data.Database.UseUUIDs)
	assert.True(t, data.Database.UseJSON)
	assert.True(t, data.Database.UseArrays)
	assert.True(t, data.Database.UseFullText)
	assert.Equal(t, "snake", data.Validation.EmitOptions.JSONTagsCaseStyle)
}

func TestAnalyticsTemplate_Generate_Basic(t *testing.T) {
	template := &templates.AnalyticsTemplate{}
	data := template.DefaultData()
	data.ProjectName = "analytics-service"

	result, err := template.Generate(data)

	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, "2", result.Version)
	assert.Len(t, result.SQL, 1)

	sqlConfig := result.SQL[0]
	assert.Equal(t, "analytics-service", sqlConfig.Name)
	assert.Equal(t, "postgresql", sqlConfig.Engine)
	assert.NotNil(t, sqlConfig.Database)
	assert.True(t, *sqlConfig.StrictFunctionChecks)
	assert.True(t, *sqlConfig.StrictOrderBy)
}

// TestingTemplate Tests
func TestTestingTemplate_Name(t *testing.T) {
	template := &templates.TestingTemplate{}
	assert.Equal(t, "testing", template.Name())
}

func TestTestingTemplate_Description(t *testing.T) {
	template := &templates.TestingTemplate{}
	assert.Contains(t, template.Description(), "test")
}

func TestTestingTemplate_DefaultData(t *testing.T) {
	template := &templates.TestingTemplate{}
	data := template.DefaultData()

	assert.Equal(t, generated.ProjectType("testing"), data.ProjectType)
	assert.Equal(t, "testdata", data.Package.Name)
	assert.Equal(t, "testdata/db", data.Package.Path)
	assert.Equal(t, generated.DatabaseTypeSQLite, data.Database.Engine)
	assert.False(t, data.Database.UseUUIDs)
	assert.False(t, data.Database.UseJSON)
	assert.False(t, data.Database.UseArrays)
	assert.False(t, data.Validation.EmitOptions.EmitJSONTags)
	assert.False(t, data.Validation.EmitOptions.EmitPreparedQueries)
}

func TestTestingTemplate_Generate_Basic(t *testing.T) {
	template := &templates.TestingTemplate{}
	data := template.DefaultData()
	data.ProjectName = "test-project"

	result, err := template.Generate(data)

	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, "2", result.Version)
	assert.Len(t, result.SQL, 1)

	sqlConfig := result.SQL[0]
	assert.Equal(t, "test-project", sqlConfig.Name)
	assert.Equal(t, "sqlite", sqlConfig.Engine)
	assert.NotNil(t, sqlConfig.Database)
	assert.False(t, *sqlConfig.StrictFunctionChecks)
	assert.False(t, *sqlConfig.StrictOrderBy)
}

// MultiTenantTemplate Tests
func TestMultiTenantTemplate_Name(t *testing.T) {
	template := &templates.MultiTenantTemplate{}
	assert.Equal(t, "multi-tenant", template.Name())
}

func TestMultiTenantTemplate_Description(t *testing.T) {
	template := &templates.MultiTenantTemplate{}
	assert.Contains(t, template.Description(), "multi-tenant")
}

func TestMultiTenantTemplate_DefaultData(t *testing.T) {
	template := &templates.MultiTenantTemplate{}
	data := template.DefaultData()

	assert.Equal(t, generated.ProjectType("multi-tenant"), data.ProjectType)
	assert.Equal(t, "db", data.Package.Name)
	assert.Equal(t, "internal/db", data.Package.Path)
	assert.Equal(t, generated.DatabaseTypePostgreSQL, data.Database.Engine)
	assert.True(t, data.Database.UseUUIDs)
	assert.True(t, data.Database.UseJSON)
	assert.True(t, data.Database.UseArrays)
	assert.True(t, data.Validation.EmitOptions.EmitJSONTags)
	assert.Equal(t, "camel", data.Validation.EmitOptions.JSONTagsCaseStyle)
}

func TestMultiTenantTemplate_Generate_Basic(t *testing.T) {
	template := &templates.MultiTenantTemplate{}
	data := template.DefaultData()
	data.ProjectName = "multi-tenant-app"

	result, err := template.Generate(data)

	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, "2", result.Version)
	assert.Len(t, result.SQL, 1)

	sqlConfig := result.SQL[0]
	assert.Equal(t, "multi-tenant-app", sqlConfig.Name)
	assert.Equal(t, "postgresql", sqlConfig.Engine)
	assert.NotNil(t, sqlConfig.Database)
	assert.True(t, *sqlConfig.StrictFunctionChecks)
	assert.True(t, *sqlConfig.StrictOrderBy)
}

// LibraryTemplate Tests
func TestLibraryTemplate_Name(t *testing.T) {
	template := &templates.LibraryTemplate{}
	assert.Equal(t, "library", template.Name())
}

func TestLibraryTemplate_Description(t *testing.T) {
	template := &templates.LibraryTemplate{}
	assert.Contains(t, template.Description(), "Library")
}

func TestLibraryTemplate_DefaultData(t *testing.T) {
	internal_testing.AssertTemplateDefaultData(t, internal_testing.TemplateTestHelper{
		Template:            &templates.LibraryTemplate{},
		ExpectedProjectType: generated.ProjectType("library"),
		ExpectUUID:          false,
		ExpectJSON:          false,
		ExpectArrays:        false,
		ExpectJSONTags:      true,
		ExpectInterface:     true,
	})
}

func TestLibraryTemplate_Generate_Basic(t *testing.T) {
	internal_testing.AssertTemplateGenerateBasic(t, internal_testing.TemplateTestHelper{
		Template:            &templates.LibraryTemplate{},
		ExpectedProjectType: generated.ProjectType("library"),
		ExpectedProjectName: "library-module",
		ExpectedEngine:      "postgresql",
		ExpectUUID:          false,
		ExpectJSON:          false,
		ExpectArrays:        false,
		ExpectJSONTags:      true,
		ExpectInterface:     true,
	})
}
