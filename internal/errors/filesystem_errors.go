// Error helpers for file system and internal errors
package errors

import (
	"os"
)

// Additional file system error codes (unique to this file)
const (
	ErrFileRead  ErrorCode = "FILE_READ"
	ErrFileWrite ErrorCode = "FILE_WRITE"
)

// Error helpers for file system (unique functions)
func FileNotFoundError(path string) *Error {
	return New(ErrFileNotFound, "file not found: "+path).
		WithDetails("path", path).
		WithCaller()
}

func PermissionDeniedError(operation, path string) *Error {
	return New(ErrPermissionDenied, "permission denied for "+operation+" on "+path).
		WithDetails("operation", operation).
		WithDetails("path", path).
		WithCaller()
}

func DirectoryExistsError(path string) *Error {
	return New(ErrDirectoryExists, "directory already exists: "+path).
		WithDetails("path", path).
		WithCaller()
}

func FileReadError(path string, err error) *Error {
	return Wrap(err, ErrFileRead, "failed to read file: "+path).
		WithDetails("path", path).
		WithCaller()
}

func FileWriteError(path string, err error) *Error {
	return Wrap(err, ErrFileWrite, "failed to write file: "+path).
		WithDetails("path", path).
		WithCaller()
}

// Helper to check file existence
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// Error helpers for internal operations (using existing functions)
func ValidationError(field, value string) *Error {
	return New(ErrInvalidType, "invalid value for "+field+": "+value).
		WithDetails("field", field).
		WithDetails("value", value).
		WithCaller()
}

func MissingFieldError(field string) *Error {
	return New(ErrMissingField, "required field is missing: "+field).
		WithDetails("field", field).
		WithCaller()
}

func StateValidationError(state string) *Error {
	return New(ErrInvalidState, "invalid state encountered: "+state).
		WithDetails("state", state).
		WithCaller()
}

func InternalError(message string) *Error {
	return New(ErrInternal, "internal error: "+message).
		WithCaller()
}

func ExecutionError(operation string, err error) *Error {
	return Wrap(err, ErrExecution, "execution failed for "+operation).
		WithDetails("operation", operation).
		WithCaller()
}

// Template helper for tests (using existing ErrTemplateNotFound)
func TemplateNotFoundError(template string) *Error {
	return New(ErrTemplateNotFound, "template not found: "+template).
		WithDetails("template", template).
		WithCaller()
}
