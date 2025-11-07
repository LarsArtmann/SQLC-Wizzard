// Package domain provides domain models for business logic
package domain

import (
	"strings"
	"github.com/LarsArtmann/SQLC-Wizzard/generated"
)

// EmitOptions is a type alias for the generated EmitOptions
// This ensures we use the single source of truth from TypeSpec generation
type EmitOptions = generated.EmitOptions

// SafetyRules represents CEL-based validation rules
type SafetyRules struct {
	NoSelectStar bool   `json:"no_select_star"`
	RequireWhere  bool   `json:"require_where"`
	NoDropTable  bool   `json:"no_drop_table"`
	NoTruncate   bool   `json:"no_truncate"`
	RequireLimit  bool   `json:"require_limit"`
	// Additional rules as needed
	Rules []string `json:"rules"`
}

// ToRuleConfigs converts safety rules to configuration format
func (s *SafetyRules) ToRuleConfigs() []RuleConfig {
	var rules []RuleConfig
	
	if s.NoSelectStar {
		rules = append(rules, RuleConfig{
			Name:  "no-select-star",
			Rule:  "!query.contains('SELECT *')",
			Message: "SELECT * is not allowed",
		})
	}
	
	if s.RequireWhere {
		rules = append(rules, RuleConfig{
			Name:  "require-where", 
			Rule:  "query.type in ('SELECT', 'UPDATE', 'DELETE') && query.hasWhereClause()",
			Message: "WHERE clause is required for this query type",
		})
	}
	
	if s.RequireLimit {
		rules = append(rules, RuleConfig{
			Name:  "require-limit",
			Rule:  "query.type == 'SELECT' && !query.hasLimitClause()",
			Message: "LIMIT clause is required for SELECT queries",
		})
	}
	
	// Add custom rules
	for _, rule := range s.Rules {
		rules = append(rules, RuleConfig{
			Name:    "custom-rule",
			Rule:    rule,
			Message: "Custom validation rule triggered",
		})
	}
	
	return rules
}

// RuleConfig is a type alias for the generated RuleConfig
// This ensures we use the single source of truth from TypeSpec generation
type RuleConfig = generated.RuleConfig

// DefaultEmitOptions returns safe defaults for code generation
func DefaultEmitOptions() EmitOptions {
	return EmitOptions{
		EmitJSONTags:           true,
		EmitPreparedQueries:    true,
		EmitInterface:          true,
		EmitEmptySlices:        true,
		EmitResultStructPointers: false,
		EmitParamsStructPointers: false,
		EmitEnumValidMethod:    true,
		EmitAllEnumValues:      true,
		JSONTagsCaseStyle:      "camel",
	}
}

// DefaultSafetyRules returns safe defaults for query validation
func DefaultSafetyRules() SafetyRules {
	return SafetyRules{
		NoSelectStar: true,
		RequireWhere:  true,
		NoDropTable:  true,
		NoTruncate:   true,
		RequireLimit:  false, // Not too restrictive by default
		Rules:         []string{},
	}
}

// Validate checks if safety rules are valid
func (s *SafetyRules) Validate() error {
	var errors []string
	
	if s.NoSelectStar && !s.RequireWhere {
		errors = append(errors, "no_select_star requires require_where")
	}
	
	// Validate custom rule syntax
	for _, rule := range s.Rules {
		if strings.TrimSpace(rule) == "" {
			errors = append(errors, "empty custom rule found")
		}
		if len(rule) > 1000 {
			errors = append(errors, "custom rule too long (>1000 chars)")
		}
	}
	
	if len(errors) > 0 {
		return &ValidationError{Errors: errors}
	}
	
	return nil
}

// ValidationError represents domain validation errors
type ValidationError struct {
	Errors []string
}

func (v *ValidationError) Error() string {
	return "validation failed: " + strings.Join(v.Errors, ", ")
}