package wizard_test

import (
	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/wizard"
	"github.com/LarsArtmann/SQLC-Wizzard/pkg/config"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Wizard Result", func() {
	var result *wizard.WizardResult

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
