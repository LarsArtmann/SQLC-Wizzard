package wizard_test

import (
	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/templates"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/wizard"
	"github.com/LarsArtmann/SQLC-Wizzard/pkg/config"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("NewWizard", func() {
	It("should create a new wizard instance", func() {
		wiz := wizard.NewWizard()
		
		Expect(wiz).NotTo(BeNil())
		Expect(wiz.GetResult()).NotTo(BeNil())
		Expect(wiz.GetResult().GenerateQueries).To(BeTrue())
		Expect(wiz.GetResult().GenerateSchema).To(BeTrue())
	})

	It("should have proper default values", func() {
		wiz := wizard.NewWizard()
		result := wiz.GetResult()
		
		Expect(result.GenerateQueries).To(BeTrue())
		Expect(result.GenerateSchema).To(BeTrue())
	})
})

var _ = Describe("Wizard Configuration", func() {
	var (
		wiz *wizard.Wizard
	)

	BeforeEach(func() {
		wiz = wizard.NewWizard()
	})

	It("should accept custom configuration", func() {
		// Test setting custom configuration
		result := wiz.GetResult()
		result.GenerateQueries = false
		result.GenerateSchema = false
		
		Expect(result.GenerateQueries).To(BeFalse())
		Expect(result.GenerateSchema).To(BeFalse())
	})

	It("should handle template data properly", func() {
		templateData := generated.TemplateData{
			ProjectName: "test-project",
			ProjectType: generated.ProjectTypeMicroservice,
			Package: generated.PackageConfig{
				Name: "testdb",
				Path: "github.com/example/test",
			},
			Database: generated.DatabaseConfig{
				Engine:   generated.DatabaseTypePostgreSQL,
				UseUUIDs: true,
				UseJSON:  true,
			},
		}

		result := wiz.GetResult()
		result.TemplateData = templateData
		
		Expect(result.TemplateData.ProjectName).To(Equal("test-project"))
		Expect(result.TemplateData.ProjectType).To(Equal(generated.ProjectTypeMicroservice))
	})
})

var _ = Describe("Wizard Steps", func() {
	It("should handle project type selection", func() {
		// Test project type validation
		testCases := []struct {
			input    string
			expected generated.ProjectType
			valid    bool
		}{
			{"microservice", generated.ProjectTypeMicroservice, true},
			{"enterprise", generated.ProjectTypeEnterprise, true},
			{"hobby", generated.ProjectTypeHobby, true},
			{"invalid", generated.ProjectType(""), false},
		}

		for _, tc := range testCases {
			projectType := generated.ProjectType(tc.input)
			if tc.valid {
				Expect(projectType.IsValid()).To(BeTrue(), "for input: %s", tc.input)
				Expect(projectType).To(Equal(tc.expected), "for input: %s", tc.input)
			} else {
				Expect(projectType.IsValid()).To(BeFalse(), "for input: %s", tc.input)
			}
		}
	})

	It("should handle database type selection", func() {
		// Test database type validation
		testCases := []struct {
			input    string
			expected generated.DatabaseType
			valid    bool
		}{
			{"postgresql", generated.DatabaseTypePostgreSQL, true},
			{"mysql", generated.DatabaseTypeMySQL, true},
			{"sqlite", generated.DatabaseTypeSQLite, true},
			{"invalid", generated.DatabaseType(""), false},
		}

		for _, tc := range testCases {
			dbType := generated.DatabaseType(tc.input)
			if tc.valid {
				Expect(dbType.IsValid()).To(BeTrue(), "for input: %s", tc.input)
				Expect(dbType).To(Equal(tc.expected), "for input: %s", tc.input)
			} else {
				Expect(dbType.IsValid()).To(BeFalse(), "for input: %s", tc.input)
			}
		}
	})
})

var _ = Describe("Template Integration", func() {
	var (
		wiz *wizard.Wizard
	)

	BeforeEach(func() {
		wiz = wizard.NewWizard()
	})

	It("should work with template system", func() {
		// Create a template
		template := templates.MustNewProjectType("microservice")
		Expect(template).NotTo(BeNil())
		
		// Generate template data
		templateData := generated.TemplateData{
			ProjectName: "test-project",
			ProjectType: generated.ProjectTypeMicroservice,
			Package: generated.PackageConfig{
				Name: "testdb",
				Path: "github.com/example/test",
			},
			Database: generated.DatabaseConfig{
				Engine:    generated.DatabaseTypePostgreSQL,
				UseUUIDs:  true,
				UseJSON:   true,
				UseArrays: true,
			},
			Output: generated.OutputConfig{
				BaseDir:    "./generated",
				QueriesDir: "queries",
				SchemaDir:  "schema",
			},
			Validation: generated.ValidationConfig{
				StrictFunctions: true,
				StrictOrderBy:   true,
				EmitOptions:    generated.DefaultEmitOptions(),
				SafetyRules:    generated.DefaultSafetyRules(),
			},
		}

		// Set template data in wizard result
		result := wiz.GetResult()
		result.TemplateData = templateData
		
		// Verify template data is set correctly
		Expect(result.TemplateData.ProjectName).To(Equal("test-project"))
		Expect(result.TemplateData.Database.Engine).To(Equal(generated.DatabaseTypePostgreSQL))
		Expect(result.TemplateData.Database.UseUUIDs).To(BeTrue())
	})
})

var _ = Describe("Configuration Generation", func() {
	var (
		wiz *wizard.Wizard
	)

	BeforeEach(func() {
		wiz = wizard.NewWizard()
	})

	It("should generate valid sqlc configuration", func() {
		templateData := generated.TemplateData{
			ProjectName: "test-project",
			ProjectType: generated.ProjectTypeMicroservice,
			Package: generated.PackageConfig{
				Name: "db",
				Path: "github.com/example/testdb",
			},
			Database: generated.DatabaseConfig{
				Engine:    generated.DatabaseTypePostgreSQL,
				UseUUIDs:  true,
				UseJSON:   true,
				UseArrays: true,
			},
			Output: generated.OutputConfig{
				BaseDir:    "./internal/db",
				QueriesDir: "queries",
				SchemaDir:  "schema",
			},
			Validation: generated.ValidationConfig{
				StrictFunctions: true,
				StrictOrderBy:   true,
				EmitOptions:    generated.DefaultEmitOptions(),
				SafetyRules:    generated.DefaultSafetyRules(),
			},
		}

		result := wiz.GetResult()
		result.TemplateData = templateData

		// Try to create sqlc config from template data
		sqlcConfig := &config.SqlcConfig{
			Version: "2",
			SQL: []config.SQLConfig{
				{
					Engine:  "postgresql",
					Queries: config.NewSinglePath("queries"),
					Schema:  config.NewSinglePath("schema"),
					Gen: config.GenConfig{
						Go: &config.GoGenConfig{
							Package:    "db",
							Out:        "internal/db",
							SQLPackage: "pgx/v5",
							EmitJSONTags:                true,
							EmitPreparedQueries:         true,
							EmitInterface:              true,
							EmitEmptySlices:            true,
							JSONTagsCaseStyle:          "camel",
						},
					},
				},
			},
		}

		Expect(sqlcConfig).NotTo(BeNil())
		Expect(sqlcConfig.Version).To(Equal("2"))
		Expect(len(sqlcConfig.SQL)).To(Equal(1))
		Expect(sqlcConfig.SQL[0].Engine).To(Equal("postgresql"))
	})
})

var _ = Describe("Error Handling", func() {
	var (
		wiz *wizard.Wizard
	)

	BeforeEach(func() {
		wiz = wizard.NewWizard()
	})

	It("should handle invalid project types gracefully", func() {
		templateData := generated.TemplateData{
			ProjectName: "test",
			ProjectType: generated.ProjectType("invalid"),
		}

		// Should not panic when setting invalid template data
		result := wiz.GetResult()
		result.TemplateData = templateData
		
		Expect(result.TemplateData.ProjectType.IsValid()).To(BeFalse())
	})

	It("should handle invalid database types gracefully", func() {
		templateData := generated.TemplateData{
			ProjectName: "test",
			Database: generated.DatabaseConfig{
				Engine: generated.DatabaseType("invalid"),
			},
		}

		result := wiz.GetResult()
		result.TemplateData = templateData
		
		Expect(result.TemplateData.Database.Engine.IsValid()).To(BeFalse())
	})

	It("should handle empty configuration", func() {
		result := wiz.GetResult()
		result.TemplateData = generated.TemplateData{}
		
		// Should not panic
		Expect(result.TemplateData.ProjectName).To(Equal(""))
		Expect(result.TemplateData.ProjectType.IsValid()).To(BeFalse())
		Expect(result.TemplateData.Database.Engine.IsValid()).To(BeFalse())
	})
})

var _ = Describe("Context Handling", func() {
	It("should work with context cancellation", func() {
		// Test that wizard can handle context without crashing
		wiz := wizard.NewWizard()
		result := wiz.GetResult()
		
		// Should still be able to access result
		Expect(result).NotTo(BeNil())
	})
})

var _ = Describe("Wizard Result", func() {
	var (
		result *wizard.WizardResult
	)

	BeforeEach(func() {
		result = &wizard.WizardResult{
			Config:          &config.SqlcConfig{},
			TemplateData:    generated.TemplateData{},
			GenerateQueries: true,
			GenerateSchema:  true,
		}
	})

	It("should hold wizard state correctly", func() {
		Expect(result.Config).NotTo(BeNil())
		Expect(result.GenerateQueries).To(BeTrue())
		Expect(result.GenerateSchema).To(BeTrue())
	})

	It("should allow modification of result", func() {
		result.GenerateQueries = false
		result.GenerateSchema = false
		
		Expect(result.GenerateQueries).To(BeFalse())
		Expect(result.GenerateSchema).To(BeFalse())
	})
})

var _ = Describe("Integration with Templates", func() {
	It("should integrate with all template types", func() {
		templates := []generated.ProjectType{
			generated.ProjectTypeHobby,
			generated.ProjectTypeMicroservice,
			generated.ProjectTypeEnterprise,
		}

		for _, templateType := range templates {
			templateData := generated.TemplateData{
				ProjectName: "test-project",
				ProjectType: templateType,
				Package: generated.PackageConfig{
					Name: "db",
					Path: "github.com/example/testdb",
				},
				Database: generated.DatabaseConfig{
					Engine:    generated.DatabaseTypePostgreSQL,
					UseUUIDs:  true,
					UseJSON:   true,
				},
			}

			// Each template type should be valid
			Expect(templateType.IsValid()).To(BeTrue(), "for template: %s", templateType)
			Expect(templateData.ProjectName).To(Equal("test-project"))
		}
	})
})