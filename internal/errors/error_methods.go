package errors

import (
	"encoding/json"
	"fmt"
)

// WithDetails adds typed details to error
// TODO: Add validation for detail fields
// TODO: Add support for multiple details
// TODO: Add method to update existing details.
func (e *Error) WithDetails(field string, value, expected, actual any) *Error {
	// TODO: Add validation for field name
	// TODO: Add validation for value types
	e.Details = &ErrorDetails{
		Field:     field,
		Value:     value,
		Expected:  expected,
		Actual:    actual,
		Component: e.Component,
	}
	return e
}

// WithMessage adds detailed message to error details
// TODO: Add validation for message content
// TODO: Add support for message formatting
// TODO: Add message length limits.
func (e *Error) WithMessage(message string) *Error {
	// TODO: Add message validation
	if e.Details == nil {
		e.Details = &ErrorDetails{}
	}
	e.Details.Message = message
	return e
}

// WithComponent sets component for error
// TODO: Add component validation
// TODO: Add component name normalization
// TODO: Add component hierarchy support.
func (e *Error) WithComponent(component string) *Error {
	// TODO: Add validation for component name
	e.Component = component
	if e.Details != nil {
		e.Details.Component = component
	}
	return e
}

// WithRequestID sets request ID for error tracking
// TODO: Add request ID format validation
// TODO: Add support for multiple tracking IDs.
func (e *Error) WithRequestID(requestID string) *Error {
	// TODO: Add request ID validation
	e.RequestID = requestID
	return e
}

// WithUserID sets user ID for error tracking
// TODO: Add user ID format validation
// TODO: Add support for user context.
func (e *Error) WithUserID(userID string) *Error {
	// TODO: Add user ID validation
	e.UserID = userID
	return e
}

// WithRetryable sets whether error is retryable
// TODO: Add validation for retryable logic
// TODO: Add automatic retryable determination based on error code.
func (e *Error) WithRetryable(retryable bool) *Error {
	e.Retryable = retryable
	return e
}

// WithSeverity sets error severity
// TODO: Add severity validation
// TODO: Add automatic severity assignment based on error code.
func (e *Error) WithSeverity(severity ErrorSeverity) *Error {
	// TODO: Validate severity is valid
	e.Severity = severity
	return e
}

// WithDescription adds detailed description
// TODO: Add description validation
// TODO: Add description length limits
// TODO: Add support for structured descriptions.
func (e *Error) WithDescription(description string) *Error {
	// TODO: Add description validation
	e.Description = description
	return e
}

// WithCause wraps the error with a cause for unwrapping
// TODO: Add cause validation
// TODO: Add support for multiple causes
// TODO: Add circular reference detection.
func (e *Error) WithCause(cause error) *Error {
	// TODO: Add nil check
	// TODO: Add circular reference check
	e.cause = cause
	return e
}

// Unwrap returns the underlying cause error
// TODO: Add support for chain unwrapping
// TODO: Add unwrapping limits.
func (e *Error) Unwrap() error {
	return e.cause
}

// Error implements error interface
// TODO: Add support for structured error output
// TODO: Add internationalization support
// TODO: Add configurable error formats.
func (e *Error) Error() string {
	if e.Description != "" {
		return fmt.Sprintf("[%s] %s: %s", string(e.Code), e.Message, e.Description)
	}
	return fmt.Sprintf("[%s] %s", string(e.Code), e.Message)
}

// IsRetryable returns whether error is retryable
// TODO: Add logic based on error code
// TODO: Add consideration for error history.
func (e *Error) IsRetryable() bool {
	return e.Retryable
}

// IsCritical returns whether error has critical severity
// TODO: Add logic based on severity and error type
// TODO: Add configurable criticality thresholds.
func (e *Error) IsCritical() bool {
	return e.Severity == ErrorSeverityCritical
}

// ToJSON returns error as JSON string
// TODO: Add pretty printing options
// TODO: Add field filtering
// TODO: Add sensitive data masking.
func (e *Error) ToJSON() (string, error) {
	// TODO: Add JSON marshaling options
	// TODO: Add error handling for special characters
	data, err := json.MarshalIndent(e, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal error to JSON: %w", err)
	}
	return string(data), nil
}

// Clone creates a copy of the error
// TODO: Implement deep copying for details
// TODO: Handle cause chain copying.
func (e *Error) Clone() *Error {
	clone := &Error{
		Code:        e.Code,
		Message:     e.Message,
		Description: e.Description,
		cause:       e.cause, // TODO: Consider deep copy
		Timestamp:   e.Timestamp,
		RequestID:   e.RequestID,
		UserID:      e.UserID,
		Component:   e.Component,
		Retryable:   e.Retryable,
		Severity:    e.Severity,
	}

	if e.Details != nil {
		clone.Details = &ErrorDetails{
			Field:     e.Details.Field,
			Value:     e.Details.Value,
			Expected:  e.Details.Expected,
			Actual:    e.Details.Actual,
			Message:   e.Details.Message,
			Component: e.Details.Component,
			Rule:      e.Details.Rule,
			Context:   e.Details.Context,
		}
	}

	return clone
}
