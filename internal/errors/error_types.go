package errors

import (
	"fmt"
	"time"
)

// ErrorCode represents strongly-typed error codes
// TODO: Consider grouping error codes by category
// TODO: Add validation for error code patterns
// TODO: Add documentation for when to use each code.
type ErrorCode string

const (
	// Migration Errors.
	ErrorCodeMigrationFailed   ErrorCode = "MIGRATION_FAILED"
	ErrorCodeMigrationNotFound ErrorCode = "MIGRATION_NOT_FOUND"
	ErrorCodeTooManyMigrations ErrorCode = "TOO_MANY_MIGRATIONS"

	// Schema Errors.
	ErrorCodeSchemaNotFound   ErrorCode = "SCHEMA_NOT_FOUND"
	ErrorCodeSchemaValidation ErrorCode = "SCHEMA_VALIDATION"
	ErrorCodeTableNotFound    ErrorCode = "TABLE_NOT_FOUND"
	ErrorCodeColumnNotFound   ErrorCode = "COLUMN_NOT_FOUND"

	// Event Errors.
	ErrorCodeEventValidation  ErrorCode = "EVENT_VALIDATION"
	ErrorCodeEventNotFound    ErrorCode = "EVENT_NOT_FOUND"
	ErrorCodeInvalidEventType ErrorCode = "INVALID_EVENT_TYPE"
	ErrorCodeEmptyAggregateID ErrorCode = "EMPTY_AGGREGATE_ID"

	// Configuration Errors.
	ErrorCodeConfigValidation   ErrorCode = "CONFIG_VALIDATION"
	ErrorCodeConfigNotFound     ErrorCode = "CONFIG_NOT_FOUND"
	ErrorCodeConfigParseFailed  ErrorCode = "CONFIG_PARSE_FAILED"
	ErrorCodeInvalidProjectType ErrorCode = "INVALID_PROJECT_TYPE"
	ErrorCodeInvalidValue       ErrorCode = "INVALID_VALUE"

	// File Errors.
	ErrorCodeFileNotFound     ErrorCode = "FILE_NOT_FOUND"
	ErrorCodeFileReadError    ErrorCode = "FILE_READ_ERROR"
	ErrorCodeTemplateNotFound ErrorCode = "TEMPLATE_NOT_FOUND"

	// Validation Errors.
	ErrorCodeValidationError ErrorCode = "VALIDATION_ERROR"

	// System Errors.
	ErrorCodeInternalServer   ErrorCode = "INTERNAL_SERVER"
	ErrorCodeTimeout          ErrorCode = "TIMEOUT"
	ErrorCodePermissionDenied ErrorCode = "PERMISSION_DENIED"
	ErrorCodeNotFound         ErrorCode = "NOT_FOUND"
)

// ErrorSeverity represents error severity levels
// TODO: Add validation for severity transitions
// TODO: Add methods for severity comparison
// TODO: Add severity-based error handling.
type ErrorSeverity string

const (
	ErrorSeverityInfo     ErrorSeverity = "info"
	ErrorSeverityWarning  ErrorSeverity = "warning"
	ErrorSeverityError    ErrorSeverity = "error"
	ErrorSeverityCritical ErrorSeverity = "critical"
)

// IsValid returns true if the severity is valid
// TODO: Add validation logic.
func (e ErrorSeverity) IsValid() bool {
	switch e {
	case ErrorSeverityInfo, ErrorSeverityWarning, ErrorSeverityError, ErrorSeverityCritical:
		return true
	default:
		return false
	}
}

// Priority returns the numeric priority of the severity (higher = more severe)
// TODO: Define priority values
// TODO: Add comparison methods.
func (e ErrorSeverity) Priority() int {
	switch e {
	case ErrorSeverityInfo:
		return 1
	case ErrorSeverityWarning:
		return 2
	case ErrorSeverityError:
		return 3
	case ErrorSeverityCritical:
		return 4
	default:
		return 0
	}
}

// ErrorDetails represents structured error details
// TODO: Add validation for detail fields
// TODO: Add methods for detail building
// TODO: Add support for nested details.
type ErrorDetails struct {
	Field     string `json:"field,omitempty"`
	Value     any    `json:"value,omitempty"`
	Expected  any    `json:"expected,omitempty"`
	Actual    any    `json:"actual,omitempty"`
	Message   string `json:"message,omitempty"`
	Component string `json:"component,omitempty"`
	Rule      string `json:"rule,omitempty"`
	Context   string `json:"context,omitempty"`
}

// Validate validates the error details
// TODO: Implement validation logic
// TODO: Add field-specific validation.
func (ed *ErrorDetails) Validate() bool {
	// TODO: Add validation rules
	return true
}

// Error represents a structured application error
// TODO: Add thread safety considerations
// TODO: Add methods for error chaining
// TODO: Add support for error aggregation.
type Error struct {
	Code        ErrorCode     `json:"code"`
	Message     string        `json:"message"`
	Description string        `json:"description,omitempty"`
	cause       error         `json:"-"` // Hidden from JSON, used for Unwrap()
	Details     *ErrorDetails `json:"details,omitempty"`
	Timestamp   int64         `json:"timestamp"`
	RequestID   string        `json:"request_id,omitempty"`
	UserID      string        `json:"user_id,omitempty"`
	Component   string        `json:"component,omitempty"`
	Retryable   bool          `json:"retryable"`
	Severity    ErrorSeverity `json:"severity"`
}

// NewError creates a new error with validation
// TODO: Add input validation
// TODO: Add default component resolution
// TODO: Add automatic severity assignment based on error code.
func NewError(code ErrorCode, message string) *Error {
	// TODO: Add validation for code and message
	return &Error{
		Code:      code,
		Message:   message,
		Timestamp: time.Now().Unix(),
		Component: "application",
		Retryable: false,
		Severity:  ErrorSeverityError,
	}
}

// Newf creates a new error with formatted message
// TODO: Add validation for format string
// TODO: Add protection against format string injection.
func Newf(code ErrorCode, format string, args ...any) *Error {
	message := fmt.Sprintf(format, args...)
	return NewError(code, message)
}
