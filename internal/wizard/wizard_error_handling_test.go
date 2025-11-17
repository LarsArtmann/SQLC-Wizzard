package wizard_test

import (
	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/wizard"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Error Handling", func() {
	It("should handle invalid project types gracefully", func() {
		wiz := wizard.NewWizard()

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
		wiz := wizard.NewWizard()

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
		result := wizard.NewWizard().GetResult()
		result.TemplateData = generated.TemplateData{}

		// Should not panic
		Expect(result.TemplateData.ProjectName).To(Equal(""))
		Expect(result.TemplateData.ProjectType.IsValid()).To(BeFalse())
		Expect(result.TemplateData.Database.Engine.IsValid()).To(BeFalse())
	})
})
