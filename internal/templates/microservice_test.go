package templates_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/LarsArtmann/SQLC-Wizzard/internal/domain"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/templates"
)

var _ = Describe("MicroserviceTemplate", func() {
	var tmpl *templates.MicroserviceTemplate

	BeforeEach(func() {
		tmpl = templates.NewMicroserviceTemplate()
	})

	Describe("NewMicroserviceTemplate", func() {
		It("should create a template instance", func() {
			Expect(tmpl).ToNot(BeNil())
		})
	})

	Describe("Name", func() {
		It("should return microservice", func() {
			Expect(tmpl.Name()).To(Equal("microservice"))
		})
	})

	Describe("Description", func() {
		It("should return a descriptive string", func() {
			desc := tmpl.Description()

			Expect(desc).ToNot(BeEmpty())
			Expect(desc).To(ContainSubstring("microservice"))
		})
	})

	Describe("DefaultData", func() {
		It("should return valid template data", func() {
			data := tmpl.DefaultData()

			Expect(data.ProjectType).To(Equal(templates.MustNewProjectType("microservice")))
			Expect(data.Database.Engine).To(Equal(templates.MustNewDatabaseType("postgresql")))
			Expect(data.Database.UseManaged).To(BeTrue())
			Expect(data.Package.Name).To(Equal("db"))
			Expect(data.Output.BaseDir).To(Equal("internal/db"))
		})

		It("should include default emit options", func() {
			data := tmpl.DefaultData()

			Expect(data.Validation.EmitOptions.EmitInterface).To(BeTrue())
			Expect(data.Validation.EmitOptions.PreparedQueries).To(BeTrue())
			Expect(data.Validation.EmitOptions.JSONTags).To(BeTrue())
		})

		It("should include default safety rules", func() {
			data := tmpl.DefaultData()

			Expect(data.Validation.SafetyRules.NoSelectStar).To(BeTrue())
			Expect(data.Validation.SafetyRules.RequireWhere).To(BeTrue())
			Expect(data.Validation.SafetyRules.NoDropTable).To(BeTrue())
		})
	})

	Describe("RequiredFeatures", func() {
		It("should return required feature list", func() {
			features := tmpl.RequiredFeatures()

			Expect(features).ToNot(BeEmpty())
			Expect(features).To(ContainElement("emit_interface"))
			Expect(features).To(ContainElement("prepared_queries"))
			Expect(features).To(ContainElement("json_tags"))
		})
	})

	Describe("Generate", func() {
		var templateData templates.TemplateData

		BeforeEach(func() {
			templateData = templates.TemplateData{
				ProjectName: "test-service",
				ProjectType: templates.MustNewProjectType("microservice"),
				
				Package: templates.PackageConfig{
					Name: "db",
				},
				
				Database: templates.DatabaseConfig{
					Engine:     templates.MustNewDatabaseType("postgresql"),
					UseManaged:  true,
					UseUUIDs:    true,
					UseJSON:     true,
					UseArrays:   false,
					UseFullText: false,
				},
				
				Output: templates.OutputConfig{
					BaseDir:    "internal/db",
					QueriesDir:  "internal/db/queries",
					SchemaDir:   "internal/db/schema",
				},
				
				Validation: templates.ValidationConfig{
					StrictFunctions: false,
					StrictOrderBy:   false,
					EmitOptions:    domain.DefaultEmitOptions(),
					SafetyRules:     domain.DefaultSafetyRules(),
				},
			}
		})

		Context("with PostgreSQL", func() {
			It("should generate valid config", func() {
				cfg, err := tmpl.Generate(templateData)

				Expect(err).ToNot(HaveOccurred())
				Expect(cfg).ToNot(BeNil())
				Expect(cfg.Version).To(Equal("2"))
				Expect(cfg.SQL).To(HaveLen(1))
			})

			It("should set correct SQL config", func() {
				cfg, err := tmpl.Generate(templateData)

				Expect(err).ToNot(HaveOccurred())
				sql := cfg.SQL[0]
				Expect(sql.Name).To(Equal("test-service"))
				Expect(sql.Engine).To(Equal("postgresql"))
				Expect(sql.Queries.Strings()).To(ContainElement("internal/db/queries"))
				Expect(sql.Schema.Strings()).To(ContainElement("internal/db/schema"))
			})

			It("should configure Go generation", func() {
				cfg, err := tmpl.Generate(templateData)

				Expect(err).ToNot(HaveOccurred())
				goCfg := cfg.SQL[0].Gen.Go
				Expect(goCfg).ToNot(BeNil())
				Expect(goCfg.Package).To(Equal("db"))
				Expect(goCfg.Out).To(Equal("internal/db"))
				Expect(goCfg.SQLPackage).To(Equal("pgx/v5"))
			})

			It("should include UUID type overrides", func() {
				cfg, err := tmpl.Generate(templateData)

				Expect(err).ToNot(HaveOccurred())
				goCfg := cfg.SQL[0].Gen.Go
				Expect(goCfg.Overrides).ToNot(BeEmpty())

				hasUUID := false
				for _, override := range goCfg.Overrides {
					if override.DBType == "uuid" {
						hasUUID = true
						Expect(override.GoType).To(Equal("UUID"))
						Expect(override.GoImportPath).To(Equal("github.com/google/uuid"))
					}
				}
				Expect(hasUUID).To(BeTrue())
			})

			It("should include JSONB type overrides", func() {
				cfg, err := tmpl.Generate(templateData)

				Expect(err).ToNot(HaveOccurred())
				goCfg := cfg.SQL[0].Gen.Go
				hasJSON := false
				for _, override := range goCfg.Overrides {
					if override.DBType == "jsonb" {
						hasJSON = true
						Expect(override.GoType).To(Equal("RawMessage"))
						Expect(override.GoImportPath).To(Equal("encoding/json"))
					}
				}
				Expect(hasJSON).To(BeTrue())
			})

			It("should apply emit options", func() {
				cfg, err := tmpl.Generate(templateData)

				Expect(err).ToNot(HaveOccurred())
				goCfg := cfg.SQL[0].Gen.Go
				Expect(goCfg.EmitInterface).To(BeTrue())
				Expect(goCfg.EmitPreparedQueries).To(BeTrue())
				Expect(goCfg.EmitJSONTags).To(BeTrue())
				Expect(goCfg.JSONTagsCaseStyle).To(Equal("camel"))
			})

			It("should include safety rules", func() {
				cfg, err := tmpl.Generate(templateData)

				Expect(err).ToNot(HaveOccurred())
				rules := cfg.SQL[0].Rules
				Expect(rules).ToNot(BeEmpty())

				ruleNames := make([]string, len(rules))
				for i, rule := range rules {
					ruleNames[i] = rule.Name
				}

				Expect(ruleNames).To(ContainElement("no-select-star"))
				Expect(ruleNames).To(ContainElement("require-where-delete"))
				Expect(ruleNames).To(ContainElement("no-drop-table"))
			})
		})

		Context("with MySQL", func() {
			BeforeEach(func() {
				templateData.Database.Engine = templates.MustNewDatabaseType("mysql")
			})

			It("should use correct SQL package", func() {
				cfg, err := tmpl.Generate(templateData)

				Expect(err).ToNot(HaveOccurred())
				goCfg := cfg.SQL[0].Gen.Go
				Expect(goCfg.SQLPackage).To(Equal("database/sql"))
				Expect(goCfg.BuildTags).To(Equal("mysql"))
			})

			It("should include JSON type override", func() {
				cfg, err := tmpl.Generate(templateData)

				Expect(err).ToNot(HaveOccurred())
				goCfg := cfg.SQL[0].Gen.Go
				hasJSON := false
				for _, override := range goCfg.Overrides {
					if override.DBType == "json" {
						hasJSON = true
						Expect(override.GoType).To(Equal("RawMessage"))
					}
				}
				Expect(hasJSON).To(BeTrue())
			})
		})

		Context("with SQLite", func() {
			BeforeEach(func() {
				templateData.Database.Engine = templates.MustNewDatabaseType("sqlite")
			})

			It("should use correct SQL package", func() {
				cfg, err := tmpl.Generate(templateData)

				Expect(err).ToNot(HaveOccurred())
				goCfg := cfg.SQL[0].Gen.Go
				Expect(goCfg.SQLPackage).To(Equal("database/sql"))
				Expect(goCfg.BuildTags).To(Equal("sqlite"))
			})
		})

		Context("with defaults", func() {
			It("should apply default values for empty fields", func() {
				minimalData := templates.TemplateData{
					ProjectType: templates.MustNewProjectType("microservice"),
					
					Database: templates.DatabaseConfig{
						Engine: templates.MustNewDatabaseType("postgresql"),
					},
					
					Validation: templates.ValidationConfig{
						EmitOptions: domain.DefaultEmitOptions(),
						SafetyRules:  domain.DefaultSafetyRules(),
					},
				}

				cfg, err := tmpl.Generate(minimalData)

				Expect(err).ToNot(HaveOccurred())
				Expect(cfg.SQL[0].Gen.Go.Package).To(Equal("db"))
				Expect(cfg.SQL[0].Gen.Go.Out).To(Equal("internal/db"))
			})
		})
	})
})
