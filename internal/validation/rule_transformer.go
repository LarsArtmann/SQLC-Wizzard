package validation

import (
	"strconv"

	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/domain"
)

// RuleTransformer consolidates rule transformation logic from duplicate implementations
// This eliminates the split brain between generated/types.go and internal/domain/
// and provides a single source of truth for rule configuration generation.
type RuleTransformer struct{}

// NewRuleTransformer creates a new rule transformer.
func NewRuleTransformer() *RuleTransformer {
	return &RuleTransformer{}
}

// TransformSafetyRules converts safety rules to configuration format
// This is the single source of truth for rule transformation logic.
func (rt *RuleTransformer) TransformSafetyRules(rules *generated.SafetyRules) []generated.RuleConfig {
	if rules == nil {
		return []generated.RuleConfig{}
	}

	var configRules []generated.RuleConfig

	// Transform NoSelectStar rule
	if rules.NoSelectStar {
		configRules = append(configRules, generated.RuleConfig{
			Name:    "no-select-star",
			Rule:    "!query.contains('SELECT *')",
			Message: "SELECT * is not allowed",
		})
	}

	// Transform RequireWhere rule
	if rules.RequireWhere {
		configRules = append(configRules, generated.RuleConfig{
			Name:    "require-where",
			Rule:    "query.type in ('SELECT', 'UPDATE', 'DELETE') && query.hasWhereClause()",
			Message: "WHERE clause is required for this query type",
		})
	}

	// Transform RequireLimit rule
	if rules.RequireLimit {
		configRules = append(configRules, generated.RuleConfig{
			Name:    "require-limit",
			Rule:    "query.type == 'SELECT' && !query.hasLimitClause()",
			Message: "LIMIT clause is required for SELECT queries",
		})
	}

	// Transform additional rules
	for _, customRule := range rules.Rules {
		configRules = append(configRules, generated.RuleConfig(customRule))
	}

	return configRules
}

// TransformDomainSafetyRules converts domain safety rules to configuration format
// Provides compatibility with domain package while using consolidated logic
//
// DEPRECATED: Use TransformTypeSafeSafetyRules for new code. This method is kept
// for backward compatibility with the old boolean-heavy SafetyRules structure.
func (rt *RuleTransformer) TransformDomainSafetyRules(rules *domain.SafetyRules) []generated.RuleConfig {
	// Since domain.SafetyRules is an alias for generated.SafetyRules,
	// we can directly use the consolidated transformation logic
	return rt.TransformSafetyRules(rules)
}

// TransformTypeSafeSafetyRules converts NEW type-safe safety rules to configuration format
// This is the preferred method for new code as it leverages semantic groupings and enums.
func (rt *RuleTransformer) TransformTypeSafeSafetyRules(rules *domain.TypeSafeSafetyRules) []generated.RuleConfig {
	var configRules []generated.RuleConfig

	// ========== STYLE RULES (Code Quality) ==========

	// Transform SelectStarPolicy rule
	if rules.StyleRules.SelectStarPolicy.ForbidsSelectStar() {
		configRules = append(configRules, generated.RuleConfig{
			Name:    "no-select-star",
			Rule:    "!query.contains('SELECT *')",
			Message: "SELECT * is not allowed - use explicit column names",
		})
	}

	// Transform RequireExplicitColumns rule (NEW!)
	if rules.StyleRules.ColumnExplicitness.RequiresExplicitColumns() {
		configRules = append(configRules, generated.RuleConfig{
			Name:    "require-explicit-columns",
			Rule:    "query.type == 'SELECT' && query.hasExplicitColumns()",
			Message: "All columns must be explicitly named",
		})
	}

	// ========== SAFETY RULES (Prevent Bugs) ==========

	// Transform RequireWhere rule
	if rules.SafetyRules.WhereRequirement.RequiresOnDestructive() {
		configRules = append(configRules, generated.RuleConfig{
			Name:    "require-where",
			Rule:    "query.type in ('SELECT', 'UPDATE', 'DELETE') && query.hasWhereClause()",
			Message: "WHERE clause is required for SELECT/UPDATE/DELETE queries to prevent accidental full-table operations",
		})
	}

	// Transform RequireLimit rule
	if rules.SafetyRules.LimitRequirement.RequiresOnSelect() {
		configRules = append(configRules, generated.RuleConfig{
			Name:    "require-limit",
			Rule:    "query.type == 'SELECT' && !query.hasLimitClause()",
			Message: "LIMIT clause is required for SELECT queries to prevent unbounded result sets",
		})
	}

	// Transform MaxRowsWithoutLimit rule (NEW!)
	if rules.SafetyRules.MaxRowsWithoutLimit > 0 {
		limitStr := uintToString(rules.SafetyRules.MaxRowsWithoutLimit)
		configRules = append(configRules, generated.RuleConfig{
			Name:    "max-rows-without-limit",
			Rule:    "query.type == 'SELECT' && (!query.hasLimitClause() || query.limitValue() > " + limitStr + ")",
			Message: "SELECT queries without LIMIT or with LIMIT > " + limitStr + " are not allowed",
		})
	}

	// ========== DESTRUCTIVE OPERATIONS (Policy-Based) ==========

	// Transform DestructiveOps policy
	switch rules.DestructiveOps {
	case domain.DestructiveForbidden:
		// Forbid DROP TABLE
		configRules = append(configRules, generated.RuleConfig{
			Name:    "no-drop-table",
			Rule:    "!query.contains('DROP TABLE')",
			Message: "DROP TABLE is forbidden by safety policy",
		})
		// Forbid TRUNCATE
		configRules = append(configRules, generated.RuleConfig{
			Name:    "no-truncate",
			Rule:    "!query.contains('TRUNCATE')",
			Message: "TRUNCATE is forbidden by safety policy",
		})

	case domain.DestructiveWithConfirmation:
		// Require confirmation for DROP TABLE
		configRules = append(configRules, generated.RuleConfig{
			Name:    "drop-table-requires-confirmation",
			Rule:    "query.contains('DROP TABLE') && query.hasComment('CONFIRMED')",
			Message: "DROP TABLE requires explicit confirmation (add comment: -- CONFIRMED)",
		})
		// Require confirmation for TRUNCATE
		configRules = append(configRules, generated.RuleConfig{
			Name:    "truncate-requires-confirmation",
			Rule:    "query.contains('TRUNCATE') && query.hasComment('CONFIRMED')",
			Message: "TRUNCATE requires explicit confirmation (add comment: -- CONFIRMED)",
		})

	case domain.DestructiveAllowed:
		// No rules needed - destructive operations are allowed
	}

	// ========== CUSTOM RULES ==========

	// Transform custom CEL rules
	for _, customRule := range rules.CustomRules {
		configRules = append(configRules, generated.RuleConfig(customRule))
	}

	return configRules
}

// uintToString converts uint to string for rule generation.
func uintToString(n uint) string {
	return strconv.FormatUint(uint64(n), 10)
}
