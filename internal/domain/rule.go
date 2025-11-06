package domain

import "github.com/LarsArtmann/SQLC-Wizzard/pkg/config"

// SafetyRules represents the boolean flags for enabling validation rules.
// This is the SINGLE SOURCE OF TRUTH for safety rules.
//
// Use ToRuleConfigs() to convert to sqlc's CEL-based RuleConfig format.
type SafetyRules struct {
	NoSelectStar bool
	RequireWhere bool
	RequireLimit bool
	NoDropTable  bool
}

// ToRuleConfigs converts SafetyRules to sqlc's RuleConfig format.
// This is the conversion method that eliminates split brain duplication.
//
// Each boolean flag is converted to a CEL (Common Expression Language) rule
// that sqlc can validate at compile time.
func (sr SafetyRules) ToRuleConfigs() []config.RuleConfig {
	var rules []config.RuleConfig

	if sr.NoSelectStar {
		rules = append(rules, config.RuleConfig{
			Name:    "no-select-star",
			Rule:    `!query.sql.contains("SELECT *")`,
			Message: "SELECT * is not allowed for security and performance reasons",
		})
	}

	if sr.RequireWhere {
		rules = append(rules, config.RuleConfig{
			Name:    "require-where-delete",
			Rule:    `query.cmd != "exec" || !query.sql.contains("DELETE") || query.sql.contains("WHERE")`,
			Message: "DELETE statements must include a WHERE clause",
		})
	}

	if sr.NoDropTable {
		rules = append(rules, config.RuleConfig{
			Name:    "no-drop-table",
			Rule:    `!query.sql.contains("DROP TABLE")`,
			Message: "DROP TABLE statements are not allowed",
		})
	}

	if sr.RequireLimit {
		rules = append(rules, config.RuleConfig{
			Name:    "require-limit-select",
			Rule:    `query.cmd == "many" implies query.sql.contains("LIMIT")`,
			Message: "SELECT queries that return multiple rows should include a LIMIT clause",
		})
	}

	return rules
}

// DefaultSafetyRules returns recommended safety rules for production use.
func DefaultSafetyRules() SafetyRules {
	return SafetyRules{
		NoSelectStar: true,
		RequireWhere: true,
		RequireLimit: false, // Optional - can be noisy
		NoDropTable:  true,
	}
}
