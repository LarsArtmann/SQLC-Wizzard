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
	internal_testing.AssertTemplateGenerateBasic(t, internal_testing.TemplateTestHelper{
		Template:            &templates.MicroserviceTemplate{},
		ExpectedProjectType: generated.ProjectTypeMicroservice,
		ExpectedProjectName: "test-service",
		ExpectedEngine:      "postgresql",
		ExpectUUID:          true,
		ExpectJSON:          true,
		ExpectArrays:        false,
		ExpectJSONTags:      true,
		ExpectInterface:     true,
		ExpectStrictChecks:  false,
	})
}

// HobbyTemplate Tests.
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
		Template:            &templates.HobbyTemplate{},
		ExpectedProjectType: generated.ProjectTypeHobby,
		ExpectedPackagePath: "db",
		ExpectedDatabaseType: generated.DatabaseTypeSQLite,
		ExpectUUID:          false,
		ExpectJSON:          false,
		ExpectArrays:        false,
		ExpectJSONTags:      false,
		ExpectInterface:     false,
	})
}

func TestHobbyTemplate_Generate_Basic(t *testing.T) {
	internal_testing.AssertTemplateGenerateBasic(t, internal_testing.TemplateTestHelper{
		Template:            &templates.HobbyTemplate{},
		ExpectedProjectType: generated.ProjectTypeHobby,
		ExpectedProjectName: "my-hobby-project",
		ExpectedEngine:      "sqlite",
		ExpectUUID:          false,
		ExpectJSON:          false,
		ExpectArrays:        false,
		ExpectJSONTags:      false,
		ExpectInterface:     false,
		ExpectStrictChecks:  false,
	})
}

// EnterpriseTemplate Tests.
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
		Template:            &templates.EnterpriseTemplate{},
		ExpectedProjectType: generated.ProjectTypeEnterprise,
		ExpectUUID:          true,
		ExpectJSON:          true,
		ExpectArrays:        true,
		ExpectFullText:      true,
		ExpectJSONTags:      true,
		ExpectInterface:     true,
	})
}

func TestEnterpriseTemplate_Generate_Basic(t *testing.T) {
	internal_testing.AssertTemplateGenerateBasic(t, internal_testing.TemplateTestHelper{
		Template:            &templates.EnterpriseTemplate{},
		ExpectedProjectType: generated.ProjectType("enterprise"),
		ExpectedProjectName: "enterprise-service",
		ExpectedEngine:      "postgresql",
		ExpectUUID:          true,
		ExpectJSON:          true,
		ExpectArrays:        true,
		ExpectJSONTags:      true,
		ExpectInterface:     true,
		ExpectStrictChecks:  true,
	})
}

// APIFirstTemplate Tests.
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
}

// AnalyticsTemplate Tests.
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
		Template:              &templates.AnalyticsTemplate{},
		ExpectedProjectType:   generated.ProjectType("analytics"),
		ExpectedPackageName:   "analytics",
		ExpectedPackagePath:   "internal/analytics",
		ExpectUUID:            false,
		ExpectJSON:            true,
		ExpectArrays:          true,
		ExpectFullText:        true,
		ExpectJSONTags:        true,
		ExpectInterface:       true,
		ExpectStrictChecks:    true,
		ExpectedJSONTagsCaseStyle: "snake",
	})
}

func TestAnalyticsTemplate_Generate_Basic(t *testing.T) {
	internal_testing.AssertTemplateGenerateBasic(t, internal_testing.TemplateTestHelper{
		Template:            &templates.AnalyticsTemplate{},
		ExpectedProjectType: generated.ProjectType("analytics"),
		ExpectedProjectName: "analytics-service",
		ExpectedEngine:      "postgresql",
		ExpectUUID:          false,
		ExpectJSON:          true,
		ExpectArrays:        true,
		ExpectJSONTags:      true,
		ExpectInterface:     true,
		ExpectStrictChecks:  true,
	})
}

// TestingTemplate Tests.
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
	internal_testing.AssertTemplateGenerateBasic(t, internal_testing.TemplateTestHelper{
		Template:            &templates.TestingTemplate{},
		ExpectedProjectType: generated.ProjectType("testing"),
		ExpectedProjectName: "test-project",
		ExpectedEngine:      "sqlite",
		ExpectUUID:          false,
		ExpectJSON:          false,
		ExpectArrays:        false,
		ExpectJSONTags:      false,
		ExpectInterface:     false,
		ExpectStrictChecks:  false,
	})
}

// MultiTenantTemplate Tests.
func TestMultiTenantTemplate_Name(t *testing.T) {
	template := &templates.MultiTenantTemplate{}
	assert.Equal(t, "multi-tenant", template.Name())
}

func TestMultiTenantTemplate_Description(t *testing.T) {
	template := &templates.MultiTenantTemplate{}
	assert.Contains(t, template.Description(), "multi-tenant")
}

func TestMultiTenantTemplate_DefaultData(t *testing.T) {
	internal_testing.AssertTemplateDefaultData(t, internal_testing.TemplateTestHelper{
		Template:            &templates.MultiTenantTemplate{},
		ExpectedProjectType: generated.ProjectType("multi-tenant"),
		ExpectUUID:          true,
		ExpectJSON:          true,
		ExpectArrays:        true,
		ExpectJSONTags:      true,
		ExpectInterface:     true,
		ExpectStrictChecks:  true,
	})
}

func TestMultiTenantTemplate_Generate_Basic(t *testing.T) {
	internal_testing.AssertTemplateGenerateBasic(t, internal_testing.TemplateTestHelper{
		Template:            &templates.MultiTenantTemplate{},
		ExpectedProjectType: generated.ProjectType("multi-tenant"),
		ExpectedProjectName: "multi-tenant-app",
		ExpectedEngine:      "postgresql",
		ExpectUUID:          true,
		ExpectJSON:          true,
		ExpectArrays:        true,
		ExpectJSONTags:      true,
		ExpectInterface:     true,
		ExpectStrictChecks:  true,
	})
}

// LibraryTemplate Tests.
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
