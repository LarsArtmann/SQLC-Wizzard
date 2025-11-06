// Package errors provides centralized, professional error handling
// All errors in the application should use these standardized types
package errors

import (
	"fmt"
	"runtime"
	"strings"
)

// ErrorCode represents standardized error codes
type ErrorCode string

const (
	// Configuration errors
	ErrConfigNotFound     ErrorCode = "CONFIG_NOT_FOUND"
	ErrConfigParseFailed  ErrorCode = "CONFIG_PARSE_FAILED"
	ErrConfigValidation   ErrorCode = "CONFIG_VALIDATION"
	ErrValidationFailed    ErrorCode = "VALIDATION_FAILED"

	// Validation errors
	ErrInvalidConfig     ErrorCode = "INVALID_CONFIG"
	ErrMissingField      ErrorCode = "MISSING_FIELD"
	ErrInvalidType      ErrorCode = "INVALID_TYPE"
	ErrInvalidState     ErrorCode = "INVALID_STATE"
	ErrInvalidValue     ErrorCode = "INVALID_VALUE"

	// File system errors
	ErrFileNotFound     ErrorCode = "FILE_NOT_FOUND"
	ErrPermissionDenied ErrorCode = "PERMISSION_DENIED"
	ErrDirectoryExists ErrorCode = "DIRECTORY_EXISTS"

	// SQLC errors
	ErrSQLCNotFound     ErrorCode = "SQLC_NOT_FOUND"
	ErrSQLCVersion      ErrorCode = "SQLC_VERSION"
	ErrSQLCValidation   ErrorCode = "SQLC_VALIDATION"

	// Template errors
	ErrTemplateNotFound  ErrorCode = "TEMPLATE_NOT_FOUND"
	ErrTemplateInvalid   ErrorCode = "TEMPLATE_INVALID"
	ErrTemplateRender   ErrorCode = "TEMPLATE_RENDER"

	// CLI errors
	ErrInvalidCommand    ErrorCode = "INVALID_COMMAND"
	ErrInvalidFlag      ErrorCode = "INVALID_FLAG"
	ErrExecution        ErrorCode = "EXECUTION"

	// Domain errors
	ErrDomainViolation  ErrorCode = "DOMAIN_VIOLATION"
	ErrAggregateNotFound ErrorCode = "AGGREGATE_NOT_FOUND"
	ErrCommandRejected   ErrorCode = "COMMAND_REJECTED"
)

// Error represents a standardized application error
type Error struct {
	Code       ErrorCode              `json:"code"`
	Message    string                 `json:"message"`
	Details    map[string]interface{} `json:"details,omitempty"`
	Cause      error                  `json:"-"`
	File       string                 `json:"-"`
	Line       int                    `json:"-"`
	Function   string                 `json:"-"`
}

// Error implements the error interface
func (e *Error) Error() string {
	if e.Details != nil {
		return fmt.Sprintf("%s: %s (details: %+v)", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap returns the underlying cause
func (e *Error) Unwrap() error {
	return e.Cause
}

// WithCause adds a cause to the error
func (e *Error) WithCause(err error) *Error {
	e.Cause = err
	return e
}

// WithDetails adds details to the error
func (e *Error) WithDetails(key string, value interface{}) *Error {
	if e.Details == nil {
		e.Details = make(map[string]interface{})
	}
	e.Details[key] = value
	return e
}

// WithCaller adds file, line, and function information
func (e *Error) WithCaller() *Error {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		e.File = file
		e.Line = line
		e.Function = getCallerName()
	}
	return e
}

// New creates a new standardized error
func New(code ErrorCode, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
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

// Validation helpers
func ValidationError(field, value string) *Error {
	return New(ErrInvalidType, fmt.Sprintf("invalid value for %s: %s", field, value)).
		WithDetails("field", field).
		WithDetails("value", value).
		WithCaller()
}

func MissingFieldError(field string) *Error {
	return New(ErrMissingField, fmt.Sprintf("required field is missing: %s", field)).
		WithDetails("field", field).
		WithCaller()
}

func StateValidationError(state string) *Error {
	return New(ErrInvalidState, fmt.Sprintf("invalid state encountered: %s", state)).
		WithDetails("state", state).
		WithCaller()
}

// File system helpers
func FileNotFoundError(path string) *Error {
	return New(ErrFileNotFound, fmt.Sprintf("file not found: %s", path)).
		WithDetails("path", path).
		WithCaller()
}

func ConfigNotFoundError(path string) *Error {
	return New(ErrFileNotFound, fmt.Sprintf("config file not found: %s", path)).
		WithDetails("path", path).
		WithCaller()
}

func FileReadError(path string, err error) *Error {
	return Wrap(err, ErrFileNotFound, fmt.Sprintf("failed to read file: %s", path)).
		WithDetails("path", path).
		WithCaller()
}

func ConfigParseError(path string, err error) *Error {
	return Wrap(err, ErrInvalidConfig, fmt.Sprintf("failed to parse config: %s", path)).
		WithDetails("path", path).
		WithCaller()
}

func PermissionDeniedError(operation, path string) *Error {
	return New(ErrPermissionDenied, fmt.Sprintf("permission denied for %s on %s", operation, path)).
		WithDetails("operation", operation).
		WithDetails("path", path).
		WithCaller()
}

func DirectoryExistsError(path string) *Error {
	return New(ErrDirectoryExists, fmt.Sprintf("directory already exists: %s", path)).
		WithDetails("path", path).
		WithCaller()
}

// Formatting helpers
func Newf(code ErrorCode, format string, args ...interface{}) *Error {
	return New(code, fmt.Sprintf(format, args...))
}

func Wrapf(err error, code ErrorCode, format string, args ...interface{}) *Error {
	if err == nil {
		return nil
	}
	
	return Wrap(err, code, fmt.Sprintf(format, args...))
}

// Template helpers
func TemplateNotFoundError(template string) *Error {
	return New(ErrTemplateNotFound, fmt.Sprintf("template not found: %s", template)).
		WithDetails("template", template).
		WithCaller()
}

func TemplateRenderError(template string, err error) *Error {
	return Wrap(err, ErrTemplateRender, fmt.Sprintf("failed to render template: %s", template)).
		WithDetails("template", template).
		WithCaller()
}

// CLI helpers
func InvalidCommandError(command string) *Error {
	return New(ErrInvalidCommand, fmt.Sprintf("invalid command: %s", command)).
		WithDetails("command", command).
		WithCaller()
}

func InvalidFlagError(flag, value string) *Error {
	return New(ErrInvalidFlag, fmt.Sprintf("invalid flag %s with value: %s", flag, value)).
		WithDetails("flag", flag).
		WithDetails("value", value).
		WithCaller()
}

// Domain helpers
func DomainViolationError(violation string) *Error {
	return New(ErrDomainViolation, fmt.Sprintf("domain violation: %s", violation)).
		WithDetails("violation", violation).
		WithCaller()
}

func AggregateNotFoundError(aggregate, id string) *Error {
	return New(ErrAggregateNotFound, fmt.Sprintf("%s not found: %s", aggregate, id)).
		WithDetails("aggregate", aggregate).
		WithDetails("id", id).
		WithCaller()
}

func CommandRejectedError(command string, reason string) *Error {
	return New(ErrCommandRejected, fmt.Sprintf("command rejected: %s, reason: %s", command, reason)).
		WithDetails("command", command).
		WithDetails("reason", reason).
		WithCaller()
}

// getCallerName extracts the calling function name
func getCallerName() string {
	pc, _, _, ok := runtime.Caller(3)
	if !ok {
		return "unknown"
	}
	
	name := runtime.FuncForPC(pc).Name()
	// Extract just the function name (remove package path)
	parts := strings.Split(name, ".")
	return parts[len(parts)-1]
}