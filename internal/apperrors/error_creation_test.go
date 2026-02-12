package apperrors

import (
	stderrors "errors"

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
			cause := stderrors.New("database connection failed")
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

			expectError(err, ErrorCodeNotFound, []string{"User", "12345"})

			// TODO: Add validation for empty ID handling
			// TODO: Add validation for special characters in ID
		})

		It("should create permission denied error", func() {
			err := NewPermissionDenied("database", "write")

			expectError(err, ErrorCodePermissionDenied, []string{"database", "write"})

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

			expectError(err, ErrorCodeFileNotFound, []string{"/path/to/file.txt"}, "filesystem")

			// TODO: Add validation for empty path
			// TODO: Add tests for relative vs absolute paths
		})

		It("should create file read error", func() {
			cause := stderrors.New("permission denied")
			err := FileReadError("/path/to/file.txt", cause)

			expectError(err, ErrorCodeFileReadError, []string{"/path/to/file.txt"}, "filesystem")

			// TODO: Add validation for nil cause
			// TODO: Add tests for different file error scenarios
		})

		It("should create config parse error", func() {
			cause := stderrors.New("invalid YAML")
			err := ConfigParseError("/path/to/config.yaml", cause)

			expectError(err, ErrorCodeConfigParseFailed, []string{"/path/to/config.yaml"}, "config")

			// TODO: Add validation for different config formats
			// TODO: Add tests for malformed paths
		})

		It("should create template not found error", func() {
			err := TemplateNotFoundError("missing-template")

			expectError(err, ErrorCodeTemplateNotFound, []string{"missing-template"}, "templates")

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

// expectError is a generic helper function for verifying common Error properties.
// It optionally validates message substrings and component.
func expectError(err *Error, expectedCode ErrorCode, expectedMessageSubstrings []string, expectedComponent ...string) {
	Expect(err).To(HaveOccurred())
	Expect(err.Code).To(Equal(expectedCode))
	for _, substr := range expectedMessageSubstrings {
		Expect(err.Message).To(ContainSubstring(substr))
	}
	if len(expectedComponent) > 0 {
		Expect(err.Component).To(Equal(expectedComponent[0]))
	}
}
