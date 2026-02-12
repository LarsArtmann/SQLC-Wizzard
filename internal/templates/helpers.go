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
