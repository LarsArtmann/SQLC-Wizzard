package templates

import (
	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/pkg/config"
)

// LibraryTemplate generates sqlc config for reusable Go libraries.
type LibraryTemplate struct {
	ConfiguredTemplate
}

// NewLibraryTemplate creates a new library template.
func NewLibraryTemplate() *LibraryTemplate {
	t := &LibraryTemplate{}
	t.ConfiguredTemplate = ConfiguredTemplate{
		TemplateName:        "library",
		TemplateDescription: "Library package configuration for reusable Go library development",
		DefaultPackageName:  "library",
		DefaultProjectName:  "library",
		StrictMode:          false,
		ProjectType:          "library",
		DbEngine:            "postgresql",
		PackagePath:         "internal/db",
		BaseOutput:          "internal/db",
		UseManaged:          false,
		UseUUIDs:            false,
		UseJSON:             false,
		UseArrays:           false,
		UseFullText:         false,
		EmitJSONTags:        true,
		EmitPreparedQueries: false,
		EmitInterface:       true,
		EmitEmptySlices:     false,
		EmitResultStructPointers: false,
		EmitParamsStructPointers: false,
		EmitEnumValidMethod:      true,
		EmitAllEnumValues:        false,
		JSONTagsCaseStyle:        "camel",
		StrictFunctions:         false,
		StrictOrderBy:           false,
		NoSelectStar:            false,
		RequireWhere:             false,
		NoDropTable:             false,
		NoTruncate:              false,
		RequireLimit:            false,
		Features:                []string{"emit_interface", "json_tags", "enum_valid_method"},
		CustomRenameRules: map[string]string{
			"id":   "ID",
			"json": "JSON",
			"api":  "API",
		},
	}
	return t
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
		Data:               data,
		DefaultName:        "library",
		DefaultDatabaseURL: "${DATABASE_URL}",
		Strict:             false,
	}
	cfg, _ := builder.Build()

	// Customize Go config with template-specific settings (including custom rename rules)
	cfg.SQL[0].Gen.Go = t.BuildGoConfigWithOverrides(data)

	// Apply validation rules
	return t.ApplyValidationRules(cfg, data)
}

// DefaultData returns default TemplateData for library template.
func (t *LibraryTemplate) DefaultData() generated.TemplateData {
	return t.BuildDefaultData(
		"library",
		"postgresql",
		"${DATABASE_URL}",
		"internal/db",
		"internal/db",
		false, // useManaged
		false, // useUUIDs
		false, // useJSON
		false, // useArrays
		false, // useFullText
		true,  // emitJSONTags
		false, // emitPreparedQueries
		true,  // emitInterface
		false, // emitEmptySlices
		false, // emitResultStructPointers
		false, // emitParamsStructPointers
		true,  // emitEnumValidMethod
		false, // emitAllEnumValues
		"camel", // jsonTagsCaseStyle
		false, // strictFunctions
		false, // strictOrderBy
		false, // noSelectStar
		false, // requireWhere
		false, // noDropTable
		false, // noTruncate
		false, // requireLimit
	)
}

// RequiredFeatures returns which features this template requires.
func (t *LibraryTemplate) RequiredFeatures() []string {
	return []string{"emit_interface", "json_tags", "enum_valid_method"}
}