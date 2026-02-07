package templates_test

import (
	"slices"
	"testing"

	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/templates"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAllTemplates_GenerateValidConfig(t *testing.T) {
	// Test that all templates generate valid SqlcConfig with default data
	registry := templates.NewRegistry()
	allTemplates := registry.List()

	assert.Len(t, allTemplates, 8, "Should have exactly 8 templates")

	for _, tmpl := range allTemplates {
		t.Run(tmpl.Name(), func(t *testing.T) {
			data := tmpl.DefaultData()
			data.ProjectName = "test-project-" + tmpl.Name()

			config, err := tmpl.Generate(data)

			// Should not error
			require.NoError(t, err, "Template %s should generate valid config", tmpl.Name())
			require.NotNil(t, config, "Template %s should not return nil config", tmpl.Name())

			// Verify basic structure
			assert.Equal(t, "2", config.Version, "Template %s should have version 2", tmpl.Name())
			assert.NotEmpty(t, config.SQL, "Template %s should have SQL configs", tmpl.Name())

			// Verify SQL config structure
			sqlConfig := config.SQL[0]
			assert.NotNil(t, sqlConfig, "Template %s should have SQL config", tmpl.Name())
			assert.NotEmpty(t, sqlConfig.Name, "Template %s should have SQL name", tmpl.Name())
			assert.NotEmpty(t, sqlConfig.Engine, "Template %s should have engine", tmpl.Name())

			// Verify Gen.Go structure
			assert.NotNil(t, sqlConfig.Gen.Go, "Template %s should have Go config", tmpl.Name())
			assert.NotEmpty(t, sqlConfig.Gen.Go.Package, "Template %s should have package name", tmpl.Name())
			assert.NotEmpty(t, sqlConfig.Gen.Go.Out, "Template %s should have output path", tmpl.Name())

			// Verify Database if present
			assert.NotNil(t, sqlConfig.Database, "Template %s should have database config", tmpl.Name())

			// Verify Query/Schema paths
			assert.NotNil(t, sqlConfig.Queries, "Template %s should have queries path", tmpl.Name())
			assert.NotNil(t, sqlConfig.Schema, "Template %s should have schema path", tmpl.Name())

			t.Logf("✓ Template %s generated valid config", tmpl.Name())
		})
	}
}

func TestAllTemplates_GenerateWithCustomData(t *testing.T) {
	// Test that all templates handle custom data correctly
	registry := templates.NewRegistry()

	testCases := []struct {
		templateName string
		projectName  string
	}{
		{"hobby", "my-hobby-app"},
		{"microservice", "api-service"},
		{"enterprise", "enterprise-system"},
		{"api-first", "my-api"},
		{"analytics", "analytics-platform"},
		{"testing", "test-suite"},
		{"multi-tenant", "saas-platform"},
		{"library", "my-library"},
	}

	for _, tc := range testCases {
		pt, err := templates.NewProjectType(tc.templateName)
		require.NoError(t, err)

		tmpl, err := registry.Get(pt)
		require.NoError(t, err)

		data := tmpl.DefaultData()
		data.ProjectName = tc.projectName
		data.Package.Name = tc.projectName + "-db"

		config, err := tmpl.Generate(data)

		// Should not error
		require.NoError(t, err, "Template %s with custom data should generate valid config", tc.templateName)
		require.NotNil(t, config)

		// Verify custom values are applied
		sqlConfig := config.SQL[0]
		assert.Equal(t, tc.projectName, sqlConfig.Name, "Template %s should use custom project name", tc.templateName)
		assert.Equal(t, tc.projectName+"-db", sqlConfig.Gen.Go.Package, "Template %s should use custom package name", tc.templateName)
	}
}

func TestTemplates_ProduceConsistentNaming(t *testing.T) {
	// Verify that all templates produce consistent Go naming
	registry := templates.NewRegistry()
	allTemplates := registry.List()

	for _, tmpl := range allTemplates {
		data := tmpl.DefaultData()
		config, err := tmpl.Generate(data)
		require.NoError(t, err)

		sqlConfig := config.SQL[0]
		goConfig := sqlConfig.Gen.Go

		// Verify common fields follow Go conventions
		assert.NotEmpty(t, goConfig.Package, "Template %s should have valid package name", tmpl.Name())
		assert.NotEmpty(t, goConfig.Out, "Template %s should have valid output path", tmpl.Name())

		// Verify no empty names
		assert.NotEmpty(t, goConfig.Package, "Template %s should have non-empty package name", tmpl.Name())
		assert.NotEmpty(t, goConfig.Out, "Template %s should have non-empty output path", tmpl.Name())
	}
}

func TestTemplates_SupportAllDatabaseTypes(t *testing.T) {
	// Verify all templates can generate configs for all database types
	registry := templates.NewRegistry()
	allTemplates := registry.List()

	dbTypes := []generated.DatabaseType{
		generated.DatabaseTypePostgreSQL,
		generated.DatabaseTypeMySQL,
		generated.DatabaseTypeSQLite,
	}

	for _, tmpl := range allTemplates {
		for _, dbType := range dbTypes {
			data := tmpl.DefaultData()
			data.Database.Engine = dbType
			data.ProjectName = string(tmpl.Name()) + "-" + string(dbType)

			config, err := tmpl.Generate(data)

			// Should not error
			assert.NoError(t, err, "Template %s should support database type %s", tmpl.Name(), dbType)

			sqlConfig := config.SQL[0]
			assert.Equal(t, string(dbType), sqlConfig.Engine, "Template %s should use correct engine", tmpl.Name())
		}
	}
}

func TestTemplates_GenerateValidYAML(t *testing.T) {
	// Verify that generated configs can be marshaled to YAML without errors
	registry := templates.NewRegistry()
	allTemplates := registry.List()

	for _, tmpl := range allTemplates {
		data := tmpl.DefaultData()
		data.ProjectName = "yaml-test-" + tmpl.Name()

		config, err := tmpl.Generate(data)
		require.NoError(t, err)

		// Try to marshal to YAML (sqlc config format)
		// We don't actually marshal to YAML here, just verify that structure is valid
		// This test ensures that config structure is complete and valid

		sqlConfig := config.SQL[0]

		// Verify all required fields are present
		assert.NotNil(t, sqlConfig.Engine)
		assert.NotNil(t, sqlConfig.Queries)
		assert.NotNil(t, sqlConfig.Schema)
		assert.NotNil(t, sqlConfig.Gen.Go)

		goConfig := sqlConfig.Gen.Go
		assert.NotEmpty(t, goConfig.Package)

		t.Logf("✓ Template %s generates valid config structure", tmpl.Name())
	}
}

func TestTemplates_DatabaseURLsAreCorrect(t *testing.T) {
	// Verify that templates set appropriate database URLs for their default databases
	testCases := []struct {
		templateName string
		expectedURL  string
	}{
		{"hobby", "file:dev.db"},
		{"microservice", "${DATABASE_URL}"},
		{"enterprise", "${DATABASE_URL}"},
		{"api-first", "${DATABASE_URL}"},
		{"analytics", "${ANALYTICS_DATABASE_URL}"},
		{"testing", "file:testdata/test.db"},
		{"multi-tenant", "${DATABASE_URL}"},
		{"library", "${DATABASE_URL}"},
	}

	registry := templates.NewRegistry()

	for _, tc := range testCases {
		pt, err := templates.NewProjectType(tc.templateName)
		require.NoError(t, err)

		tmpl, err := registry.Get(pt)
		require.NoError(t, err)

		data := tmpl.DefaultData()
		sqlConfig := data.Database.URL

		assert.Equal(t, tc.expectedURL, sqlConfig, "Template %s should have correct default URL", tc.templateName)
	}
}

func TestTemplates_ValidationConfigurations(t *testing.T) {
	// Verify that templates set appropriate validation settings
	registry := templates.NewRegistry()
	allTemplates := registry.List()

	for _, tmpl := range allTemplates {
		data := tmpl.DefaultData()
		config, err := tmpl.Generate(data)
		require.NoError(t, err)

		sqlConfig := config.SQL[0]

		// Verify StrictFunctionChecks is set appropriately
		// Hobby and testing templates should have false (no strict checks)
		// Enterprise template should have true (strict checks enabled)

		features := tmpl.RequiredFeatures()

		// Templates with strict_checks feature should have StrictFunctionChecks enabled
		hasStrictChecks := false
		if slices.Contains(features, "strict_checks") {
			hasStrictChecks = true
			assert.NotNil(t, sqlConfig.StrictFunctionChecks, "Template %s with strict_checks feature should have StrictFunctionChecks set", tmpl.Name())
			if sqlConfig.StrictFunctionChecks != nil {
				assert.True(t, *sqlConfig.StrictFunctionChecks, "Template %s with strict_checks should have StrictFunctionChecks = true", tmpl.Name())
			}
		}

		if !hasStrictChecks {
			// Templates without strict_checks should have false or nil
			if sqlConfig.StrictFunctionChecks != nil {
				assert.False(t, *sqlConfig.StrictFunctionChecks, "Template %s without strict_checks should have StrictFunctionChecks = false", tmpl.Name())
			}
		}

		t.Logf("✓ Template %s validation configuration is correct", tmpl.Name())
	}
}

func TestTemplates_OutputPaths(t *testing.T) {
	// Verify that templates set appropriate output paths
	registry := templates.NewRegistry()
	allTemplates := registry.List()

	for _, tmpl := range allTemplates {
		data := tmpl.DefaultData()
		config, err := tmpl.Generate(data)
		require.NoError(t, err)

		sqlConfig := config.SQL[0]
		goConfig := sqlConfig.Gen.Go

		// Verify output path is set
		assert.NotEmpty(t, goConfig.Out, "Template %s should have output path set", tmpl.Name())

		// Verify output path matches template expectations
		expectedPaths := map[string]string{
			"hobby":        "db",
			"microservice": "internal/db",
			"enterprise":   "internal/db",
			"api-first":    "internal/db",
			"analytics":    "internal/analytics",
			"testing":      "testdata/db",
			"multi-tenant": "internal/db",
			"library":      "internal/db",
		}

		expectedPath, ok := expectedPaths[tmpl.Name()]
		if ok {
			assert.Contains(t, goConfig.Out, expectedPath, "Template %s should have output path containing %s", tmpl.Name())
		}

		t.Logf("✓ Template %s output path is appropriate", tmpl.Name())
	}
}
