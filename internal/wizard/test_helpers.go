package wizard

import (
	"github.com/LarsArtmann/SQLC-Wizzard/generated"
)

// CreateTemplateData creates a basic template data structure for testing.
func CreateTemplateData() generated.TemplateData {
	return generated.TemplateData{
		Package: generated.PackageConfig{
			Name: "db",
			Path: "github.com/example/test",
		},
		Database: generated.DatabaseConfig{
			Engine: generated.DatabaseTypePostgreSQL,
		},
		Output: generated.OutputConfig{
			BaseDir:    "./internal/db",
			QueriesDir: "./internal/db/queries",
			SchemaDir:  "./internal/db/schema",
		},
		Validation: generated.ValidationConfig{
			EmitOptions: generated.DefaultEmitOptions(),
			SafetyRules: generated.DefaultSafetyRules(),
		},
	}
}

// CreateTemplateDataWithFeatures creates template data with feature flags enabled.
func CreateTemplateDataWithFeatures(projectName string, projectType generated.ProjectType) generated.TemplateData {
	data := CreateTemplateData()
	data.ProjectName = projectName
	data.ProjectType = projectType
	data.Database.UseUUIDs = true
	data.Database.UseJSON = true
	data.Database.UseArrays = true
	return data
}

// CreateTemplateDataWithCustomOutput creates template data with custom output directories.
func CreateTemplateDataWithCustomOutput(baseDir, queriesDir, schemaDir string) generated.TemplateData {
	data := CreateTemplateData()
	data.Output.BaseDir = baseDir
	data.Output.QueriesDir = queriesDir
	data.Output.SchemaDir = schemaDir
	return data
}

// CreateTemplateDataWithAllFeatures creates template data with all features enabled or disabled.
func CreateTemplateDataWithAllFeatures(enabled bool) *generated.TemplateData {
	data := CreateTemplateData()
	data.Database.UseUUIDs = enabled
	data.Database.UseJSON = enabled
	data.Database.UseArrays = enabled
	data.Database.UseFullText = enabled
	data.Validation.StrictFunctions = enabled
	data.Validation.StrictOrderBy = enabled
	return &data
}

// createTemplateDataWithAllFeatures creates template data with all features enabled or disabled.
// This is an unexported wrapper for use by internal test files.
func createTemplateDataWithAllFeatures(enabled bool) *generated.TemplateData {
	return CreateTemplateDataWithAllFeatures(enabled)
}