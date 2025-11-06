package errors

import (
	"errors"
	"fmt"
)

// ErrorCode represents a type-safe error categorization
type ErrorCode string

const (
	// Validation errors
	ErrCodeInvalidInput      ErrorCode = "INVALID_INPUT"
	ErrCodeMissingField      ErrorCode = "MISSING_FIELD"
	ErrCodeInvalidValue      ErrorCode = "INVALID_VALUE"
	ErrCodeValidationFailed  ErrorCode = "VALIDATION_FAILED"

	// Configuration errors
	ErrCodeConfigNotFound    ErrorCode = "CONFIG_NOT_FOUND"
	ErrCodeConfigParseFailed ErrorCode = "CONFIG_PARSE_FAILED"
	ErrCodeConfigInvalid     ErrorCode = "CONFIG_INVALID"

	// Template errors
	ErrCodeTemplateNotFound  ErrorCode = "TEMPLATE_NOT_FOUND"
	ErrCodeTemplateInvalid   ErrorCode = "TEMPLATE_INVALID"
	ErrCodeGenerationFailed  ErrorCode = "GENERATION_FAILED"

	// Wizard errors
	ErrCodeWizardCancelled   ErrorCode = "WIZARD_CANCELLED"
	ErrCodeWizardFailed      ErrorCode = "WIZARD_FAILED"

	// File system errors
	ErrCodeFileNotFound      ErrorCode = "FILE_NOT_FOUND"
	ErrCodeFileAlreadyExists ErrorCode = "FILE_ALREADY_EXISTS"
	ErrCodeFileWriteFailed   ErrorCode = "FILE_WRITE_FAILED"
	ErrCodeFileReadFailed    ErrorCode = "FILE_READ_FAILED"

	// Internal errors
	ErrCodeInternal          ErrorCode = "INTERNAL_ERROR"
	ErrCodeNotImplemented    ErrorCode = "NOT_IMPLEMENTED"
)

// Error represents a structured application error with context
type Error struct {
	Code    ErrorCode
	Message string
	Cause   error
	Context map[string]interface{}
}

// Error implements the error interface
func (e *Error) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap implements error unwrapping for errors.Is/As
func (e *Error) Unwrap() error {
	return e.Cause
}

// New creates a new Error with the given code and message
func New(code ErrorCode, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
		Context: make(map[string]interface{}),
	}
}

// Newf creates a new Error with formatted message
func Newf(code ErrorCode, format string, args ...interface{}) *Error {
	return &Error{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
		Context: make(map[string]interface{}),
	}
}

// Wrap wraps an existing error with additional context
func Wrap(code ErrorCode, message string, cause error) *Error {
	return &Error{
		Code:    code,
		Message: message,
		Cause:   cause,
		Context: make(map[string]interface{}),
	}
}

// Wrapf wraps an existing error with formatted message
func Wrapf(code ErrorCode, cause error, format string, args ...interface{}) *Error {
	return &Error{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
		Cause:   cause,
		Context: make(map[string]interface{}),
	}
}

// WithContext adds context information to the error
func (e *Error) WithContext(key string, value interface{}) *Error {
	if e.Context == nil {
		e.Context = make(map[string]interface{})
	}
	e.Context[key] = value
	return e
}

// GetContext retrieves context value by key
func (e *Error) GetContext(key string) (interface{}, bool) {
	if e.Context == nil {
		return nil, false
	}
	val, ok := e.Context[key]
	return val, ok
}

// Is checks if the error matches the target error
func (e *Error) Is(target error) bool {
	if t, ok := target.(*Error); ok {
		return e.Code == t.Code
	}
	return false
}

// HasCode checks if the error has the given code
func (e *Error) HasCode(code ErrorCode) bool {
	return e.Code == code
}

// HasCode checks if any error in the chain has the given code
func HasCode(err error, code ErrorCode) bool {
	var appErr *Error
	if errors.As(err, &appErr) {
		return appErr.Code == code
	}
	return false
}

// Validation Errors

// ValidationError represents a validation error
type ValidationError struct {
	*Error
	Field string
}

// NewValidationError creates a validation error for a specific field
func NewValidationError(field, message string) *ValidationError {
	return &ValidationError{
		Error: New(ErrCodeValidationFailed, message),
		Field: field,
	}
}

// NewValidationErrorf creates a validation error with formatted message
func NewValidationErrorf(field, format string, args ...interface{}) *ValidationError {
	return &ValidationError{
		Error: Newf(ErrCodeValidationFailed, format, args...),
		Field: field,
	}
}

// Config Errors

// NewConfigNotFoundError creates a config not found error
func NewConfigNotFoundError(path string) *Error {
	return New(ErrCodeConfigNotFound, fmt.Sprintf("config file not found: %s", path)).
		WithContext("path", path)
}

// NewConfigParseError wraps a parse error
func NewConfigParseError(path string, cause error) *Error {
	return Wrap(ErrCodeConfigParseFailed, fmt.Sprintf("failed to parse config: %s", path), cause).
		WithContext("path", path)
}

// NewConfigInvalidError creates a config validation error
func NewConfigInvalidError(message string) *Error {
	return New(ErrCodeConfigInvalid, message)
}

// Template Errors

// NewTemplateNotFoundError creates a template not found error
func NewTemplateNotFoundError(templateName string) *Error {
	return New(ErrCodeTemplateNotFound, fmt.Sprintf("template not found: %s", templateName)).
		WithContext("template", templateName)
}

// NewTemplateInvalidError creates a template invalid error
func NewTemplateInvalidError(templateName, message string) *Error {
	return New(ErrCodeTemplateInvalid, fmt.Sprintf("template %s is invalid: %s", templateName, message)).
		WithContext("template", templateName)
}

// NewGenerationFailedError wraps a generation error
func NewGenerationFailedError(message string, cause error) *Error {
	return Wrap(ErrCodeGenerationFailed, message, cause)
}

// Wizard Errors

// NewWizardCancelledError creates a wizard cancelled error
func NewWizardCancelledError() *Error {
	return New(ErrCodeWizardCancelled, "wizard was cancelled by user")
}

// NewWizardFailedError wraps a wizard failure
func NewWizardFailedError(step string, cause error) *Error {
	return Wrap(ErrCodeWizardFailed, fmt.Sprintf("wizard failed at step: %s", step), cause).
		WithContext("step", step)
}

// File System Errors

// NewFileNotFoundError creates a file not found error
func NewFileNotFoundError(path string) *Error {
	return New(ErrCodeFileNotFound, fmt.Sprintf("file not found: %s", path)).
		WithContext("path", path)
}

// NewFileAlreadyExistsError creates a file already exists error
func NewFileAlreadyExistsError(path string) *Error {
	return New(ErrCodeFileAlreadyExists, fmt.Sprintf("file already exists: %s", path)).
		WithContext("path", path)
}

// NewFileWriteError wraps a file write error
func NewFileWriteError(path string, cause error) *Error {
	return Wrap(ErrCodeFileWriteFailed, fmt.Sprintf("failed to write file: %s", path), cause).
		WithContext("path", path)
}

// NewFileReadError wraps a file read error
func NewFileReadError(path string, cause error) *Error {
	return Wrap(ErrCodeFileReadFailed, fmt.Sprintf("failed to read file: %s", path), cause).
		WithContext("path", path)
}
