package templates

import (
	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/pkg/config"
)

// HobbyTemplate generates sqlc config for hobby/personal projects.
type HobbyTemplate struct {
	BaseTemplate
}

// NewHobbyTemplate creates a new hobby template.
func NewHobbyTemplate() *HobbyTemplate {
	return &HobbyTemplate{}
}

// Name returns the template name.
func (t *HobbyTemplate) Name() string {
	return "hobby"
}

// Description returns a human-readable description.
func (t *HobbyTemplate) Description() string {
	return "Lightweight hobby configuration for personal projects and learning"
}

// Generate creates a SqlcConfig from template data.
func (t *HobbyTemplate) Generate(data generated.TemplateData) (*config.SqlcConfig, error) {
	return t.BaseTemplate.GenerateWithDefaults(
		data,
		"db",                          // packageName
		"db",                          // packagePath
		"db",                          // baseDir
		"db/queries",                  // queriesDir
		"db/schema",                   // schemaDir
		"file:dev.db",                 // databaseURL
		"hobby",                       // projectName
		false,                         // strict
	)
}

// DefaultData returns default TemplateData for hobby template.
func (t *HobbyTemplate) DefaultData() TemplateData {
	return generated.TemplateData{
		ProjectName: "",
		ProjectType: MustNewProjectType("hobby"),

		Package: generated.PackageConfig{
			Name: "db",
			Path: "db",
		},

		Database: generated.DatabaseConfig{
			Engine:      MustNewDatabaseType("sqlite"),
			URL:         "file:dev.db",
			UseManaged:  false,
			UseUUIDs:    false,
			UseJSON:     false,
			UseArrays:   false,
			UseFullText: false,
		},

		Output: generated.OutputConfig{
			BaseDir:    "db",
			QueriesDir: "db/queries",
			SchemaDir:  "db/schema",
		},

		Validation: generated.ValidationConfig{
			StrictFunctions: false,
			StrictOrderBy:   false,
			EmitOptions: generated.EmitOptions{
				EmitJSONTags:             false,
				EmitPreparedQueries:      false,
				EmitInterface:            false,
				EmitEmptySlices:          true,
				EmitResultStructPointers: false,
				EmitParamsStructPointers: false,
				EmitEnumValidMethod:      false,
				EmitAllEnumValues:        false,
				JSONTagsCaseStyle:        "snake",
			},
			SafetyRules: generated.SafetyRules{
				NoSelectStar: false,
				RequireWhere: false,
				NoDropTable:  false,
				NoTruncate:   false,
				RequireLimit: false,
				Rules:        []generated.SafetyRule{},
			},
		},
	}
}

// RequiredFeatures returns which features this template requires.
func (t *HobbyTemplate) RequiredFeatures() []string {
	return []string{}
}

