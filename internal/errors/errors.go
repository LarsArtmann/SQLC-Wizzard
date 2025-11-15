package errors

import (
	"encoding/json"
	"fmt"
	"time"
)

// ErrorCode represents strongly-typed error codes
// Replaces string-based error codes with type safety
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
	ErrorCodeInvalidProjectType ErrorCode = "INVALID_PROJECT_TYPE"
	
	// System Errors
	ErrorCodeInternalServer     ErrorCode = "INTERNAL_SERVER"
	ErrorCodeTimeout           ErrorCode = "TIMEOUT"
	ErrorCodePermissionDenied   ErrorCode = "PERMISSION_DENIED"
	ErrorCodeNotFound          ErrorCode = "NOT_FOUND"
)

// IsValid returns true if ErrorCode is valid
func (ec ErrorCode) IsValid() bool {
	switch ec {
	case ErrorCodeMigrationFailed, ErrorCodeMigrationNotFound, ErrorCodeTooManyMigrations,
		 ErrorCodeSchemaNotFound, ErrorCodeSchemaValidation, ErrorCodeTableNotFound, ErrorCodeColumnNotFound,
		 ErrorCodeEventValidation, ErrorCodeEventNotFound, ErrorCodeInvalidEventType, ErrorCodeEmptyAggregateID,
		 ErrorCodeConfigValidation, ErrorCodeConfigNotFound, ErrorCodeInvalidProjectType,
		 ErrorCodeInternalServer, ErrorCodeTimeout, ErrorCodePermissionDenied, ErrorCodeNotFound:
		return true
	default:
		return false
	}
}

// ErrorSeverity represents error severity levels
type ErrorSeverity string

const (
	ErrorSeverityInfo    ErrorSeverity = "info"
	ErrorSeverityWarning ErrorSeverity = "warning"
	ErrorSeverityError   ErrorSeverity = "error"
	ErrorSeverityCritical ErrorSeverity = "critical"
)

// IsValid returns true if ErrorSeverity is valid
func (es ErrorSeverity) IsValid() bool {
	switch es {
	case ErrorSeverityInfo, ErrorSeverityWarning, ErrorSeverityError, ErrorSeverityCritical:
		return true
	default:
		return false
	}
}

// ErrorDetails represents structured error details
// Replaces map[string]any with typed fields
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
// Replaces any type usage with strongly typed structure
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
func NewError(code ErrorCode, message string) (*Error, error) {
	if !code.IsValid() {
		return nil, &Error{
			Code:        ErrorCodeInternalServer,
			Message:     "Invalid error code provided",
			Description: fmt.Sprintf("Error code '%s' is not a valid error code", string(code)),
			Timestamp:   time.Now().Unix(),
			Component:   "errors",
			Retryable:   false,
			Severity:    ErrorSeverityCritical,
		}
	}
	
	if message == "" {
		return nil, &Error{
			Code:        ErrorCodeEventValidation,
			Message:     "Error message cannot be empty",
			Description: "Error messages must provide meaningful information",
			Timestamp:   time.Now().Unix(),
			Component:   "errors",
			Retryable:   false,
			Severity:    ErrorSeverityError,
		}
	}
	
	return &Error{
		Code:        code,
		Message:     message,
		Timestamp:   time.Now().Unix(),
		Component:   "application",
		Retryable:   false,
		Severity:    ErrorSeverityError,
	}, nil
}

// Newf creates a new error with formatted message
// Replaces any type variadic parameters with type-safe formatting
func Newf(code ErrorCode, format string, args ...interface{}) (*Error, error) {
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
	if severity.IsValid() {
		e.Severity = severity
	}
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

// FromJSON creates error from JSON string
func FromJSON(jsonStr string) (*Error, error) {
	var err Error
	if parseErr := json.Unmarshal([]byte(jsonStr), &err); parseErr != nil {
		return nil, fmt.Errorf("failed to parse error from JSON: %w", parseErr)
	}
	
	if !err.Code.IsValid() {
		return nil, fmt.Errorf("invalid error code in JSON: %s", err.Code)
	}
	
	return &err, nil
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

// Add adds an error to the list
func (el *ErrorList) Add(err *Error) {
	if err != nil {
		el.Errors = append(el.Errors, err)
	}
}

// AddError adds an error using NewError
func (el *ErrorList) AddError(code ErrorCode, message string) error {
	err, createErr := NewError(code, message)
	if createErr != nil {
		return createErr
	}
	el.Add(err)
	return nil
}

// HasErrors returns true if list contains errors
func (el *ErrorList) HasErrors() bool {
	return len(el.Errors) > 0
}

// GetCount returns number of errors
func (el *ErrorList) GetCount() int {
	return len(el.Errors)
}

// GetByCode returns errors with specific code
func (el *ErrorList) GetByCode(code ErrorCode) []*Error {
	var matching []*Error
	for _, err := range el.Errors {
		if err.Code == code {
			matching = append(matching, err)
		}
	}
	return matching
}

// GetCritical returns all critical errors
func (el *ErrorList) GetCritical() []*Error {
	var critical []*Error
	for _, err := range el.Errors {
		if err.IsCritical() {
			critical = append(critical, err)
		}
	}
	return critical
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

// ToJSON returns error list as JSON
func (el *ErrorList) ToJSON() (string, error) {
	data, err := json.MarshalIndent(el, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal error list to JSON: %w", err)
	}
	return string(data), nil
}

// Wrap wraps an existing error with additional context
func Wrap(original error, code ErrorCode, component string) (*Error, error) {
	if original == nil {
		return nil, NewError(ErrorCodeInternalServer, "Cannot wrap nil error")
	}
	
	err, createErr := NewError(code, original.Error())
	if createErr != nil {
		return nil, createErr
	}
	
	return err.WithComponent(component).WithDescription(fmt.Sprintf("Wrapped error: %v", original)), nil
}

// WrapWithRequestID wraps an error with request ID tracking
func WrapWithRequestID(original error, code ErrorCode, requestID, component string) (*Error, error) {
	err, wrapErr := Wrap(original, code, component)
	if wrapErr != nil {
		return nil, wrapErr
	}
	
	return err.WithRequestID(requestID), nil
}

// WrapWithUserID wraps an error with user ID tracking
func WrapWithUserID(original error, code ErrorCode, userID, component string) (*Error, error) {
	err, wrapErr := Wrap(original, code, component)
	if wrapErr != nil {
		return nil, wrapErr
	}
	
	return err.WithUserID(userID), nil
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
			wrapped, _ := NewInternal("unknown", "error handling", err)
			errorList.Add(wrapped)
		}
	}
	return errorList
}

// NewInternal creates an internal server error with context
func NewInternal(component, operation string, cause error) (*Error, error) {
	message := fmt.Sprintf("Internal error in %s during %s", component, operation)
	err, createErr := NewError(ErrorCodeInternalServer, message)
	if createErr != nil {
		return nil, createErr
	}
	
	if cause != nil {
		err = err.WithDescription(cause.Error())
	}
	
	return err.WithComponent(component), nil
}

// NewNotFound creates a not found error with context
func NewNotFound(resource, id string) (*Error, error) {
	message := fmt.Sprintf("%s with ID '%s' not found", resource, id)
	return NewError(ErrorCodeNotFound, message)
}

// NewPermissionDenied creates a permission denied error
func NewPermissionDenied(resource, operation string) (*Error, error) {
	message := fmt.Sprintf("Permission denied for %s on resource: %s", operation, resource)
	return NewError(ErrorCodePermissionDenied, message)
}

// NewTimeout creates a timeout error
func NewTimeout(operation string, timeoutMs int64) (*Error, error) {
	message := fmt.Sprintf("Operation '%s' timed out after %dms", operation, timeoutMs)
	err, createErr := NewError(ErrorCodeTimeout, message)
	if createErr != nil {
		return nil, createErr
	}
	
	return err.WithRetryable(false), nil
}