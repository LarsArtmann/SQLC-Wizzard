package errors

import (
	"fmt"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Error Creation", func() {
	Context("NewError", func() {
		It("should create a new error with defaults", func() {
			err := NewError(ErrorCodeInternalServer, "test error")

			Expect(err).To(HaveOccurred())
			Expect(err.Code).To(Equal(ErrorCodeInternalServer))
			Expect(err.Message).To(Equal("test error"))
			Expect(err.Component).To(Equal("application"))
			Expect(err.Retryable).To(BeFalse())
			Expect(err.Severity).To(Equal(ErrorSeverityError))
			Expect(err.Timestamp).To(BeNumerically(">", 0))

			// TODO: Add timestamp range validation
			// TODO: Add validation for all default fields
		})
	})

	Context("Newf", func() {
		It("should create error with formatted message", func() {
			err := Newf(ErrorCodeValidationError, "field %s has value %d", "age", 150)

			Expect(err).To(HaveOccurred())
			Expect(err.Message).To(Equal("field age has value 150"))

			// TODO: Add tests for complex formatting scenarios
			// TODO: Add tests for edge cases in formatting
		})
	})

	Context("Helper Constructors", func() {
		It("should create internal error with context", func() {
			cause := errors.New("database connection failed")
			err := NewInternal("database", "connect", cause)

			Expect(err).To(HaveOccurred())
			Expect(err.Code).To(Equal(ErrorCodeInternalServer))
			Expect(err.Component).To(Equal("database"))
			Expect(err.Message).To(ContainSubstring("database"))
			Expect(err.Message).To(ContainSubstring("connect"))
			Expect(err.Description).To(Equal("database connection failed"))

			// TODO: Add validation for cause chaining
			// TODO: Add tests for nil cause handling
		})

		It("should create not found error", func() {
			err := NewNotFound("User", "12345")

			Expect(err).To(HaveOccurred())
			Expect(err.Code).To(Equal(ErrorCodeNotFound))
			Expect(err.Message).To(ContainSubstring("User"))
			Expect(err.Message).To(ContainSubstring("12345"))

			// TODO: Add validation for empty ID handling
			// TODO: Add validation for special characters in ID
		})

		It("should create permission denied error", func() {
			err := NewPermissionDenied("database", "write")

			Expect(err).To(HaveOccurred())
			Expect(err.Code).To(Equal(ErrorCodePermissionDenied))
			Expect(err.Message).To(ContainSubstring("write"))
			Expect(err.Message).To(ContainSubstring("database"))

			// TODO: Add validation for empty resource/operation
			// TODO: Add tests for special characters
		})

		It("should create timeout error", func() {
			err := NewTimeout("fetchData", 5000)

			Expect(err).To(HaveOccurred())
			Expect(err.Code).To(Equal(ErrorCodeTimeout))
			Expect(err.Message).To(ContainSubstring("fetchData"))
			Expect(err.Message).To(ContainSubstring("5000"))
			Expect(err.Retryable).To(BeFalse())

			// TODO: Add validation for negative timeout values
			// TODO: Add tests for zero timeout
		})

		It("should create file not found error", func() {
			err := FileNotFoundError("/path/to/file.txt")

			Expect(err).To(HaveOccurred())
			Expect(err.Code).To(Equal(ErrorCodeFileNotFound))
			Expect(err.Message).To(ContainSubstring("/path/to/file.txt"))
			Expect(err.Component).To(Equal("filesystem"))

			// TODO: Add validation for empty path
			// TODO: Add tests for relative vs absolute paths
		})

		It("should create file read error", func() {
			cause := errors.New("permission denied")
			err := FileReadError("/path/to/file.txt", cause)

			Expect(err).To(HaveOccurred())
			Expect(err.Code).To(Equal(ErrorCodeFileReadError))
			Expect(err.Message).To(ContainSubstring("/path/to/file.txt"))
			Expect(err.Component).To(Equal("filesystem"))

			// TODO: Add validation for nil cause
			// TODO: Add tests for different file error scenarios
		})

		It("should create config parse error", func() {
			cause := errors.New("invalid YAML")
			err := ConfigParseError("/path/to/config.yaml", cause)

			Expect(err).To(HaveOccurred())
			Expect(err.Code).To(Equal(ErrorCodeConfigParseFailed))
			Expect(err.Message).To(ContainSubstring("/path/to/config.yaml"))
			Expect(err.Component).To(Equal("config"))

			// TODO: Add validation for different config formats
			// TODO: Add tests for malformed paths
		})

		It("should create template not found error", func() {
			err := TemplateNotFoundError("missing-template")

			Expect(err).To(HaveOccurred())
			Expect(err.Code).To(Equal(ErrorCodeTemplateNotFound))
			Expect(err.Message).To(ContainSubstring("missing-template"))
			Expect(err.Component).To(Equal("templates"))

			// TODO: Add validation for empty template names
			// TODO: Add tests for template name patterns
		})

		It("should create validation error", func() {
			err := ValidationError("email", "invalid format")

			Expect(err).To(HaveOccurred())
			Expect(err.Code).To(Equal(ErrorCodeValidationError))
			Expect(err.Message).To(ContainSubstring("email"))
			Expect(err.Message).To(ContainSubstring("invalid format"))
			Expect(err.Component).To(Equal("validation"))

			// TODO: Add validation for empty field names
			// TODO: Add tests for special characters in field names
		})
	})

	// TODO: Add tests for error creation with nil inputs
	// TODO: Add tests for error creation with special characters
	// TODO: Add tests for error creation with unicode characters
	// TODO: Add performance tests for error creation
})

func TestErrorCreation(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Error Creation Suite")
}
