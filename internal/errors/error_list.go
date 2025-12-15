package errors

import "fmt"

// ErrorList represents multiple errors
// TODO: Add thread safety for concurrent access
// TODO: Add methods for error filtering
// TODO: Add support for error aggregation
type ErrorList struct {
	Errors []*Error `json:"errors"`
}

// NewErrorList creates a new error list
// TODO: Add initial capacity configuration
func NewErrorList() *ErrorList {
	// TODO: Add configurable initial capacity
	return &ErrorList{
		Errors: make([]*Error, 0),
	}
}

// Add adds an error to list
// TODO: Add duplicate error detection
// TODO: Add error validation
// TODO: Add error deduplication
func (el *ErrorList) Add(err *Error) {
	// TODO: Add nil check
	// TODO: Add validation
	// TODO: Add deduplication logic
	if err != nil {
		el.Errors = append(el.Errors, err)
	}
}

// AddError adds an error using NewError
// TODO: Add message validation
// TODO: Add code validation
func (el *ErrorList) AddError(code ErrorCode, message string) {
	// TODO: Add validation
	el.Add(NewError(code, message))
}

// HasErrors returns true if list contains errors
// TODO: Add error severity filtering
func (el *ErrorList) HasErrors() bool {
	return len(el.Errors) > 0
}

// GetCount returns number of errors
func (el *ErrorList) GetCount() int {
	return len(el.Errors)
}

// Error implements error interface for ErrorList
// TODO: Add configurable error formats
// TODO: Add severity-based formatting
// TODO: Add error limiting for large lists
func (el *ErrorList) Error() string {
	if !el.HasErrors() {
		return "no errors"
	}

	if len(el.Errors) == 1 {
		return el.Errors[0].Error()
	}

	// TODO: Add configurable summary length
	return fmt.Sprintf("%d errors occurred (first: %s)", len(el.Errors), el.Errors[0].Error())
}

// Filter returns errors matching the given criteria
// TODO: Implement filtering by code, severity, component
func (el *ErrorList) Filter(predicate func(*Error) bool) *ErrorList {
	filtered := NewErrorList()
	for _, err := range el.Errors {
		if predicate(err) {
			filtered.Add(err)
		}
	}
	return filtered
}

// GroupByCode groups errors by error code
// TODO: Implement grouping functionality
func (el *ErrorList) GroupByCode() map[ErrorCode][]*Error {
	groups := make(map[ErrorCode][]*Error)
	for _, err := range el.Errors {
		groups[err.Code] = append(groups[err.Code], err)
	}
	return groups
}

// GetByCode returns errors with specific code
// TODO: Implement code-based retrieval
func (el *ErrorList) GetByCode(code ErrorCode) []*Error {
	return el.Filter(func(err *Error) bool {
		return err.Code == code
	}).Errors
}

// GetCritical returns critical errors
// TODO: Implement severity filtering
func (el *ErrorList) GetCritical() []*Error {
	return el.Filter(func(err *Error) bool {
		return err.IsCritical()
	}).Errors
}

// GetRetryable returns retryable errors
// TODO: Implement retryable filtering
func (el *ErrorList) GetRetryable() []*Error {
	return el.Filter(func(err *Error) bool {
		return err.IsRetryable()
	}).Errors
}

// Clear removes all errors from the list
// TODO: Add memory management considerations
func (el *ErrorList) Clear() {
	// TODO: Consider memory clearing for sensitive data
	el.Errors = el.Errors[:0]
}

// Clone creates a deep copy of the error list
// TODO: Implement proper deep cloning
func (el *ErrorList) Clone() *ErrorList {
	clone := NewErrorList()
	for _, err := range el.Errors {
		clone.Add(err.Clone())
	}
	return clone
}
