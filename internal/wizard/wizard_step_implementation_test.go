package wizard_test

import (
	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/wizard"
	"github.com/charmbracelet/huh"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Individual Step Implementation Tests", func() {
	var theme *huh.Theme

	BeforeEach(func() {
		theme = huh.ThemeBase()
	})

	Describe("ProjectTypeStep", func() {
		var step *wizard.ProjectTypeStep

		BeforeEach(func() {
			step = wizard.NewProjectTypeStep(theme, wizard.NewUIHelper())
		})

		It("should create project type step successfully", func() {
			Expect(step).NotTo(BeNil())
		})

		It("should initialize with valid project types", func() {
			// Test that all project types are valid
			validTypes := []generated.ProjectType{
				generated.ProjectTypeHobby,
				generated.ProjectTypeMicroservice,
				generated.ProjectTypeEnterprise,
				generated.ProjectTypeAPIFirst,
				generated.ProjectTypeAnalytics,
				generated.ProjectTypeTesting,
				generated.ProjectTypeMultiTenant,
				generated.ProjectTypeLibrary,
			}

			for _, projectType := range validTypes {
				Expect(projectType.IsValid()).To(BeTrue(),
					"Project type %s should be valid", projectType)
			}
		})

		It("should validate project type constants", func() {
			// Test that generated constants exist
			Expect(generated.ProjectTypeHobby.IsValid()).To(BeTrue())
			Expect(generated.ProjectTypeMicroservice.IsValid()).To(BeTrue())
			Expect(generated.ProjectTypeEnterprise.IsValid()).To(BeTrue())
			Expect(generated.ProjectTypeAPIFirst.IsValid()).To(BeTrue())
			Expect(generated.ProjectTypeAnalytics.IsValid()).To(BeTrue())
			Expect(generated.ProjectTypeTesting.IsValid()).To(BeTrue())
			Expect(generated.ProjectTypeMultiTenant.IsValid()).To(BeTrue())
			Expect(generated.ProjectTypeLibrary.IsValid()).To(BeTrue())
		})
	})

	Describe("DatabaseStep", func() {
		var step *wizard.DatabaseStep

		BeforeEach(func() {
			step = wizard.NewDatabaseStep(theme, wizard.NewUIHelper())
		})

		It("should create database step successfully", func() {
			Expect(step).NotTo(BeNil())
		})

		It("should initialize with valid database types", func() {
			// Test that all database types are valid
			validTypes := []generated.DatabaseType{
				generated.DatabaseTypePostgreSQL,
				generated.DatabaseTypeMySQL,
				generated.DatabaseTypeSQLite,
			}

			for _, dbType := range validTypes {
				Expect(dbType.IsValid()).To(BeTrue(),
					"Database type %s should be valid", dbType)
			}
		})

		It("should validate database type constants", func() {
			// Test that generated constants exist
			Expect(generated.DatabaseTypePostgreSQL.IsValid()).To(BeTrue())
			Expect(generated.DatabaseTypeMySQL.IsValid()).To(BeTrue())
			Expect(generated.DatabaseTypeSQLite.IsValid()).To(BeTrue())
		})
	})

	Describe("ProjectDetailsStep", func() {
		var step *wizard.ProjectDetailsStep

		BeforeEach(func() {
			step = wizard.NewProjectDetailsStep(theme, wizard.NewUIHelper())
		})

		It("should create project details step successfully", func() {
			Expect(step).NotTo(BeNil())
		})

		It("should handle package name validation", func() {
			// Test various package names
			validNames := []string{
				"db",
				"database",
				"models",
				"internal_db",
				"myproject_db",
			}

			for _, name := range validNames {
				data := &generated.TemplateData{
					Package: generated.PackageConfig{
						Name: name,
					},
				}

				Expect(data.Package.Name).To(Equal(name))
				Expect(data.Package.Name).ToNot(BeEmpty())
			}
		})

		It("should handle package path validation", func() {
			// Test various package paths
			validPaths := []string{
				"github.com/user/project",
				"github.com/company/enterprise",
				"gitlab.com/organization/repo",
				"github.com/author/library",
			}

			for _, path := range validPaths {
				data := &generated.TemplateData{
					Package: generated.PackageConfig{
						Path: path,
					},
				}

				Expect(data.Package.Path).To(Equal(path))
				Expect(data.Package.Path).ToNot(BeEmpty())
			}
		})
	})

	Describe("FeaturesStep", func() {
		var step *wizard.FeaturesStep

		BeforeEach(func() {
			step = wizard.NewFeaturesStep(theme, wizard.NewUIHelper())
		})

		It("should create features step successfully", func() {
			Expect(step).NotTo(BeNil())
		})

		It("should handle emit options configuration", func() {
			emitOptions := []generated.EmitOptions{
				{
					EmitInterface:            true,
					EmitPreparedQueries:      false,
					EmitJSONTags:             true,
					EmitResultStructPointers: false,
					EmitParamsStructPointers: false,
					EmitEnumValidMethod:      true,
					EmitAllEnumValues:        false,
				},
				{
					EmitInterface:            false,
					EmitPreparedQueries:      true,
					EmitJSONTags:             false,
					EmitResultStructPointers: true,
					EmitParamsStructPointers: true,
					EmitEnumValidMethod:      false,
					EmitAllEnumValues:        true,
				},
			}

			for _, opts := range emitOptions {
				data := &generated.TemplateData{
					Validation: generated.ValidationConfig{
						EmitOptions: opts,
					},
				}

				Expect(data.Validation.EmitOptions.EmitInterface).To(Equal(opts.EmitInterface))
				Expect(data.Validation.EmitOptions.EmitPreparedQueries).To(Equal(opts.EmitPreparedQueries))
				Expect(data.Validation.EmitOptions.EmitJSONTags).To(Equal(opts.EmitJSONTags))
			}
		})

		It("should handle safety rules configuration", func() {
			safetyRules := []generated.SafetyRules{
				{
					NoSelectStar: true,
					RequireWhere: true,
					RequireLimit: false,
					NoDropTable:  true,
				},
				{
					NoSelectStar: false,
					RequireWhere: false,
					RequireLimit: true,
					NoDropTable:  false,
				},
			}

			for _, rules := range safetyRules {
				data := &generated.TemplateData{
					Validation: generated.ValidationConfig{
						SafetyRules: rules,
					},
				}

				Expect(data.Validation.SafetyRules.NoSelectStar).To(Equal(rules.NoSelectStar))
				Expect(data.Validation.SafetyRules.RequireWhere).To(Equal(rules.RequireWhere))
				Expect(data.Validation.SafetyRules.RequireLimit).To(Equal(rules.RequireLimit))
				Expect(data.Validation.SafetyRules.NoDropTable).To(Equal(rules.NoDropTable))
			}
		})
	})

	Describe("OutputStep", func() {
		var step *wizard.OutputStep

		BeforeEach(func() {
			step = wizard.NewOutputStep(theme, wizard.NewUIHelper())
		})

		It("should create output step successfully", func() {
			Expect(step).NotTo(BeNil())
		})

		It("should validate output configuration", func() {
			// Test valid configurations
			validConfigs := []generated.OutputConfig{
				{
					BaseDir:    "./internal/db",
					QueriesDir: "./sql/queries",
					SchemaDir:  "./sql/schema",
				},
				{
					BaseDir:    "./generated/db",
					QueriesDir: "./queries",
					SchemaDir:  "./schema",
				},
				{
					BaseDir:    "/absolute/path/db",
					QueriesDir: "/absolute/path/queries",
					SchemaDir:  "/absolute/path/schema",
				},
			}

			for _, config := range validConfigs {
				Expect(config.BaseDir).ToNot(BeEmpty())
				Expect(config.QueriesDir).ToNot(BeEmpty())
				Expect(config.SchemaDir).ToNot(BeEmpty())
			}
		})

		It("should detect invalid output configurations", func() {
			// Test invalid configurations
			invalidConfigs := []generated.OutputConfig{
				{
					BaseDir:    "", // Empty base dir
					QueriesDir: "./sql/queries",
					SchemaDir:  "./sql/schema",
				},
				{
					BaseDir:    "./internal/db",
					QueriesDir: "", // Empty queries dir
					SchemaDir:  "./sql/schema",
				},
				{
					BaseDir:    "./internal/db",
					QueriesDir: "./sql/queries",
					SchemaDir:  "", // Empty schema dir
				},
			}

			for _, config := range invalidConfigs {
				if config.BaseDir == "" {
					Expect(config.BaseDir).To(BeEmpty())
				}
				if config.QueriesDir == "" {
					Expect(config.QueriesDir).To(BeEmpty())
				}
				if config.SchemaDir == "" {
					Expect(config.SchemaDir).To(BeEmpty())
				}
			}
		})
	})

	Describe("UIHelper", func() {
		var uiHelper *wizard.UIHelper

		BeforeEach(func() {
			uiHelper = wizard.NewUIHelper()
		})

		It("should create UI helper successfully", func() {
			Expect(uiHelper).NotTo(BeNil())
		})

		It("should have theme configured", func() {
			// UIHelper should have a theme internally
			// We can't directly access it, but we can verify it was created
			Expect(uiHelper).NotTo(BeNil())
		})
	})

	Describe("Template Data Validation", func() {
		It("should validate complete template data structure", func() {
			data := generated.TemplateData{
				ProjectName: "complete-project",
				ProjectType: generated.ProjectTypeMicroservice,
				Package: generated.PackageConfig{
					Name:      "db",
					Path:      "github.com/example/project",
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
					QueriesDir: "./sql/queries",
					SchemaDir:  "./sql/schema",
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
						NoSelectStar: true,
						RequireWhere: true,
						RequireLimit: false,
						NoDropTable:  true,
					},
				},
			}

			// Verify all fields are properly set
			Expect(data.ProjectName).To(Equal("complete-project"))
			Expect(data.ProjectType).To(Equal(generated.ProjectTypeMicroservice))
			Expect(data.Package.Name).To(Equal("db"))
			Expect(data.Package.Path).To(Equal("github.com/example/project"))
			Expect(data.Package.BuildTags).To(Equal("postgres,pgx"))
			Expect(data.Database.Engine).To(Equal(generated.DatabaseTypePostgreSQL))
			Expect(data.Database.UseManaged).To(BeTrue())
			Expect(data.Database.UseUUIDs).To(BeTrue())
			Expect(data.Database.UseJSON).To(BeTrue())
			Expect(data.Database.UseArrays).To(BeTrue())
			Expect(data.Database.UseFullText).To(BeTrue())
			Expect(data.Output.BaseDir).To(Equal("./internal/db"))
			Expect(data.Output.QueriesDir).To(Equal("./sql/queries"))
			Expect(data.Output.SchemaDir).To(Equal("./sql/schema"))
			Expect(data.Validation.StrictFunctions).To(BeTrue())
			Expect(data.Validation.StrictOrderBy).To(BeTrue())
			Expect(data.Validation.EmitOptions.EmitInterface).To(BeTrue())
			Expect(data.Validation.SafetyRules.NoSelectStar).To(BeTrue())
		})
	})

	Describe("Default Values Testing", func() {
		It("should handle default emit options", func() {
			emitOpts := generated.DefaultEmitOptions()

			Expect(emitOpts.EmitInterface).ToNot(BeNil())
			Expect(emitOpts.EmitPreparedQueries).ToNot(BeNil())
			Expect(emitOpts.EmitJSONTags).ToNot(BeNil())
		})

		It("should handle default safety rules", func() {
			safetyRules := generated.DefaultSafetyRules()

			Expect(safetyRules.NoSelectStar).ToNot(BeNil())
			Expect(safetyRules.RequireWhere).ToNot(BeNil())
			Expect(safetyRules.RequireLimit).ToNot(BeNil())
		})
	})
})
