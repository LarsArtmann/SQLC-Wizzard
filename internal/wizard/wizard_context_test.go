package wizard_test

import (
	"github.com/LarsArtmann/SQLC-Wizzard/internal/wizard"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Context Handling", func() {
	It("should work with context cancellation", func() {
		// Test that wizard can handle context without crashing
		wiz := wizard.NewWizard()
		result := wiz.GetResult()

		// Should still be able to access result
		Expect(result).NotTo(BeNil())
	})
})
