package templates_test

import (
	"testing"

	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/templates"
	"github.com/stretchr/testify/assert"
)

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
