package validation

import (
	"strconv"

	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/domain"
)

// RuleTransformer consolidates rule transformation logic from duplicate implementations.
// This eliminates the split brain between generated/types.go and internal/domain/
// and provides a single source of truth for rule configuration generation.
type RuleTransformer struct{}

// NewRuleTransformer creates a new rule transformer.
func NewRuleTransformer() *RuleTransformer {
	return &RuleTransformer{}
}

// ruleAppender is a helper to reduce duplication in rule creation patterns.
type ruleAppender struct {
	rules []generated.RuleConfig
}

// add appends a new rule with the given name, rule expression, and message.
func (a *ruleAppender) add(name, rule, message string) {
	a.rules = append(a.rules, generated.RuleConfig{
		Name:    name,
		Rule:    rule,
		Message: message,
	})
}

// addSafetyRules appends safety rules to the rule appender.
func (a *ruleAppender) addSafetyRules(safetyRules []generated.SafetyRule) {
	for _, rule := range safetyRules {
		a.rules = append(a.rules, generated.RuleConfig(rule))
	}
}

// TransformSafetyRules converts safety rules to configuration format
// This is the single source of truth for rule transformation logic.
func (rt *RuleTransformer) TransformSafetyRules(
	rules *generated.SafetyRules,
) []generated.RuleConfig {
	if rules == nil {
		return []generated.RuleConfig{}
	}

	a := &ruleAppender{}

	if rules.NoSelectStar {
		a.add("no-select-star", "!query.contains('SELECT *')", "SELECT * is not allowed")
	}

	if rules.RequireWhere {
		a.add("require-where", "query.type in ('SELECT', 'UPDATE', 'DELETE') && !query.hasWhereClause()", "WHERE clause is required for this query type")
	}

	if rules.RequireLimit {
		a.add("require-limit", "query.type == 'SELECT' && !query.hasLimitClause()", "LIMIT clause is required for SELECT queries")
	}

	a.addSafetyRules(rules.Rules)

	return a.rules
}

// TransformDomainSafetyRules converts domain safety rules to configuration format
// Provides compatibility with domain package while using consolidated logic
//
// DEPRECATED: Use TransformTypeSafeSafetyRules for new code. This method is kept
// for backward compatibility with the old boolean-heavy SafetyRules structure.
func (rt *RuleTransformer) TransformDomainSafetyRules(
	rules *domain.SafetyRules,
) []generated.RuleConfig {
	// Since domain.SafetyRules is an alias for generated.SafetyRules,
	// we can directly use the consolidated transformation logic
	return rt.TransformSafetyRules(rules)
}

// TransformTypeSafeSafetyRules converts NEW type-safe safety rules to configuration format
// This is the preferred method for new code as it leverages semantic groupings and enums.
func (rt *RuleTransformer) TransformTypeSafeSafetyRules(
	rules *domain.TypeSafeSafetyRules,
) []generated.RuleConfig {
	a := &ruleAppender{}

	// ========== STYLE RULES (Code Quality) ==========

	if rules.StyleRules.SelectStarPolicy.ForbidsSelectStar() {
		a.add("no-select-star", "!query.contains('SELECT *')", "SELECT * is not allowed - use explicit column names")
	}

	if rules.StyleRules.ColumnExplicitness.RequiresExplicitColumns() {
		a.add("require-explicit-columns", "query.type == 'SELECT' && query.hasExplicitColumns()", "All columns must be explicitly named")
	}

	// ========== SAFETY RULES (Prevent Bugs) ==========

	if rules.SafetyRules.WhereRequirement.RequiresOnDestructive() {
		a.add("require-where", "query.type in ('SELECT', 'UPDATE', 'DELETE') && !query.hasWhereClause()", "WHERE clause is required for SELECT/UPDATE/DELETE queries to prevent accidental full-table operations")
	}

	if rules.SafetyRules.LimitRequirement.RequiresOnSelect() {
		a.add("require-limit", "query.type == 'SELECT' && !query.hasLimitClause()", "LIMIT clause is required for SELECT queries to prevent unbounded result sets")
	}

	if rules.SafetyRules.MaxRowsWithoutLimit > 0 {
		limitStr := uintToString(rules.SafetyRules.MaxRowsWithoutLimit)
		a.add("max-rows-without-limit", "query.type == 'SELECT' && (!query.hasLimitClause() || query.limitValue() > "+limitStr+")", "SELECT queries without LIMIT or with LIMIT > "+limitStr+" are not allowed")
	}

	// ========== DESTRUCTIVE OPERATIONS (Policy-Based) ==========

	switch rules.DestructiveOps {
	case domain.DestructiveForbidden:
		a.add("no-drop-table", "!query.contains('DROP TABLE')", "DROP TABLE is forbidden by safety policy")
		a.add("no-truncate", "!query.contains('TRUNCATE')", "TRUNCATE is forbidden by safety policy")

	case domain.DestructiveWithConfirmation:
		a.add("drop-table-requires-confirmation", "query.contains('DROP TABLE') && query.hasComment('CONFIRMED')", "DROP TABLE requires explicit confirmation (add comment: -- CONFIRMED)")
		a.add("truncate-requires-confirmation", "query.contains('TRUNCATE') && query.hasComment('CONFIRMED')", "TRUNCATE requires explicit confirmation (add comment: -- CONFIRMED)")

	case domain.DestructiveAllowed:
	}

	// ========== CUSTOM RULES ==========

	a.addSafetyRules(rules.CustomRules)

	return a.rules
}

// uintToString converts uint to string for rule generation.
func uintToString(n uint) string {
	return strconv.FormatUint(uint64(n), 10)
}
