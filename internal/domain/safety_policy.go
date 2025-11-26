package domain

import (
	"strconv"

	"github.com/LarsArtmann/SQLC-Wizzard/generated"
)

// This file contains type-safe safety policy configuration
// Replaces the boolean-heavy SafetyRules from generated/types.go with
// semantic groupings that clarify the purpose of each rule category

// SelectStarPolicy defines how SELECT * queries are handled
// Replaces NoSelectStar bool with type-safe policy that prevents invalid states
type SelectStarPolicy string

const (
	// SelectStarAllowed permits SELECT * queries
	SelectStarAllowed SelectStarPolicy = "allowed"

	// SelectStarForbidden blocks all SELECT * queries
	SelectStarForbidden SelectStarPolicy = "forbidden"

	// SelectStarExplicit requires explicit column listing in all cases
	// (stricter than forbidden - affects SELECT col1, col2 too)
	SelectStarExplicit SelectStarPolicy = "explicit"
)

// IsValid returns true if select star policy is recognized
func (s SelectStarPolicy) IsValid() bool {
	switch s {
	case SelectStarAllowed, SelectStarForbidden, SelectStarExplicit:
		return true
	default:
		return false
	}
}

// String returns string representation of policy
func (s SelectStarPolicy) String() string {
	return string(s)
}

// ForbidsSelectStar returns true if this policy forbids SELECT * queries
func (s SelectStarPolicy) ForbidsSelectStar() bool {
	return s == SelectStarForbidden || s == SelectStarExplicit
}

// RequiresExplicitColumns returns true if this policy requires all columns to be named
func (s SelectStarPolicy) RequiresExplicitColumns() bool {
	return s == SelectStarExplicit
}

// ColumnExplicitnessPolicy defines how explicit column selection is enforced
type ColumnExplicitnessPolicy string

const (
	// ColumnExplicitnessDefault allows any valid column selection
	ColumnExplicitnessDefault ColumnExplicitnessPolicy = "default"

	// ColumnExplicitnessRequired requires all columns to be explicitly named
	// (affects even SELECT col1, col2 queries)
	ColumnExplicitnessRequired ColumnExplicitnessPolicy = "required"

	// ColumnExplicitnessNamed requires columns to have aliases
	ColumnExplicitnessNamed ColumnExplicitnessPolicy = "named"
)

// IsValid returns true if column explicitness policy is recognized
func (c ColumnExplicitnessPolicy) IsValid() bool {
	switch c {
	case ColumnExplicitnessDefault, ColumnExplicitnessRequired, ColumnExplicitnessNamed:
		return true
	default:
		return false
	}
}

// String returns string representation of policy
func (c ColumnExplicitnessPolicy) String() string {
	return string(c)
}

// QueryStyleRules defines rules about SQL query style and best practices
// Uses type-safe enums to prevent invalid state combinations
type QueryStyleRules struct {
	// SelectStarPolicy defines how SELECT * queries are handled
	SelectStarPolicy SelectStarPolicy `json:"selectStarPolicy"`

	// ColumnExplicitness defines how explicit column selection is enforced
	ColumnExplicitness ColumnExplicitnessPolicy `json:"columnExplicitness"`
}

// WhereClauseRequirement defines when WHERE clauses are required
// Replaces RequireWhere bool with type-safe policy that prevents invalid states
type WhereClauseRequirement string

const (
	// WhereClauseNever allows queries without WHERE clauses
	WhereClauseNever WhereClauseRequirement = "never"

	// WhereClauseOnDestructive requires WHERE on UPDATE/DELETE operations only
	WhereClauseOnDestructive WhereClauseRequirement = "destructive"

	// WhereClauseOnSelect requires WHERE on SELECT queries with no LIMIT
	WhereClauseOnSelect WhereClauseRequirement = "select_unlimited"

	// WhereClauseAlways requires WHERE clauses on all SELECT/UPDATE/DELETE queries
	// (strictest setting for production safety)
	WhereClauseAlways WhereClauseRequirement = "always"
)

// IsValid returns true if where clause requirement is recognized
func (w WhereClauseRequirement) IsValid() bool {
	switch w {
	case WhereClauseNever, WhereClauseOnDestructive, WhereClauseOnSelect, WhereClauseAlways:
		return true
	default:
		return false
	}
}

// String returns string representation of requirement
func (w WhereClauseRequirement) String() string {
	return string(w)
}

// RequiresOnSelect returns true if this policy requires WHERE on SELECT queries
func (w WhereClauseRequirement) RequiresOnSelect() bool {
	return w == WhereClauseAlways || w == WhereClauseOnSelect
}

// RequiresOnDestructive returns true if this policy requires WHERE on UPDATE/DELETE
func (w WhereClauseRequirement) RequiresOnDestructive() bool {
	return w == WhereClauseAlways || w == WhereClauseOnDestructive
}

// LimitClauseRequirement defines when LIMIT clauses are required
// Replaces RequireLimit bool with type-safe policy for fine-grained control
type LimitClauseRequirement string

const (
	// LimitClauseNever allows queries without LIMIT clauses
	LimitClauseNever LimitClauseRequirement = "never"

	// LimitClauseOnSelect requires LIMIT on SELECT queries
	LimitClauseOnSelect LimitClauseRequirement = "select"

	// LimitClauseOnSelectWithoutWhere requires LIMIT only when no WHERE clause
	LimitClauseOnSelectWithoutWhere LimitClauseRequirement = "select_without_where"

	// LimitClauseAlways requires LIMIT clauses on all SELECT queries
	// (strictest setting for production environments)
	LimitClauseAlways LimitClauseRequirement = "always"
)

// IsValid returns true if limit clause requirement is recognized
func (l LimitClauseRequirement) IsValid() bool {
	switch l {
	case LimitClauseNever, LimitClauseOnSelect, LimitClauseOnSelectWithoutWhere, LimitClauseAlways:
		return true
	default:
		return false
	}
}

// String returns string representation of requirement
func (l LimitClauseRequirement) String() string {
	return string(l)
}

// RequiresOnSelect returns true if this policy requires LIMIT on SELECT queries
func (l LimitClauseRequirement) RequiresOnSelect() bool {
	return l == LimitClauseOnSelect || l == LimitClauseOnSelectWithoutWhere || l == LimitClauseAlways
}

// RequiresWithoutWhere returns true if this policy requires LIMIT only without WHERE
func (l LimitClauseRequirement) RequiresWithoutWhere() bool {
	return l == LimitClauseOnSelectWithoutWhere
}

// QuerySafetyRules defines rules that prevent common SQL safety issues
// Uses type-safe enums to prevent invalid state combinations
type QuerySafetyRules struct {
	// WhereRequirement defines when WHERE clauses are required
	WhereRequirement WhereClauseRequirement `json:"whereRequirement"`

	// LimitRequirement defines when LIMIT clauses are required
	LimitRequirement LimitClauseRequirement `json:"limitRequirement"`

	// MaxRowsWithoutLimit specifies max rows that can be returned without explicit LIMIT
	// 0 means no limit, >0 means enforce this limit if no LIMIT clause present
	MaxRowsWithoutLimit uint `json:"maxRowsWithoutLimit"`
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
				Message: "Custom rule at index " + strconv.Itoa(i) + " has empty name",
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
			SelectStarPolicy:   SelectStarForbidden,       // Enforce explicit columns
			ColumnExplicitness: ColumnExplicitnessDefault, // Not too strict by default
		},
		SafetyRules: QuerySafetyRules{
			WhereRequirement:    WhereClauseOnDestructive, // Require WHERE on UPDATE/DELETE only
			LimitRequirement:    LimitClauseNever,         // Not enforced by default (too restrictive)
			MaxRowsWithoutLimit: 1000,                     // Soft limit when no LIMIT clause
		},
		DestructiveOps: DestructiveForbidden, // Production-safe default
		CustomRules:    []generated.SafetyRule{},
	}
}

// NewDevelopmentSafetyRules returns relaxed rules for development environments
func NewDevelopmentSafetyRules() TypeSafeSafetyRules {
	return TypeSafeSafetyRules{
		StyleRules: QueryStyleRules{
			SelectStarPolicy:   SelectStarAllowed, // Allow for rapid prototyping
			ColumnExplicitness: ColumnExplicitnessDefault,
		},
		SafetyRules: QuerySafetyRules{
			WhereRequirement:    WhereClauseNever, // Allow flexible querying
			LimitRequirement:    LimitClauseNever, // Not enforced in development
			MaxRowsWithoutLimit: 0,                // No limit enforcement
		},
		DestructiveOps: DestructiveAllowed, // Allow migrations and schema changes
		CustomRules:    []generated.SafetyRule{},
	}
}

// NewProductionSafetyRules returns strict rules for production environments
func NewProductionSafetyRules() TypeSafeSafetyRules {
	return TypeSafeSafetyRules{
		StyleRules: QueryStyleRules{
			SelectStarPolicy:   SelectStarExplicit,         // Extra strict
			ColumnExplicitness: ColumnExplicitnessRequired, // Extra strict
		},
		SafetyRules: QuerySafetyRules{
			WhereRequirement:    WhereClauseAlways, // Always require WHERE
			LimitRequirement:    LimitClauseAlways, // Always require LIMIT
			MaxRowsWithoutLimit: 100,               // Very conservative
		},
		DestructiveOps: DestructiveForbidden, // Never allow destructive ops
		CustomRules:    []generated.SafetyRule{},
	}
}
