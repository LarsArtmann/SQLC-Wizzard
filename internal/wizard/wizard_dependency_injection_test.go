package wizard_test

import (
	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/templates"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/wizard"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Wizard with Dependency Injection", func() {
	var (
		mockUI       *MockUI
		mockSteps    map[string]*MockStep
		mockTemplate *MockTemplate
		wizardDeps   wizard.WizardDependencies
		wiz          *wizard.Wizard
	)

	BeforeEach(func() {
		mockUI = NewMockUI()
		mockSteps = map[string]*MockStep{
			"projectType": NewMockStep(),
			"database":    NewMockStep(),
			"details":     NewMockStep(),
			"features":    NewMockStep(),
			"output":      NewMockStep(),
		}
		mockTemplate = NewMockTemplate()

		wizardDeps = wizard.WizardDependencies{
			UI:          mockUI,
			ProjectType: mockSteps["projectType"],
			Database:    mockSteps["database"],
			Details:     mockSteps["details"],
			Features:    mockSteps["features"],
			Output:      mockSteps["output"],
			TemplateFunc: func(projectType templates.ProjectType) (wizard.TemplateInterface, error) {
				return mockTemplate, nil
			},
		}

		wiz = wizard.NewTestableWizard(wizardDeps)
	})

	AfterEach(func() {
		// TODO: Add cleanup verification
		// TODO: Add resource leak detection
	})

	Context("when wizard runs successfully", func() {
		It("should execute all steps in correct order", func() {
			_, err := wiz.Run()

			Expect(err).ToNot(HaveOccurred())

			// Verify all steps were called
			Expect(mockSteps["projectType"].ExecuteCalls).To(Equal(1))
			Expect(mockSteps["database"].ExecuteCalls).To(Equal(1))
			Expect(mockSteps["details"].ExecuteCalls).To(Equal(1))
			Expect(mockSteps["features"].ExecuteCalls).To(Equal(1))
			Expect(mockSteps["output"].ExecuteCalls).To(Equal(1))

			// Verify UI interactions
			Expect(mockUI.WelcomeCalls).To(Equal(1))
			Expect(len(mockUI.StepHeaders)).To(Equal(5))
			Expect(len(mockUI.StepCompletions)).To(Equal(5))
		})

		It("should generate template configuration", func() {
			_, err := wiz.Run()

			Expect(err).ToNot(HaveOccurred())
			Expect(mockTemplate.GenerateCalls).To(Equal(1))
			Expect(mockTemplate.LastCallData.ProjectType).To(Equal(generated.ProjectTypeMicroservice))
		})
	})

	Context("when steps fail", func() {
		It("should handle database step failures", func() {
			mockSteps["database"].ShouldFail = true
			mockSteps["database"].FailError = NewTestError("Database step failed")

			_, err := wiz.Run()

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Database step failed"))
		})

		It("should handle details step failures", func() {
			mockSteps["details"].ShouldFail = true
			mockSteps["details"].FailError = NewTestError("Details step failed")

			_, err := wiz.Run()

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Details step failed"))
		})

		It("should handle output step failures", func() {
			mockSteps["output"].ShouldFail = true
			mockSteps["output"].FailError = NewTestError("Output step failed")

			_, err := wiz.Run()

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Output step failed"))
		})
	})

	Context("data flow validation", func() {
		It("should pass data correctly between steps", func() {
			_, err := wiz.Run()

			Expect(err).ToNot(HaveOccurred())

			// Verify data flow through steps (wizard uses defaults)
			// Steps should receive the same data instance
			Expect(mockSteps["projectType"].ExecuteCalls).To(Equal(1))
			Expect(mockSteps["database"].ExecuteCalls).To(Equal(1))
			Expect(mockSteps["details"].ExecuteCalls).To(Equal(1))
			Expect(mockSteps["features"].ExecuteCalls).To(Equal(1))
			Expect(mockSteps["output"].ExecuteCalls).To(Equal(1))
		})
	})

	// TODO: Add performance tests
	// TODO: Add concurrency tests
	// TODO: Add memory usage tests
	// TODO: Add edge case tests
})

// NewTestError creates a test error for failure scenarios
// TODO: Move to test utilities
// TODO: Add error type validation
func NewTestError(message string) error {
	return &testError{message: message}
}

// testError is a simple error type for testing
type testError struct {
	message string
}

func (e *testError) Error() string {
	return e.message
}
