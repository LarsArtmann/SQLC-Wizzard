// Core error types and helper functions - keep minimal
package errors

import (
	"fmt"
	"runtime"
	"strings"
)

// ErrorCode represents standardized error codes
type ErrorCode string

// Error represents a standardized application error
type Error struct {
	Code    ErrorCode      `json:"code"`
	Message string         `json:"message"`
	Cause   error          `json:"cause,omitempty"`
	Details map[string]any `json:"details,omitempty"`
	Stack   []string       `json:"stack,omitempty"`
}

// Error implements the error interface
func (e *Error) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap returns the underlying cause
func (e *Error) Unwrap() error {
	return e.Cause
}

// WithCaller adds caller information to the error
func (e *Error) WithCaller() *Error {
	if e.Details == nil {
		e.Details = make(map[string]any)
	}

	if _, file, line, ok := runtime.Caller(1); ok {
		e.Details["file"] = file
		e.Details["line"] = line
		e.Details["function"] = functionName(file, line)
	}

	return e
}

// WithDetails adds details to the error
func (e *Error) WithDetails(key string, value any) *Error {
	if e.Details == nil {
		e.Details = make(map[string]any)
	}
	e.Details[key] = value
	return e
}

// WithCause adds a cause to the error
func (e *Error) WithCause(err error) *Error {
	e.Cause = err
	return e
}

// New creates a new standardized error
func New(code ErrorCode, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
		Details: make(map[string]any),
	}
}

// Wrap wraps an existing error with context
func Wrap(err error, code ErrorCode, message string) *Error {
	if err == nil {
		return nil
	}

	e := &Error{
		Code:    code,
		Message: message,
		Cause:   err,
		Details: make(map[string]any),
	}
	return e.WithCaller()
}

// Is checks if error matches specific code
func Is(err error, code ErrorCode) bool {
	if appErr, ok := err.(*Error); ok {
		return appErr.Code == code
	}
	return false
}

// HasCode checks if error contains specific code
func HasCode(err error, code ErrorCode) bool {
	return Is(err, code)
}

// GetCode extracts error code from error
func GetCode(err error) ErrorCode {
	if appErr, ok := err.(*Error); ok {
		return appErr.Code
	}
	return ErrExecution
}

// GetMessage extracts message from error
func GetMessage(err error) string {
	if appErr, ok := err.(*Error); ok {
		return appErr.Message
	}
	return err.Error()
}

// GetDetails extracts details from error
func GetDetails(err error) map[string]any {
	if appErr, ok := err.(*Error); ok {
		return appErr.Details
	}
	return make(map[string]any)
}

// functionName extracts function name from file and line
func functionName(file string, line int) string {
	// Simplified function name extraction
	parts := strings.Split(file, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return "unknown"
}

// Additional core error codes (preserve from original)
const (
	// File system errors
	ErrFileNotFound     ErrorCode = "FILE_NOT_FOUND"
	ErrPermissionDenied ErrorCode = "PERMISSION_DENIED"
	ErrDirectoryExists  ErrorCode = "DIRECTORY_EXISTS"

	// Validation errors
	ErrInvalidConfig ErrorCode = "INVALID_CONFIG"
	ErrMissingField  ErrorCode = "MISSING_FIELD"
	ErrInvalidType   ErrorCode = "INVALID_TYPE"
	ErrInvalidState  ErrorCode = "INVALID_STATE"
	ErrInvalidValue  ErrorCode = "INVALID_VALUE"

	// Internal errors
	ErrInternal         ErrorCode = "INTERNAL"
	ErrExecution        ErrorCode = "EXECUTION"
	ErrValidationFailed ErrorCode = "VALIDATION_FAILED"

	// SQLC errors
	ErrSQLCNotFound   ErrorCode = "SQLC_NOT_FOUND"
	ErrSQLCVersion    ErrorCode = "SQLC_VERSION"
	ErrSQLCValidation ErrorCode = "SQLC_VALIDATION"

	// Template errors
	ErrTemplateNotFound ErrorCode = "TEMPLATE_NOT_FOUND"
	ErrTemplateInvalid  ErrorCode = "TEMPLATE_INVALID"
	ErrTemplateRender   ErrorCode = "TEMPLATE_RENDER"

	// CLI errors
	ErrInvalidCommand ErrorCode = "INVALID_COMMAND"
	ErrInvalidFlag    ErrorCode = "INVALID_FLAG"

	// Domain errors
	ErrDomainViolation   ErrorCode = "DOMAIN_VIOLATION"
	ErrAggregateNotFound ErrorCode = "AGGREGATE_NOT_FOUND"
	ErrCommandRejected   ErrorCode = "COMMAND_REJECTED"
)

// Formatting helpers (preserve from original)
func Newf(code ErrorCode, format string, args ...any) *Error {
	return New(code, fmt.Sprintf(format, args...))
}

func Wrapf(err error, code ErrorCode, format string, args ...any) *Error {
	if err == nil {
		return nil
	}

	return Wrap(err, code, fmt.Sprintf(format, args...))
}
