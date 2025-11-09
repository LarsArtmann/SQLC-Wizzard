// Error codes for external system errors
package errors

// External system error codes
const (
	ErrExternalNotFound    ErrorCode = "EXTERNAL_NOT_FOUND"
	ErrExternalUnavailable ErrorCode = "EXTERNAL_UNAVAILABLE"
	ErrExternalTimeout     ErrorCode = "EXTERNAL_TIMEOUT"
	ErrExternalPermission  ErrorCode = "EXTERNAL_PERMISSION"
)

// Error helpers for external systems
func ExternalNotFoundError(system, resource string) *Error {
	return New(ErrExternalNotFound, "external resource not found in "+system+": "+resource).
		WithDetails("system", system).
		WithDetails("resource", resource).
		WithCaller()
}

func ExternalUnavailableError(system string) *Error {
	return New(ErrExternalUnavailable, "external system unavailable: "+system).
		WithDetails("system", system).
		WithCaller()
}

func ExternalTimeoutError(system, operation string) *Error {
	return New(ErrExternalTimeout, "external system timeout: "+system+" during "+operation).
		WithDetails("system", system).
		WithDetails("operation", operation).
		WithCaller()
}

func ExternalPermissionError(system, operation string) *Error {
	return New(ErrExternalPermission, "external system permission denied: "+system+" for "+operation).
		WithDetails("system", system).
		WithDetails("operation", operation).
		WithCaller()
}
