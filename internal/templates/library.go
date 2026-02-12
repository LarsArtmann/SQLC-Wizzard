package templates

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

// RequiredFeatures returns which features this template requires.
func (t *LibraryTemplate) RequiredFeatures() []string {
	return []string{"emit_interface", "json_tags", "enum_valid_method"}
}