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
			rules := &generated.SafetyRules{
				NoSelectStar: true,
			}

			configRules := transformer.TransformSafetyRules(rules)

			Expect(configRules).To(HaveLen(1))
			Expect(configRules[0].Name).To(Equal("no-select-star"))
			Expect(configRules[0].Rule).To(Equal("!query.contains('SELECT *')"))
			Expect(configRules[0].Message).To(Equal("SELECT * is not allowed"))
		})

		It("should transform RequireWhere rule", func() {
			rules := &generated.SafetyRules{
				RequireWhere: true,
			}

			configRules := transformer.TransformSafetyRules(rules)

			Expect(configRules).To(HaveLen(1))
			Expect(configRules[0].Name).To(Equal("require-where"))
			Expect(configRules[0].Rule).To(Equal("query.type in ('SELECT', 'UPDATE', 'DELETE') && query.hasWhereClause()"))
			Expect(configRules[0].Message).To(Equal("WHERE clause is required for this query type"))
		})

		It("should transform RequireLimit rule", func() {
			rules := &generated.SafetyRules{
				RequireLimit: true,
			}

			configRules := transformer.TransformSafetyRules(rules)

			Expect(configRules).To(HaveLen(1))
			Expect(configRules[0].Name).To(Equal("require-limit"))
			Expect(configRules[0].Rule).To(Equal("query.type == 'SELECT' && !query.hasLimitClause()"))
			Expect(configRules[0].Message).To(Equal("LIMIT clause is required for SELECT queries"))
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
				Rules:      customRules,
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
				Rules:      []generated.SafetyRule{customRule},
			}

			configRules := transformer.TransformDomainSafetyRules(rules)

			Expect(configRules).To(HaveLen(1))
			Expect(configRules[0].Name).To(Equal("domain-custom"))
		})
	})
})