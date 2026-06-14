package templates_test

import (
	"testing"

	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/templates"
	internal_testing "github.com/LarsArtmann/SQLC-Wizzard/internal/testing"
	"github.com/stretchr/testify/assert"
)

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
	internal_testing.AssertTemplateGenerateBasicWithDefaults(
		t,
		&templates.MultiTenantTemplate{},
		generated.ProjectTypeMultiTenant,
		"multi-tenant-app",
	)
}
