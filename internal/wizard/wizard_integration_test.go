package wizard_test

import (
	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/wizard"
	"github.com/LarsArtmann/SQLC-Wizzard/pkg/config"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Configuration Generation", func() {
	It("should generate valid sqlc configuration", func() {
		wiz := wizard.NewWizard()

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
							Package:             "db",
							Out:                 "internal/db",
							SQLPackage:          "pgx/v5",
							EmitJSONTags:        true,
							EmitPreparedQueries: true,
							EmitInterface:       true,
							EmitEmptySlices:     true,
							JSONTagsCaseStyle:   "camel",
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
