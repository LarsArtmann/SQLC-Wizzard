// Error codes for domain-related errors
package errors

// Domain error codes
const (
	ErrDomainNotFound   ErrorCode = "DOMAIN_NOT_FOUND"
	ErrDomainValidation ErrorCode = "DOMAIN_VALIDATION"
	ErrDomainState      ErrorCode = "DOMAIN_STATE"
	ErrDomainInvariant  ErrorCode = "DOMAIN_INVARIANT"
)

// Error helpers for domain operations
func DomainNotFoundError(aggregateID string) *Error {
	return New(ErrDomainNotFound, "domain aggregate not found: "+aggregateID).
		WithDetails("aggregate_id", aggregateID).
		WithCaller()
}

func DomainValidationError(field, value string) *Error {
	return New(ErrDomainValidation, "domain validation failed: "+field).
		WithDetails("field", field).
		WithDetails("value", value).
		WithCaller()
}

func DomainStateError(state string) *Error {
	return New(ErrDomainState, "invalid domain state: "+state).
		WithDetails("state", state).
		WithCaller()
}

func DomainInvariantError(invariant string) *Error {
	return New(ErrDomainInvariant, "domain invariant violated: "+invariant).
		WithDetails("invariant", invariant).
		WithCaller()
}
