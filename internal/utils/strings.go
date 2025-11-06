package utils

// Contains checks if a string slice contains a specific value
func Contains(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}