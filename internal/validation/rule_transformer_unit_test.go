package validation_test

import (
	"testing"

	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/domain"
	testingHelper "github.com/LarsArtmann/SQLC-Wizzard/internal/testing"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/validation"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

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
			runMultiRuleTransformationTest(
				transformer,
				func() *generated.SafetyRules {
					return &generated.SafetyRules{
						NoSelectStar: true,
						RequireWhere: true,
						RequireLimit: false,
					}
				},
				func(t *validation.RuleTransformer, r *generated.SafetyRules) []generated.RuleConfig {
					return t.TransformSafetyRules(r)
				},
				[]string{"no-select-star", "require-where"},
			)
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

			runMultiRuleTransformationTest(
				transformer,
				func() *generated.SafetyRules {
					return &generated.SafetyRules{
						NoSelectStar: false,
						RequireWhere: false,
						RequireLimit: false,
						Rules:        customRules,
					}
				},
				func(t *validation.RuleTransformer, r *generated.SafetyRules) []generated.RuleConfig {
					return t.TransformSafetyRules(r)
				},
				[]string{"custom-rule-1", "custom-rule-2"},
			)
		})

		It("should preserve rule order", func() {
			runMultiRuleTransformationTest(
				transformer,
				func() *generated.SafetyRules {
					return &generated.SafetyRules{
						NoSelectStar: true,
						RequireWhere: true,
						RequireLimit: true,
					}
				},
				func(t *validation.RuleTransformer, r *generated.SafetyRules) []generated.RuleConfig {
					return t.TransformSafetyRules(r)
				},
				[]string{"no-select-star", "require-where", "require-limit"},
			)
		})
	})

	Context("TransformDomainSafetyRules", func() {
		It("should transform domain safety rules", func() {
			runMultiRuleTransformationTest(
				transformer,
				func() *generated.SafetyRules {
					rules := &domain.SafetyRules{
						NoSelectStar: true,
						RequireWhere: false,
						RequireLimit: true,
					}

					return toGeneratedSafetyRules(rules)
				},
				func(t *validation.RuleTransformer, r *generated.SafetyRules) []generated.RuleConfig {
					return t.TransformDomainSafetyRules(r)
				},
				[]string{"no-select-star", "require-limit"},
			)
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
			rules := testingHelper.CreateGeneratedSafetyRulesAllowedWithCustomRules(
				[]generated.SafetyRule{customRule},
			)

			configRules := transformer.TransformDomainSafetyRules(rules)

			Expect(configRules).To(HaveLen(1))
			Expect(configRules[0].Name).To(Equal("domain-custom"))
		})
	})

	Context("TransformTypeSafeSafetyRules", func() {
		It("should produce empty rules when all type-safe policies are disabled", func() {
			rules := testingHelper.CreateBaseTypeSafeSafetyRules()

			configRules := transformer.TransformTypeSafeSafetyRules(rules)

			Expect(configRules).To(BeEmpty())
		})

		It("should transform SelectStarPolicy.ForbidsSelectStar() correctly", func() {
			rules := testingHelper.CreateBaseTypeSafeSafetyRules()
			rules.StyleRules.SelectStarPolicy = domain.SelectStarForbidden

			configRules := transformer.TransformTypeSafeSafetyRules(rules)

			Expect(configRules).To(HaveLen(1))
			Expect(configRules[0].Name).To(Equal("no-select-star"))
			Expect(configRules[0].Rule).To(Equal("!query.contains('SELECT *')"))
			Expect(configRules[0].Message).To(ContainSubstring("explicit column names"))
		})

		It("should transform WhereRequirement.RequiresOnDestructive() correctly", func() {
			rules := testingHelper.CreateBaseTypeSafeSafetyRules()
			rules.SafetyRules.WhereRequirement = domain.WhereClauseOnDestructive

			configRules := transformer.TransformTypeSafeSafetyRules(rules)

			Expect(configRules).To(HaveLen(1))
			Expect(configRules[0].Name).To(Equal("require-where"))
			Expect(
				configRules[0].Rule,
			).To(Equal("query.type in ('SELECT', 'UPDATE', 'DELETE') && !query.hasWhereClause()"))
		})

		It("should transform LimitRequirement.RequiresOnSelect() correctly", func() {
			rules := testingHelper.CreateBaseTypeSafeSafetyRules()
			rules.SafetyRules.LimitRequirement = domain.LimitClauseOnSelect

			expectSingleRule(
				transformer,
				rules,
				"require-limit",
				"query.type == 'SELECT' && !query.hasLimitClause()",
			)
		})

		It("should transform MaxRowsWithoutLimit correctly", func() {
			rules := testingHelper.CreateBaseTypeSafeSafetyRules()
			rules.SafetyRules.MaxRowsWithoutLimit = 100

			expectSingleRule(
				transformer,
				rules,
				"max-rows-without-limit",
				"query.type == 'SELECT' && (!query.hasLimitClause() || query.limitValue() > 100)",
			)
		})
	})
})
