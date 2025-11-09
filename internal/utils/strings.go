package utils

import "slices"

// Contains checks if a string slice contains a specific value
func Contains(slice []string, value string) bool {
	return slices.Contains(slice, value)
}
