package errors_test

import (
	stderrors "errors"
	"fmt"

	"github.com/LarsArtmann/SQLC-Wizzard/internal/errors"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Error Wrapping and Combining", func() {
	Context("Wrap", func() {
		It("should wrap error with code and component", func() {
			original := errors.New("database connection failed")
			wrapped := errors.Wrap(original, errors.ErrorCodeInternalServer, "database")

			Expect(wrapped).To(HaveOccurred())
			Expect(wrapped.Code).To(Equal(errors.ErrorCodeInternalServer))
			Expect(wrapped.Component).To(Equal("database"))
			Expect(wrapped.Message).To(Equal("database connection failed"))
			Expect(stderrors.Unwrap(wrapped)).To(Equal(original))

			// TODO: Add validation for nil original error
			// TODO: Add tests for wrapping application errors
		})
	})

	Context("WrapWithRequestID", func() {
		It("should wrap error with request ID", func() {
			original := errors.New("operation failed")
			wrapped := errors.WrapWithRequestID(original, errors.ErrorCodeInternalServer, "req-123", "api")

			Expect(wrapped).To(HaveOccurred())
			Expect(wrapped.RequestID).To(Equal("req-123"))
			Expect(wrapped.Component).To(Equal("api"))
			Expect(stderrors.Unwrap(wrapped)).To(Equal(original))

			// TODO: Add validation for empty request ID
			// TODO: Add tests for request ID format validation
		})
	})

	Context("WrapWithUserID", func() {
		It("should wrap error with user ID", func() {
			original := errors.New("permission denied")
			wrapped := errors.WrapWithUserID(original, errors.ErrorCodePermissionDenied, "user-456", "auth")

			Expect(wrapped).To(HaveOccurred())
			Expect(wrapped.UserID).To(Equal("user-456"))
			Expect(wrapped.Component).To(Equal("auth"))
			Expect(stderrors.Unwrap(wrapped)).To(Equal(original))

			// TODO: Add validation for empty user ID
			// TODO: Add tests for user ID format validation
		})
	})

	Context("Wrapf", func() {
		It("should wrap error with formatted message", func() {
			original := errors.New("field validation failed")
			baseErr := errors.NewError(errors.ErrorCodeValidationError, "validation error")
			wrapped := errors.Wrapf(original, baseErr, "user %s has invalid %s", "john", "email")

			Expect(wrapped).To(HaveOccurred())
			Expect(wrapped.Code).To(Equal(errors.ErrorCodeValidationError))
			Expect(wrapped.Message).To(Equal("user john has invalid email"))
			Expect(wrapped.Description).To(ContainSubstring("field validation failed"))
			Expect(stderrors.Unwrap(wrapped)).To(Equal(original))

			// TODO: Add tests for complex formatting scenarios
			// TODO: Add validation for format string errors
		})

		It("should handle nil original error", func() {
			baseErr := errors.NewError(errors.ErrorCodeValidationError, "base error")
			wrapped := errors.Wrapf(nil, baseErr, "formatted message")

			Expect(wrapped).To(HaveOccurred())
			Expect(wrapped.Message).To(Equal("Cannot wrap nil error"))
			Expect(wrapped.Code).To(Equal(errors.ErrorCodeInternalServer))

			// TODO: Add tests for nil base error
			// TODO: Add validation for both errors being nil
		})
	})

	Context("errors.Combine", func() {
		It("should combine multiple application errors", func() {
			err1 := errors.NewError(errors.ErrorCodeValidationError, "error 1")
			err2 := errors.NewError(errors.ErrorCodeValidationError, "error 2")
			err3 := errors.NewError(errors.ErrorCodeInternalServer, "error 3")

			list := errors.Combine(err1, err2, err3)

			Expect(list.GetCount()).To(Equal(3))
			Expect(list.Errors[0].Message).To(Equal("error 1"))
			Expect(list.Errors[1].Message).To(Equal("error 2"))
			Expect(list.Errors[2].Message).To(Equal("error 3"))

			// TODO: Add tests for duplicate errors
			// TODO: Add tests for mixed error types
		})

		It("should handle empty error list", func() {
			list := errors.Combine()

			Expect(list.GetCount()).To(Equal(0))
			Expect(list.HasErrors()).To(BeFalse())
		})
	})

	Context("CombineErrors", func() {
		It("should combine application errors", func() {
			err1 := errors.NewError(errors.ErrorCodeValidationError, "error 1")
			err2 := errors.NewError(errors.ErrorCodeValidationError, "error 2")

			list := errors.CombineErrors(err1, err2)

			Expect(list.GetCount()).To(Equal(2))
		})

		It("should wrap non-application errors", func() {
			err1 := errors.New("standard error")
			err2 := errors.NewError(errors.ErrorCodeValidationError, "app error")

			list := errors.CombineErrors(err1, err2)

			Expect(list.GetCount()).To(Equal(2))
			Expect(list.Errors[0].Code).To(Equal(errors.ErrorCodeInternalServer))
			Expect(list.Errors[0].Component).To(Equal("unknown"))
			Expect(list.Errors[0].Message).To(ContainSubstring("standard error"))
			Expect(list.Errors[1].Code).To(Equal(errors.ErrorCodeValidationError))

			// TODO: Add validation for wrapping behavior
			// TODO: Add tests for error preservation
		})

		It("should skip nil errors", func() {
			appErr := errors.NewError(errors.ErrorCodeValidationError, "app error")

			list := errors.CombineErrors(nil, appErr, nil)

			Expect(list.GetCount()).To(Equal(1))
			Expect(list.Errors[0].Code).To(Equal(errors.ErrorCodeValidationError))
			Expect(list.Errors[0].Message).To(Equal("app error"))
		})

		It("should handle all nil errors", func() {
			list := errors.CombineErrors(nil, nil, nil)

			Expect(list.GetCount()).To(Equal(0))
			Expect(list.HasErrors()).To(BeFalse())

			// TODO: Add tests for memory allocation behavior
			// TODO: Add performance tests for large error lists
		})
	})

	// TODO: Add tests for deep error wrapping chains
	// TODO: Add tests for circular error references
	// TODO: Add tests for error unwrapping limits
})
