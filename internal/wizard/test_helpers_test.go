package wizard_test

import (
	"github.com/LarsArtmann/SQLC-Wizzard/generated"
)

// createTemplateData creates a basic template data structure for testing.
func createTemplateData() generated.TemplateData {
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

// createTemplateDataWithFeatures creates template data with feature flags enabled.
func createTemplateDataWithFeatures(projectName string, projectType generated.ProjectType) generated.TemplateData {
	data := createTemplateData()
	data.ProjectName = projectName
	data.ProjectType = projectType
	data.Database.UseUUIDs = true
	data.Database.UseJSON = true
	data.Database.UseArrays = true
	return data
}

// createTemplateDataWithCustomOutput creates template data with custom output directories.
func createTemplateDataWithCustomOutput(baseDir, queriesDir, schemaDir string) generated.TemplateData {
	data := createTemplateData()
	data.Output.BaseDir = baseDir
	data.Output.QueriesDir = queriesDir
	data.Output.SchemaDir = schemaDir
	return data
}
