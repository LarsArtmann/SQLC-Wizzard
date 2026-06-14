package templates

import (
	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/validation"
	"github.com/LarsArtmann/SQLC-Wizzard/pkg/config"
	"github.com/samber/lo"
)

// TransformSafetyRulesToConfig converts generated safety rules to config rule configs.
// This centralized function eliminates duplicated rule transformation code across templates.
func TransformSafetyRulesToConfig(safetyRules *generated.SafetyRules) []config.RuleConfig {
	transformer := validation.NewRuleTransformer()
	rules := transformer.TransformSafetyRules(safetyRules)

	return lo.Map(rules, func(r generated.RuleConfig, _ int) config.RuleConfig {
		return config.RuleConfig{
			Name:    r.Name,
			Rule:    r.Rule,
			Message: r.Message,
		}
	})
}
