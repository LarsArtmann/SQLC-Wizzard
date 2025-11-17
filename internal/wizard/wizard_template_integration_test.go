package wizard_test

import (
	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

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
					Engine:   generated.DatabaseTypePostgreSQL,
					UseUUIDs: true,
					UseJSON:  true,
				},
			}

			// Each template type should be valid
			Expect(templateType.IsValid()).To(BeTrue(), "for template: %s", templateType)
			Expect(templateData.ProjectName).To(Equal("test-project"))
		}
	})
})
