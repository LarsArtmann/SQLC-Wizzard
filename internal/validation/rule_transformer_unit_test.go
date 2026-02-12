package validation_test

import (
	"testing"

	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/domain"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/validation"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// runRuleTransformationTest runs a generic test for safety rule transformation.
func runRuleTransformationTest(transformer *validation.RuleTransformer, ruleName, expectedRule, expectedMessage string, setupRules func() *generated.SafetyRules) {
	rules := setupRules()
	configRules := transformer.TransformSafetyRules(rules)

	Expect(configRules).To(HaveLen(1))
	Expect(configRules[0].Name).To(Equal(ruleName))
	Expect(configRules[0].Rule).To(Equal(expectedRule))
	Expect(configRules[0].Message).To(Equal(expectedMessage))
}

// createBaseTypeSafeSafetyRules creates a default TypeSafeSafetyRules structure for test reuse.
func createBaseTypeSafeSafetyRules() *domain.TypeSafeSafetyRules {
	return &domain.TypeSafeSafetyRules{
		StyleRules: domain.QueryStyleRules{
			SelectStarPolicy:   domain.SelectStarAllowed,
			ColumnExplicitness: domain.ColumnExplicitnessDefault,
		},
		SafetyRules: domain.QuerySafetyRules{
			WhereRequirement:    domain.WhereClauseNever,
			LimitRequirement:    domain.LimitClauseNever,
			MaxRowsWithoutLimit: 0,
		},
		DestructiveOps: domain.DestructiveAllowed,
		CustomRules:    []generated.SafetyRule{},
	}
}

// expectSingleRule verifies that a single rule with the expected name and rule expression is generated.
func expectSingleRule(transformer *validation.RuleTransformer, rules *domain.TypeSafeSafetyRules, expectedName, expectedRule string) {
	configRules := transformer.TransformTypeSafeSafetyRules(rules)

	Expect(configRules).To(HaveLen(1))
	Expect(configRules[0].Name).To(Equal(expectedName))
	Expect(configRules[0].Rule).To(Equal(expectedRule))
}

// assertIdenticalRules checks that boolean-based and type-safe rule transformations produce identical results.
func assertIdenticalRules(boolResult, typeSafeResult []generated.RuleConfig, expectedRule string) {
	Expect(boolResult).To(HaveLen(1))
	Expect(typeSafeResult).To(HaveLen(1))
	Expect(boolResult[0].Rule).To(Equal(typeSafeResult[0].Rule))
	Expect(boolResult[0].Name).To(Equal(typeSafeResult[0].Name))
	Expect(boolResult[0].Rule).To(Equal(expectedRule))
}

// runParityTest runs a parity test between boolean and type-safe safety rules.
func runParityTest(
	transformer *validation.RuleTransformer,
	boolRules *generated.SafetyRules,
	setupTypeSafeRules func() *domain.TypeSafeSafetyRules,
	expectedExpression string,
) {
	boolResult := transformer.TransformSafetyRules(boolRules)
	typeSafeRules := setupTypeSafeRules()
	typeSafeResult := transformer.TransformTypeSafeSafetyRules(typeSafeRules)
	assertIdenticalRules(boolResult, typeSafeResult, expectedExpression)
}

// runClauseRequirementParityTest is a generic helper for clause requirement parity tests.
// It reduces duplication for tests that verify boolean flags (RequireWhere, RequireLimit)
// produce identical expressions to type-safe enum values (WhereClauseOnDestructive, LimitClauseOnSelect).
func runClauseRequirementParityTest(
	transformer *validation.RuleTransformer,
	boolRules *generated.SafetyRules,
	setTypeSafeRequirement func(*domain.TypeSafeSafetyRules),
	expectedExpression string,
) {
	runParityTest(
		transformer,
		boolRules,
		func() *domain.TypeSafeSafetyRules {
			rules := createBaseTypeSafeSafetyRules()
			setTypeSafeRequirement(rules)
			return rules
		},
		expectedExpression,
	)
}

func TestValidation(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Validation Unit Suite")
}

var _ = Describe("RuleTransformer Unit Tests", func() {
	var transformer *validation.RuleTransformer

	BeforeEach(func() {
		transformer = validation.NewRuleTransformer()
	})

	Context("NewRuleTransformer", func() {
		It("should create a new transformer", func() {
			Expect(transformer).NotTo(BeNil())
		})
	})

	Context("TransformSafetyRules", func() {
		It("should return empty rules when all flags are false", func() {
			rules := &generated.SafetyRules{
				NoSelectStar: false,
				RequireWhere: false,
				RequireLimit: false,
			}

			configRules := transformer.TransformSafetyRules(rules)

			Expect(configRules).To(BeEmpty())
		})

		It("should transform NoSelectStar rule", func() {
			runRuleTransformationTest(
				transformer,
				"no-select-star",
				"!query.contains('SELECT *')",
				"SELECT * is not allowed",
				func() *generated.SafetyRules {
					return &generated.SafetyRules{NoSelectStar: true}
				},
			)
		})

		It("should transform RequireWhere rule", func() {
			runRuleTransformationTest(
				transformer,
				"require-where",
				"query.type in ('SELECT', 'UPDATE', 'DELETE') && !query.hasWhereClause()",
				"WHERE clause is required for this query type",
				func() *generated.SafetyRules {
					return &generated.SafetyRules{RequireWhere: true}
				},
			)
		})

		It("should transform RequireLimit rule", func() {
			runRuleTransformationTest(
				transformer,
				"require-limit",
				"query.type == 'SELECT' && !query.hasLimitClause()",
				"LIMIT clause is required for SELECT queries",
				func() *generated.SafetyRules {
					return &generated.SafetyRules{RequireLimit: true}
				},
			)
		})

		It("should transform multiple rules", func() {
			rules := &generated.SafetyRules{
				NoSelectStar: true,
				RequireWhere: true,
				RequireLimit: false,
			}

			configRules := transformer.TransformSafetyRules(rules)

			Expect(configRules).To(HaveLen(2))
			Expect(configRules[0].Name).To(Equal("no-select-star"))
			Expect(configRules[1].Name).To(Equal("require-where"))
		})

		It("should transform custom rules", func() {
			customRules := []generated.SafetyRule{
				{
					Name:    "custom-rule-1",
					Rule:    "query.contains('FOR UPDATE')",
					Message: "FOR UPDATE not allowed",
				},
				{
					Name:    "custom-rule-2",
					Rule:    "query.contains('INSERT INTO temp_table')",
					Message: "temp table operations not allowed",
				},
			}
			rules := &generated.SafetyRules{
				NoSelectStar: false,
				RequireWhere: false,
				RequireLimit: false,
				Rules:        customRules,
			}

			configRules := transformer.TransformSafetyRules(rules)

			Expect(configRules).To(HaveLen(2))
			Expect(configRules[0].Name).To(Equal("custom-rule-1"))
			Expect(configRules[1].Name).To(Equal("custom-rule-2"))
		})

		It("should preserve rule order", func() {
			rules := &generated.SafetyRules{
				NoSelectStar: true,
				RequireWhere: true,
				RequireLimit: true,
			}

			configRules := transformer.TransformSafetyRules(rules)

			Expect(configRules).To(HaveLen(3))
			Expect(configRules[0].Name).To(Equal("no-select-star"))
			Expect(configRules[1].Name).To(Equal("require-where"))
			Expect(configRules[2].Name).To(Equal("require-limit"))
		})
	})

	Context("TransformDomainSafetyRules", func() {
		It("should transform domain safety rules", func() {
			rules := &domain.SafetyRules{
				NoSelectStar: true,
				RequireWhere: false,
				RequireLimit: true,
			}

			configRules := transformer.TransformDomainSafetyRules(rules)

			Expect(configRules).To(HaveLen(2))
			Expect(configRules[0].Name).To(Equal("no-select-star"))
			Expect(configRules[1].Name).To(Equal("require-limit"))
		})

		It("should return empty rules for false flags", func() {
			rules := &domain.SafetyRules{
				NoSelectStar: false,
				RequireWhere: false,
				RequireLimit: false,
			}

			configRules := transformer.TransformDomainSafetyRules(rules)

			Expect(configRules).To(BeEmpty())
		})

		It("should handle custom rules in domain safety", func() {
			customRule := generated.SafetyRule{
				Name:    "domain-custom",
				Rule:    "query.contains('DOMAIN_PATTERNS')",
				Message: "Domain patterns not allowed",
			}
			rules := &domain.SafetyRules{
				NoSelectStar: false,
				RequireWhere: false,
				RequireLimit: false,
				Rules:        []generated.SafetyRule{customRule},
			}

			configRules := transformer.TransformDomainSafetyRules(rules)

			Expect(configRules).To(HaveLen(1))
			Expect(configRules[0].Name).To(Equal("domain-custom"))
		})
	})

	Context("TransformTypeSafeSafetyRules", func() {
		It("should produce empty rules when all type-safe policies are disabled", func() {
			rules := &domain.TypeSafeSafetyRules{
				StyleRules: domain.QueryStyleRules{
					SelectStarPolicy:   domain.SelectStarAllowed,
					ColumnExplicitness: domain.ColumnExplicitnessDefault,
				},
				SafetyRules: domain.QuerySafetyRules{
					WhereRequirement:    domain.WhereClauseNever,
					LimitRequirement:    domain.LimitClauseNever,
					MaxRowsWithoutLimit: 0,
				},
				DestructiveOps: domain.DestructiveAllowed,
				CustomRules:    []generated.SafetyRule{},
			}

			configRules := transformer.TransformTypeSafeSafetyRules(rules)

			Expect(configRules).To(BeEmpty())
		})

		It("should transform SelectStarPolicy.ForbidsSelectStar() correctly", func() {
			rules := createBaseTypeSafeSafetyRules()
			rules.StyleRules.SelectStarPolicy = domain.SelectStarForbidden

			configRules := transformer.TransformTypeSafeSafetyRules(rules)

			Expect(configRules).To(HaveLen(1))
			Expect(configRules[0].Name).To(Equal("no-select-star"))
			Expect(configRules[0].Rule).To(Equal("!query.contains('SELECT *')"))
			Expect(configRules[0].Message).To(ContainSubstring("explicit column names"))
		})

		It("should transform WhereRequirement.RequiresOnDestructive() correctly", func() {
			rules := createBaseTypeSafeSafetyRules()
			rules.SafetyRules.WhereRequirement = domain.WhereClauseOnDestructive

			configRules := transformer.TransformTypeSafeSafetyRules(rules)

			Expect(configRules).To(HaveLen(1))
			Expect(configRules[0].Name).To(Equal("require-where"))
			Expect(configRules[0].Rule).To(Equal("query.type in ('SELECT', 'UPDATE', 'DELETE') && !query.hasWhereClause()"))
		})

		It("should transform LimitRequirement.RequiresOnSelect() correctly", func() {
			rules := createBaseTypeSafeSafetyRules()
			rules.SafetyRules.LimitRequirement = domain.LimitClauseOnSelect

			expectSingleRule(transformer, rules, "require-limit", "query.type == 'SELECT' && !query.hasLimitClause()")
		})

		It("should transform MaxRowsWithoutLimit correctly", func() {
			rules := createBaseTypeSafeSafetyRules()
			rules.SafetyRules.MaxRowsWithoutLimit = 100

			expectSingleRule(transformer, rules, "max-rows-without-limit", "query.type == 'SELECT' && (!query.hasLimitClause() || query.limitValue() > 100)")
		})
	})

	Context("Transformer Parity Tests", func() {
		// These tests verify that equivalent configurations produce identical expressions
		// following the "violation when true" convention consistently across both transformers

		It("should produce identical expressions for NoSelectStar vs ForbidsSelectStar", func() {
			runParityTest(
				transformer,
				&generated.SafetyRules{
					NoSelectStar: true,
					RequireWhere: false,
					RequireLimit: false,
				},
				func() *domain.TypeSafeSafetyRules {
					rules := createBaseTypeSafeSafetyRules()
					rules.StyleRules.SelectStarPolicy = domain.SelectStarForbidden
					return rules
				},
				"!query.contains('SELECT *')",
			)
		})

		It("should produce identical expressions for RequireWhere vs RequiresOnDestructive", func() {
			runClauseRequirementParityTest(
				transformer,
				&generated.SafetyRules{
					NoSelectStar: false,
					RequireWhere: true,
					RequireLimit: false,
				},
				func(rules *domain.TypeSafeSafetyRules) {
					rules.SafetyRules.WhereRequirement = domain.WhereClauseOnDestructive
				},
				"query.type in ('SELECT', 'UPDATE', 'DELETE') && !query.hasWhereClause()",
			)
		})

		It("should produce identical expressions for RequireLimit vs RequiresOnSelect", func() {
			runParityTest(
				transformer,
				&generated.SafetyRules{
					NoSelectStar: false,
					RequireWhere: false,
					RequireLimit: true,
				},
				func() *domain.TypeSafeSafetyRules {
					rules := createBaseTypeSafeSafetyRules()
					rules.SafetyRules.LimitRequirement = domain.LimitClauseOnSelect
					return rules
				},
				"query.type == 'SELECT' && !query.hasLimitClause()",
			)
		})

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

			typeSafeRules := &domain.TypeSafeSafetyRules{
				StyleRules: domain.QueryStyleRules{
					SelectStarPolicy:   domain.SelectStarForbidden,
					ColumnExplicitness: domain.ColumnExplicitnessDefault,
				},
				SafetyRules: domain.QuerySafetyRules{
					WhereRequirement:    domain.WhereClauseOnDestructive,
					LimitRequirement:    domain.LimitClauseOnSelect,
					MaxRowsWithoutLimit: 100,
				},
				DestructiveOps: domain.DestructiveForbidden,
				CustomRules:    []generated.SafetyRule{},
			}
			result := transformer.TransformTypeSafeSafetyRules(typeSafeRules)

			// Should have 4 rules: no-select-star, require-where, require-limit, max-rows-without-limit, no-drop-table, no-truncate
			Expect(result).To(HaveLen(6))

			// Find and verify each rule's expression follows violation convention
			ruleMap := make(map[string]generated.RuleConfig)
			for _, rule := range result {
				ruleMap[rule.Name] = rule
			}

			// NoSelectStar: violation when SELECT * exists → !contains('SELECT *')
			Expect(ruleMap["no-select-star"].Rule).To(MatchRegexp(`!query\.contains\('SELECT \*'\)`))

			// RequireWhere: violation when WHERE is missing → hasWhereClause() in expression
			Expect(ruleMap["require-where"].Rule).To(ContainSubstring("hasWhereClause()"))

			// RequireLimit: violation when LIMIT is missing → !hasLimitClause()
			Expect(ruleMap["require-limit"].Rule).To(MatchRegexp(`!query\.hasLimitClause\(\)`))

			// MaxRowsWithoutLimit: violation when limit is missing or too small
			Expect(ruleMap["max-rows-without-limit"].Rule).To(ContainSubstring("hasLimitClause()"))
			Expect(ruleMap["max-rows-without-limit"].Rule).To(ContainSubstring("limitValue()"))

			// NoDropTable: violation when DROP TABLE exists
			Expect(ruleMap["no-drop-table"].Rule).To(MatchRegexp(`!query\.contains\('DROP TABLE'\)`))

			// NoTruncate: violation when TRUNCATE exists
			Expect(ruleMap["no-truncate"].Rule).To(MatchRegexp(`!query\.contains\('TRUNCATE'\)`))
		})
	})
})
