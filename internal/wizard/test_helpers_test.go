package wizard_test

import (
	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/wizard"
	. "github.com/onsi/gomega"
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

// createTemplateDataWithAllFeatures creates template data with all features enabled or disabled.
func createTemplateDataWithAllFeatures(enabled bool) generated.TemplateData {
	data := createTemplateData()
	data.Database.UseUUIDs = enabled
	data.Database.UseJSON = enabled
	data.Database.UseArrays = enabled
	data.Database.UseFullText = enabled
	data.Validation.StrictFunctions = enabled
	data.Validation.StrictOrderBy = enabled
	return data
}

// testOutputPathConfiguration tests output directory configuration with custom paths.
func testOutputPathConfiguration(wiz *wizard.Wizard, baseDir, queriesDir, schemaDir string) {
	result := wiz.GetResult()
	result.TemplateData.Output = createTemplateDataWithCustomOutput(
		baseDir,
		queriesDir,
		schemaDir,
	).Output

	Expect(result.TemplateData.Output.BaseDir).To(Equal(baseDir))
	Expect(result.TemplateData.Output.QueriesDir).To(Equal(queriesDir))
	Expect(result.TemplateData.Output.SchemaDir).To(Equal(schemaDir))
}
