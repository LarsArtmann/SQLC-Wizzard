package errors

// Is checks if an error matches a specific error type
func Is(err, target error) bool {
	if appErr, ok := err.(*Error); ok {
		if targetErr, ok := target.(*Error); ok {
			return appErr.Code == targetErr.Code
		}
	}
	return false
}
