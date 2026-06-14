package templates_test

import (
	"testing"

	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/templates"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestConfiguredTemplate_ZeroValueInitialization verifies that templates work correctly when created
// as zero-valued structs (empty struct literal). This is important because:
// - Users may create templates via reflection or deserialization
// - Tests use empty struct literals (&templates.EnterpriseTemplate{})
// - Ensures DefaultData() provides sensible defaults for all fields.
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
		assert.True(
			t,
			data.Database.UseFullText,
			"UseFullText should default to true for enterprise",
		)
		assert.True(
			t,
			data.Validation.EmitOptions.EmitJSONTags,
			"EmitJSONTags should be true for enterprise",
		)
		assert.True(
			t,
			data.Validation.EmitOptions.EmitInterface,
			"EmitInterface should be true for enterprise",
		)
		assert.True(
			t,
			data.Validation.StrictFunctions,
			"StrictFunctions should be true for enterprise",
		)
		assert.True(t, data.Validation.StrictOrderBy, "StrictOrderBy should be true for enterprise")
		assert.Equal(
			t,
			"camel",
			data.Validation.EmitOptions.JSONTagsCaseStyle,
			"JSONTagsCaseStyle should be camel for enterprise",
		)
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

		assert.True(
			t,
			*result.SQL[0].StrictFunctionChecks,
			"StrictFunctionChecks should come from DefaultData",
		)
		assert.True(t, *result.SQL[0].StrictOrderBy, "StrictOrderBy should come from DefaultData")
	})
}
