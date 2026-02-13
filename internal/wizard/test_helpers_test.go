package wizard_test

import (
	"github.com/LarsArtmann/SQLC-Wizzard/internal/wizard"
	. "github.com/onsi/gomega"
)

// Test output path constants for reuse across test cases.
const (
	CustomBaseDir    = "./custom/db"
	CustomQueriesDir = "./custom/queries"
	CustomSchemaDir  = "./custom/schema"

	AbsoluteBaseDir    = "/absolute/path/db"
	RelativeQueriesDir = "./relative/queries"
	RelativeSchemaDir  = "../schema"
)

// testOutputPathConfiguration tests output directory configuration with custom paths.
func testOutputPathConfiguration(wiz *wizard.Wizard, baseDir, queriesDir, schemaDir string) {
	result := wiz.GetResult()
	result.TemplateData.Output = wizard.CreateTemplateDataWithCustomOutput(
		baseDir,
		queriesDir,
		schemaDir,
	).Output

	Expect(result.TemplateData.Output.BaseDir).To(Equal(baseDir))
	Expect(result.TemplateData.Output.QueriesDir).To(Equal(queriesDir))
	Expect(result.TemplateData.Output.SchemaDir).To(Equal(schemaDir))
}
