package validation_test

import (
	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/domain"
	testingHelper "github.com/LarsArtmann/SQLC-Wizzard/internal/testing"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("RuleTransformer Parity Tests", func() {
	transformer := newRuleTransformerForParity()

	// These tests verify that equivalent configurations produce identical expressions
	// following the "violation when true" convention consistently across both transformers
	runBooleanFlagParityTest(
		transformer,
		"should produce identical expressions for NoSelectStar vs ForbidsSelectStar",
		true, false, false,
		func(rules *domain.TypeSafeSafetyRules) {
			rules.StyleRules.SelectStarPolicy = domain.SelectStarForbidden
		},
		"!query.contains('SELECT *')",
	)

	runBooleanFlagParityTest(
		transformer,
		"should produce identical expressions for RequireWhere vs RequiresOnDestructive",
		false, true, false,
		func(rules *domain.TypeSafeSafetyRules) {
			rules.SafetyRules.WhereRequirement = domain.WhereClauseOnDestructive
		},
		"query.type in ('SELECT', 'UPDATE', 'DELETE') && !query.hasWhereClause()",
	)

	runBooleanFlagParityTest(
		transformer,
		"should produce identical expressions for RequireLimit vs RequiresOnSelect",
		false, false, true,
		func(rules *domain.TypeSafeSafetyRules) {
			rules.SafetyRules.LimitRequirement = domain.LimitClauseOnSelect
		},
		"query.type == 'SELECT' && !query.hasLimitClause()",
	)

	It("should follow consistent violation polarity convention", func() {
		// Verify the convention: rules express "violation when true"
		//
		// For forbidden patterns (NoSelectStar, NoDropTable, NoTruncate):
		//   Rule uses "violation when pattern exists" → !query.contains('PATTERN')
		//
		// For required patterns (RequireWhere, RequireLimit):
		//   Rule uses "violation when pattern is missing" → !query.hasClause()
		//
		// This is consistent: the boolean flag + expression together mean "flag is true → violation condition"
		result := transformer.TransformTypeSafeSafetyRules(
			testingHelper.CreateTypeSafeSafetyRules(nil),
		)

		// Should have 4 rules: no-select-star, require-where, require-limit, max-rows-without-limit, no-drop-table, no-truncate
		Expect(result).To(HaveLen(6))

		// Find and verify each rule's expression follows violation convention
		ruleMap := make(map[string]generated.RuleConfig)
		for _, rule := range result {
			ruleMap[rule.Name] = rule
		}

		// NoSelectStar: violation when SELECT * exists → !contains('SELECT *')
		Expect(
			ruleMap["no-select-star"].Rule,
		).To(MatchRegexp(`!query\.contains\('SELECT \*'\)`))

		// RequireWhere: violation when WHERE is missing → hasWhereClause() in expression
		Expect(ruleMap["require-where"].Rule).To(ContainSubstring("hasWhereClause()"))

		// RequireLimit: violation when LIMIT is missing → !hasLimitClause()
		Expect(ruleMap["require-limit"].Rule).To(MatchRegexp(`!query\.hasLimitClause\(\)`))

		// MaxRowsWithoutLimit: violation when limit is missing or too small
		Expect(ruleMap["max-rows-without-limit"].Rule).To(ContainSubstring("hasLimitClause()"))
		Expect(ruleMap["max-rows-without-limit"].Rule).To(ContainSubstring("limitValue()"))

		// NoDropTable: violation when DROP TABLE exists
		Expect(
			ruleMap["no-drop-table"].Rule,
		).To(MatchRegexp(`!query\.contains\('DROP TABLE'\)`))

		// NoTruncate: violation when TRUNCATE exists
		Expect(ruleMap["no-truncate"].Rule).To(MatchRegexp(`!query\.contains\('TRUNCATE'\)`))
	})
})
