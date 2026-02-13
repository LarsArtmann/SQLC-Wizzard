package templates

import (
	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/pkg/config"
)

// HobbyTemplate generates sqlc config for hobby/personal projects.
type HobbyTemplate struct {
	ConfiguredTemplate
}

// NewHobbyTemplate creates a new hobby template.
func NewHobbyTemplate() *HobbyTemplate {
	base := NewConfiguredTemplate(
		"hobby",
		"Lightweight hobby configuration for personal projects and learning",
		"db",
		"hobby",
		false, // strictMode
		"hobby",
		"sqlite",
	)

	// Override hobby-specific settings: SQLite, minimal features
	base.UseManaged = false
	base.UseUUIDs = false
	base.UseJSON = false
	base.UseArrays = false
	base.UseFullText = false
	base.EmitJSONTags = false
	base.EmitInterface = false
	base.EmitEmptySlices = true
	base.JSONTagsCaseStyle = "camel"
	base.StrictFunctions = false
	base.StrictOrderBy = false
	base.NoSelectStar = false
	base.RequireWhere = false
	base.Features = []string{} // Hobby has no required features

	return &HobbyTemplate{ConfiguredTemplate: base}
}

// Name returns the template name.
func (t *HobbyTemplate) Name() string {
	return "hobby"
}

// Description returns the template description.
func (t *HobbyTemplate) Description() string {
	return "Lightweight hobby configuration for personal projects and learning"
}

// DefaultData returns default TemplateData for hobby template.
// Uses explicit BuildDefaultData for "file:dev.db" SQLite URL.
func (t *HobbyTemplate) DefaultData() generated.TemplateData {
	return t.BuildDefaultData(
		"hobby",       // projectType
		"sqlite",      // dbEngine
		"file:dev.db", // databaseURL
		"db",          // packagePath
		"db",          // baseOutputDir
		false,         // useManaged
		false,         // useUUIDs
		false,         // useJSON
		false,         // useArrays
		false,         // useFullText
		false,         // emitJSONTags
		false,         // emitPreparedQueries
		false,         // emitInterface
		true,          // emitEmptySlices
		false,         // emitResultStructPointers
		false,         // emitParamsStructPointers
		false,         // emitEnumValidMethod
		false,         // emitAllEnumValues
		"camel",       // jsonTagsCaseStyle
		false,         // strictFunctions
		false,         // strictOrderBy
		false,         // noSelectStar
		false,         // requireWhere
		false,         // noDropTable
		false,         // noTruncate
		false,         // requireLimit
	)
}

// Generate creates a SqlcConfig from template data.
func (t *HobbyTemplate) Generate(data generated.TemplateData) (*config.SqlcConfig, error) {
	return t.ConfiguredTemplate.Generate(data)
}

// RequiredFeatures returns which features this template requires.
func (t *HobbyTemplate) RequiredFeatures() []string {
	return t.Features
}
