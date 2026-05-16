package apperrors

import "fmt"

// MultiError represents multiple errors
// TODO: Add thread safety for concurrent access
// TODO: Add methods for error filtering
// TODO: Add support for error aggregation.
type MultiError struct {
	Errors []*Error `json:"errors"`
}

// NewMultiError creates a new error list
// TODO: Add initial capacity configuration.
func NewMultiError() *MultiError {
	// TODO: Add configurable initial capacity
	return &MultiError{
		Errors: make([]*Error, 0),
	}
}

// Add adds an error to list
// TODO: Add duplicate error detection
// TODO: Add error validation
// TODO: Add error deduplication.
func (me *MultiError) Add(err *Error) {
	// TODO: Add nil check
	// TODO: Add validation
	// TODO: Add deduplication logic
	if err != nil {
		me.Errors = append(me.Errors, err)
	}
}

// AddError adds an error using NewError
// TODO: Add message validation
// TODO: Add code validation.
func (me *MultiError) AddError(code ErrorCode, message string) {
	// TODO: Add validation
	me.Add(NewError(code, message))
}

// HasErrors returns true if list contains errors
// TODO: Add error severity filtering.
func (me *MultiError) HasErrors() bool {
	return len(me.Errors) > 0
}

// GetCount returns number of errors.
func (me *MultiError) GetCount() int {
	return len(me.Errors)
}

// Error implements error interface for MultiError
// TODO: Add configurable error formats
// TODO: Add severity-based formatting
// TODO: Add error limiting for large lists.
func (me *MultiError) Error() string {
	if !me.HasErrors() {
		return "no errors"
	}

	if len(me.Errors) == 1 {
		return me.Errors[0].Error()
	}

	// TODO: Add configurable summary length
	return fmt.Sprintf("%d errors occurred (first: %s)", len(me.Errors), me.Errors[0].Error())
}

// Filter returns errors matching the given criteria
// TODO: Implement filtering by code, severity, component.
func (me *MultiError) Filter(predicate func(*Error) bool) *MultiError {
	filtered := NewMultiError()

	for _, err := range me.Errors {
		if predicate(err) {
			filtered.Add(err)
		}
	}

	return filtered
}

// GroupByCode groups errors by error code
// TODO: Implement grouping functionality.
func (me *MultiError) GroupByCode() map[ErrorCode][]*Error {
	groups := make(map[ErrorCode][]*Error)
	for _, err := range me.Errors {
		groups[err.Code] = append(groups[err.Code], err)
	}

	return groups
}

// GetByCode returns errors with specific code
// TODO: Implement code-based retrieval.
func (me *MultiError) GetByCode(code ErrorCode) []*Error {
	return me.Filter(func(err *Error) bool {
		return err.Code == code
	}).Errors
}

// GetCritical returns critical errors
// TODO: Implement severity filtering.
func (me *MultiError) GetCritical() []*Error {
	return me.Filter(func(err *Error) bool {
		return err.IsCritical()
	}).Errors
}

// GetRetryable returns retryable errors
// TODO: Implement retryable filtering.
func (me *MultiError) GetRetryable() []*Error {
	return me.Filter(func(err *Error) bool {
		return err.IsRetryable()
	}).Errors
}

// Clear removes all errors from the list
// TODO: Add memory management considerations.
func (me *MultiError) Clear() {
	// TODO: Consider memory clearing for sensitive data
	me.Errors = me.Errors[:0]
}

// Clone creates a deep copy of the error list
// TODO: Implement proper deep cloning.
func (me *MultiError) Clone() *MultiError {
	clone := NewMultiError()
	for _, err := range me.Errors {
		clone.Add(err.Clone())
	}

	return clone
}
