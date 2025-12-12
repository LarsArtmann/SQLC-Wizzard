package wizard_test

import (
	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/wizard"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Wizard Execution Tests", func() {
	var wiz *wizard.Wizard

	BeforeEach(func() {
		wiz = wizard.NewWizard()
	})

	Describe("Wizard Run Method Execution", func() {
		It("should have proper default values", func() {
			result := wiz.GetResult()

			Expect(result).NotTo(BeNil())
			Expect(result.GenerateQueries).To(BeTrue())
			Expect(result.GenerateSchema).To(BeTrue())
		})

		It("should allow setting template data", func() {
			result := wiz.GetResult()

			// Verify we can set project name
			result.TemplateData.ProjectName = "test-project"
			Expect(result.TemplateData.ProjectName).To(Equal("test-project"))

			// Verify we can set project type
			result.TemplateData.ProjectType = generated.ProjectTypeMicroservice
			Expect(result.TemplateData.ProjectType).To(Equal(generated.ProjectTypeMicroservice))
		})

		It("should handle wizard result updates", func() {
			result := wiz.GetResult()

			// Test setting generate flags
			result.GenerateQueries = false
			result.GenerateSchema = false

			Expect(result.GenerateQueries).To(BeFalse())
			Expect(result.GenerateSchema).To(BeFalse())

			// Reset to defaults
			result.GenerateQueries = true
			result.GenerateSchema = true

			Expect(result.GenerateQueries).To(BeTrue())
			Expect(result.GenerateSchema).To(BeTrue())
		})

		It("should maintain wizard state", func() {
			// Get initial result
			result := wiz.GetResult()
			initialTemplateData := result.TemplateData

			// Modify template data
			result.TemplateData.ProjectName = "modified-project"
			result.TemplateData.Package.Name = "modifieddb"

			// Get result again and verify changes persist
			result2 := wiz.GetResult()
			Expect(result2.TemplateData.ProjectName).To(Equal("modified-project"))
			Expect(result2.TemplateData.Package.Name).To(Equal("modifieddb"))

			// Verify other fields are unchanged
			Expect(result2.TemplateData.ProjectType).To(Equal(initialTemplateData.ProjectType))
		})

		It("should handle complex template configurations", func() {
			result := wiz.GetResult()

			// Set up complex configuration
			result.TemplateData.ProjectName = "complex-project"
			result.TemplateData.Package.Name = "db"
			result.TemplateData.Package.Path = "github.com/company/complex"
			result.TemplateData.Package.BuildTags = "postgres,pgx,custom"

			// Verify complex configuration
			Expect(result.TemplateData.ProjectName).To(Equal("complex-project"))
			Expect(result.TemplateData.Package.Name).To(Equal("db"))
			Expect(result.TemplateData.Package.Path).To(Equal("github.com/company/complex"))
			Expect(result.TemplateData.Package.BuildTags).To(Equal("postgres,pgx,custom"))
		})
	})

	Describe("Wizard Configuration Updates", func() {
		It("should handle database configuration changes", func() {
			result := wiz.GetResult()

			// Test PostgreSQL configuration
			result.TemplateData.Database.Engine = generated.DatabaseTypePostgreSQL
			result.TemplateData.Database.UseUUIDs = true
			result.TemplateData.Database.UseJSON = true
			result.TemplateData.Database.UseArrays = true
			result.TemplateData.Database.UseFullText = true

			Expect(result.TemplateData.Database.Engine).To(Equal(generated.DatabaseTypePostgreSQL))
			Expect(result.TemplateData.Database.UseUUIDs).To(BeTrue())
			Expect(result.TemplateData.Database.UseJSON).To(BeTrue())
			Expect(result.TemplateData.Database.UseArrays).To(BeTrue())
			Expect(result.TemplateData.Database.UseFullText).To(BeTrue())

			// Test SQLite configuration
			result.TemplateData.Database.Engine = generated.DatabaseTypeSQLite
			result.TemplateData.Database.UseUUIDs = false
			result.TemplateData.Database.UseJSON = false
			result.TemplateData.Database.UseArrays = false
			result.TemplateData.Database.UseFullText = false

			Expect(result.TemplateData.Database.Engine).To(Equal(generated.DatabaseTypeSQLite))
			Expect(result.TemplateData.Database.UseUUIDs).To(BeFalse())
			Expect(result.TemplateData.Database.UseJSON).To(BeFalse())
			Expect(result.TemplateData.Database.UseArrays).To(BeFalse())
			Expect(result.TemplateData.Database.UseFullText).To(BeFalse())
		})

		It("should handle validation configuration changes", func() {
			result := wiz.GetResult()

			// Test emit options
			result.TemplateData.Validation.EmitOptions.EmitJSONTags = true
			result.TemplateData.Validation.EmitOptions.EmitPreparedQueries = true
			result.TemplateData.Validation.EmitOptions.EmitInterface = true
			result.TemplateData.Validation.EmitOptions.EmitEmptySlices = true
			result.TemplateData.Validation.EmitOptions.JSONTagsCaseStyle = "camel"

			Expect(result.TemplateData.Validation.EmitOptions.EmitJSONTags).To(BeTrue())
			Expect(result.TemplateData.Validation.EmitOptions.EmitPreparedQueries).To(BeTrue())
			Expect(result.TemplateData.Validation.EmitOptions.EmitInterface).To(BeTrue())
			Expect(result.TemplateData.Validation.EmitOptions.EmitEmptySlices).To(BeTrue())
			Expect(result.TemplateData.Validation.EmitOptions.JSONTagsCaseStyle).To(Equal("camel"))

			// Test safety rules
			result.TemplateData.Validation.SafetyRules.NoSelectStar = true
			result.TemplateData.Validation.SafetyRules.RequireWhere = true
			result.TemplateData.Validation.SafetyRules.RequireLimit = false
			result.TemplateData.Validation.SafetyRules.NoDropTable = true
			result.TemplateData.Validation.SafetyRules.NoTruncate = true

			Expect(result.TemplateData.Validation.SafetyRules.NoSelectStar).To(BeTrue())
			Expect(result.TemplateData.Validation.SafetyRules.RequireWhere).To(BeTrue())
			Expect(result.TemplateData.Validation.SafetyRules.RequireLimit).To(BeFalse())
			Expect(result.TemplateData.Validation.SafetyRules.NoDropTable).To(BeTrue())
			Expect(result.TemplateData.Validation.SafetyRules.NoTruncate).To(BeTrue())
		})

		It("should handle output configuration changes", func() {
			result := wiz.GetResult()

			// Test custom output directories
			result.TemplateData.Output.BaseDir = "./custom/output"
			result.TemplateData.Output.QueriesDir = "./custom/queries"
			result.TemplateData.Output.SchemaDir = "./custom/schema"

			Expect(result.TemplateData.Output.BaseDir).To(Equal("./custom/output"))
			Expect(result.TemplateData.Output.QueriesDir).To(Equal("./custom/queries"))
			Expect(result.TemplateData.Output.SchemaDir).To(Equal("./custom/schema"))

			// Test absolute paths
			result.TemplateData.Output.BaseDir = "/absolute/path/output"
			result.TemplateData.Output.QueriesDir = "/absolute/path/queries"
			result.TemplateData.Output.SchemaDir = "/absolute/path/schema"

			Expect(result.TemplateData.Output.BaseDir).To(Equal("/absolute/path/output"))
			Expect(result.TemplateData.Output.QueriesDir).To(Equal("/absolute/path/queries"))
			Expect(result.TemplateData.Output.SchemaDir).To(Equal("/absolute/path/schema"))
		})
	})

	Describe("Wizard State Management", func() {
		It("should maintain separate wizard instances", func() {
			// Create two wizard instances
			wiz1 := wizard.NewWizard()
			wiz2 := wizard.NewWizard()

			// Modify first wizard
			result1 := wiz1.GetResult()
			result1.TemplateData.ProjectName = "wizard1-project"

			// Modify second wizard
			result2 := wiz2.GetResult()
			result2.TemplateData.ProjectName = "wizard2-project"

			// Verify they're independent
			Expect(result1.TemplateData.ProjectName).To(Equal("wizard1-project"))
			Expect(result2.TemplateData.ProjectName).To(Equal("wizard2-project"))
		})

		It("should handle wizard result modifications", func() {
			result := wiz.GetResult()

			// Store initial state
			initialProjectName := result.TemplateData.ProjectName
			initialGenerateQueries := result.GenerateQueries

			// Modify state
			result.TemplateData.ProjectName = "new-project"
			result.GenerateQueries = false

			// Verify changes
			Expect(result.TemplateData.ProjectName).To(Equal("new-project"))
			Expect(result.GenerateQueries).To(BeFalse())

			// Verify initial values are different
			Expect(initialProjectName).NotTo(Equal(result.TemplateData.ProjectName))
			Expect(initialGenerateQueries).NotTo(Equal(result.GenerateQueries))
		})

		It("should handle template data resets", func() {
			result := wiz.GetResult()

			// Set up complex state
			result.TemplateData.ProjectName = "complex-project"
			result.TemplateData.ProjectType = "enterprise"
			result.TemplateData.Database.Engine = "postgresql"
			result.TemplateData.Package.Name = "db"
			result.TemplateData.Package.Path = "github.com/company/enterprise"

			// Reset to simple state
			result.TemplateData.ProjectName = "simple-project"
			result.TemplateData.ProjectType = generated.ProjectTypeHobby
			result.TemplateData.Database.Engine = generated.DatabaseTypeSQLite
			result.TemplateData.Package.Name = "simpledb"
			result.TemplateData.Package.Path = "github.com/user/simple"

			// Verify reset
			Expect(result.TemplateData.ProjectName).To(Equal("simple-project"))
			Expect(result.TemplateData.ProjectType).To(Equal(generated.ProjectTypeHobby))
			Expect(result.TemplateData.Database.Engine).To(Equal(generated.DatabaseTypeSQLite))
			Expect(result.TemplateData.Package.Name).To(Equal("simpledb"))
			Expect(result.TemplateData.Package.Path).To(Equal("github.com/user/simple"))
		})
	})
})
