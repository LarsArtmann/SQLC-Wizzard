// Package domain provides domain models for business logic
package domain

import "strings"

// EmitOptions defines SQL code generation options
type EmitOptions struct {
	EmitJSONTags                bool   `json:"emit_json_tags"`
	EmitPreparedQueries         bool   `json:"emit_prepared_queries"`
	EmitInterface              bool   `json:"emit_interface"`
	EmitEmptySlices            bool   `json:"emit_empty_slices"`
	EmitResultStructPointers    bool   `json:"emit_result_struct_pointers"`
	EmitParamsStructPointers    bool   `json:"emit_params_struct_pointers"`
	EmitEnumValidMethod        bool   `json:"emit_enum_valid_method"`
	EmitAllEnumValues         bool   `json:"emit_all_enum_values"`
	JSONTagsCaseStyle          string `json:"json_tags_case_style"`
}

// ApplyToGoGenConfig applies emit options to a GoGenConfig
// This eliminates field-by-field copying (DRY principle)
func (e *EmitOptions) ApplyToGoGenConfig(cfg interface{}) {
	// Type assertion for flexibility
	type HasEmitJSONTags interface{ SetEmitJSONTags(bool) }
	type HasEmitPreparedQueries interface{ SetEmitPreparedQueries(bool) }
	type HasEmitInterface interface{ SetEmitInterface(bool) }
	// ... more as needed
	
	if has, ok := cfg.(HasEmitJSONTags); ok {
		has.SetEmitJSONTags(e.EmitJSONTags)
	}
	if has, ok := cfg.(HasEmitPreparedQueries); ok {
		has.SetEmitPreparedQueries(e.EmitPreparedQueries)
	}
	if has, ok := cfg.(HasEmitInterface); ok {
		has.SetEmitInterface(e.EmitInterface)
	}
}

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

// RuleConfig represents a validation rule configuration
type RuleConfig struct {
	Name    string `json:"name"`
	Rule    string `json:"rule"`
	Message string `json:"message"`
}

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