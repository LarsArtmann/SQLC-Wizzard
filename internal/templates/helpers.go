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
	emitPreparedQueries bool,
	emitResultPointers bool,
	emitParamsPointers bool,
	emitEnumValidMethod bool,
	jsonTagsCaseStyle string,
	noSelectStar bool,
	requireWhere bool,
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

		Validation: generated.ValidationConfig{
			StrictFunctions: strictFunctions,
			StrictOrderBy:   strictOrderBy,
			EmitOptions: generated.EmitOptions{
				EmitJSONTags:             true,
				EmitPreparedQueries:      emitPreparedQueries,
				EmitInterface:            true,
				EmitEmptySlices:          true,
				EmitResultStructPointers: emitResultPointers,
				EmitParamsStructPointers: emitParamsPointers,
				EmitEnumValidMethod:      emitEnumValidMethod,
				EmitAllEnumValues:        true,
				JSONTagsCaseStyle:        jsonTagsCaseStyle,
			},
			SafetyRules: generated.SafetyRules{
				NoSelectStar: noSelectStar,
				RequireWhere: requireWhere,
				NoDropTable:  true,
				NoTruncate:   true,
				RequireLimit: requireLimit,
				Rules:        []generated.SafetyRule{},
			},
		},
	}
}