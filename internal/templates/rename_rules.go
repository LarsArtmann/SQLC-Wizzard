package templates

// CommonRenameRules returns common rename rules for better Go naming.
// Used as a centralized, single source of truth for sqlc rename mappings.
func CommonRenameRules() map[string]string {
	return map[string]string{
		"id":   "ID",
		"uuid": "UUID",
		"url":  "URL",
		"uri":  "URI",
		"api":  "API",
		"http": "HTTP",
		"json": "JSON",
		"db":   "DB",
		"otp":  "OTP",
	}
}
