package errors

import (
	"encoding/json"
	"fmt"
	"time"
)

// ErrorCode represents strongly-typed error codes
type ErrorCode string

const (
	// Migration Errors
	ErrorCodeMigrationFailed    ErrorCode = "MIGRATION_FAILED"
	ErrorCodeMigrationNotFound  ErrorCode = "MIGRATION_NOT_FOUND"
	ErrorCodeTooManyMigrations ErrorCode = "TOO_MANY_MIGRATIONS"
	
	// Schema Errors
	ErrorCodeSchemaNotFound      ErrorCode = "SCHEMA_NOT_FOUND"
	ErrorCodeSchemaValidation   ErrorCode = "SCHEMA_VALIDATION"
	ErrorCodeTableNotFound      ErrorCode = "TABLE_NOT_FOUND"
	ErrorCodeColumnNotFound     ErrorCode = "COLUMN_NOT_FOUND"
	
	// Event Errors
	ErrorCodeEventValidation     ErrorCode = "EVENT_VALIDATION"
	ErrorCodeEventNotFound       ErrorCode = "EVENT_NOT_FOUND"
	ErrorCodeInvalidEventType    ErrorCode = "INVALID_EVENT_TYPE"
	ErrorCodeEmptyAggregateID   ErrorCode = "EMPTY_AGGREGATE_ID"
	
	// Configuration Errors
	ErrorCodeConfigValidation    ErrorCode = "CONFIG_VALIDATION"
	ErrorCodeConfigNotFound     ErrorCode = "CONFIG_NOT_FOUND"
	ErrorCodeConfigParseFailed  ErrorCode = "CONFIG_PARSE_FAILED"
	ErrorCodeInvalidProjectType ErrorCode = "INVALID_PROJECT_TYPE"
	ErrorCodeInvalidValue        ErrorCode = "INVALID_VALUE"
	
	// File Errors
	ErrorCodeFileNotFound     ErrorCode = "FILE_NOT_FOUND"
	ErrorCodeFileReadError    ErrorCode = "FILE_READ_ERROR"
	ErrorCodeTemplateNotFound ErrorCode = "TEMPLATE_NOT_FOUND"
	
	// Validation Errors
	ErrorCodeValidationError ErrorCode = "VALIDATION_ERROR"
	
	// System Errors
	ErrorCodeInternalServer     ErrorCode = "INTERNAL_SERVER"
	ErrorCodeTimeout           ErrorCode = "TIMEOUT"
	ErrorCodePermissionDenied   ErrorCode = "PERMISSION_DENIED"
	ErrorCodeNotFound          ErrorCode = "NOT_FOUND"
)

// ErrorSeverity represents error severity levels
type ErrorSeverity string

const (
	ErrorSeverityInfo    ErrorSeverity = "info"
	ErrorSeverityWarning ErrorSeverity = "warning"
	ErrorSeverityError   ErrorSeverity = "error"
	ErrorSeverityCritical ErrorSeverity = "critical"
)

// ErrorDetails represents structured error details
type ErrorDetails struct {
	Field       string      `json:"field,omitempty"`
	Value       interface{} `json:"value,omitempty"`
	Expected    interface{} `json:"expected,omitempty"`
	Actual      interface{} `json:"actual,omitempty"`
	Message     string      `json:"message,omitempty"`
	Component   string      `json:"component,omitempty"`
	Rule        string      `json:"rule,omitempty"`
	Context     string      `json:"context,omitempty"`
}

// Error represents a structured application error
type Error struct {
	Code        ErrorCode    `json:"code"`
	Message     string       `json:"message"`
	Description string       `json:"description,omitempty"`
	Details     *ErrorDetails `json:"details,omitempty"`
	Timestamp   int64        `json:"timestamp"`
	RequestID   string        `json:"request_id,omitempty"`
	UserID      string        `json:"user_id,omitempty"`
	Component   string        `json:"component,omitempty"`
	Retryable   bool          `json:"retryable"`
	Severity    ErrorSeverity `json:"severity"`
}

// NewError creates a new error with validation
func NewError(code ErrorCode, message string) *Error {
	return &Error{
		Code:        code,
		Message:     message,
		Timestamp:   time.Now().Unix(),
		Component:   "application",
		Retryable:   false,
		Severity:    ErrorSeverityError,
	}
}

// Newf creates a new error with formatted message
func Newf(code ErrorCode, format string, args ...interface{}) *Error {
	message := fmt.Sprintf(format, args...)
	return NewError(code, message)
}

// WithDetails adds typed details to error
func (e *Error) WithDetails(field string, value, expected, actual interface{}) *Error {
	e.Details = &ErrorDetails{
		Field:    field,
		Value:    value,
		Expected: expected,
		Actual:   actual,
		Component: e.Component,
	}
	return e
}

// WithMessage adds detailed message to error details
func (e *Error) WithMessage(message string) *Error {
	if e.Details == nil {
		e.Details = &ErrorDetails{}
	}
	e.Details.Message = message
	return e
}

// WithComponent sets component for error
func (e *Error) WithComponent(component string) *Error {
	e.Component = component
	return e
}

// WithRequestID sets request ID for error tracking
func (e *Error) WithRequestID(requestID string) *Error {
	e.RequestID = requestID
	return e
}

// WithUserID sets user ID for error tracking
func (e *Error) WithUserID(userID string) *Error {
	e.UserID = userID
	return e
}

// WithRetryable sets whether error is retryable
func (e *Error) WithRetryable(retryable bool) *Error {
	e.Retryable = retryable
	return e
}

// WithSeverity sets error severity
func (e *Error) WithSeverity(severity ErrorSeverity) *Error {
	e.Severity = severity
	return e
}

// WithDescription adds detailed description
func (e *Error) WithDescription(description string) *Error {
	e.Description = description
	return e
}

// Error implements error interface
func (e *Error) Error() string {
	if e.Description != "" {
		return fmt.Sprintf("[%s] %s: %s", string(e.Code), e.Message, e.Description)
	}
	return fmt.Sprintf("[%s] %s", string(e.Code), e.Message)
}

// IsRetryable returns whether error is retryable
func (e *Error) IsRetryable() bool {
	return e.Retryable
}

// IsCritical returns whether error has critical severity
func (e *Error) IsCritical() bool {
	return e.Severity == ErrorSeverityCritical
}

// ToJSON returns error as JSON string
func (e *Error) ToJSON() (string, error) {
	data, err := json.MarshalIndent(e, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal error to JSON: %w", err)
	}
	return string(data), nil
}

// ErrorList represents multiple errors
type ErrorList struct {
	Errors []*Error `json:"errors"`
}

// NewErrorList creates a new error list
func NewErrorList() *ErrorList {
	return &ErrorList{
		Errors: make([]*Error, 0),
	}
}

// Add adds an error to list
func (el *ErrorList) Add(err *Error) {
	if err != nil {
		el.Errors = append(el.Errors, err)
	}
}

// AddError adds an error using NewError
func (el *ErrorList) AddError(code ErrorCode, message string) {
	el.Add(NewError(code, message))
}

// HasErrors returns true if list contains errors
func (el *ErrorList) HasErrors() bool {
	return len(el.Errors) > 0
}

// GetCount returns number of errors
func (el *ErrorList) GetCount() int {
	return len(el.Errors)
}

// Error implements error interface for ErrorList
func (el *ErrorList) Error() string {
	if !el.HasErrors() {
		return "no errors"
	}
	
	if len(el.Errors) == 1 {
		return el.Errors[0].Error()
	}
	
	return fmt.Sprintf("%d errors occurred (first: %s)", len(el.Errors), el.Errors[0].Error())
}

// Wrap wraps an existing error with additional context
func Wrap(original error, code ErrorCode, component string) *Error {
	if original == nil {
		return NewError(ErrorCodeInternalServer, "Cannot wrap nil error")
	}
	
	err := NewError(code, original.Error())
	return err.WithComponent(component).WithDescription(fmt.Sprintf("Wrapped error: %v", original))
}

// WrapWithRequestID wraps an error with request ID tracking
func WrapWithRequestID(original error, code ErrorCode, requestID, component string) *Error {
	err := Wrap(original, code, component)
	return err.WithRequestID(requestID)
}

// WrapWithUserID wraps an error with user ID tracking
func WrapWithUserID(original error, code ErrorCode, userID, component string) *Error {
	err := Wrap(original, code, component)
	return err.WithUserID(userID)
}

// Combine combines multiple errors into an ErrorList
func Combine(errors ...*Error) *ErrorList {
	errorList := NewErrorList()
	for _, err := range errors {
		errorList.Add(err)
	}
	return errorList
}

// CombineErrors combines error interface types into ErrorList
func CombineErrors(errs ...error) *ErrorList {
	errorList := NewErrorList()
	for _, err := range errs {
		if appErr, ok := err.(*Error); ok {
			errorList.Add(appErr)
		} else {
			// Wrap non-application errors
			wrapped := Wrap(err, ErrorCodeInternalServer, "unknown")
			errorList.Add(wrapped)
		}
	}
	return errorList
}

// NewInternal creates an internal server error with context
func NewInternal(component, operation string, cause error) *Error {
	message := fmt.Sprintf("Internal error in %s during %s", component, operation)
	err := NewError(ErrorCodeInternalServer, message)
	if cause != nil {
		err = err.WithDescription(cause.Error())
	}
	return err.WithComponent(component)
}

// NewNotFound creates a not found error with context
func NewNotFound(resource, id string) *Error {
	message := fmt.Sprintf("%s with ID '%s' not found", resource, id)
	return NewError(ErrorCodeNotFound, message)
}

// NewPermissionDenied creates a permission denied error
func NewPermissionDenied(resource, operation string) *Error {
	message := fmt.Sprintf("Permission denied for %s on resource: %s", operation, resource)
	return NewError(ErrorCodePermissionDenied, message)
}

// NewTimeout creates a timeout error
func NewTimeout(operation string, timeoutMs int64) *Error {
	message := fmt.Sprintf("Operation '%s' timed out after %dms", operation, timeoutMs)
	return NewError(ErrorCodeTimeout, message).WithRetryable(false)
}

// FileNotFoundError creates a file not found error
func FileNotFoundError(path string) *Error {
	message := fmt.Sprintf("File not found: %s", path)
	return NewError(ErrorCodeFileNotFound, message).WithComponent("filesystem")
}

// FileReadError creates a file read error
func FileReadError(path string, err error) *Error {
	message := fmt.Sprintf("Failed to read file: %s", path)
	appErr := NewError(ErrorCodeFileReadError, message).WithComponent("filesystem")
	if err != nil {
		appErr = appErr.WithDescription(err.Error())
	}
	return appErr
}

// ConfigParseError creates a config parse error
func ConfigParseError(path string, err error) *Error {
	message := fmt.Sprintf("Failed to parse config file: %s", path)
	appErr := NewError(ErrorCodeConfigParseFailed, message).WithComponent("config")
	if err != nil {
		appErr = appErr.WithDescription(err.Error())
	}
	return appErr
}

// ErrConfigParseFailed is a base error for config parsing failures
var ErrConfigParseFailed = NewError(ErrorCodeConfigParseFailed, "Config parse failed")

// ErrInvalidValue is a base error for invalid values
var ErrInvalidValue = NewError(ErrorCodeInvalidValue, "Invalid value")

// Wrapf wraps an error with formatted message and base error
func Wrapf(err error, baseErr *Error, format string, args ...interface{}) *Error {
	if err == nil {
		return NewError(ErrorCodeInternalServer, "Cannot wrap nil error")
	}
	
	message := fmt.Sprintf(format, args...)
	wrapped := NewError(baseErr.Code, message).WithComponent(baseErr.Component)
	if baseErr.Description != "" {
		wrapped = wrapped.WithDescription(baseErr.Description)
	}
	return wrapped.WithDescription(fmt.Sprintf("Original error: %v", err))
}

// TemplateNotFoundError creates a template not found error
func TemplateNotFoundError(template string) *Error {
	message := fmt.Sprintf("Template not found: %s", template)
	return NewError(ErrorCodeTemplateNotFound, message).WithComponent("templates")
}

// ValidationError creates a validation error
func ValidationError(field, message string) *Error {
	errMsg := fmt.Sprintf("Validation failed for field '%s': %s", field, message)
	return NewError(ErrorCodeValidationError, errMsg).WithComponent("validation")
}

// Is checks if an error matches a specific error type
func Is(err, target error) bool {
	if appErr, ok := err.(*Error); ok {
		if targetErr, ok := target.(*Error); ok {
			return appErr.Code == targetErr.Code
		}
	}
	return false
}

// ErrFileNotFound is a base error for file not found
var ErrFileNotFound = NewError(ErrorCodeFileNotFound, "File not found")

// ErrInvalidType is a base error for invalid type errors
var ErrInvalidType = NewError(ErrorCodeValidationError, "Invalid type")