package wizard_test

import (
	"github.com/LarsArtmann/SQLC-Wizzard/generated"
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
			Template:    mockTemplate,
		}

		wiz = wizard.NewWizard(wizardDeps)
	})

	AfterEach(func() {
		// TODO: Add cleanup verification
		// TODO: Add resource leak detection
	})

	Context("when wizard runs successfully", func() {
		It("should execute all steps in correct order", func() {
			data := &generated.TemplateData{}

			err := wiz.Run(data)

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
			data := &generated.TemplateData{
				ProjectType: generated.ProjectTypeMicroservice,
			}

			err := wiz.Run(data)

			Expect(err).ToNot(HaveOccurred())
			Expect(mockTemplate.GenerateCalls).To(Equal(1))
			Expect(mockTemplate.LastCallData.ProjectType).To(Equal(generated.ProjectTypeMicroservice))
		})
	})

	Context("when steps fail", func() {
		It("should handle step execution failures", func() {
			mockSteps["database"].ShouldFail = true
			mockSteps["database"].FailError = NewTestError("Database step failed")

			data := &generated.TemplateData{}
			err := wiz.Run(data)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Database step failed"))
		})

		It("should handle validation failures", func() {
			mockSteps["details"].ShouldValidateFail = true
			mockSteps["details"].ValidateError = NewTestError("Validation failed")

			data := &generated.TemplateData{}
			err := wiz.Run(data)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Validation failed"))
		})
	})

	Context("when UI operations fail", func() {
		It("should handle UI failures gracefully", func() {
			mockUI.ShouldFailRun = true
			mockUI.FailErrorMessage = "UI initialization failed"

			data := &generated.TemplateData{}
			err := wiz.Run(data)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("UI initialization failed"))
		})
	})

	Context("data flow validation", func() {
		It("should pass data correctly between steps", func() {
			testData := &generated.TemplateData{
				ProjectType: generated.ProjectTypeAPIFirst,
			}

			err := wiz.Run(testData)

			Expect(err).ToNot(HaveOccurred())

			// Verify data flow through steps
			Expect(mockSteps["projectType"].LastCallData.ProjectType).To(Equal(generated.ProjectTypeAPIFirst))
			Expect(mockSteps["database"].LastCallData.ProjectType).To(Equal(generated.ProjectTypeAPIFirst))
			Expect(mockSteps["details"].LastCallData.ProjectType).To(Equal(generated.ProjectTypeAPIFirst))
			Expect(mockSteps["features"].LastCallData.ProjectType).To(Equal(generated.ProjectTypeAPIFirst))
			Expect(mockSteps["output"].LastCallData.ProjectType).To(Equal(generated.ProjectTypeAPIFirst))
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
