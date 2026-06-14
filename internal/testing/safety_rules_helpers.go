package testing

import (
	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/domain"
)

// newTypeSafeSafetyRules creates a base TypeSafeSafetyRules structure.
// This helper eliminates duplicate struct initialization code across helper functions.
func newTypeSafeSafetyRules() *domain.TypeSafeSafetyRules {
	return &domain.TypeSafeSafetyRules{
		StyleRules: domain.QueryStyleRules{
			SelectStarPolicy:   domain.SelectStarForbidden,
			ColumnExplicitness: domain.ColumnExplicitnessDefault,
		},
		SafetyRules: domain.QuerySafetyRules{
			WhereRequirement:    domain.WhereClauseOnDestructive,
			LimitRequirement:    domain.LimitClauseOnSelect,
			MaxRowsWithoutLimit: domain.MaxRowsWithoutLimitProduction,
		},
		DestructiveOps: domain.DestructiveForbidden,
		CustomRules:    []generated.SafetyRule{},
	}
}

// CreateBaseTypeSafeSafetyRules creates a default TypeSafeSafetyRules structure for test reuse.
// This helper eliminates duplicate fixture code across test files.
func CreateBaseTypeSafeSafetyRules() *domain.TypeSafeSafetyRules {
	rules := newTypeSafeSafetyRules()
	rules.StyleRules.SelectStarPolicy = domain.SelectStarAllowed
	rules.SafetyRules.WhereRequirement = domain.WhereClauseNever
	rules.SafetyRules.LimitRequirement = domain.LimitClauseNever
	rules.SafetyRules.MaxRowsWithoutLimit = 0
	rules.DestructiveOps = domain.DestructiveAllowed

	return rules
}

// CreateTypeSafeSafetyRules creates a TypeSafeSafetyRules with a common test configuration.
// This helper provides a default for most test scenarios and accepts an optional configuration callback.
//
// Example usage:
//
//	rules := testing.CreateTypeSafeSafetyRules(func(r *domain.TypeSafeSafetyRules) {
//	    r.SafetyRules.WhereRequirement = domain.WhereClauseAlways
//	})
func CreateTypeSafeSafetyRules(
	configure func(*domain.TypeSafeSafetyRules),
) *domain.TypeSafeSafetyRules {
	rules := newTypeSafeSafetyRules()
	if configure != nil {
		configure(rules)
	}

	return rules
}

// CreateRestrictiveTypeSafeSafetyRulesWithCustomRules creates a TypeSafeSafetyRules with
// restrictive safety settings and accepts custom rules for testing custom rule preservation.
// This helper eliminates duplicate fixture code for testing custom safety rules.
func CreateRestrictiveTypeSafeSafetyRulesWithCustomRules(
	customRules []generated.SafetyRule,
) *domain.TypeSafeSafetyRules {
	return CreateTypeSafeSafetyRules(func(r *domain.TypeSafeSafetyRules) {
		r.SafetyRules.WhereRequirement = domain.WhereClauseAlways
		r.SafetyRules.LimitRequirement = domain.LimitClauseNever
		r.SafetyRules.MaxRowsWithoutLimit = domain.MaxRowsWithoutLimitDefault
		r.CustomRules = customRules
	})
}

// CreateStrictTypeSafeSafetyRules creates a strictly configured TypeSafeSafetyRules for testing.
// All safety rules are set to their most restrictive values.
func CreateStrictTypeSafeSafetyRules() *domain.TypeSafeSafetyRules {
	rules := newTypeSafeSafetyRules()
	rules.StyleRules.ColumnExplicitness = domain.ColumnExplicitnessRequired
	rules.SafetyRules.WhereRequirement = domain.WhereClauseAlways
	rules.SafetyRules.LimitRequirement = domain.LimitClauseAlways

	return rules
}

// CreateQueryStyleRulesForbiddenSelectStar creates QueryStyleRules that forbids SELECT *.
// This helper eliminates duplicate fixture code for testing select star policies.
func CreateQueryStyleRulesForbiddenSelectStar() domain.QueryStyleRules {
	return domain.QueryStyleRules{
		SelectStarPolicy:   domain.SelectStarForbidden,
		ColumnExplicitness: domain.ColumnExplicitnessDefault,
	}
}

// CreateQuerySafetyRulesStrict creates QuerySafetyRules with strict safety settings.
// All safety rules are set to their most restrictive values.
func CreateQuerySafetyRulesStrict() domain.QuerySafetyRules {
	return domain.QuerySafetyRules{
		WhereRequirement:    domain.WhereClauseAlways,
		LimitRequirement:    domain.LimitClauseNever,
		MaxRowsWithoutLimit: domain.MaxRowsWithoutLimitDefault,
	}
}

// CreateGeneratedSafetyRulesForbidden creates a SafetyRules with all forbidden flags set.
// This helper eliminates duplicate fixture code across test files.
func CreateGeneratedSafetyRulesForbidden() generated.SafetyRules {
	return generated.SafetyRules{
		NoSelectStar: true,
		RequireWhere: true,
		NoDropTable:  true,
		NoTruncate:   true,
		RequireLimit: false,
		Rules:        []generated.SafetyRule{},
	}
}

// CreateGeneratedSafetyRulesAllowed creates a SafetyRules with all forbidden flags cleared.
// This helper eliminates duplicate fixture code across test files.
func CreateGeneratedSafetyRulesAllowed() generated.SafetyRules {
	return generated.SafetyRules{
		NoSelectStar: false,
		RequireWhere: false,
		NoDropTable:  false,
		NoTruncate:   false,
		RequireLimit: true,
		Rules:        []generated.SafetyRule{},
	}
}

// CreateGeneratedSafetyRulesAllowedWithCustomRules creates a SafetyRules with all forbidden flags
// cleared and accepts custom rules for testing custom rule preservation.
// This helper eliminates duplicate fixture code for testing custom safety rules.
//
// Deprecated: domain.SafetyRules is the legacy boolean version. New tests should use the
// type-safe helpers (CreateBaseTypeSafeSafetyRules, etc.) instead.
func CreateGeneratedSafetyRulesAllowedWithCustomRules(
	customRules []generated.SafetyRule,
) *domain.SafetyRules {
	return &domain.SafetyRules{
		NoSelectStar: false,
		RequireWhere: false,
		NoDropTable:  false,
		NoTruncate:   false,
		RequireLimit: false,
		Rules:        customRules,
	}
}
