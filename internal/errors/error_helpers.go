package errors

import (
	"fmt"
)

// Wrap wraps an existing error with additional context
// TODO: Add validation for nil error handling
// TODO: Add automatic component detection
// TODO: Add context preservation
func Wrap(original error, code ErrorCode, component string) *Error {
	if original == nil {
		return NewError(ErrorCodeInternalServer, "Cannot wrap nil error")
	}

	// TODO: Add validation for component
	err := NewError(code, original.Error())
	return err.WithCause(original).WithComponent(component).WithDescription(fmt.Sprintf("Wrapped error: %v", original))
}

// WrapWithRequestID wraps an error with request ID tracking
// TODO: Add request ID validation
// TODO: Add automatic request ID extraction from context
func WrapWithRequestID(original error, code ErrorCode, requestID, component string) *Error {
	err := Wrap(original, code, component)
	return err.WithRequestID(requestID)
}

// WrapWithUserID wraps an error with user ID tracking
// TODO: Add user ID validation
// TODO: Add automatic user ID extraction from context
func WrapWithUserID(original error, code ErrorCode, userID, component string) *Error {
	err := Wrap(original, code, component)
	return err.WithUserID(userID)
}

// Combine combines multiple errors into an ErrorList
// TODO: Add nil error filtering
// TODO: Add duplicate error detection
func Combine(errors ...*Error) *ErrorList {
	errorList := NewErrorList()
	for _, err := range errors {
		// TODO: Add error validation
		errorList.Add(err)
	}
	return errorList
}

// CombineErrors combines error interface types into ErrorList
// TODO: Add better handling of wrapped errors
// TODO: Add error type detection and categorization
func CombineErrors(errs ...error) *ErrorList {
	errorList := NewErrorList()
	for _, err := range errs {
		if err == nil {
			continue // Skip nil errors
		}
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

// Helper constructors for common error types
// TODO: Add more specific error constructors
// TODO: Add context-aware constructors
// TODO: Add validation constructors

// NewInternal creates an internal server error with context
// TODO: Add automatic component detection from call stack
// TODO: Add operation validation
func NewInternal(component, operation string, cause error) *Error {
	message := fmt.Sprintf("Internal error in %s during %s", component, operation)
	err := NewError(ErrorCodeInternalServer, message)
	if cause != nil {
		err = err.WithDescription(cause.Error())
	}
	return err.WithComponent(component)
}

// NewNotFound creates a not found error with context
// TODO: Add resource type validation
// TODO: Add ID format validation
func NewNotFound(resource, id string) *Error {
	// TODO: Add validation for resource and ID
	message := fmt.Sprintf("%s with ID '%s' not found", resource, id)
	return NewError(ErrorCodeNotFound, message)
}

// NewPermissionDenied creates a permission denied error
// TODO: Add operation validation
// TODO: Add resource validation
func NewPermissionDenied(resource, operation string) *Error {
	// TODO: Add validation for inputs
	message := fmt.Sprintf("Permission denied for %s on resource: %s", operation, resource)
	return NewError(ErrorCodePermissionDenied, message)
}

// NewTimeout creates a timeout error
// TODO: Add timeout value validation
// TODO: Add operation validation
func NewTimeout(operation string, timeoutMs int64) *Error {
	// TODO: Add validation for timeout value
	message := fmt.Sprintf("Operation '%s' timed out after %dms", operation, timeoutMs)
	return NewError(ErrorCodeTimeout, message).WithRetryable(false)
}

// FileNotFoundError creates a file not found error
// TODO: Add path validation
// TODO: Add path normalization
func FileNotFoundError(path string) *Error {
	// TODO: Add path validation
	message := fmt.Sprintf("File not found: %s", path)
	return NewError(ErrorCodeFileNotFound, message).WithComponent("filesystem")
}

// FileReadError creates a file read error
// TODO: Add path validation
// TODO: Add error wrapping validation
func FileReadError(path string, err error) *Error {
	// TODO: Add path validation
	message := fmt.Sprintf("Failed to read file: %s", path)
	appErr := NewError(ErrorCodeFileReadError, message).WithComponent("filesystem")
	if err != nil {
		appErr = appErr.WithCause(err)
	}
	return appErr
}

// ConfigParseError creates a config parse error
// TODO: Add path validation
// TODO: Add config type detection
func ConfigParseError(path string, err error) *Error {
	// TODO: Add path validation
	message := fmt.Sprintf("Failed to parse config file: %s", path)
	appErr := NewError(ErrorCodeConfigParseFailed, message).WithComponent("config")
	if err != nil {
		appErr = appErr.WithDescription(err.Error())
	}
	return appErr
}

// TemplateNotFoundError creates a template not found error
// TODO: Add template name validation
// TODO: Add template type checking
func TemplateNotFoundError(template string) *Error {
	// TODO: Add template name validation
	message := fmt.Sprintf("Template not found: %s", template)
	return NewError(ErrorCodeTemplateNotFound, message).WithComponent("templates")
}

// ValidationError creates a validation error
// TODO: Add field name validation
// TODO: Add validation rule checking
func ValidationError(field, message string) *Error {
	// TODO: Add field validation
	errMsg := fmt.Sprintf("Validation failed for field '%s': %s", field, message)
	return NewError(ErrorCodeValidationError, errMsg).WithComponent("validation")
}

// Wrapf wraps an error with formatted message and base error
// TODO: Add format string validation
// TODO: Add protection against format injection
func Wrapf(err error, baseErr *Error, format string, args ...any) *Error {
	if err == nil {
		return NewError(ErrorCodeInternalServer, "Cannot wrap nil error")
	}

	// TODO: Add format validation
	message := fmt.Sprintf(format, args...)
	wrapped := NewError(baseErr.Code, message).WithComponent(baseErr.Component)

	// Merge descriptions instead of clobbering
	var finalDesc string
	if baseErr.Description != "" {
		finalDesc = baseErr.Description + " | " + fmt.Sprintf("Original error: %v", err)
	} else {
		finalDesc = fmt.Sprintf("Original error: %v", err)
	}

	return wrapped.WithCause(err).WithDescription(finalDesc)
}

// Base error instances for common cases
// TODO: Consider making these immutable
// TODO: Add context-specific base errors
var (
	ErrConfigParseFailed = NewError(ErrorCodeConfigParseFailed, "Config parse failed")
	ErrInvalidValue      = NewError(ErrorCodeInvalidValue, "Invalid value")
	ErrFileNotFound      = NewError(ErrorCodeFileNotFound, "File not found")
	ErrInvalidType       = NewError(ErrorCodeValidationError, "Invalid type")
)
