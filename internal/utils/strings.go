package utils

import (
	"strings"
	"unicode"
)

// StringToCamelCase converts snake_case to CamelCase
func StringToCamelCase(s string) string {
	if s == "" {
		return ""
	}

	words := strings.Split(s, "_")
	for i, word := range words {
		if word == "" {
			continue
		}
		if len(word) > 0 {
			words[i] = strings.ToUpper(word[:1]) + strings.ToLower(word[1:])
		}
	}

	return strings.Join(words, "")
}

// stringToCase converts CamelCase to specified case with separator
func stringToCase(s string, separator rune) string {
	if s == "" {
		return ""
	}

	var result []rune
	for i, r := range s {
		if unicode.IsUpper(r) && i > 0 && !unicode.IsUpper(rune(s[i-1])) {
			result = append(result, separator)
		}
		result = append(result, unicode.ToLower(r))
	}

	return string(result)
}

// StringToSnakeCase converts CamelCase to snake_case
func StringToSnakeCase(s string) string {
	return stringToCase(s, '_')
}

// StringToKebabCase converts CamelCase to kebab-case
func StringToKebabCase(s string) string {
	return stringToCase(s, '-')
}

// Pluralize converts a noun to its plural form (simplified)
func Pluralize(word string) string {
	if word == "" {
		return ""
	}

	// Handle common irregular nouns
	switch strings.ToLower(word) {
	case "person":
		return "people"
	case "child":
		return "children"
	case "index":
		return "indices"
	case "status":
		return "statuses"
	case "match":
		return "matches"
	case "box":
		return "boxes"
	case "sheep":
		return "sheep"
	case "data":
		return "data" // Already plural or irregular
	}

	// Handle words ending in specific patterns
	lowered := strings.ToLower(word)
	if strings.HasSuffix(lowered, "s") {
		return word // Already plural
	}

	if strings.HasSuffix(lowered, "y") {
		return word[:len(word)-1] + "ies"
	}

	if strings.HasSuffix(lowered, "s") || strings.HasSuffix(lowered, "sh") ||
		strings.HasSuffix(lowered, "ch") || strings.HasSuffix(lowered, "x") ||
		strings.HasSuffix(lowered, "z") {
		return word + "es"
	}

	return word + "s"
}

// Singularize converts a plural noun to its singular form (simplified)
func Singularize(word string) string {
	if word == "" {
		return ""
	}

	lowered := strings.ToLower(word)

	// Handle common irregular nouns
	switch lowered {
	case "people":
		return "person"
	case "children":
		return "child"
	case "indices":
		return "index"
	case "statuses":
		return "status"
	case "matches":
		return "match"
	case "boxes":
		return "box"
	case "data":
		return "data" // Already singular or irregular
	case "sheep":
		return "sheep"
	}

	// Handle words ending in specific patterns
	if strings.HasSuffix(lowered, "ies") {
		return word[:len(word)-3] + "y"
	}

	if strings.HasSuffix(lowered, "es") {
		return word[:len(word)-2]
	}

	if strings.HasSuffix(lowered, "s") {
		return word[:len(word)-1]
	}

	return word
}

// IsValidIdentifier checks if a string is a valid Go identifier
func IsValidIdentifier(s string) bool {
	if s == "" {
		return false
	}

	for i, r := range s {
		if i == 0 {
			if !unicode.IsLetter(r) && r != '_' {
				return false
			}
		} else {
			if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' {
				return false
			}
		}
	}

	return true
}

// EscapeSQLIdentifier escapes a SQL identifier with double quotes
func EscapeSQLIdentifier(s string) string {
	return `"` + s + `"`
}

// GetFileExtension returns the file extension including the dot
func GetFileExtension(path string) string {
	dot := strings.LastIndex(path, ".")
	if dot == -1 || dot == len(path)-1 {
		return ""
	}

	return path[dot:]
}

// HasExtension checks if a path has a specific extension
func HasExtension(path, ext string) bool {
	if !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}

	return strings.HasSuffix(strings.ToLower(path), strings.ToLower(ext))
}

// TruncateString truncates a string to a specific length
func TruncateString(s string, length int) string {
	if len(s) <= length {
		return s
	}

	if length <= 3 {
		return s[:length]
	}

	return s[:length-3] + "..."
}

// PadString pads a string to a specific length with the given character
func PadString(s string, length int, char string) string {
	if len(s) >= length {
		return s
	}

	if char == "" {
		char = " "
	}

	padding := strings.Repeat(char, length-len(s))
	return s + padding
}
