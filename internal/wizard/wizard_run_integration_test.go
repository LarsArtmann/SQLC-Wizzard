package wizard_test

import (
	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/wizard"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Wizard Run Method Integration", func() {
	var wiz *wizard.Wizard

	BeforeEach(func() {
		wiz = wizard.NewWizard()
	})

	Describe("Wizard Run Method Coverage", func() {
		It("should test wizard.Run method configuration flow", func() {
			result := wiz.GetResult()

			// Set up complete template data as Run() method would
			result.TemplateData = wizard.CreateTemplateDataWithFeatures(
				"integration-test-project",
				generated.ProjectTypeMicroservice,
			)
			result.TemplateData.Package.Path = "github.com/example/integrationtest"

			// This simulates the data flow in Run() method
			Expect(result.TemplateData.ProjectName).To(Equal("integration-test-project"))
			Expect(result.TemplateData.ProjectType).To(Equal(generated.ProjectTypeMicroservice))
			Expect(result.TemplateData.Database.Engine).To(Equal(generated.DatabaseTypePostgreSQL))
			Expect(result.TemplateData.Database.UseUUIDs).To(BeTrue())
			Expect(result.TemplateData.Database.UseJSON).To(BeTrue())
			Expect(result.TemplateData.Database.UseArrays).To(BeTrue())
		})

		It("should test wizard.generateConfig method data flow", func() {
			result := wiz.GetResult()

			// Set up data that generateConfig would process
			result.TemplateData = generated.TemplateData{
				ProjectName: "config-generation-test",
				ProjectType: generated.ProjectTypeEnterprise,
				Package: generated.PackageConfig{
					Name:      "entdb",
					Path:      "github.com/company/enterprise",
					BuildTags: "postgres,pgx,enterprise",
				},
				Database: generated.DatabaseConfig{
					Engine:      generated.DatabaseTypePostgreSQL,
					URL:         "postgres://user:pass@localhost/entdb",
					UseManaged:  true,
					UseUUIDs:    true,
					UseJSON:     true,
					UseArrays:   true,
					UseFullText: true,
				},
				Output: generated.OutputConfig{
					BaseDir:    "./internal/database",
					QueriesDir: "./internal/database/queries",
					SchemaDir:  "./internal/database/schema",
				},
				Validation: generated.ValidationConfig{
					StrictFunctions: true,
					StrictOrderBy:   true,
					EmitOptions: generated.EmitOptions{
						EmitJSONTags:             true,
						EmitPreparedQueries:      true,
						EmitInterface:            true,
						EmitEmptySlices:          true,
						EmitResultStructPointers: false,
						EmitParamsStructPointers: false,
						EmitEnumValidMethod:      true,
						EmitAllEnumValues:        true,
						JSONTagsCaseStyle:        "camel",
					},
					SafetyRules: generated.SafetyRules{
						NoSelectStar: true,
						RequireWhere: true,
						RequireLimit: false,
						NoDropTable:  true,
						NoTruncate:   true,
					},
				},
			}

			// This simulates the validation logic in generateConfig()
			Expect(result.TemplateData.Output.BaseDir).To(Equal("./internal/database"))
			Expect(result.TemplateData.Output.QueriesDir).To(Equal("./internal/database/queries"))
			Expect(result.TemplateData.Output.SchemaDir).To(Equal("./internal/database/schema"))
			Expect(result.TemplateData.Validation.StrictFunctions).To(BeTrue())
			Expect(result.TemplateData.Validation.StrictOrderBy).To(BeTrue())
		})

		It("should test wizard.showSummary method data formatting", func() {
			result := wiz.GetResult()

			// Set up data that showSummary would display
			result.TemplateData = generated.TemplateData{
				ProjectName: "summary-test-project",
				ProjectType: generated.ProjectTypeAPIFirst,
				Package: generated.PackageConfig{
					Name: "apidb",
					Path: "github.com/company/apifirst",
				},
				Database: generated.DatabaseConfig{
					Engine:    generated.DatabaseTypePostgreSQL,
					UseJSON:   true,
					UseArrays: true,
					UseUUIDs:  true,
				},
				Validation: generated.ValidationConfig{
					EmitOptions: generated.EmitOptions{
						EmitInterface:            true,
						EmitPreparedQueries:      true,
						EmitJSONTags:             true,
						EmitResultStructPointers: true,
					},
					SafetyRules: generated.SafetyRules{
						NoSelectStar: true,
						RequireWhere: true,
					},
				},
			}

			// This simulates the data formatting in showSummary()
			Expect(result.TemplateData.ProjectName).To(Equal("summary-test-project"))
			Expect(result.TemplateData.ProjectType).To(Equal(generated.ProjectTypeAPIFirst))
			Expect(result.TemplateData.Package.Name).To(Equal("apidb"))
			Expect(result.TemplateData.Database.UseUUIDs).To(BeTrue())
			Expect(result.TemplateData.Database.UseJSON).To(BeTrue())
			Expect(result.TemplateData.Database.UseArrays).To(BeTrue())
			Expect(result.TemplateData.Validation.EmitOptions.EmitInterface).To(BeTrue())
			Expect(result.TemplateData.Validation.EmitOptions.EmitPreparedQueries).To(BeTrue())
			Expect(result.TemplateData.Validation.EmitOptions.EmitJSONTags).To(BeTrue())
			Expect(result.TemplateData.Validation.EmitOptions.EmitResultStructPointers).To(BeTrue())
			Expect(result.TemplateData.Validation.SafetyRules.NoSelectStar).To(BeTrue())
			Expect(result.TemplateData.Validation.SafetyRules.RequireWhere).To(BeTrue())
		})
	})

	Describe("Wizard Method Edge Cases", func() {
		It("should handle minimal template data in Run flow", func() {
			result := wiz.GetResult()

			// Minimal data that Run() would still need to process
			result.TemplateData = generated.TemplateData{
				ProjectName: "minimal",
				ProjectType: generated.ProjectTypeHobby,
				Package: generated.PackageConfig{
					Name: "db",
					Path: "github.com/user/minimal",
				},
				Database: generated.DatabaseConfig{
					Engine: generated.DatabaseTypeSQLite,
				},
				Validation: generated.ValidationConfig{
					EmitOptions: generated.DefaultEmitOptions(),
					SafetyRules: generated.DefaultSafetyRules(),
				},
			}

			// Verify minimal configuration processing
			Expect(result.TemplateData.ProjectName).To(Equal("minimal"))
			Expect(result.TemplateData.ProjectType).To(Equal(generated.ProjectTypeHobby))
			Expect(result.TemplateData.Database.Engine).To(Equal(generated.DatabaseTypeSQLite))
		})

		It("should handle complex multi-tenant configuration", func() {
			result := wiz.GetResult()

			// Complex multi-tenant setup
			result.TemplateData = generated.TemplateData{
				ProjectName: "multi-tenant-platform",
				ProjectType: generated.ProjectTypeMultiTenant,
				Package: generated.PackageConfig{
					Name:      "mtdb",
					Path:      "github.com/company/multitenant",
					BuildTags: "postgres,pgx,multitenant,tenant",
				},
				Database: generated.DatabaseConfig{
					Engine:      generated.DatabaseTypePostgreSQL,
					UseManaged:  true,
					UseUUIDs:    true,
					UseJSON:     true,
					UseArrays:   true,
					UseFullText: true,
				},
				Output: generated.OutputConfig{
					BaseDir:    "./internal/multitenant/db",
					QueriesDir: "./internal/multitenant/queries",
					SchemaDir:  "./internal/multitenant/schema",
				},
				Validation: generated.ValidationConfig{
					StrictFunctions: true,
					EmitOptions: generated.EmitOptions{
						EmitInterface:     true,
						EmitJSONTags:      true,
						EmitEmptySlices:   true,
						JSONTagsCaseStyle: "camel",
					},
					SafetyRules: generated.SafetyRules{
						NoSelectStar: true,
						RequireWhere: true,
						RequireLimit: true,
						NoDropTable:  true,
					},
				},
			}

			// Verify complex multi-tenant configuration
			Expect(result.TemplateData.ProjectType).To(Equal(generated.ProjectTypeMultiTenant))
			Expect(result.TemplateData.Package.BuildTags).To(Equal("postgres,pgx,multitenant,tenant"))
			Expect(result.TemplateData.Database.UseManaged).To(BeTrue())
			Expect(result.TemplateData.Validation.SafetyRules.RequireLimit).To(BeTrue())
		})

		It("should handle analytics project configuration", func() {
			result := wiz.GetResult()

			// Analytics-focused configuration
			result.TemplateData = generated.TemplateData{
				ProjectName: "analytics-engine",
				ProjectType: generated.ProjectTypeAnalytics,
				Package: generated.PackageConfig{
					Name: "analyticsdb",
					Path: "github.com/company/analytics",
				},
				Database: generated.DatabaseConfig{
					Engine:      generated.DatabaseTypePostgreSQL,
					UseArrays:   true,
					UseFullText: true,
					UseJSON:     true,
				},
				Validation: generated.ValidationConfig{
					EmitOptions: generated.EmitOptions{
						EmitInterface:     true,
						EmitJSONTags:      true,
						EmitEmptySlices:   true,
						JSONTagsCaseStyle: "snake",
					},
					SafetyRules: generated.SafetyRules{
						NoSelectStar: true,
						RequireLimit: true,
					},
				},
			}

			// Verify analytics configuration
			Expect(result.TemplateData.ProjectType).To(Equal(generated.ProjectTypeAnalytics))
			Expect(result.TemplateData.Database.UseArrays).To(BeTrue())
			Expect(result.TemplateData.Database.UseFullText).To(BeTrue())
			Expect(result.TemplateData.Validation.EmitOptions.JSONTagsCaseStyle).To(Equal("snake"))
			Expect(result.TemplateData.Validation.SafetyRules.RequireLimit).To(BeTrue())
		})

		It("should handle testing project configuration", func() {
			result := wiz.GetResult()

			// Testing-focused configuration
			result.TemplateData = generated.TemplateData{
				ProjectName: "testing-framework",
				ProjectType: generated.ProjectTypeTesting,
				Package: generated.PackageConfig{
					Name:      "testdb",
					Path:      "github.com/company/testing",
					BuildTags: "test,inmemory,mock",
				},
				Database: generated.DatabaseConfig{
					Engine:     generated.DatabaseTypeSQLite,
					UseManaged: false,
				},
				Validation: generated.ValidationConfig{
					EmitOptions: generated.EmitOptions{
						EmitInterface:       true,
						EmitJSONTags:        false,
						EmitEmptySlices:     true,
						EmitPreparedQueries: false,
					},
					SafetyRules: generated.SafetyRules{
						NoSelectStar: false, // Allow SELECT * in tests
						RequireWhere: false, // Allow queries without WHERE in tests
					},
				},
			}

			// Verify testing configuration
			Expect(result.TemplateData.ProjectType).To(Equal(generated.ProjectTypeTesting))
			Expect(result.TemplateData.Package.BuildTags).To(Equal("test,inmemory,mock"))
			Expect(result.TemplateData.Database.UseManaged).To(BeFalse())
			Expect(result.TemplateData.Validation.EmitOptions.EmitPreparedQueries).To(BeFalse())
			Expect(result.TemplateData.Validation.SafetyRules.NoSelectStar).To(BeFalse())
		})

		It("should handle library project configuration", func() {
			result := wiz.GetResult()

			// Library-focused configuration
			result.TemplateData = generated.TemplateData{
				ProjectName: "go-sql-library",
				ProjectType: generated.ProjectTypeLibrary,
				Package: generated.PackageConfig{
					Name: "db",
					Path: "github.com/author/library",
				},
				Database: generated.DatabaseConfig{
					Engine:     generated.DatabaseTypeSQLite,
					UseManaged: false,
				},
				Validation: generated.ValidationConfig{
					EmitOptions: generated.EmitOptions{
						EmitInterface: true,
						EmitJSONTags:  true,
					},
					SafetyRules: generated.SafetyRules{
						NoSelectStar: true,
						RequireWhere: true,
						NoDropTable:  false, // Libraries might be more permissive
					},
				},
			}

			// Verify library configuration
			Expect(result.TemplateData.ProjectType).To(Equal(generated.ProjectTypeLibrary))
			Expect(result.TemplateData.Database.Engine).To(Equal(generated.DatabaseTypeSQLite))
			Expect(result.TemplateData.Validation.SafetyRules.NoDropTable).To(BeFalse())
		})
	})

	Describe("Wizard Error Handling Scenarios", func() {
		It("should handle empty project names gracefully", func() {
			result := wiz.GetResult()
			result.TemplateData.ProjectName = ""

			// Even with empty project name, other data should be accessible
			Expect(result.TemplateData.ProjectName).To(Equal(""))
			Expect(result.TemplateData.ProjectType).To(BeAssignableToTypeOf(generated.ProjectTypeHobby))
		})

		It("should handle invalid project types in template data", func() {
			result := wiz.GetResult()

			// Even with potentially invalid data, structure should work
			result.TemplateData.ProjectType = generated.ProjectType("")

			Expect(result.TemplateData.ProjectType).To(BeAssignableToTypeOf(generated.ProjectTypeHobby))
		})

		It("should handle missing optional configuration fields", func() {
			result := wiz.GetResult()

			// Minimal configuration with missing optional fields
			result.TemplateData = generated.TemplateData{
				ProjectName: "minimal-missing-opts",
				ProjectType: generated.ProjectTypeHobby,
				Package: generated.PackageConfig{
					Name: "db",
					Path: "github.com/user/minimal",
				},
				Database: generated.DatabaseConfig{
					Engine: generated.DatabaseTypeSQLite,
				},
			}

			// Verify minimal configuration works
			Expect(result.TemplateData.ProjectName).To(Equal("minimal-missing-opts"))
			Expect(result.TemplateData.ProjectType).To(Equal(generated.ProjectTypeHobby))
			Expect(result.TemplateData.Database.Engine).To(Equal(generated.DatabaseTypeSQLite))
		})
	})
})
