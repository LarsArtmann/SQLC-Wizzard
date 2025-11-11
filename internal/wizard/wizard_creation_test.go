package wizard_test

import (
	"github.com/LarsArtmann/SQLC-Wizzard/internal/wizard"
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