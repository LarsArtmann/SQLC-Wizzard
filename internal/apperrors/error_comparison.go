package apperrors

import "errors"

// Is checks if an error matches a specific error type.
func Is(err, target error) bool {
	appErr := &Error{}
	if errors.As(err, &appErr) {
		targetErr := &Error{}
		if errors.As(target, &targetErr) {
			return appErr.Code == targetErr.Code
		}
	}
	return false
}
