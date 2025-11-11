package wizard_test

import (
	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/wizard"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Wizard Configuration", func() {
	It("should accept custom configuration", func() {
		wiz := wizard.NewWizard()
		
		// Test setting custom configuration
		result := wiz.GetResult()
		result.GenerateQueries = false
		result.GenerateSchema = false
		
		Expect(result.GenerateQueries).To(BeFalse())
		Expect(result.GenerateSchema).To(BeFalse())
	})

	It("should handle template data properly", func() {
		wiz := wizard.NewWizard()
		
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