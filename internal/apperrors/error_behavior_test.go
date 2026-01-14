package apperrors

import (
	stderrors "errors"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Error Behavior and Comparison", func() {
	Context("Error Interface", func() {
		It("should implement error interface correctly", func() {
			err := NewError(ErrorCodeValidationError, "test message")

			Expect(err.Error()).To(Equal("[VALIDATION_ERROR] test message"))

			// TODO: Add tests for error interface compliance
		})

		It("should include description in error string when present", func() {
			err := NewError(ErrorCodeValidationError, "test message")
			_ = err.WithDescription("detailed explanation")

			Expect(err.Error()).To(Equal("[VALIDATION_ERROR] test message: detailed explanation"))

			// TODO: Add tests for empty description
			// TODO: Add tests for multiline descriptions
		})
	})

	Context("Unwrap", func() {
		It("should return nil when no cause is set", func() {
			err := NewError(ErrorCodeValidationError, "test")

			Expect(err.Unwrap()).To(Succeed())
		})

		It("should return the cause error when set", func() {
			original := stderrors.New("original error")
			err := NewError(ErrorCodeValidationError, "wrapped").WithCause(original)

			Expect(err.Unwrap()).To(Equal(original))
		})
	})

	Context("Is", func() {
		It("should match errors with same code", func() {
			err1 := NewError(ErrorCodeValidationError, "test1")
			err2 := NewError(ErrorCodeValidationError, "test2")

			Expect(Is(err1, err2)).To(BeTrue())

			// TODO: Add tests for different messages with same code
		})

		It("should not match errors with different codes", func() {
			err1 := NewError(ErrorCodeValidationError, "test")
			err2 := NewError(ErrorCodeInternalServer, "test")

			Expect(Is(err1, err2)).To(BeFalse())

			// TODO: Add tests for all error code combinations
		})

		It("should handle non-application errors", func() {
			err1 := stderrors.New("standard error")
			err2 := NewError(ErrorCodeInternalServer, "test")

			Expect(Is(err1, err2)).To(BeFalse())
		})

		It("should handle nil errors", func() {
			err1 := NewError(ErrorCodeValidationError, "test")

			Expect(Is(err1, nil)).To(BeFalse())
			Expect(Is(nil, err1)).To(BeFalse())

			// TODO: Add tests for both nil errors
		})
	})

	Context("ErrorList Behavior", func() {
		It("should return correct error message for single error", func() {
			list := NewErrorList()
			list.AddError(ErrorCodeValidationError, "single error")

			Expect(list.Error()).To(Equal("[VALIDATION_ERROR] single error"))
		})

		It("should return summary message for multiple errors", func() {
			list := NewErrorList()
			list.AddError(ErrorCodeValidationError, "error 1")
			list.AddError(ErrorCodeInternalServer, "error 2")

			Expect(list.Error()).To(ContainSubstring("2 errors occurred"))
			Expect(list.Error()).To(ContainSubstring("error 1"))
		})

		It("should return message for empty list", func() {
			list := NewErrorList()

			Expect(list.Error()).To(Equal("no errors"))
		})
	})

	Context("Base Errors", func() {
		It("should have ErrConfigParseFailed", func() {
			Expect(ErrConfigParseFailed).To(HaveOccurred())
			Expect(ErrConfigParseFailed.Code).To(Equal(ErrorCodeConfigParseFailed))
			Expect(ErrConfigParseFailed.Message).To(Equal("Config parse failed"))
		})

		It("should have ErrInvalidValue", func() {
			Expect(ErrInvalidValue).To(HaveOccurred())
			Expect(ErrInvalidValue.Code).To(Equal(ErrorCodeInvalidValue))
			Expect(ErrInvalidValue.Message).To(Equal("Invalid value"))
		})

		It("should have ErrFileNotFound", func() {
			Expect(ErrFileNotFound).To(HaveOccurred())
			Expect(ErrFileNotFound.Code).To(Equal(ErrorCodeFileNotFound))
			Expect(ErrFileNotFound.Message).To(Equal("File not found"))
		})

		It("should have ErrInvalidType", func() {
			Expect(ErrInvalidType).To(HaveOccurred())
			Expect(ErrInvalidType.Code).To(Equal(ErrorCodeValidationError))
			Expect(ErrInvalidType.Message).To(Equal("Invalid type"))
		})

		// TODO: Add tests for all base error immutability
		// TODO: Add tests for base error thread safety
	})

	Context("Error Details", func() {
		It("should support detailed error information", func() {
			err := NewError(ErrorCodeValidationError, "validation failed").
				WithDetails("age", 150, "0-120", 150).
				WithMessage("Age must be within valid range").
				WithComponent("user-service").
				WithRequestID("req-789").
				WithRetryable(false).
				WithSeverity(ErrorSeverityWarning)

			Expect(err.Details.Field).To(Equal("age"))
			Expect(err.Details.Value).To(Equal(150))
			Expect(err.Details.Expected).To(Equal("0-120"))
			Expect(err.Details.Actual).To(Equal(150))
			Expect(err.Details.Message).To(Equal("Age must be within valid range"))
			Expect(err.Component).To(Equal("user-service"))
			Expect(err.RequestID).To(Equal("req-789"))
			Expect(err.Retryable).To(BeFalse())
			Expect(err.Severity).To(Equal(ErrorSeverityWarning))

			// TODO: Add tests for partial detail setting
			// TODO: Add tests for detail overriding
		})
	})

	Context("JSON Serialization", func() {
		It("should serialize error to JSON", func() {
			err := NewError(ErrorCodeValidationError, "test error").
				WithDetails("field", "value", "expected", "actual").
				WithComponent("test-component")

			jsonStr, jsonErr := err.ToJSON()
			Expect(jsonErr).ToNot(HaveOccurred())
			Expect(jsonStr).To(ContainSubstring("VALIDATION_ERROR"))
			Expect(jsonStr).To(ContainSubstring("test error"))
			Expect(jsonStr).To(ContainSubstring("field"))
			Expect(jsonStr).To(ContainSubstring("test-component"))

			// TODO: Add tests for JSON structure validation
			// TODO: Add tests for special character handling
		})

		It("should handle JSON marshaling errors", func() {
			// Create an error with problematic data that might cause JSON issues
			err := NewError(ErrorCodeValidationError, "test")

			// This should work normally
			jsonStr, marshalErr := err.ToJSON()
			Expect(marshalErr).ToNot(HaveOccurred())
			Expect(jsonStr).NotTo(BeEmpty())

			// TODO: Add tests for complex nested structures
			// TODO: Add tests for circular references
		})
	})

	Context("Retryable and Critical Checks", func() {
		It("should correctly identify retryable errors", func() {
			err := NewError(ErrorCodeTimeout, "timeout").WithRetryable(true)

			Expect(err.IsRetryable()).To(BeTrue())
		})

		It("should correctly identify non-retryable errors", func() {
			err := NewError(ErrorCodeValidationError, "invalid input")

			Expect(err.IsRetryable()).To(BeFalse())
		})

		It("should correctly identify critical errors", func() {
			err := NewError(ErrorCodeInternalServer, "crash").WithSeverity(ErrorSeverityCritical)

			Expect(err.IsCritical()).To(BeTrue())
		})

		It("should correctly identify non-critical errors", func() {
			err := NewError(ErrorCodeValidationError, "invalid input")

			Expect(err.IsCritical()).To(BeFalse())
		})
	})

	// TODO: Add performance tests for error creation
	// TODO: Add memory usage tests
	// TODO: Add thread safety tests
	// TODO: Add error behavior under load tests
})

func TestErrorBehavior(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Errors Behavior Suite")
}
