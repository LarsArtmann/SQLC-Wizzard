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
		Template:              &templates.EnterpriseTemplate{},
		ExpectedProjectType:   generated.ProjectTypeEnterprise,
		ExpectUUID:            true,
		ExpectJSON:            true,
		ExpectArrays:          true,
		ExpectFullText:        true,
		ExpectJSONTags:       true,
		ExpectInterface:       true,
		ExpectPreparedQueries: true,
	})
}

func TestEnterpriseTemplate_Generate_Basic(t *testing.T) {
	internal_testing.AssertTemplateGenerateBasicWithDefaults(t,
		&templates.EnterpriseTemplate{},
		generated.ProjectTypeEnterprise,
		"enterprise-service",
	)
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
		ExpectStrictChecks:  false,
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
		Template:                  &templates.AnalyticsTemplate{},
		ExpectedProjectType:       generated.ProjectType("analytics"),
		ExpectedPackageName:       "analytics",
		ExpectedPackagePath:       "internal/analytics",
		ExpectUUID:               false,
		ExpectJSON:               true,
		ExpectArrays:             true,
		ExpectFullText:           true,
		ExpectJSONTags:           true,
		ExpectInterface:          true,
		ExpectStrictChecks:       true,
		ExpectedJSONTagsCaseStyle: "camel",
		ExpectPreparedQueries:    false,
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
	internal_testing.AssertTemplateDefaultData(t, internal_testing.NewTemplateTestHelper(
		&templates.MultiTenantTemplate{},
		internal_testing.WithProjectType(generated.ProjectType("multi-tenant")),
		internal_testing.WithEngine("postgresql"),
		internal_testing.WithUUID(true),
		internal_testing.WithJSON(true),
		internal_testing.WithArrays(true),
		internal_testing.WithJSONTags(true),
		internal_testing.WithInterface(true),
		internal_testing.WithStrictChecks(true),
		internal_testing.WithPreparedQueries(true),
	))
}

func TestMultiTenantTemplate_Generate_Basic(t *testing.T) {
	internal_testing.AssertTemplateGenerateBasicWithDefaults(t,
		&templates.MultiTenantTemplate{},
		generated.ProjectTypeMultiTenant,
		"multi-tenant-app",
	)
}

// LibraryTemplate Tests.
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
	internal_testing.AssertTemplateGenerateBasic(t, internal_testing.TemplateTestHelper{
		Template:            templates.NewLibraryTemplate(),
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

// ZeroValueInitializationTests verify that templates work correctly when created
// as zero-valued structs (empty struct literal). This is important because:
// - Users may create templates via reflection or deserialization
// - Tests use empty struct literals (&templates.EnterpriseTemplate{})
// - Ensures DefaultData() provides sensible defaults for all fields
func TestConfiguredTemplate_ZeroValueInitialization(t *testing.T) {
	// Test EnterpriseTemplate with zero values
	t.Run("EnterpriseTemplate zero initialization", func(t *testing.T) {
		tmpl := &templates.EnterpriseTemplate{} // zero values
		data := tmpl.DefaultData()

		assert.Equal(t, generated.ProjectTypeEnterprise, data.ProjectType)
		assert.Equal(t, generated.DatabaseTypePostgreSQL, data.Database.Engine)
		assert.True(t, data.Database.UseUUIDs, "UseUUIDs should default to true")
		assert.True(t, data.Database.UseJSON, "UseJSON should default to true")
		assert.True(t, data.Database.UseArrays, "UseArrays should default to true")
		assert.True(t, data.Database.UseFullText, "UseFullText should default to true for enterprise")
		assert.True(t, data.Validation.EmitOptions.EmitJSONTags, "EmitJSONTags should be true for enterprise")
		assert.True(t, data.Validation.EmitOptions.EmitInterface, "EmitInterface should be true for enterprise")
		assert.True(t, data.Validation.StrictFunctions, "StrictFunctions should be true for enterprise")
		assert.True(t, data.Validation.StrictOrderBy, "StrictOrderBy should be true for enterprise")
		assert.Equal(t, "camel", data.Validation.EmitOptions.JSONTagsCaseStyle, "JSONTagsCaseStyle should be camel for enterprise")
	})

	t.Run("APIFirstTemplate zero initialization", func(t *testing.T) {
		tmpl := &templates.APIFirstTemplate{} // zero values
		data := tmpl.DefaultData()

		assert.Equal(t, generated.ProjectTypeAPIFirst, data.ProjectType)
		assert.Equal(t, generated.DatabaseTypePostgreSQL, data.Database.Engine)
		assert.True(t, data.Database.UseUUIDs)
		assert.True(t, data.Database.UseJSON)
		assert.True(t, data.Database.UseArrays)
		assert.True(t, data.Validation.EmitOptions.EmitJSONTags)
		assert.True(t, data.Validation.EmitOptions.EmitInterface)
		assert.Equal(t, "camel", data.Validation.EmitOptions.JSONTagsCaseStyle)
	})

	t.Run("ConfiguredTemplate strict mode fallback", func(t *testing.T) {
		tmpl := &templates.EnterpriseTemplate{} // zero values
		data := tmpl.DefaultData()
		data.ProjectName = "test"

		// Generate should use strict settings from DefaultData, not from zero struct
		result, err := tmpl.Generate(data)
		require.NoError(t, err)
		require.NotNil(t, result)

		assert.True(t, *result.SQL[0].StrictFunctionChecks, "StrictFunctionChecks should come from DefaultData")
		assert.True(t, *result.SQL[0].StrictOrderBy, "StrictOrderBy should come from DefaultData")
	})
}

// TestBuildDefaultDataFromOptions verifies the new BuildOptions-based API works correctly.
func TestBuildDefaultDataFromOptions(t *testing.T) {
	t.Run("with defaults from NewBuildOptions", func(t *testing.T) {
		tmpl := &templates.MicroserviceTemplate{}
		opts := templates.NewBuildOptions("microservice", "postgresql")
		data := tmpl.BuildDefaultDataFromOptions(opts)

		assert.Equal(t, generated.ProjectTypeMicroservice, data.ProjectType)
		assert.Equal(t, generated.DatabaseTypePostgreSQL, data.Database.Engine)
		assert.True(t, data.Database.UseUUIDs, "UseUUIDs should default to true")
		assert.True(t, data.Database.UseJSON, "UseJSON should default to true")
		assert.True(t, data.Database.UseArrays, "UseArrays should default to true")
		assert.Equal(t, "internal/db", data.Package.Path)
	})

	t.Run("with custom options", func(t *testing.T) {
		tmpl := &templates.MicroserviceTemplate{}
		opts := templates.NewBuildOptions("enterprise", "postgresql")
		opts.UseArrays = false
		opts.EmitJSONTags = true
		opts.EmitInterface = true
		opts.JSONTagsCaseStyle = "camel"
		opts.StrictFunctions = true
		opts.StrictOrderBy = true
		data := tmpl.BuildDefaultDataFromOptions(opts)

		assert.Equal(t, generated.ProjectTypeEnterprise, data.ProjectType)
		assert.False(t, data.Database.UseArrays, "UseArrays should be false")
		assert.True(t, data.Validation.EmitOptions.EmitJSONTags, "EmitJSONTags should be true")
		assert.True(t, data.Validation.EmitOptions.EmitInterface, "EmitInterface should be true")
		assert.Equal(t, "camel", data.Validation.EmitOptions.JSONTagsCaseStyle)
		assert.True(t, data.Validation.StrictFunctions, "StrictFunctions should be true")
		assert.True(t, data.Validation.StrictOrderBy, "StrictOrderBy should be true")
	})

	t.Run("sqlite with hobby defaults", func(t *testing.T) {
		tmpl := &templates.MicroserviceTemplate{}
		opts := templates.NewBuildOptions("hobby", "sqlite")
		opts.UseManaged = false
		opts.UseUUIDs = false
		opts.UseJSON = false
		data := tmpl.BuildDefaultDataFromOptions(opts)

		assert.Equal(t, generated.ProjectTypeHobby, data.ProjectType)
		assert.Equal(t, generated.DatabaseTypeSQLite, data.Database.Engine)
		assert.False(t, data.Database.UseManaged, "UseManaged should be false for hobby")
		assert.False(t, data.Database.UseUUIDs, "UseUUIDs should be false for hobby")
		assert.False(t, data.Database.UseJSON, "UseJSON should be false for hobby")
	})
}
