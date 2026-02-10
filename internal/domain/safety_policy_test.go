package domain_test

import (
	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/domain"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/testing"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// Test cases for TypeSafeSafetyRules and related types
// Run via TestDomain in domain_test.go

// safetyRulesExpectations defines the expected values for safety rules testing.
type safetyRulesExpectations struct {
	description            string
	noSelectStar           bool
	requireExplicitColumns bool
	requireWhere           bool
	requireLimit           bool
	maxRowsWithoutLimit    uint
	destructiveOps         domain.DestructiveOperationPolicy
}

// assertSafetyRules validates safety rules against expected values.
func assertSafetyRules(rules any, expectations safetyRulesExpectations) {
	switch r := rules.(type) {
	case domain.TypeSafeSafetyRules:
		Expect(r.StyleRules.SelectStarPolicy.ForbidsSelectStar()).To(Equal(expectations.noSelectStar), expectations.description+": SelectStarPolicy")
		Expect(r.StyleRules.ColumnExplicitness.RequiresExplicitColumns()).To(Equal(expectations.requireExplicitColumns), expectations.description+": ColumnExplicitness")
		Expect(r.SafetyRules.WhereRequirement.RequiresOnDestructive()).To(Equal(expectations.requireWhere), expectations.description+": WhereRequirement")
		Expect(r.SafetyRules.LimitRequirement.RequiresOnSelect()).To(Equal(expectations.requireLimit), expectations.description+": LimitRequirement")
		Expect(r.SafetyRules.MaxRowsWithoutLimit).To(Equal(expectations.maxRowsWithoutLimit), expectations.description+": MaxRowsWithoutLimit")
		Expect(r.DestructiveOps).To(Equal(expectations.destructiveOps), expectations.description+": DestructiveOps")

		err := r.IsValid()
		Expect(err).NotTo(HaveOccurred(), expectations.description+": Should be valid")

	default:
		Fail("Unsupported safety rules type")
	}
}

var _ = Describe("QueryStyleRules", func() {
	It("should allow independent configuration of style rules", func() {
		rules := domain.QueryStyleRules{
			SelectStarPolicy:   domain.SelectStarForbidden,
			ColumnExplicitness: domain.ColumnExplicitnessDefault,
		}

		Expect(rules.SelectStarPolicy.ForbidsSelectStar()).To(BeTrue())
		Expect(rules.ColumnExplicitness.RequiresExplicitColumns()).To(BeFalse())
	})

	It("should support both rules enabled", func() {
		rules := domain.QueryStyleRules{
			SelectStarPolicy:   domain.SelectStarForbidden,
			ColumnExplicitness: domain.ColumnExplicitnessRequired,
		}

		Expect(rules.SelectStarPolicy.ForbidsSelectStar()).To(BeTrue())
		Expect(rules.ColumnExplicitness.RequiresExplicitColumns()).To(BeTrue())
	})
})

var _ = Describe("QuerySafetyRules", func() {
	It("should allow independent configuration of safety rules", func() {
		rules := domain.QuerySafetyRules{
			WhereRequirement:    domain.WhereClauseAlways,
			LimitRequirement:    domain.LimitClauseNever,
			MaxRowsWithoutLimit: 1000,
		}

		Expect(rules.WhereRequirement.RequiresOnDestructive()).To(BeTrue())
		Expect(rules.LimitRequirement.RequiresOnSelect()).To(BeFalse())
		Expect(rules.MaxRowsWithoutLimit).To(Equal(uint(1000)))
	})

	It("should support uint for MaxRowsWithoutLimit", func() {
		rules := domain.QuerySafetyRules{
			MaxRowsWithoutLimit: 500,
		}

		// This ensures we're using uint (no negative values possible)
		Expect(rules.MaxRowsWithoutLimit).To(BeNumerically(">=", 0))
		Expect(rules.MaxRowsWithoutLimit).To(Equal(uint(500)))
	})

	It("should allow zero to mean no limit", func() {
		rules := domain.QuerySafetyRules{
			MaxRowsWithoutLimit: 0,
		}

		Expect(rules.MaxRowsWithoutLimit).To(Equal(uint(0)))
	})
})

// DestructiveOperationPolicy validation test suite.
type DestructiveOperationPolicyTestSuite struct{}

func (DestructiveOperationPolicyTestSuite) GetValidValues() []domain.DestructiveOperationPolicy {
	return []domain.DestructiveOperationPolicy{
		domain.DestructiveAllowed,
		domain.DestructiveWithConfirmation,
		domain.DestructiveForbidden,
	}
}

func (DestructiveOperationPolicyTestSuite) GetInvalidValues() []domain.DestructiveOperationPolicy {
	return []domain.DestructiveOperationPolicy{
		"invalid",
		"",
		"sometimes",
	}
}

func (DestructiveOperationPolicyTestSuite) GetTypeName() string {
	return "DestructiveOperationPolicy"
}

// destructivePolicyTestCase defines a test case for DestructiveOperationPolicy behavior.
type destructivePolicyTestCase struct {
	description      string
	allowedPolicy    domain.DestructiveOperationPolicy
	forbiddenPolicies []domain.DestructiveOperationPolicy
	checker          func(domain.DestructiveOperationPolicy) bool
}

// testDestructivePolicyBehavior tests a destructive policy behavior against all policy values.
func testDestructivePolicyBehavior(tc destructivePolicyTestCase) {
	It(tc.description, func() {
		Expect(tc.checker(tc.allowedPolicy)).To(BeTrue())
	})

	It(tc.description+" for other policies", func() {
		for _, policy := range tc.forbiddenPolicies {
			Expect(tc.checker(policy)).To(BeFalse())
		}
	})
}

var _ = Describe("DestructiveOperationPolicy", func() {
	// Use generic validation test suite
	testing.TestValidationSuite(DestructiveOperationPolicyTestSuite{})

	Context("AllowsDropTable", func() {
		testDestructivePolicyBehavior(destructivePolicyTestCase{
			description:       "should allow DROP TABLE only when policy is 'allowed'",
			allowedPolicy:     domain.DestructiveAllowed,
			forbiddenPolicies: []domain.DestructiveOperationPolicy{domain.DestructiveWithConfirmation, domain.DestructiveForbidden},
			checker:           func(p domain.DestructiveOperationPolicy) bool { return p.AllowsDropTable() },
		})
	})

	Context("AllowsTruncate", func() {
		testDestructivePolicyBehavior(destructivePolicyTestCase{
			description:       "should allow TRUNCATE only when policy is 'allowed'",
			allowedPolicy:     domain.DestructiveAllowed,
			forbiddenPolicies: []domain.DestructiveOperationPolicy{domain.DestructiveWithConfirmation, domain.DestructiveForbidden},
			checker:           func(p domain.DestructiveOperationPolicy) bool { return p.AllowsTruncate() },
		})
	})

	Context("RequiresConfirmation", func() {
		testDestructivePolicyBehavior(destructivePolicyTestCase{
			description:       "should require confirmation only for 'with_confirmation' policy",
			allowedPolicy:     domain.DestructiveWithConfirmation,
			forbiddenPolicies: []domain.DestructiveOperationPolicy{domain.DestructiveAllowed, domain.DestructiveForbidden},
			checker:           func(p domain.DestructiveOperationPolicy) bool { return p.RequiresConfirmation() },
		})
	})

	Context("String", func() {
		It("should return correct string representation", func() {
			Expect(domain.DestructiveAllowed.String()).To(Equal("allowed"))
			Expect(domain.DestructiveWithConfirmation.String()).To(Equal("with_confirmation"))
			Expect(domain.DestructiveForbidden.String()).To(Equal("forbidden"))
		})
	})
})

var _ = Describe("TypeSafeSafetyRules", func() {
	Context("IsValid", func() {
		It("should validate rules with valid policy", func() {
			rules := domain.TypeSafeSafetyRules{
				StyleRules: domain.QueryStyleRules{
					SelectStarPolicy: domain.SelectStarForbidden,
				},
				SafetyRules: domain.QuerySafetyRules{
					WhereRequirement: domain.WhereClauseAlways,
				},
				DestructiveOps: domain.DestructiveForbidden,
				CustomRules:    []generated.SafetyRule{},
			}

			err := rules.IsValid()
			Expect(err).NotTo(HaveOccurred())
		})

		It("should reject invalid destructive operation policy", func() {
			rules := domain.TypeSafeSafetyRules{
				DestructiveOps: "invalid",
			}

			err := rules.IsValid()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("DestructiveOps"))
		})

		It("should reject custom rule with empty name", func() {
			rules := domain.TypeSafeSafetyRules{
				DestructiveOps: domain.DestructiveForbidden,
				CustomRules: []generated.SafetyRule{
					{
						Name:    "",
						Rule:    "query.type == 'SELECT'",
						Message: "Test rule",
					},
				},
			}

			err := rules.IsValid()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("empty name"))
		})

		It("should reject custom rule with empty rule expression", func() {
			rules := domain.TypeSafeSafetyRules{
				DestructiveOps: domain.DestructiveForbidden,
				CustomRules: []generated.SafetyRule{
					{
						Name:    "test-rule",
						Rule:    "",
						Message: "Test rule",
					},
				},
			}

			err := rules.IsValid()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("empty rule expression"))
		})

		It("should validate custom rules with all fields", func() {
			rules := domain.TypeSafeSafetyRules{
				DestructiveOps: domain.DestructiveForbidden,
				CustomRules: []generated.SafetyRule{
					{
						Name:    "no-complex-joins",
						Rule:    "query.joins().count() <= 3",
						Message: "Maximum 3 joins allowed",
					},
				},
			}

			err := rules.IsValid()
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Context("NewTypeSafeSafetyRules", func() {
		It("should return valid default rules", func() {
			defaults := domain.NewTypeSafeSafetyRules()

			assertSafetyRules(defaults, safetyRulesExpectations{
				description:            "Default safety rules",
				noSelectStar:           true,
				requireExplicitColumns: false,
				requireWhere:           true,
				requireLimit:           false,
				maxRowsWithoutLimit:    1000,
				destructiveOps:         domain.DestructiveForbidden,
			})
		})

		It("should have safe production defaults", func() {
			defaults := domain.NewTypeSafeSafetyRules()

			// Production-safe: No SELECT *
			Expect(defaults.StyleRules.SelectStarPolicy.ForbidsSelectStar()).To(BeTrue())

			// Production-safe: Require WHERE
			Expect(defaults.SafetyRules.WhereRequirement.RequiresOnDestructive()).To(BeTrue())

			// Production-safe: Block destructive ops
			Expect(defaults.DestructiveOps).To(Equal(domain.DestructiveForbidden))
		})
	})

	Context("NewDevelopmentSafetyRules", func() {
		It("should return relaxed rules for development", func() {
			devRules := domain.NewDevelopmentSafetyRules()

			Expect(devRules.StyleRules.SelectStarPolicy.ForbidsSelectStar()).To(BeFalse())
			Expect(devRules.SafetyRules.WhereRequirement.RequiresOnDestructive()).To(BeFalse())
			Expect(devRules.SafetyRules.LimitRequirement.RequiresOnSelect()).To(BeFalse())
			Expect(devRules.SafetyRules.MaxRowsWithoutLimit).To(Equal(uint(0)))
			Expect(devRules.DestructiveOps).To(Equal(domain.DestructiveAllowed))

			err := devRules.IsValid()
			Expect(err).NotTo(HaveOccurred())
		})

		It("should allow destructive operations in dev", func() {
			devRules := domain.NewDevelopmentSafetyRules()

			Expect(devRules.DestructiveOps.AllowsDropTable()).To(BeTrue())
			Expect(devRules.DestructiveOps.AllowsTruncate()).To(BeTrue())
		})
	})

	Context("NewProductionSafetyRules", func() {
		It("should return strict rules for production", func() {
			prodRules := domain.NewProductionSafetyRules()

			assertSafetyRules(prodRules, safetyRulesExpectations{
				description:            "Production safety rules",
				noSelectStar:           true,
				requireExplicitColumns: true,
				requireWhere:           true,
				requireLimit:           true,
				maxRowsWithoutLimit:    100,
				destructiveOps:         domain.DestructiveForbidden,
			})
		})

		It("should forbid all destructive operations in production", func() {
			prodRules := domain.NewProductionSafetyRules()

			Expect(prodRules.DestructiveOps.AllowsDropTable()).To(BeFalse())
			Expect(prodRules.DestructiveOps.AllowsTruncate()).To(BeFalse())
		})

		It("should enforce both WHERE and LIMIT in production", func() {
			prodRules := domain.NewProductionSafetyRules()

			Expect(prodRules.SafetyRules.WhereRequirement.RequiresOnDestructive()).To(BeTrue())
			Expect(prodRules.SafetyRules.LimitRequirement.RequiresOnSelect()).To(BeTrue())
		})
	})
})

var _ = Describe("Type Safety Benefits", func() {
	Context("Split Brain Elimination", func() {
		It("should prevent boolean flag confusion", func() {
			// Before: NoDropTable and NoTruncate as separate booleans
			// What if NoDropTable=true but NoTruncate=false? Inconsistent policy!

			// After: Single DestructiveOperationPolicy enum
			policy := domain.DestructiveForbidden

			// Clear, consistent policy for ALL destructive operations
			Expect(policy.AllowsDropTable()).To(BeFalse())
			Expect(policy.AllowsTruncate()).To(BeFalse())
		})

		It("should group related safety rules semantically", func() {
			// Style rules (code quality)
			styleRules := domain.QueryStyleRules{
				SelectStarPolicy: domain.SelectStarForbidden,
			}

			// Safety rules (prevent bugs)
			safetyRules := domain.QuerySafetyRules{
				WhereRequirement: domain.WhereClauseAlways,
			}

			// Clear separation of concerns
			Expect(styleRules.SelectStarPolicy.ForbidsSelectStar()).To(BeTrue())
			Expect(safetyRules.WhereRequirement.RequiresOnDestructive()).To(BeTrue())
		})
	})

	Context("Environment-Specific Presets", func() {
		It("should provide different rules for different environments", func() {
			dev := domain.NewDevelopmentSafetyRules()
			prod := domain.NewProductionSafetyRules()

			// Dev is permissive
			Expect(dev.DestructiveOps).To(Equal(domain.DestructiveAllowed))
			Expect(dev.SafetyRules.WhereRequirement.RequiresOnDestructive()).To(BeFalse())

			// Prod is strict
			Expect(prod.DestructiveOps).To(Equal(domain.DestructiveForbidden))
			Expect(prod.SafetyRules.WhereRequirement.RequiresOnDestructive()).To(BeTrue())
			Expect(prod.SafetyRules.LimitRequirement.RequiresOnSelect()).To(BeTrue())
		})
	})

	Context("Uint Type Safety", func() {
		It("should prevent negative row limits", func() {
			rules := domain.QuerySafetyRules{
				MaxRowsWithoutLimit: 100,
			}

			// Before: int could be negative (-100 doesn't make sense!)
			// After: uint cannot be negative (type-safe at compile time)
			Expect(rules.MaxRowsWithoutLimit).To(BeNumerically(">=", 0))
		})
	})
})
