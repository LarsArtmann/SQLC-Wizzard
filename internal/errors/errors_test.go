package errors_test

import (
	stderrors "errors"
	"fmt"
	"testing"

	"github.com/LarsArtmann/SQLC-Wizzard/internal/errors"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestErrors(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Errors Suite")
}

// testErrorCreation runs generic tests for error creation functions
func testErrorCreation(err *errors.Error, expectedCode, expectedComponent, expectedMessage string) {
	Expect(err).NotTo(BeNil())
	Expect(err.Error()).To(Equal(expectedCode + ": " + expectedMessage))
	Expect(err.Code).To(Equal(errors.ErrorCode(expectedCode)))
	Expect(err.Component).To(Equal(expectedComponent))
	if expectedMessage != "" {
		Expect(err.Description).To(Equal(expectedMessage))
	}
}

// testErrorWithCause runs generic tests for error creation functions with cause
func testErrorWithCause(err error, expectedCode errors.ErrorCode, expectedComponent, expectedMessage string, causeExists bool) {
	Expect(err).NotTo(BeNil())
	Expect(err.Error()).To(Equal("[" + string(expectedCode) + "] " + expectedMessage))
	if appErr, ok := err.(*errors.Error); ok {
		Expect(appErr.Code).To(Equal(expectedCode))
		Expect(appErr.Component).To(Equal(expectedComponent))
	}
	if causeExists {
		Expect(stderrors.Unwrap(err).Error()).To(Equal(expectedMessage))
	}
}

var _ = Describe("Error Creation", func() {
	Context("NewError", func() {
		It("should create a new error with defaults", func() {
			err := errors.NewError(errors.ErrorCodeInternalServer, "test error")

			Expect(err).NotTo(BeNil())
			Expect(err.Code).To(Equal(errors.ErrorCodeInternalServer))
			Expect(err.Message).To(Equal("test error"))
			Expect(err.Component).To(Equal("application"))
			Expect(err.Retryable).To(BeFalse())
			Expect(err.Severity).To(Equal(errors.ErrorSeverityError))
			Expect(err.Timestamp).To(BeNumerically(">", 0))
		})
	})

	Context("Newf", func() {
		It("should create error with formatted message", func() {
			err := errors.Newf(errors.ErrorCodeValidationError, "field %s has value %d", "age", 150)

			Expect(err).NotTo(BeNil())
			Expect(err.Message).To(Equal("field age has value 150"))
		})
	})

	Context("Helper Constructors", func() {
		It("should create internal error with context", func() {
			cause := fmt.Errorf("database connection failed")
			err := errors.NewInternal("database", "connect", cause)

			Expect(err).NotTo(BeNil())
			Expect(err.Code).To(Equal(errors.ErrorCodeInternalServer))
			Expect(err.Component).To(Equal("database"))
			Expect(err.Message).To(ContainSubstring("database"))
			Expect(err.Message).To(ContainSubstring("connect"))
			Expect(err.Description).To(Equal("database connection failed"))
		})

		It("should create not found error", func() {
			err := errors.NewNotFound("User", "12345")

			Expect(err).NotTo(BeNil())
			Expect(err.Code).To(Equal(errors.ErrorCodeNotFound))
			Expect(err.Message).To(ContainSubstring("User"))
			Expect(err.Message).To(ContainSubstring("12345"))
		})

		It("should create permission denied error", func() {
			err := errors.NewPermissionDenied("database", "write")

			Expect(err).NotTo(BeNil())
			Expect(err.Code).To(Equal(errors.ErrorCodePermissionDenied))
			Expect(err.Message).To(ContainSubstring("write"))
			Expect(err.Message).To(ContainSubstring("database"))
		})

		It("should create timeout error", func() {
			err := errors.NewTimeout("fetchData", 5000)

			Expect(err).NotTo(BeNil())
			Expect(err.Code).To(Equal(errors.ErrorCodeTimeout))
			Expect(err.Message).To(ContainSubstring("fetchData"))
			Expect(err.Message).To(ContainSubstring("5000"))
			Expect(err.Retryable).To(BeFalse())
		})

		It("should create file not found error", func() {
			err := errors.FileNotFoundError("/path/to/file.txt")

			Expect(err).NotTo(BeNil())
			Expect(err.Code).To(Equal(errors.ErrorCodeFileNotFound))
			Expect(err.Component).To(Equal("filesystem"))
			Expect(err.Message).To(ContainSubstring("/path/to/file.txt"))
		})

		It("should create file read error", func() {
			cause := fmt.Errorf("permission denied")
			fileErr := errors.FileReadError("/path/to/file.txt", cause)
			
			// Check main error format
			Expect(fileErr).NotTo(BeNil())
			Expect(fileErr.Error()).To(Equal("[FILE_READ_ERROR] Failed to read file: /path/to/file.txt"))
			
			// Check that error has cause
			Expect(stderrors.Unwrap(fileErr)).NotTo(BeNil())
			Expect(stderrors.Unwrap(fileErr).Error()).To(Equal("permission denied"))
		})

		It("should create config parse error", func() {
			cause := fmt.Errorf("invalid YAML")
			err := errors.ConfigParseError("/config/app.yaml", cause)

			Expect(err).NotTo(BeNil())
			Expect(err.Code).To(Equal(errors.ErrorCodeConfigParseFailed))
			Expect(err.Component).To(Equal("config"))
			Expect(err.Description).To(Equal("invalid YAML"))
		})

		It("should create template not found error", func() {
			err := errors.TemplateNotFoundError("microservice")

			Expect(err).NotTo(BeNil())
			Expect(err.Code).To(Equal(errors.ErrorCodeTemplateNotFound))
			Expect(err.Component).To(Equal("templates"))
			Expect(err.Message).To(ContainSubstring("microservice"))
		})

		It("should create validation error", func() {
			err := errors.ValidationError("email", "invalid format")

			Expect(err).NotTo(BeNil())
			Expect(err.Code).To(Equal(errors.ErrorCodeValidationError))
			Expect(err.Component).To(Equal("validation"))
			Expect(err.Message).To(ContainSubstring("email"))
			Expect(err.Message).To(ContainSubstring("invalid format"))
		})
	})
})

var _ = Describe("Error Builder Pattern", func() {
	Context("WithDetails", func() {
		It("should add typed details to error", func() {
			err := errors.NewError(errors.ErrorCodeValidationError, "invalid input").
				WithDetails("age", 150, "1-120", 150)

			Expect(err.Details).NotTo(BeNil())
			Expect(err.Details.Field).To(Equal("age"))
			Expect(err.Details.Value).To(Equal(150))
			Expect(err.Details.Expected).To(Equal("1-120"))
			Expect(err.Details.Actual).To(Equal(150))
		})
	})

	Context("WithMessage", func() {
		It("should create details when nil and set message", func() {
			err := errors.NewError(errors.ErrorCodeValidationError, "test")

			// Initially details should be nil
			Expect(err.Details).To(BeNil())

			// WithMessage should create details and set message
			err = err.WithMessage("test message")

			Expect(err.Details).NotTo(BeNil())
			Expect(err.Details.Message).To(Equal("test message"))
		})
	})

	Context("WithComponent", func() {
		It("should set component", func() {
			err := errors.NewError(errors.ErrorCodeInternalServer, "test").
				WithComponent("database")

			Expect(err.Component).To(Equal("database"))
		})
	})

	Context("WithRequestID", func() {
		It("should set request ID", func() {
			err := errors.NewError(errors.ErrorCodeInternalServer, "test").
				WithRequestID("req-12345")

			Expect(err.RequestID).To(Equal("req-12345"))
		})
	})

	Context("WithUserID", func() {
		It("should set user ID", func() {
			err := errors.NewError(errors.ErrorCodeInternalServer, "test").
				WithUserID("user-67890")

			Expect(err.UserID).To(Equal("user-67890"))
		})
	})

	Context("WithRetryable", func() {
		It("should set retryable flag", func() {
			err := errors.NewError(errors.ErrorCodeInternalServer, "test").
				WithRetryable(true)

			Expect(err.Retryable).To(BeTrue())
		})
	})

	Context("WithSeverity", func() {
		It("should set severity level", func() {
			err := errors.NewError(errors.ErrorCodeInternalServer, "test").
				WithSeverity(errors.ErrorSeverityCritical)

			Expect(err.Severity).To(Equal(errors.ErrorSeverityCritical))
		})
	})

	Context("WithDescription", func() {
		It("should add detailed description", func() {
			err := errors.NewError(errors.ErrorCodeInternalServer, "test").
				WithDescription("This is a detailed explanation")

			Expect(err.Description).To(Equal("This is a detailed explanation"))
		})
	})

	Context("Chaining", func() {
		It("should allow method chaining", func() {
			err := errors.NewError(errors.ErrorCodeInternalServer, "test").
				WithComponent("api").
				WithRequestID("req-123").
				WithUserID("user-456").
				WithRetryable(true).
				WithSeverity(errors.ErrorSeverityWarning).
				WithDescription("detailed info")

			Expect(err.Component).To(Equal("api"))
			Expect(err.RequestID).To(Equal("req-123"))
			Expect(err.UserID).To(Equal("user-456"))
			Expect(err.Retryable).To(BeTrue())
			Expect(err.Severity).To(Equal(errors.ErrorSeverityWarning))
			Expect(err.Description).To(Equal("detailed info"))
		})
	})
})

var _ = Describe("Error Interface", func() {
	Context("Error()", func() {
		It("should return formatted message without description", func() {
			err := errors.NewError(errors.ErrorCodeValidationError, "invalid input")

			Expect(err.Error()).To(Equal("[VALIDATION_ERROR] invalid input"))
		})

		It("should include description when present", func() {
			err := errors.NewError(errors.ErrorCodeValidationError, "invalid input").
				WithDescription("age must be positive")

			Expect(err.Error()).To(Equal("[VALIDATION_ERROR] invalid input: age must be positive"))
		})
	})

	Context("IsRetryable()", func() {
		It("should return retryable status", func() {
			err := errors.NewError(errors.ErrorCodeInternalServer, "test")
			Expect(err.IsRetryable()).To(BeFalse())

			err = err.WithRetryable(true)
			Expect(err.IsRetryable()).To(BeTrue())
		})
	})

	Context("IsCritical()", func() {
		It("should return false for non-critical errors", func() {
			err := errors.NewError(errors.ErrorCodeInternalServer, "test")
			Expect(err.IsCritical()).To(BeFalse())
		})

		It("should return true for critical severity", func() {
			err := errors.NewError(errors.ErrorCodeInternalServer, "test").
				WithSeverity(errors.ErrorSeverityCritical)
			Expect(err.IsCritical()).To(BeTrue())
		})
	})

	Context("ToJSON()", func() {
		It("should marshal error to JSON", func() {
			err := errors.NewError(errors.ErrorCodeValidationError, "test error")

			jsonStr, jsonErr := err.ToJSON()

			Expect(jsonErr).NotTo(HaveOccurred())
			Expect(jsonStr).To(ContainSubstring("VALIDATION_ERROR"))
			Expect(jsonStr).To(ContainSubstring("test error"))
		})
	})
})

var _ = Describe("ErrorList", func() {
	Context("NewErrorList", func() {
		It("should create empty error list", func() {
			list := errors.NewErrorList()

			Expect(list).NotTo(BeNil())
			Expect(list.Errors).NotTo(BeNil())
			Expect(list.Errors).To(BeEmpty())
		})
	})

	Context("Add", func() {
		It("should add error to list", func() {
			list := errors.NewErrorList()
			err := errors.NewError(errors.ErrorCodeValidationError, "test")

			list.Add(err)

			Expect(list.Errors).To(HaveLen(1))
			Expect(list.Errors[0]).To(Equal(err))
		})

		It("should ignore nil errors", func() {
			list := errors.NewErrorList()

			list.Add(nil)

			Expect(list.Errors).To(BeEmpty())
		})
	})

	Context("AddError", func() {
		It("should create and add error", func() {
			list := errors.NewErrorList()

			list.AddError(errors.ErrorCodeValidationError, "test error")

			Expect(list.Errors).To(HaveLen(1))
			Expect(list.Errors[0].Code).To(Equal(errors.ErrorCodeValidationError))
			Expect(list.Errors[0].Message).To(Equal("test error"))
		})
	})

	Context("HasErrors", func() {
		It("should return false for empty list", func() {
			list := errors.NewErrorList()

			Expect(list.HasErrors()).To(BeFalse())
		})

		It("should return true when errors exist", func() {
			list := errors.NewErrorList()
			list.AddError(errors.ErrorCodeValidationError, "test")

			Expect(list.HasErrors()).To(BeTrue())
		})
	})

	Context("GetCount", func() {
		It("should return 0 for empty list", func() {
			list := errors.NewErrorList()

			Expect(list.GetCount()).To(Equal(0))
		})

		It("should return correct count", func() {
			list := errors.NewErrorList()
			list.AddError(errors.ErrorCodeValidationError, "test1")
			list.AddError(errors.ErrorCodeValidationError, "test2")
			list.AddError(errors.ErrorCodeValidationError, "test3")

			Expect(list.GetCount()).To(Equal(3))
		})
	})

	Context("Error()", func() {
		It("should return 'no errors' for empty list", func() {
			list := errors.NewErrorList()

			Expect(list.Error()).To(Equal("no errors"))
		})

		It("should return single error message", func() {
			list := errors.NewErrorList()
			list.AddError(errors.ErrorCodeValidationError, "test error")

			Expect(list.Error()).To(Equal("[VALIDATION_ERROR] test error"))
		})

		It("should return count and first error for multiple", func() {
			list := errors.NewErrorList()
			list.AddError(errors.ErrorCodeValidationError, "first error")
			list.AddError(errors.ErrorCodeValidationError, "second error")

			errMsg := list.Error()
			Expect(errMsg).To(ContainSubstring("2 errors occurred"))
			Expect(errMsg).To(ContainSubstring("first error"))
		})
	})
})

var _ = Describe("Error Wrapping", func() {
	Context("Wrap", func() {
		It("should wrap error with context", func() {
			original := fmt.Errorf("database connection failed")
			wrapped := errors.Wrap(original, errors.ErrorCodeInternalServer, "database")

			Expect(wrapped).NotTo(BeNil())
			Expect(wrapped.Code).To(Equal(errors.ErrorCodeInternalServer))
			Expect(wrapped.Component).To(Equal("database"))
			Expect(wrapped.Message).To(Equal("database connection failed"))
			Expect(wrapped.Description).To(ContainSubstring("Wrapped error"))
		})

		It("should handle nil error", func() {
			wrapped := errors.Wrap(nil, errors.ErrorCodeInternalServer, "test")

			Expect(wrapped).NotTo(BeNil())
			Expect(wrapped.Message).To(ContainSubstring("Cannot wrap nil error"))
			Expect(wrapped.Code).To(Equal(errors.ErrorCodeInternalServer))
			Expect(wrapped.Component).To(Equal("application")) // Default component from NewError
		})
	})

	Context("WrapWithRequestID", func() {
		It("should wrap error with request ID", func() {
			original := fmt.Errorf("test error")
			wrapped := errors.WrapWithRequestID(original, errors.ErrorCodeInternalServer, "req-123", "api")

			Expect(wrapped.RequestID).To(Equal("req-123"))
			Expect(wrapped.Component).To(Equal("api"))
		})
	})

	Context("WrapWithUserID", func() {
		It("should wrap error with user ID", func() {
			original := fmt.Errorf("test error")
			wrapped := errors.WrapWithUserID(original, errors.ErrorCodeInternalServer, "user-456", "auth")

			Expect(wrapped.UserID).To(Equal("user-456"))
			Expect(wrapped.Component).To(Equal("auth"))
		})
	})

	Context("Wrapf", func() {
		It("should wrap error with formatted message", func() {
			original := fmt.Errorf("connection timeout")
			baseErr := errors.NewError(errors.ErrorCodeTimeout, "timeout")
			wrapped := errors.Wrapf(original, baseErr, "failed to connect to %s", "database")

			Expect(wrapped.Code).To(Equal(errors.ErrorCodeTimeout))
			Expect(wrapped.Message).To(Equal("failed to connect to database"))
			Expect(wrapped.Description).To(ContainSubstring("connection timeout"))
		})

		It("should handle nil error", func() {
			baseErr := errors.NewError(errors.ErrorCodeInternalServer, "test")
			wrapped := errors.Wrapf(nil, baseErr, "test")

			Expect(wrapped.Message).To(ContainSubstring("Cannot wrap nil error"))
		})
	})
})

var _ = Describe("Error Combination", func() {
	Context("Combine", func() {
		It("should combine multiple errors", func() {
			err1 := errors.NewError(errors.ErrorCodeValidationError, "error 1")
			err2 := errors.NewError(errors.ErrorCodeValidationError, "error 2")
			err3 := errors.NewError(errors.ErrorCodeValidationError, "error 3")

			list := errors.Combine(err1, err2, err3)

			Expect(list.GetCount()).To(Equal(3))
			Expect(list.Errors[0]).To(Equal(err1))
			Expect(list.Errors[1]).To(Equal(err2))
			Expect(list.Errors[2]).To(Equal(err3))
		})

		It("should handle empty input", func() {
			list := errors.Combine()

			Expect(list.GetCount()).To(Equal(0))
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
			err1 := fmt.Errorf("standard error")
			err2 := errors.NewError(errors.ErrorCodeValidationError, "app error")

			list := errors.CombineErrors(err1, err2)

			Expect(list.GetCount()).To(Equal(2))
			Expect(list.Errors[0].Code).To(Equal(errors.ErrorCodeInternalServer))
			Expect(list.Errors[0].Component).To(Equal("unknown"))
			Expect(list.Errors[0].Message).To(ContainSubstring("standard error"))
			Expect(list.Errors[1].Code).To(Equal(errors.ErrorCodeValidationError))
		})

		It("should skip nil errors", func() {
			appErr := errors.NewError(errors.ErrorCodeValidationError, "app error")
			
			list := errors.CombineErrors(nil, appErr, nil)

			Expect(list.GetCount()).To(Equal(1))
			Expect(list.Errors[0].Code).To(Equal(errors.ErrorCodeValidationError))
			Expect(list.Errors[0].Message).To(Equal("app error"))
		})
	})
})

var _ = Describe("Error Comparison", func() {
	Context("Is", func() {
		It("should match errors with same code", func() {
			err1 := errors.NewError(errors.ErrorCodeValidationError, "test1")
			err2 := errors.NewError(errors.ErrorCodeValidationError, "test2")

			Expect(errors.Is(err1, err2)).To(BeTrue())
		})

		It("should not match errors with different codes", func() {
			err1 := errors.NewError(errors.ErrorCodeValidationError, "test")
			err2 := errors.NewError(errors.ErrorCodeInternalServer, "test")

			Expect(errors.Is(err1, err2)).To(BeFalse())
		})

		It("should handle non-application errors", func() {
			err1 := fmt.Errorf("standard error")
			err2 := errors.NewError(errors.ErrorCodeInternalServer, "test")

			Expect(errors.Is(err1, err2)).To(BeFalse())
		})
	})
})

var _ = Describe("Base Errors", func() {
	It("should have ErrConfigParseFailed", func() {
		Expect(errors.ErrConfigParseFailed).NotTo(BeNil())
		Expect(errors.ErrConfigParseFailed.Code).To(Equal(errors.ErrorCodeConfigParseFailed))
		Expect(errors.ErrConfigParseFailed.Message).To(Equal("Config parse failed"))
	})

	It("should have ErrInvalidValue", func() {
		Expect(errors.ErrInvalidValue).NotTo(BeNil())
		Expect(errors.ErrInvalidValue.Code).To(Equal(errors.ErrorCodeInvalidValue))
		Expect(errors.ErrInvalidValue.Message).To(Equal("Invalid value"))
	})

	It("should have ErrFileNotFound", func() {
		Expect(errors.ErrFileNotFound).NotTo(BeNil())
		Expect(errors.ErrFileNotFound.Code).To(Equal(errors.ErrorCodeFileNotFound))
		Expect(errors.ErrFileNotFound.Message).To(Equal("File not found"))
	})

	It("should have ErrInvalidType", func() {
		Expect(errors.ErrInvalidType).NotTo(BeNil())
		Expect(errors.ErrInvalidType.Code).To(Equal(errors.ErrorCodeValidationError))
		Expect(errors.ErrInvalidType.Message).To(Equal("Invalid type"))
	})
})
