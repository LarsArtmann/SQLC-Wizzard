package templates

import (
	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/pkg/config"
)

// LibraryTemplate generates sqlc config for reusable Go libraries.
type LibraryTemplate struct {
	BaseTemplate
}

// NewLibraryTemplate creates a new library template.
func NewLibraryTemplate() *LibraryTemplate {
	return &LibraryTemplate{}
}

// Name returns the template name.
func (t *LibraryTemplate) Name() string {
	return "library"
}

// Description returns a human-readable description.
func (t *LibraryTemplate) Description() string {
	return "Library package configuration for reusable Go library development"
}

// Generate creates a SqlcConfig from template data.
func (t *LibraryTemplate) Generate(data generated.TemplateData) (*config.SqlcConfig, error) {
	// Apply default values
	t.ApplyDefaultValues(&data)

	// Build config using shared builder
	builder := &ConfigBuilder{
		Data:             data,
		DefaultName:      "library",
		DefaultDatabaseURL: "${DATABASE_URL}",
		Strict:           false,
	}
	cfg, _ := builder.Build()

	// Customize Go config with template-specific settings
	cfg.SQL[0].Gen.Go = t.buildGoGenConfig(data)

	// Apply validation rules
	return t.ApplyValidationRules(cfg, data)
}

// DefaultData returns default TemplateData for library template.
func (t *LibraryTemplate) DefaultData() TemplateData {
	return generated.TemplateData{
		ProjectName: "",
		ProjectType: MustNewProjectType("library"),

		Package: generated.PackageConfig{
			Name: "db",
			Path: "internal/db",
		},

		Database: generated.DatabaseConfig{
			Engine:      MustNewDatabaseType("postgresql"),
			URL:         "${DATABASE_URL}",
			UseManaged:  false,
			UseUUIDs:    false,
			UseJSON:     false,
			UseArrays:   false,
			UseFullText: false,
		},

		Output: generated.OutputConfig{
			BaseDir:    "internal/db",
			QueriesDir: "internal/db/queries",
			SchemaDir:  "internal/db/schema",
		},

		Validation: generated.ValidationConfig{
			StrictFunctions: false,
			StrictOrderBy:   false,
			EmitOptions: generated.EmitOptions{
				EmitJSONTags:             true,
				EmitPreparedQueries:      false,
				EmitInterface:            true,
				EmitEmptySlices:          true,
				EmitResultStructPointers: false,
				EmitParamsStructPointers: false,
				EmitEnumValidMethod:      true,
				EmitAllEnumValues:        true,
				JSONTagsCaseStyle:        "camel",
			},
			SafetyRules: generated.SafetyRules{
				NoSelectStar: false,
				RequireWhere: false,
				NoDropTable:  true,
				NoTruncate:   true,
				RequireLimit: false,
				Rules:        []generated.SafetyRule{},
			},
		},
	}
}

// RequiredFeatures returns which features this template requires.
func (t *LibraryTemplate) RequiredFeatures() []string {
	return []string{"emit_interface", "json_tags", "enum_valid_method"}
}

// buildGoGenConfig builds the GoGenConfig from template data.
func (t *LibraryTemplate) buildGoGenConfig(data generated.TemplateData) *config.GoGenConfig {
	sqlPackage := t.GetSQLPackage(data.Database.Engine)
	cfg := t.BuildGoGenConfig(data, sqlPackage)

	// Libraries use simpler rename rules
	cfg.Rename = map[string]string{
		"id":   "ID",
		"json": "JSON",
		"api":  "API",
	}

	return cfg
}
