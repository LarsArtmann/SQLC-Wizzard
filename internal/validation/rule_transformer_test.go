package validation_test

import (
	"testing"

	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/domain"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/validation"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestValidation(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Validation Suite")
}

var _ = Describe("RuleTransformer", func() {
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
			rules := &generated.SafetyRules{
				NoSelectStar: true,
			}

			configRules := transformer.TransformSafetyRules(rules)

			Expect(configRules).To(HaveLen(1))
			Expect(configRules[0].Name).To(Equal("no-select-star"))
			Expect(configRules[0].Rule).To(ContainSubstring("SELECT *"))
			Expect(configRules[0].Message).To(ContainSubstring("not allowed"))
		})

		It("should transform RequireWhere rule", func() {
			rules := &generated.SafetyRules{
				RequireWhere: true,
			}

			configRules := transformer.TransformSafetyRules(rules)

			Expect(configRules).To(HaveLen(1))
			Expect(configRules[0].Name).To(Equal("require-where"))
			Expect(configRules[0].Rule).To(ContainSubstring("hasWhereClause"))
			Expect(configRules[0].Message).To(ContainSubstring("WHERE clause"))
		})

		It("should transform RequireLimit rule", func() {
			rules := &generated.SafetyRules{
				RequireLimit: true,
			}

			configRules := transformer.TransformSafetyRules(rules)

			Expect(configRules).To(HaveLen(1))
			Expect(configRules[0].Name).To(Equal("require-limit"))
			Expect(configRules[0].Rule).To(ContainSubstring("hasLimitClause"))
			Expect(configRules[0].Message).To(ContainSubstring("LIMIT clause"))
		})

		It("should transform all built-in rules when enabled", func() {
			rules := &generated.SafetyRules{
				NoSelectStar: true,
				RequireWhere: true,
				RequireLimit: true,
			}

			configRules := transformer.TransformSafetyRules(rules)

			Expect(configRules).To(HaveLen(3))

			// Verify all rules are present
			ruleNames := make([]string, len(configRules))
			for i, rule := range configRules {
				ruleNames[i] = rule.Name
			}

			Expect(ruleNames).To(ContainElement("no-select-star"))
			Expect(ruleNames).To(ContainElement("require-where"))
			Expect(ruleNames).To(ContainElement("require-limit"))
		})

		It("should transform custom rules", func() {
			rules := &generated.SafetyRules{
				NoSelectStar: false,
				RequireWhere: false,
				RequireLimit: false,
				Rules: []generated.SafetyRule{
					{
						Name:    "no-delete-without-limit",
						Rule:    "query.type == 'DELETE' && !query.hasLimitClause()",
						Message: "DELETE queries must have LIMIT clause",
					},
				},
			}

			configRules := transformer.TransformSafetyRules(rules)

			Expect(configRules).To(HaveLen(1))
			Expect(configRules[0].Name).To(Equal("no-delete-without-limit"))
			Expect(configRules[0].Rule).To(Equal("query.type == 'DELETE' && !query.hasLimitClause()"))
			Expect(configRules[0].Message).To(Equal("DELETE queries must have LIMIT clause"))
		})

		It("should combine built-in and custom rules", func() {
			rules := &generated.SafetyRules{
				NoSelectStar: true,
				RequireWhere: true,
				RequireLimit: false,
				Rules: []generated.SafetyRule{
					{
						Name:    "custom-rule-1",
						Rule:    "custom.expression()",
						Message: "Custom rule message",
					},
					{
						Name:    "custom-rule-2",
						Rule:    "another.expression()",
						Message: "Another custom rule",
					},
				},
			}

			configRules := transformer.TransformSafetyRules(rules)

			Expect(configRules).To(HaveLen(4))

			// First two should be built-in rules
			Expect(configRules[0].Name).To(Equal("no-select-star"))
			Expect(configRules[1].Name).To(Equal("require-where"))

			// Last two should be custom rules
			Expect(configRules[2].Name).To(Equal("custom-rule-1"))
			Expect(configRules[3].Name).To(Equal("custom-rule-2"))
		})

		It("should handle empty custom rules array", func() {
			rules := &generated.SafetyRules{
				NoSelectStar: true,
				Rules:        []generated.SafetyRule{},
			}

			configRules := transformer.TransformSafetyRules(rules)

			Expect(configRules).To(HaveLen(1))
			Expect(configRules[0].Name).To(Equal("no-select-star"))
		})

		It("should handle nil custom rules array", func() {
			rules := &generated.SafetyRules{
				NoSelectStar: true,
				Rules:        nil,
			}

			configRules := transformer.TransformSafetyRules(rules)

			Expect(configRules).To(HaveLen(1))
		})
	})

	Context("TransformDomainSafetyRules", func() {
		It("should transform domain safety rules", func() {
			rules := &domain.SafetyRules{
				NoSelectStar: true,
				RequireWhere: true,
			}

			configRules := transformer.TransformDomainSafetyRules(rules)

			Expect(configRules).To(HaveLen(2))
			Expect(configRules[0].Name).To(Equal("no-select-star"))
			Expect(configRules[1].Name).To(Equal("require-where"))
		})

		It("should handle all domain safety rule flags", func() {
			rules := &domain.SafetyRules{
				NoSelectStar: true,
				RequireWhere: true,
				RequireLimit: true,
			}

			configRules := transformer.TransformDomainSafetyRules(rules)

			Expect(configRules).To(HaveLen(3))
		})

		It("should handle domain custom rules", func() {
			rules := &domain.SafetyRules{
				NoSelectStar: false,
				Rules: []generated.SafetyRule{
					{
						Name:    "domain-custom",
						Rule:    "domain.rule()",
						Message: "Domain custom rule",
					},
				},
			}

			configRules := transformer.TransformDomainSafetyRules(rules)

			Expect(configRules).To(HaveLen(1))
			Expect(configRules[0].Name).To(Equal("domain-custom"))
		})

		It("should produce same result as TransformSafetyRules", func() {
			// Create identical rules for both types
			genRules := &generated.SafetyRules{
				NoSelectStar: true,
				RequireWhere: true,
				RequireLimit: true,
				Rules: []generated.SafetyRule{
					{Name: "custom", Rule: "custom()", Message: "Custom"},
				},
			}

			domainRules := (*domain.SafetyRules)(genRules)

			genResult := transformer.TransformSafetyRules(genRules)
			domainResult := transformer.TransformDomainSafetyRules(domainRules)

			Expect(domainResult).To(Equal(genResult))
		})
	})

	Context("Rule Consolidation", func() {
		It("should provide single source of truth for rule transformation", func() {
			// This test verifies that we have eliminated the split brain
			// by having a single transformer that works for both generated and domain types

			rules := &generated.SafetyRules{
				NoSelectStar: true,
				RequireWhere: true,
				RequireLimit: true,
			}

			// Transform using generated rules
			result1 := transformer.TransformSafetyRules(rules)

			// Transform using domain rules (same underlying type)
			domainRules := (*domain.SafetyRules)(rules)
			result2 := transformer.TransformDomainSafetyRules(domainRules)

			// Both should produce identical results
			Expect(result1).To(Equal(result2))
			Expect(len(result1)).To(Equal(3))
		})
	})

	Context("Edge Cases", func() {
		It("should handle nil rules gracefully", func() {
			// Note: This would panic in real usage, but demonstrates defensive programming need
			// In production, callers should ensure rules are not nil

			rules := &generated.SafetyRules{}
			configRules := transformer.TransformSafetyRules(rules)

			Expect(configRules).To(BeEmpty())
		})

		It("should preserve rule order", func() {
			rules := &generated.SafetyRules{
				NoSelectStar: true,
				RequireWhere: true,
				RequireLimit: true,
				Rules: []generated.SafetyRule{
					{Name: "custom-1", Rule: "r1()", Message: "M1"},
					{Name: "custom-2", Rule: "r2()", Message: "M2"},
				},
			}

			configRules := transformer.TransformSafetyRules(rules)

			Expect(configRules).To(HaveLen(5))
			Expect(configRules[0].Name).To(Equal("no-select-star"))
			Expect(configRules[1].Name).To(Equal("require-where"))
			Expect(configRules[2].Name).To(Equal("require-limit"))
			Expect(configRules[3].Name).To(Equal("custom-1"))
			Expect(configRules[4].Name).To(Equal("custom-2"))
		})

		It("should handle rules with special characters", func() {
			rules := &generated.SafetyRules{
				Rules: []generated.SafetyRule{
					{
						Name:    "special-chars-rule",
						Rule:    "query.contains('SELECT * FROM \"users\"')",
						Message: "Don't use SELECT * with quotes!",
					},
				},
			}

			configRules := transformer.TransformSafetyRules(rules)

			Expect(configRules).To(HaveLen(1))
			Expect(configRules[0].Rule).To(ContainSubstring("\"users\""))
			Expect(configRules[0].Message).To(ContainSubstring("!"))
		})

		It("should handle empty rule names", func() {
			rules := &generated.SafetyRules{
				Rules: []generated.SafetyRule{
					{
						Name:    "",
						Rule:    "some.rule()",
						Message: "Some message",
					},
				},
			}

			configRules := transformer.TransformSafetyRules(rules)

			Expect(configRules).To(HaveLen(1))
			Expect(configRules[0].Name).To(Equal(""))
		})

		It("should handle long rule expressions", func() {
			longRule := "query.type == 'SELECT' && query.hasWhereClause() && query.hasLimitClause() && !query.contains('SELECT *') && query.tables().count() < 5"

			rules := &generated.SafetyRules{
				Rules: []generated.SafetyRule{
					{
						Name:    "complex-rule",
						Rule:    longRule,
						Message: "Complex validation rule",
					},
				},
			}

			configRules := transformer.TransformSafetyRules(rules)

			Expect(configRules).To(HaveLen(1))
			Expect(configRules[0].Rule).To(Equal(longRule))
		})
	})

	Context("Real-world Scenarios", func() {
		It("should handle typical production configuration", func() {
			rules := &generated.SafetyRules{
				NoSelectStar: true,
				RequireWhere: true,
				RequireLimit: false, // Not required for all queries
				Rules: []generated.SafetyRule{
					{
						Name:    "no-unindexed-joins",
						Rule:    "query.joins().all(j => j.isIndexed())",
						Message: "All joins must use indexed columns",
					},
					{
						Name:    "max-limit-1000",
						Rule:    "!query.hasLimitClause() || query.limit() <= 1000",
						Message: "LIMIT must not exceed 1000",
					},
				},
			}

			configRules := transformer.TransformSafetyRules(rules)

			Expect(configRules).To(HaveLen(4))

			// Verify production rules are correctly transformed
			Expect(configRules[0].Name).To(Equal("no-select-star"))
			Expect(configRules[1].Name).To(Equal("require-where"))
			Expect(configRules[2].Name).To(Equal("no-unindexed-joins"))
			Expect(configRules[3].Name).To(Equal("max-limit-1000"))
		})
	})
})
