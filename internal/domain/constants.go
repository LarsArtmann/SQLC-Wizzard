package domain

// Constants for safety rules and configuration defaults.
// These values are used across the domain layer for consistent
// validation and configuration defaults.

const (
	// MaxRowsWithoutLimitDefault is the default maximum number of rows
	// that can be returned when no LIMIT clause is specified.
	// This provides a soft safety limit to prevent accidental large result sets.
	MaxRowsWithoutLimitDefault = 1000

	// MaxRowsWithoutLimitProduction is a conservative limit for production environments.
	MaxRowsWithoutLimitProduction = 100

	// MaxRowsWithoutLimitDevelopment is 0 (no limit) for development environments.
	MaxRowsWithoutLimitDevelopment = 0
)
