package domain_test

import (
	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/domain"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// Test cases for TypeSafeSafetyRules and related types
// Run via TestDomain in domain_test.go

// safetyRulesExpectations defines the expected values for safety rules testing
type safetyRulesExpectations struct {
	description                     string
	noSelectStar                   bool
	requireExplicitColumns         bool
	requireWhere                   bool
	requireLimit                   bool
	maxRowsWithoutLimit            uint
	destructiveOps                 domain.DestructiveOperationPolicy
}

// assertSafetyRules validates safety rules against expected values
func assertSafetyRules(rules interface{}, expectations safetyRulesExpectations) {
	switch r := rules.(type) {
	case domain.TypeSafeSafetyRules:
		Expect(r.StyleRules.NoSelectStar).To(Equal(expectations.noSelectStar), expectations.description+": NoSelectStar")
		Expect(r.StyleRules.RequireExplicitColumns).To(Equal(expectations.requireExplicitColumns), expectations.description+": RequireExplicitColumns")
		Expect(r.SafetyRules.RequireWhere).To(Equal(expectations.requireWhere), expectations.description+": RequireWhere")
		Expect(r.SafetyRules.RequireLimit).To(Equal(expectations.requireLimit), expectations.description+": RequireLimit")
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
			NoSelectStar:           true,
			RequireExplicitColumns: false,
		}

		Expect(rules.NoSelectStar).To(BeTrue())
		Expect(rules.RequireExplicitColumns).To(BeFalse())
	})

	It("should support both rules enabled", func() {
		rules := domain.QueryStyleRules{
			NoSelectStar:           true,
			RequireExplicitColumns: true,
		}

		Expect(rules.NoSelectStar).To(BeTrue())
		Expect(rules.RequireExplicitColumns).To(BeTrue())
	})
})

var _ = Describe("QuerySafetyRules", func() {
	It("should allow independent configuration of safety rules", func() {
		rules := domain.QuerySafetyRules{
			RequireWhere:        true,
			RequireLimit:        false,
			MaxRowsWithoutLimit: 1000,
		}

		Expect(rules.RequireWhere).To(BeTrue())
		Expect(rules.RequireLimit).To(BeFalse())
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

var _ = Describe("DestructiveOperationPolicy", func() {
	Context("IsValid", func() {
		It("should validate all defined policies", func() {
			validPolicies := []domain.DestructiveOperationPolicy{
				domain.DestructiveAllowed,
				domain.DestructiveWithConfirmation,
				domain.DestructiveForbidden,
			}

			for _, policy := range validPolicies {
				Expect(policy.IsValid()).To(BeTrue(), "Policy %s should be valid", policy)
			}
		})

		It("should reject invalid policies", func() {
			invalidPolicies := []domain.DestructiveOperationPolicy{
				"invalid",
				"",
				"sometimes",
			}

			for _, policy := range invalidPolicies {
				Expect(policy.IsValid()).To(BeFalse(), "Policy %s should be invalid", policy)
			}
		})
	})

	Context("AllowsDropTable", func() {
		It("should allow DROP TABLE only when policy is 'allowed'", func() {
			Expect(domain.DestructiveAllowed.AllowsDropTable()).To(BeTrue())
		})

		It("should not allow DROP TABLE for other policies", func() {
			Expect(domain.DestructiveWithConfirmation.AllowsDropTable()).To(BeFalse())
			Expect(domain.DestructiveForbidden.AllowsDropTable()).To(BeFalse())
		})
	})

	Context("AllowsTruncate", func() {
		It("should allow TRUNCATE only when policy is 'allowed'", func() {
			Expect(domain.DestructiveAllowed.AllowsTruncate()).To(BeTrue())
		})

		It("should not allow TRUNCATE for other policies", func() {
			Expect(domain.DestructiveWithConfirmation.AllowsTruncate()).To(BeFalse())
			Expect(domain.DestructiveForbidden.AllowsTruncate()).To(BeFalse())
		})
	})

	Context("RequiresConfirmation", func() {
		It("should require confirmation only for 'with_confirmation' policy", func() {
			Expect(domain.DestructiveWithConfirmation.RequiresConfirmation()).To(BeTrue())
		})

		It("should not require confirmation for other policies", func() {
			Expect(domain.DestructiveAllowed.RequiresConfirmation()).To(BeFalse())
			Expect(domain.DestructiveForbidden.RequiresConfirmation()).To(BeFalse())
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
					NoSelectStar: true,
				},
				SafetyRules: domain.QuerySafetyRules{
					RequireWhere: true,
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
				description:               "Default safety rules",
				noSelectStar:             true,
				requireExplicitColumns:   false,
				requireWhere:             true,
				requireLimit:             false,
				maxRowsWithoutLimit:      1000,
				destructiveOps:           domain.DestructiveForbidden,
			})
		})

		It("should have safe production defaults", func() {
			defaults := domain.NewTypeSafeSafetyRules()

			// Production-safe: No SELECT *
			Expect(defaults.StyleRules.NoSelectStar).To(BeTrue())

			// Production-safe: Require WHERE
			Expect(defaults.SafetyRules.RequireWhere).To(BeTrue())

			// Production-safe: Block destructive ops
			Expect(defaults.DestructiveOps).To(Equal(domain.DestructiveForbidden))
		})
	})

	Context("NewDevelopmentSafetyRules", func() {
		It("should return relaxed rules for development", func() {
			devRules := domain.NewDevelopmentSafetyRules()

			Expect(devRules.StyleRules.NoSelectStar).To(BeFalse())
			Expect(devRules.SafetyRules.RequireWhere).To(BeFalse())
			Expect(devRules.SafetyRules.RequireLimit).To(BeFalse())
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
				description:               "Production safety rules",
				noSelectStar:             true,
				requireExplicitColumns:   true,
				requireWhere:             true,
				requireLimit:             true,
				maxRowsWithoutLimit:      100,
				destructiveOps:           domain.DestructiveForbidden,
			})
		})

		It("should forbid all destructive operations in production", func() {
			prodRules := domain.NewProductionSafetyRules()

			Expect(prodRules.DestructiveOps.AllowsDropTable()).To(BeFalse())
			Expect(prodRules.DestructiveOps.AllowsTruncate()).To(BeFalse())
		})

		It("should enforce both WHERE and LIMIT in production", func() {
			prodRules := domain.NewProductionSafetyRules()

			Expect(prodRules.SafetyRules.RequireWhere).To(BeTrue())
			Expect(prodRules.SafetyRules.RequireLimit).To(BeTrue())
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
				NoSelectStar: true,
			}

			// Safety rules (prevent bugs)
			safetyRules := domain.QuerySafetyRules{
				RequireWhere: true,
			}

			// Clear separation of concerns
			Expect(styleRules.NoSelectStar).To(BeTrue())
			Expect(safetyRules.RequireWhere).To(BeTrue())
		})
	})

	Context("Environment-Specific Presets", func() {
		It("should provide different rules for different environments", func() {
			dev := domain.NewDevelopmentSafetyRules()
			prod := domain.NewProductionSafetyRules()

			// Dev is permissive
			Expect(dev.DestructiveOps).To(Equal(domain.DestructiveAllowed))
			Expect(dev.SafetyRules.RequireWhere).To(BeFalse())

			// Prod is strict
			Expect(prod.DestructiveOps).To(Equal(domain.DestructiveForbidden))
			Expect(prod.SafetyRules.RequireWhere).To(BeTrue())
			Expect(prod.SafetyRules.RequireLimit).To(BeTrue())
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
