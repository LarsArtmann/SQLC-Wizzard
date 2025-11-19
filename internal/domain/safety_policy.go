package domain

import "github.com/LarsArtmann/SQLC-Wizzard/generated"

// This file contains type-safe safety policy configuration
// Replaces the boolean-heavy SafetyRules from generated/types.go with
// semantic groupings that clarify the purpose of each rule category

// QueryStyleRules defines rules about SQL query style and best practices
type QueryStyleRules struct {
	// NoSelectStar prevents SELECT * queries (forces explicit column selection)
	NoSelectStar bool

	// RequireExplicitColumns requires all columns to be explicitly named
	// (stricter than NoSelectStar, affects even SELECT col1, col2)
	RequireExplicitColumns bool
}

// QuerySafetyRules defines rules that prevent common SQL safety issues
type QuerySafetyRules struct {
	// RequireWhere requires WHERE clauses on SELECT, UPDATE, DELETE queries
	RequireWhere bool

	// RequireLimit requires LIMIT clauses on SELECT queries
	RequireLimit bool

	// MaxRowsWithoutLimit specifies max rows that can be returned without explicit LIMIT
	// 0 means no limit, >0 means enforce this limit if no LIMIT clause present
	MaxRowsWithoutLimit uint
}

// DestructiveOperationPolicy defines how to handle destructive database operations
// This is an ENUM that replaces multiple boolean flags with clear policy states
type DestructiveOperationPolicy string

const (
	// DestructiveAllowed permits all destructive operations (DROP, TRUNCATE, etc.)
	DestructiveAllowed DestructiveOperationPolicy = "allowed"

	// DestructiveWithConfirmation requires explicit confirmation for destructive ops
	// (useful for production environments with safeguards)
	DestructiveWithConfirmation DestructiveOperationPolicy = "with_confirmation"

	// DestructiveForbidden completely blocks all destructive operations
	// (recommended for production code generation)
	DestructiveForbidden DestructiveOperationPolicy = "forbidden"
)

// IsValid returns true if the destructive operation policy is recognized
func (d DestructiveOperationPolicy) IsValid() bool {
	switch d {
	case DestructiveAllowed, DestructiveWithConfirmation, DestructiveForbidden:
		return true
	default:
		return false
	}
}

// String returns the string representation of the policy
func (d DestructiveOperationPolicy) String() string {
	return string(d)
}

// AllowsDropTable returns true if DROP TABLE is allowed under this policy
func (d DestructiveOperationPolicy) AllowsDropTable() bool {
	return d == DestructiveAllowed
}

// AllowsTruncate returns true if TRUNCATE is allowed under this policy
func (d DestructiveOperationPolicy) AllowsTruncate() bool {
	return d == DestructiveAllowed
}

// RequiresConfirmation returns true if destructive ops need confirmation
func (d DestructiveOperationPolicy) RequiresConfirmation() bool {
	return d == DestructiveWithConfirmation
}

// TypeSafeSafetyRules represents type-safe query validation rules
// This is the NEW type-safe version that replaces the boolean-heavy generated.SafetyRules
// with semantic groupings that clarify the purpose of each rule
//
// Migration path: Use this for new code, gradually migrate existing code from SafetyRules
type TypeSafeSafetyRules struct {
	// StyleRules defines rules about query style and best practices
	StyleRules QueryStyleRules

	// SafetyRules defines rules that prevent common safety issues
	SafetyRules QuerySafetyRules

	// DestructiveOps defines the policy for destructive operations
	DestructiveOps DestructiveOperationPolicy

	// CustomRules contains additional CEL-based validation rules
	CustomRules []generated.SafetyRule
}

// IsValid validates that all safety rules are valid
func (t *TypeSafeSafetyRules) IsValid() error {
	if !t.DestructiveOps.IsValid() {
		return &DomainValidationError{
			Field:   "DestructiveOps",
			Message: "Invalid destructive operation policy: " + string(t.DestructiveOps),
		}
	}

	// Validate custom rules
	for i, rule := range t.CustomRules {
		if rule.Name == "" {
			return &DomainValidationError{
				Field:   "CustomRules",
				Message: "Custom rule at index " + string(rune(i)) + " has empty name",
			}
		}
		if rule.Rule == "" {
			return &DomainValidationError{
				Field:   "CustomRules",
				Message: "Custom rule '" + rule.Name + "' has empty rule expression",
			}
		}
	}

	return nil
}

// NewTypeSafeSafetyRules returns production-ready defaults for safety rules
// This creates the NEW type-safe version with semantic groupings
func NewTypeSafeSafetyRules() TypeSafeSafetyRules {
	return TypeSafeSafetyRules{
		StyleRules: QueryStyleRules{
			NoSelectStar:           true, // Enforce explicit columns
			RequireExplicitColumns: false, // Not too strict by default
		},
		SafetyRules: QuerySafetyRules{
			RequireWhere:        true, // Prevent accidental full-table operations
			RequireLimit:        false, // Not enforced by default (too restrictive)
			MaxRowsWithoutLimit: 1000, // Soft limit when no LIMIT clause
		},
		DestructiveOps: DestructiveForbidden, // Production-safe default
		CustomRules:    []generated.SafetyRule{},
	}
}

// NewDevelopmentSafetyRules returns relaxed rules for development environments
func NewDevelopmentSafetyRules() TypeSafeSafetyRules {
	return TypeSafeSafetyRules{
		StyleRules: QueryStyleRules{
			NoSelectStar:           false, // Allow for rapid prototyping
			RequireExplicitColumns: false,
		},
		SafetyRules: QuerySafetyRules{
			RequireWhere:        false, // Allow flexible querying
			RequireLimit:        false,
			MaxRowsWithoutLimit: 0, // No limit enforcement
		},
		DestructiveOps: DestructiveAllowed, // Allow migrations and schema changes
		CustomRules:    []generated.SafetyRule{},
	}
}

// NewProductionSafetyRules returns strict rules for production environments
func NewProductionSafetyRules() TypeSafeSafetyRules {
	return TypeSafeSafetyRules{
		StyleRules: QueryStyleRules{
			NoSelectStar:           true,
			RequireExplicitColumns: true, // Extra strict
		},
		SafetyRules: QuerySafetyRules{
			RequireWhere:        true,
			RequireLimit:        true, // Always require LIMIT
			MaxRowsWithoutLimit: 100, // Very conservative
		},
		DestructiveOps: DestructiveForbidden, // Never allow destructive ops
		CustomRules:    []generated.SafetyRule{},
	}
}
