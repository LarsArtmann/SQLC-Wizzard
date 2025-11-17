package validation

import (
	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/domain"
)

// RuleTransformer consolidates rule transformation logic from duplicate implementations
// This eliminates the split brain between generated/types.go and internal/domain/
// and provides a single source of truth for rule configuration generation
type RuleTransformer struct{}

// NewRuleTransformer creates a new rule transformer
func NewRuleTransformer() *RuleTransformer {
	return &RuleTransformer{}
}

// TransformSafetyRules converts safety rules to configuration format
// This is the single source of truth for rule transformation logic
func (rt *RuleTransformer) TransformSafetyRules(rules *generated.SafetyRules) []generated.RuleConfig {
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
		configRules = append(configRules, generated.RuleConfig{
			Name:    customRule.Name,
			Rule:    customRule.Rule,
			Message: customRule.Message,
		})
	}

	return configRules
}

// TransformDomainSafetyRules converts domain safety rules to configuration format
// Provides compatibility with domain package while using consolidated logic
func (rt *RuleTransformer) TransformDomainSafetyRules(rules *domain.SafetyRules) []generated.RuleConfig {
	// Since domain.SafetyRules is an alias for generated.SafetyRules,
	// we can directly use the consolidated transformation logic
	generatedRules := (*generated.SafetyRules)(rules)
	return rt.TransformSafetyRules(generatedRules)
}
