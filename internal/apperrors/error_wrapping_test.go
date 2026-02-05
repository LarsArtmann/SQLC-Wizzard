package apperrors_test

import (
	stderrors "errors"

	"github.com/LarsArtmann/SQLC-Wizzard/internal/apperrors"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Error Wrapping and Combining", func() {
	Context("Wrap", func() {
		It("should wrap error with code and component", func() {
			original := stderrors.New("database connection failed")
			wrapped := apperrors.Wrap(original, apperrors.ErrorCodeInternalServer, "database")

			Expect(wrapped).To(HaveOccurred())
			Expect(wrapped.Code).To(Equal(apperrors.ErrorCodeInternalServer))
			Expect(wrapped.Component).To(Equal("database"))
			Expect(wrapped.Message).To(Equal("database connection failed"))
			Expect(stderrors.Unwrap(wrapped)).To(Equal(original))

			// TODO: Add validation for nil original error
			// TODO: Add tests for wrapping application errors
		})
	})

	Context("WrapWithRequestID and WrapWithUserID", func() {
		type wrappingTestCase struct {
			name        string
			errorMsg    string
			code        apperrors.ErrorCode
			id          string
			component   string
			expectID    string
			expectError string
		}

		DescribeTable("should wrap error with context",
			func(tc wrappingTestCase, wrapFunc func(error, apperrors.ErrorCode, string, string) *apperrors.Error) {
				original := stderrors.New(tc.errorMsg)
				wrapped := wrapFunc(original, tc.code, tc.id, tc.component)

				Expect(wrapped).To(HaveOccurred())
				Expect(stderrors.Unwrap(wrapped)).To(Equal(original))

				if tc.expectID == "requestID" {
					Expect(wrapped.RequestID).To(Equal(tc.id))
				} else {
					Expect(wrapped.UserID).To(Equal(tc.id))
				}
				Expect(wrapped.Component).To(Equal(tc.component))
			},
			Entry("WrapWithRequestID", wrappingTestCase{
				name:        "wrap with request ID",
				errorMsg:    "operation failed",
				code:        apperrors.ErrorCodeInternalServer,
				id:          "req-123",
				component:   "api",
				expectID:    "requestID",
				expectError: "operation failed",
			}, apperrors.WrapWithRequestID),

			Entry("WrapWithUserID", wrappingTestCase{
				name:        "wrap with user ID",
				errorMsg:    "permission denied",
				code:        apperrors.ErrorCodePermissionDenied,
				id:          "user-456",
				component:   "auth",
				expectID:    "userID",
				expectError: "permission denied",
			}, apperrors.WrapWithUserID),
		)
	})

	Context("Wrapf", func() {
		It("should wrap error with formatted message", func() {
			original := stderrors.New("field validation failed")
			baseErr := apperrors.NewError(apperrors.ErrorCodeValidationError, "validation error")
			wrapped := apperrors.Wrapf(original, baseErr, "user %s has invalid %s", "john", "email")

			Expect(wrapped).To(HaveOccurred())
			Expect(wrapped.Code).To(Equal(apperrors.ErrorCodeValidationError))
			Expect(wrapped.Message).To(Equal("user john has invalid email"))
			Expect(wrapped.Description).To(ContainSubstring("field validation failed"))
			Expect(stderrors.Unwrap(wrapped)).To(Equal(original))

			// TODO: Add tests for complex formatting scenarios
			// TODO: Add validation for format string errors
		})

		It("should handle nil original error", func() {
			baseErr := apperrors.NewError(apperrors.ErrorCodeValidationError, "base error")
			wrapped := apperrors.Wrapf(nil, baseErr, "formatted message")

			Expect(wrapped).To(HaveOccurred())
			Expect(wrapped.Message).To(Equal("Cannot wrap nil error"))
			Expect(wrapped.Code).To(Equal(apperrors.ErrorCodeInternalServer))

			// TODO: Add tests for nil base error
			// TODO: Add validation for both errors being nil
		})
	})

	Context("apperrors.Combine", func() {
		It("should combine multiple application errors", func() {
			err1 := apperrors.NewError(apperrors.ErrorCodeValidationError, "error 1")
			err2 := apperrors.NewError(apperrors.ErrorCodeValidationError, "error 2")
			err3 := apperrors.NewError(apperrors.ErrorCodeInternalServer, "error 3")

			list := apperrors.Combine(err1, err2, err3)

			Expect(list.GetCount()).To(Equal(3))
			Expect(list.Errors[0].Message).To(Equal("error 1"))
			Expect(list.Errors[1].Message).To(Equal("error 2"))
			Expect(list.Errors[2].Message).To(Equal("error 3"))

			// TODO: Add tests for duplicate errors
			// TODO: Add tests for mixed error types
		})

		It("should handle empty error list", func() {
			list := apperrors.Combine()

			Expect(list.GetCount()).To(Equal(0))
			Expect(list.HasErrors()).To(BeFalse())
		})
	})

	Context("CombineErrors", func() {
		It("should combine application errors", func() {
			err1 := apperrors.NewError(apperrors.ErrorCodeValidationError, "error 1")
			err2 := apperrors.NewError(apperrors.ErrorCodeValidationError, "error 2")

			list := apperrors.CombineErrors(err1, err2)

			Expect(list.GetCount()).To(Equal(2))
		})

		It("should wrap non-application errors", func() {
			err1 := stderrors.New("standard error")
			err2 := apperrors.NewError(apperrors.ErrorCodeValidationError, "app error")

			list := apperrors.CombineErrors(err1, err2)

			Expect(list.GetCount()).To(Equal(2))
			Expect(list.Errors[0].Code).To(Equal(apperrors.ErrorCodeInternalServer))
			Expect(list.Errors[0].Component).To(Equal("unknown"))
			Expect(list.Errors[0].Message).To(ContainSubstring("standard error"))
			Expect(list.Errors[1].Code).To(Equal(apperrors.ErrorCodeValidationError))

			// TODO: Add validation for wrapping behavior
			// TODO: Add tests for error preservation
		})

		It("should skip nil errors", func() {
			appErr := apperrors.NewError(apperrors.ErrorCodeValidationError, "app error")

			list := apperrors.CombineErrors(nil, appErr, nil)

			Expect(list.GetCount()).To(Equal(1))
			Expect(list.Errors[0].Code).To(Equal(apperrors.ErrorCodeValidationError))
			Expect(list.Errors[0].Message).To(Equal("app error"))
		})

		It("should handle all nil errors", func() {
			list := apperrors.CombineErrors(nil, nil, nil)

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
