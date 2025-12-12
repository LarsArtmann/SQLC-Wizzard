package wizard_test

import (
	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/wizard"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Wizard Steps", func() {
	var (
		wiz *wizard.Wizard
	)

	BeforeEach(func() {
		wiz = wizard.NewWizard()
	})

	Describe("Template Data Configuration", func() {
		It("should handle complete template data configuration", func() {
			result := wiz.GetResult()
			result.TemplateData = generated.TemplateData{
				ProjectName: "complete-test-project",
				ProjectType: generated.ProjectTypeEnterprise,
				Package: generated.PackageConfig{
					Name:      "apib",
					Path:      "github.com/company/api",
					BuildTags: "postgres,pgx",
				},
				Database: generated.DatabaseConfig{
					Engine:      generated.DatabaseTypePostgreSQL,
					URL:         "postgres://user:pass@localhost/db",
					UseManaged:  true,
					UseUUIDs:    true,
					UseJSON:     true,
					UseArrays:   true,
					UseFullText: true,
				},
				Output: generated.OutputConfig{
					BaseDir:    "./internal/db",
					QueriesDir: "./internal/db/queries",
					SchemaDir:  "./internal/db/schema",
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
						EmitAllEnumValues:        false,
					},
					SafetyRules: generated.SafetyRules{
						NoSelectStar:    true,
						RequireWhere:     true,
						RequireLimit:     false,
						NoDropTable:      true,
					},
				},
			}

			// Verify enterprise configuration
			Expect(result.TemplateData.ProjectType).To(Equal(generated.ProjectTypeEnterprise))
			Expect(result.TemplateData.Database.UseManaged).To(BeTrue())
			Expect(result.TemplateData.Database.UseFullText).To(BeTrue())
			Expect(result.TemplateData.Validation.StrictFunctions).To(BeTrue())
		})

		It("should handle minimal template data configuration", func() {
			result := wiz.GetResult()
			result.TemplateData = generated.TemplateData{
				ProjectName: "minimal-project",
				ProjectType: generated.ProjectTypeHobby,
				Package: generated.PackageConfig{
					Name: "db",
					Path: "github.com/user/minimal",
				},
				Database: generated.DatabaseConfig{
					Engine: generated.DatabaseTypeSQLite,
				},
			}

			// Verify minimal configuration
			Expect(result.TemplateData.ProjectType).To(Equal(generated.ProjectTypeHobby))
			Expect(result.TemplateData.Database.Engine).To(Equal(generated.DatabaseTypeSQLite))
			Expect(result.TemplateData.Package.Name).To(Equal("db"))
		})

		It("should handle API-first project configuration", func() {
			result := wiz.GetResult()
			result.TemplateData = generated.TemplateData{
				ProjectName: "api-first-project",
				ProjectType: generated.ProjectTypeAPIFirst,
				Package: generated.PackageConfig{
					Name: "api",
					Path: "github.com/company/api",
				},
				Database: generated.DatabaseConfig{
					Engine:    generated.DatabaseTypePostgreSQL,
					UseJSON:   true,
					UseArrays: true,
				},
				Validation: generated.ValidationConfig{
					EmitOptions: generated.EmitOptions{
						EmitJSONTags:             true,
						EmitPreparedQueries:      true,
						EmitInterface:            true,
						EmitResultStructPointers: true, // API-first typically uses pointers
					},
				},
			}

			Expect(result.TemplateData.ProjectType).To(Equal(generated.ProjectTypeAPIFirst))
			Expect(result.TemplateData.Database.UseJSON).To(BeTrue())
			Expect(result.TemplateData.Validation.EmitOptions.EmitResultStructPointers).To(BeTrue())
		})

		It("should handle multi-tenant configuration", func() {
			result := wiz.GetResult()
			result.TemplateData = generated.TemplateData{
				ProjectName: "multi-tenant-saas",
				ProjectType: generated.ProjectTypeMultiTenant,
				Package: generated.PackageConfig{
					Name: "db",
					Path: "github.com/company/multitenant",
				},
				Database: generated.DatabaseConfig{
					Engine:    generated.DatabaseTypePostgreSQL,
					UseUUIDs:  true,
					UseArrays: true,
				},
				Validation: generated.ValidationConfig{
					SafetyRules: generated.SafetyRules{
						RequireWhere:     true,
						RequireLimit:     true,
						NoDropTable:      true,
					},
				},
			}

			Expect(result.TemplateData.ProjectType).To(Equal(generated.ProjectTypeMultiTenant))
			Expect(result.TemplateData.Validation.SafetyRules.RequireLimit).To(BeTrue())
		})

		It("should handle testing project configuration", func() {
			result := wiz.GetResult()
			result.TemplateData = generated.TemplateData{
				ProjectName: "testing-project",
				ProjectType: generated.ProjectTypeTesting,
				Package: generated.PackageConfig{
					Name:      "testdb",
					Path:      "github.com/company/testing",
					BuildTags: "test,inmemory",
				},
				Database: generated.DatabaseConfig{
					Engine:   generated.DatabaseTypeSQLite,
					UseManaged: false, // Tests typically use in-memory DB
				},
				Validation: generated.ValidationConfig{
					EmitOptions: generated.EmitOptions{
						EmitPreparedQueries: false, // Tests might not need prepared queries
						EmitInterface:       true,   // Always useful for testing
					},
				},
			}

			Expect(result.TemplateData.ProjectType).To(Equal(generated.ProjectTypeTesting))
			Expect(result.TemplateData.Database.UseManaged).To(BeFalse())
			Expect(result.TemplateData.Package.BuildTags).To(Equal("test,inmemory"))
		})

		It("should handle analytics project configuration", func() {
			result := wiz.GetResult()
			result.TemplateData = generated.TemplateData{
				ProjectName: "analytics-platform",
				ProjectType: generated.ProjectTypeAnalytics,
				Package: generated.PackageConfig{
					Name: "analytics",
					Path: "github.com/company/analytics",
				},
				Database: generated.DatabaseConfig{
					Engine:      generated.DatabaseTypePostgreSQL,
					UseArrays:   true, // Analytics often need arrays
					UseFullText: true, // Analytics often needs full-text search
				},
				Validation: generated.ValidationConfig{
					EmitOptions: generated.EmitOptions{
						EmitEmptySlices: true, // Analytics often handles empty results
					},
					SafetyRules: generated.SafetyRules{
						RequireLimit: true, // Analytics queries should be limited
					},
				},
			}

			Expect(result.TemplateData.ProjectType).To(Equal(generated.ProjectTypeAnalytics))
			Expect(result.TemplateData.Database.UseArrays).To(BeTrue())
			Expect(result.TemplateData.Database.UseFullText).To(BeTrue())
			Expect(result.TemplateData.Validation.SafetyRules.RequireLimit).To(BeTrue())
		})

		It("should handle library project configuration", func() {
			result := wiz.GetResult()
			result.TemplateData = generated.TemplateData{
				ProjectName: "go-library",
				ProjectType: generated.ProjectTypeLibrary,
				Package: generated.PackageConfig{
					Name: "db",
					Path: "github.com/author/library",
				},
				Database: generated.DatabaseConfig{
					Engine:   generated.DatabaseTypeSQLite, // Libraries often use SQLite for simplicity
					UseManaged: false,
				},
				Validation: generated.ValidationConfig{
					EmitOptions: generated.EmitOptions{
						EmitInterface: true, // Libraries should expose interfaces
					},
					SafetyRules: generated.SafetyRules{
						NoDropTable: false, // Libraries might be more flexible
					},
				},
			}

			Expect(result.TemplateData.ProjectType).To(Equal(generated.ProjectTypeLibrary))
			Expect(result.TemplateData.Database.Engine).To(Equal(generated.DatabaseTypeSQLite))
			Expect(result.TemplateData.Validation.EmitOptions.EmitInterface).To(BeTrue())
		})
	})

	Describe("Database Configuration Variations", func() {
		It("should handle PostgreSQL with all features", func() {
			result := wiz.GetResult()
			result.TemplateData.Database = generated.DatabaseConfig{
				Engine:      generated.DatabaseTypePostgreSQL,
				UseUUIDs:    true,
				UseJSON:     true,
				UseArrays:   true,
				UseFullText: true,
				UseManaged:  true,
			}

			Expect(result.TemplateData.Database.Engine.IsValid()).To(BeTrue())
			Expect(result.TemplateData.Database.UseUUIDs).To(BeTrue())
			Expect(result.TemplateData.Database.UseJSON).To(BeTrue())
			Expect(result.TemplateData.Database.UseArrays).To(BeTrue())
			Expect(result.TemplateData.Database.UseFullText).To(BeTrue())
			Expect(result.TemplateData.Database.UseManaged).To(BeTrue())
		})

		It("should handle MySQL configuration", func() {
			result := wiz.GetResult()
			result.TemplateData.Database = generated.DatabaseConfig{
				Engine:    generated.DatabaseTypeMySQL,
				UseJSON:   true,
				UseArrays: false, // MySQL arrays are different
				UseManaged: true,
			}

			Expect(result.TemplateData.Database.Engine).To(Equal(generated.DatabaseTypeMySQL))
			Expect(result.TemplateData.Database.Engine.IsValid()).To(BeTrue())
			Expect(result.TemplateData.Database.UseArrays).To(BeFalse())
		})

		It("should handle SQLite configuration", func() {
			result := wiz.GetResult()
			result.TemplateData.Database = generated.DatabaseConfig{
				Engine:    generated.DatabaseTypeSQLite,
				UseArrays: false, // SQLite arrays are limited
				UseManaged: false, // SQLite is typically not managed
			}

			Expect(result.TemplateData.Database.Engine).To(Equal(generated.DatabaseTypeSQLite))
			Expect(result.TemplateData.Database.Engine.IsValid()).To(BeTrue())
			Expect(result.TemplateData.Database.UseArrays).To(BeFalse())
			Expect(result.TemplateData.Database.UseManaged).To(BeFalse())
		})
	})

	Describe("Output Configuration", func() {
		It("should handle custom output directories", func() {
			result := wiz.GetResult()
			result.TemplateData.Output = generated.OutputConfig{
				BaseDir:    "./custom/db",
				QueriesDir: "./custom/queries",
				SchemaDir:  "./custom/schema",
			}

			Expect(result.TemplateData.Output.BaseDir).To(Equal("./custom/db"))
			Expect(result.TemplateData.Output.QueriesDir).To(Equal("./custom/queries"))
			Expect(result.TemplateData.Output.SchemaDir).To(Equal("./custom/schema"))
		})

		It("should handle relative and absolute paths", func() {
			result := wiz.GetResult()
			result.TemplateData.Output = generated.OutputConfig{
				BaseDir:    "/absolute/path/db",
				QueriesDir: "./relative/queries",
				SchemaDir:  "../schema",
			}

			Expect(result.TemplateData.Output.BaseDir).To(Equal("/absolute/path/db"))
			Expect(result.TemplateData.Output.QueriesDir).To(Equal("./relative/queries"))
			Expect(result.TemplateData.Output.SchemaDir).To(Equal("../schema"))
		})
	})

	Describe("Wizard Result Configuration", func() {
		It("should handle generate flags", func() {
			result := wiz.GetResult()
			
			// Test default values
			Expect(result.GenerateQueries).To(BeTrue())
			Expect(result.GenerateSchema).To(BeTrue())
			
			// Test setting values
			result.GenerateQueries = false
			result.GenerateSchema = false
			
			Expect(result.GenerateQueries).To(BeFalse())
			Expect(result.GenerateSchema).To(BeFalse())
		})

		It("should handle template data updates", func() {
			result := wiz.GetResult()
			
			// Set initial template data
			result.TemplateData = generated.TemplateData{
				ProjectName: "initial-project",
				ProjectType: generated.ProjectTypeHobby,
			}
			
			Expect(result.TemplateData.ProjectName).To(Equal("initial-project"))
			
			// Update template data
			result.TemplateData.ProjectName = "updated-project"
			result.TemplateData.ProjectType = generated.ProjectTypeMicroservice
			
			Expect(result.TemplateData.ProjectName).To(Equal("updated-project"))
			Expect(result.TemplateData.ProjectType).To(Equal(generated.ProjectTypeMicroservice))
		})
	})
})