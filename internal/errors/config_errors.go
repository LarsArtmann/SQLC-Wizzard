// Error codes for configuration-related errors
package errors

// Configuration error codes
const (
	ErrConfigNotFound    ErrorCode = "CONFIG_NOT_FOUND"
	ErrConfigParseFailed ErrorCode = "CONFIG_PARSE_FAILED"
	ErrConfigValidation  ErrorCode = "CONFIG_VALIDATION"
)

// Error helpers for configuration
func ConfigNotFoundError(path string) *Error {
	return New(ErrConfigNotFound, "config file not found: "+path).
		WithDetails("path", path).
		WithCaller()
}

func ConfigParseError(path string, err error) *Error {
	return Wrap(err, ErrConfigParseFailed, "failed to parse config: "+path).
		WithDetails("path", path).
		WithCaller()
}

func ConfigValidationError(field, value string) *Error {
	return New(ErrConfigValidation, "config validation failed for "+field).
		WithDetails("field", field).
		WithDetails("value", value).
		WithCaller()
}