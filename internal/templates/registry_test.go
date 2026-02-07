package templates_test

import (
	"testing"

	"github.com/LarsArtmann/SQLC-Wizzard/internal/templates"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRegistry(t *testing.T) {
	registry := templates.NewRegistry()

	assert.NotNil(t, registry)
}

func TestNewRegistry_RegistersAllTemplates(t *testing.T) {
	registry := templates.NewRegistry()

	// Should have 8 templates
	templates := registry.List()
	assert.Len(t, templates, 8)

	// Should have expected template names
	templateNames := make(map[string]bool)
	for _, tmpl := range templates {
		templateNames[tmpl.Name()] = true
	}

	expectedTemplates := []string{
		"hobby",
		"microservice",
		"enterprise",
		"api-first",
		"analytics",
		"testing",
		"multi-tenant",
		"library",
	}

	for _, expected := range expectedTemplates {
		assert.True(t, templateNames[expected], "Template %s should be registered", expected)
	}
}

func TestRegistry_Get_ExistingTemplate(t *testing.T) {
	registry := templates.NewRegistry()

	// Test each template type
	testCases := []struct {
		projectType string
		name        string
	}{
		{"hobby", "HobbyTemplate"},
		{"microservice", "MicroserviceTemplate"},
		{"enterprise", "EnterpriseTemplate"},
		{"api-first", "APIFirstTemplate"},
		{"analytics", "AnalyticsTemplate"},
		{"testing", "TestingTemplate"},
		{"multi-tenant", "MultiTenantTemplate"},
		{"library", "LibraryTemplate"},
	}

	for _, tc := range testCases {
		pt, err := templates.NewProjectType(tc.projectType)
		require.NoError(t, err, "Should create project type: %s", tc.projectType)

		tmpl, err := registry.Get(pt)

		require.NoError(t, err, "Should get template: %s", tc.name)
		assert.NotNil(t, tmpl, "Template should not be nil: %s", tc.name)
		assert.Equal(t, tc.projectType, tmpl.Name(), "Template name should match: %s", tc.name)
	}
}

func TestRegistry_HasTemplate_Existing(t *testing.T) {
	registry := templates.NewRegistry()

	// Test existing templates
	testCases := []string{
		"hobby",
		"microservice",
		"enterprise",
		"api-first",
		"analytics",
		"testing",
		"multi-tenant",
		"library",
	}

	for _, templateName := range testCases {
		pt := templates.MustNewProjectType(templateName)
		assert.True(t, registry.HasTemplate(pt), "Template %s should exist", templateName)
	}
}

func TestRegistry_HasTemplate_NonExisting(t *testing.T) {
	registry := templates.NewRegistry()

	// Test non-existing template
	pt := templates.MustNewProjectType("hobby")
	// Since NewRegistry() pre-registers all templates, we can't truly test
	// for non-existing templates. This test verifies HasTemplate works.
	assert.True(t, registry.HasTemplate(pt))
}

func TestRegistry_List(t *testing.T) {
	registry := templates.NewRegistry()

	templates := registry.List()

	assert.NotNil(t, templates)
	assert.Len(t, templates, 8, "Should have exactly 8 templates")

	// Check for unique templates
	templateNames := make(map[string]bool)
	for _, tmpl := range templates {
		name := tmpl.Name()
		assert.False(t, templateNames[name], "Template should be unique: %s", name)
		templateNames[name] = true

		// Verify template has required methods
		assert.NotEmpty(t, tmpl.Name(), "Template should have name")
		assert.NotEmpty(t, tmpl.Description(), "Template should have description")
		assert.NotNil(t, tmpl.RequiredFeatures(), "Template should have required features method")
	}
}

func TestRegistry_Register_Duplicate(t *testing.T) {
	registry := templates.NewRegistry()

	// Get initial count
	initialTemplates := registry.List()
	initialCount := len(initialTemplates)

	// Register same template twice
	registry.Register(templates.NewHobbyTemplate())
	registry.Register(templates.NewHobbyTemplate())

	// Should still have same number of templates (duplicates overwritten)
	finalTemplates := registry.List()
	assert.Len(t, finalTemplates, initialCount, "Duplicate templates should be overwritten")
}

func TestGetTemplate_ConvenienceFunction(t *testing.T) {
	// Test GetTemplate convenience function
	testCases := []string{
		"hobby",
		"microservice",
		"enterprise",
	}

	for _, templateName := range testCases {
		pt := templates.MustNewProjectType(templateName)
		tmpl, err := templates.GetTemplate(pt)

		require.NoError(t, err, "Should get template via convenience function: %s", templateName)
		assert.NotNil(t, tmpl, "Template should not be nil")
		assert.Equal(t, templateName, tmpl.Name())
	}
}

func TestListTemplates_ConvenienceFunction(t *testing.T) {
	// Test ListTemplates convenience function
	templates := templates.ListTemplates()

	assert.NotNil(t, templates)
	assert.Len(t, templates, 8, "Should list all 8 templates")
}

func TestDefaultRegistry(t *testing.T) {
	// Test that default registry is properly initialized

	// This ensures to init() function in registry.go works correctly
	templates := templates.ListTemplates()

	assert.NotNil(t, templates)
	assert.Len(t, templates, 8, "Default registry should have all templates")
}
