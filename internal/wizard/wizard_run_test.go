package wizard_test

import (
	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/wizard"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Wizard Run Method", func() {
	var wiz *wizard.Wizard

	BeforeEach(func() {
		wiz = wizard.NewWizard()
	})

	It("should initialize with proper default template data", func() {
		result := wiz.GetResult()

		Expect(result).NotTo(BeNil())
		Expect(result.GenerateQueries).To(BeTrue())
		Expect(result.GenerateSchema).To(BeTrue())
	})

	It("should handle valid template data generation", func() {
		result := wiz.GetResult()
		result.TemplateData = createTemplateDataWithFeatures("test-project", generated.ProjectTypeMicroservice)
		result.TemplateData.Package.Path = "github.com/example/testdb"

		Expect(result.TemplateData.ProjectName).To(Equal("test-project"))
		Expect(result.TemplateData.ProjectType).To(Equal(generated.ProjectTypeMicroservice))
		Expect(result.TemplateData.Database.Engine).To(Equal(generated.DatabaseTypePostgreSQL))
		Expect(result.TemplateData.Database.UseUUIDs).To(BeTrue())
		Expect(result.TemplateData.Database.UseJSON).To(BeTrue())
	})

	It("should validate database configuration", func() {
		result := wiz.GetResult()
		result.TemplateData = generated.TemplateData{
			Database: generated.DatabaseConfig{
				Engine: generated.DatabaseTypeSQLite,
			},
		}

		dbConfig := result.TemplateData.Database
		Expect(dbConfig.Engine).To(Equal(generated.DatabaseTypeSQLite))
		Expect(dbConfig.Engine.IsValid()).To(BeTrue())
	})

	It("should validate project configuration", func() {
		result := wiz.GetResult()
		result.TemplateData = generated.TemplateData{
			ProjectType: generated.ProjectTypeHobby,
			Package: generated.PackageConfig{
				Name: "testpkg",
				Path: "github.com/user/test",
			},
		}

		projectConfig := result.TemplateData
		Expect(projectConfig.ProjectType).To(Equal(generated.ProjectTypeHobby))
		Expect(projectConfig.ProjectType.IsValid()).To(BeTrue())
		Expect(projectConfig.Package.Name).To(Equal("testpkg"))
		Expect(projectConfig.Package.Path).To(Equal("github.com/user/test"))
	})

	It("should handle emit options configuration", func() {
		result := wiz.GetResult()
		emitOpts := generated.DefaultEmitOptions()

		result.TemplateData.Validation.EmitOptions = emitOpts

		Expect(result.TemplateData.Validation.EmitOptions.EmitJSONTags).To(BeTrue())
		Expect(result.TemplateData.Validation.EmitOptions.EmitPreparedQueries).To(BeTrue())
		Expect(result.TemplateData.Validation.EmitOptions.EmitInterface).To(BeTrue())
	})

	It("should handle safety rules configuration", func() {
		result := wiz.GetResult()
		safetyRules := generated.DefaultSafetyRules()

		result.TemplateData.Validation.SafetyRules = safetyRules

		Expect(result.TemplateData.Validation.SafetyRules.NoSelectStar).To(BeTrue())
		Expect(result.TemplateData.Validation.SafetyRules.RequireWhere).To(BeTrue())
	})

	It("should validate all project types", func() {
		validProjectTypes := []generated.ProjectType{
			generated.ProjectTypeHobby,
			generated.ProjectTypeMicroservice,
			generated.ProjectTypeEnterprise,
			generated.ProjectTypeAPIFirst,
			generated.ProjectTypeAnalytics,
			generated.ProjectTypeTesting,
			generated.ProjectTypeMultiTenant,
			generated.ProjectTypeLibrary,
		}

		for _, projectType := range validProjectTypes {
			Expect(projectType.IsValid()).To(BeTrue(),
				"Project type %s should be valid", projectType)
		}
	})

	It("should validate all database types", func() {
		validDatabaseTypes := []generated.DatabaseType{
			generated.DatabaseTypePostgreSQL,
			generated.DatabaseTypeMySQL,
			generated.DatabaseTypeSQLite,
		}

		for _, dbType := range validDatabaseTypes {
			Expect(dbType.IsValid()).To(BeTrue(),
				"Database type %s should be valid", dbType)
		}
	})

	It("should handle invalid project types", func() {
		invalidProjectType := generated.ProjectType("invalid-type")
		Expect(invalidProjectType.IsValid()).To(BeFalse())
	})

	It("should handle invalid database types", func() {
		invalidDatabaseType := generated.DatabaseType("invalid-db")
		Expect(invalidDatabaseType.IsValid()).To(BeFalse())
	})
})

var _ = Describe("Wizard Configuration Generation", func() {
	var (
		wiz    *wizard.Wizard
		result *wizard.WizardResult
	)

	BeforeEach(func() {
		wiz = wizard.NewWizard()
		result = wiz.GetResult()
	})

	It("should generate valid configuration for microservice template", func() {
		result.TemplateData = generated.TemplateData{
			ProjectName: "test-microservice",
			ProjectType: generated.ProjectTypeMicroservice,
			Package: generated.PackageConfig{
				Name: "db",
				Path: "github.com/company/microservice",
			},
			Database: generated.DatabaseConfig{
				Engine:    generated.DatabaseTypePostgreSQL,
				UseUUIDs:  true,
				UseJSON:   true,
				UseArrays: true,
			},
		}

		Expect(result.TemplateData.ProjectType).To(Equal(generated.ProjectTypeMicroservice))
		Expect(result.TemplateData.Database.Engine).To(Equal(generated.DatabaseTypePostgreSQL))
		Expect(result.TemplateData.Package.Name).To(Equal("db"))
	})

	It("should generate valid configuration for hobby template", func() {
		result.TemplateData = generated.TemplateData{
			ProjectName: "test-hobby",
			ProjectType: generated.ProjectTypeHobby,
			Package: generated.PackageConfig{
				Name: "db",
				Path: "github.com/user/hobbyproject",
			},
			Database: generated.DatabaseConfig{
				Engine:   generated.DatabaseTypeSQLite,
				UseUUIDs: false,
				UseJSON:  false,
			},
		}

		Expect(result.TemplateData.ProjectType).To(Equal(generated.ProjectTypeHobby))
		Expect(result.TemplateData.Database.Engine).To(Equal(generated.DatabaseTypeSQLite))
		Expect(result.TemplateData.Database.UseUUIDs).To(BeFalse())
	})
})
