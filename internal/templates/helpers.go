package templates

import (
	"github.com/LarsArtmann/SQLC-Wizzard/generated"
)

// createDefaultTemplateData creates default TemplateData with customizable parameters.
// This helper reduces duplication across template implementations.
func createDefaultTemplateData(
	projectType string,
	packageName string,
	packagePath string,
	databaseURL string,
	databaseEngine string,
	useManaged bool,
	useUUIDs bool,
	useJSON bool,
	useArrays bool,
	useFullText bool,
	baseDir string,
	strictFunctions bool,
	strictOrderBy bool,
	emitJSONTags bool,
	emitInterface bool,
	emitAllEnumValues bool,
	emitPreparedQueries bool,
	emitEmptySlices bool,
	emitResultStructPointers bool,
	emitParamsStructPointers bool,
	emitEnumValidMethod bool,
	jsonTagsCaseStyle string,
	noSelectStar bool,
	requireWhere bool,
	noDropTable bool,
	noTruncate bool,
	requireLimit bool,
) TemplateData {
	return generated.TemplateData{
		ProjectName: "",
		ProjectType: MustNewProjectType(projectType),

		Package: generated.PackageConfig{
			Name: packageName,
			Path: packagePath,
		},

		Database: generated.DatabaseConfig{
			Engine:      MustNewDatabaseType(databaseEngine),
			URL:         databaseURL,
			UseManaged:  useManaged,
			UseUUIDs:    useUUIDs,
			UseJSON:     useJSON,
			UseArrays:   useArrays,
			UseFullText: useFullText,
		},

		Output: generated.OutputConfig{
			BaseDir:    baseDir,
			QueriesDir: baseDir + "/queries",
			SchemaDir:  baseDir + "/schema",
		},

		Validation: (&BaseTemplate{}).BuildValidationConfig(
			strictFunctions,
			strictOrderBy,
			emitJSONTags,
			emitPreparedQueries,
			emitInterface,
			emitEmptySlices,
			emitResultStructPointers,
			emitParamsStructPointers,
			emitEnumValidMethod,
			emitAllEnumValues,
			jsonTagsCaseStyle,
			noSelectStar,
			requireWhere,
			noDropTable,
			noTruncate,
			requireLimit,
		),
	}
}

// NewMinimalConfiguredTemplate builds a ConfiguredTemplate configured for
// lightweight projects (hobby, testing) that share the same minimal-feature
// profile: no managed/UUIDs/JSON/arrays/fulltext, no JSON tags, no interface,
// with empty slices enabled and all strict/safety checks off.
// Centralizing this eliminates the duplicated override block that hobby and
// testing templates used to repeat verbatim.
func NewMinimalConfiguredTemplate(
	name, description string,
	packageName, projectName string,
	projectType, dbEngine string,
	features []string,
) ConfiguredTemplate {
	tpl := NewConfiguredTemplate(
		name,
		description,
		packageName,
		projectName,
		false, // strictMode - minimal templates are not strict
		projectType,
		dbEngine,
	)

	// Minimal overrides: disable features common to hobby/testing workloads
	tpl.UseManaged = false
	tpl.UseUUIDs = false
	tpl.UseJSON = false
	tpl.UseArrays = false
	tpl.UseFullText = false
	tpl.EmitJSONTags = false
	tpl.EmitInterface = false
	tpl.EmitEmptySlices = true
	tpl.JSONTagsCaseStyle = CamelCaseStyle
	tpl.StrictFunctions = false
	tpl.StrictOrderBy = false
	tpl.NoSelectStar = false
	tpl.RequireWhere = false
	tpl.Features = features

	return tpl
}

// ApplyJSONInterfaceOptions configures a ConfiguredTemplate for templates that
// produce JSON-tagged structs with the interface emitter enabled. Used by
// api-first and enterprise templates that share this minimal common profile.
func ApplyJSONInterfaceOptions(tpl *ConfiguredTemplate) {
	tpl.JSONTagsCaseStyle = CamelCaseStyle
	tpl.EmitJSONTags = true
	tpl.EmitInterface = true
}
