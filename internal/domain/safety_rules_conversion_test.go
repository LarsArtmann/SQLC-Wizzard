package domain_test

import (
	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/domain"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/testing"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("SafetyRules Conversions", func() {
	Context("SafetyRulesToTypeSafe", func() {
		It("should convert forbidden destructive ops correctly", func() {
			old := testing.CreateGeneratedSafetyRulesForbidden()

			typeSafe := domain.SafetyRulesToTypeSafe(old)

			Expect(typeSafe.StyleRules.SelectStarPolicy.ForbidsSelectStar()).To(BeTrue())
			Expect(typeSafe.StyleRules.ColumnExplicitness.RequiresExplicitColumns()).To(BeFalse())
			Expect(typeSafe.SafetyRules.WhereRequirement.RequiresOnDestructive()).To(BeTrue())
			Expect(typeSafe.SafetyRules.LimitRequirement.RequiresOnSelect()).To(BeFalse())
			Expect(typeSafe.SafetyRules.MaxRowsWithoutLimit).To(Equal(uint(1000)))
			Expect(typeSafe.DestructiveOps).To(Equal(domain.DestructiveForbidden))
		})

		It("should convert allowed destructive ops correctly", func() {
			old := testing.CreateGeneratedSafetyRulesAllowed()

			typeSafe := domain.SafetyRulesToTypeSafe(old)

			Expect(typeSafe.StyleRules.SelectStarPolicy.ForbidsSelectStar()).To(BeFalse())
			Expect(typeSafe.SafetyRules.WhereRequirement.RequiresOnDestructive()).To(BeFalse())
			Expect(typeSafe.SafetyRules.LimitRequirement.RequiresOnSelect()).To(BeTrue())
			Expect(typeSafe.DestructiveOps).To(Equal(domain.DestructiveAllowed))
		})

		It("should handle mixed destructive ops by defaulting to forbidden", func() {
			old := generated.SafetyRules{
				NoSelectStar: true,
				RequireWhere: true,
				NoDropTable:  true,
				NoTruncate:   false,
				RequireLimit: false,
				Rules:        []generated.SafetyRule{},
			}

			typeSafe := domain.SafetyRulesToTypeSafe(old)

			// Mixed state should map to forbidden for safety
			Expect(typeSafe.DestructiveOps).To(Equal(domain.DestructiveForbidden))
		})

		It("should preserve custom rules", func() {
			customRules := []generated.SafetyRule{
				{
					Name:    "no-complex-joins",
					Rule:    "query.joins().count() <= 3",
					Message: "Maximum 3 joins allowed",
				},
				{
					Name:    "require-index",
					Rule:    "query.usesIndex()",
					Message: "Query must use an index",
				},
			}

			old := generated.SafetyRules{
				NoSelectStar: true,
				RequireWhere: true,
				NoDropTable:  true,
				NoTruncate:   true,
				RequireLimit: false,
				Rules:        customRules,
			}

			typeSafe := domain.SafetyRulesToTypeSafe(old)

			Expect(typeSafe.CustomRules).To(Equal(customRules))
			Expect(typeSafe.CustomRules).To(HaveLen(2))
			Expect(typeSafe.CustomRules[0].Name).To(Equal("no-complex-joins"))
			Expect(typeSafe.CustomRules[1].Name).To(Equal("require-index"))
		})
	})

	Context("TypeSafeSafetyRules.ToLegacy", func() {
		It("should convert back to legacy format correctly", func() {
			typeSafe := *testing.CreateStrictTypeSafeSafetyRules()

			legacy := typeSafe.ToLegacy()

			Expect(legacy.NoSelectStar).To(BeTrue())
			Expect(legacy.RequireWhere).To(BeTrue())
			Expect(legacy.NoDropTable).To(BeTrue())
			Expect(legacy.NoTruncate).To(BeTrue())
			Expect(legacy.RequireLimit).To(BeTrue())
		})

		It("should convert allowed destructive ops correctly", func() {
			typeSafe := *testing.CreateBaseTypeSafeSafetyRules()

			legacy := typeSafe.ToLegacy()

			Expect(legacy.NoDropTable).To(BeFalse())
			Expect(legacy.NoTruncate).To(BeFalse())
		})

		It("should preserve custom rules", func() {
			customRules := []generated.SafetyRule{
				{
					Name:    "test-rule",
					Rule:    "query.type == 'SELECT'",
					Message: "Only SELECT allowed",
				},
			}

			typeSafe := testing.CreateRestrictiveTypeSafeSafetyRulesWithCustomRules(customRules)

			legacy := typeSafe.ToLegacy()

			Expect(legacy.Rules).To(Equal(customRules))
		})
	})

	Context("Roundtrip Conversions", func() {
		It("should preserve data through old→new→old conversion", func() {
			original := generated.SafetyRules{
				NoSelectStar: true,
				RequireWhere: true,
				NoDropTable:  true,
				NoTruncate:   true,
				RequireLimit: false,
				Rules: []generated.SafetyRule{
					{
						Name:    "test",
						Rule:    "true",
						Message: "test message",
					},
				},
			}

			typeSafe := domain.SafetyRulesToTypeSafe(original)
			roundtrip := typeSafe.ToLegacy()

			Expect(roundtrip).To(Equal(original))
		})

		It("should preserve core data through new→old→new conversion", func() {
			original := *testing.CreateRestrictiveTypeSafeSafetyRulesWithCustomRules(nil)

			legacy := original.ToLegacy()
			roundtrip := domain.SafetyRulesToTypeSafe(legacy)

			Expect(
				roundtrip.StyleRules.SelectStarPolicy.ForbidsSelectStar(),
			).To(Equal(original.StyleRules.SelectStarPolicy.ForbidsSelectStar()))
			Expect(
				roundtrip.SafetyRules.WhereRequirement.RequiresOnDestructive(),
			).To(Equal(original.SafetyRules.WhereRequirement.RequiresOnDestructive()))
			Expect(
				roundtrip.SafetyRules.LimitRequirement.RequiresOnSelect(),
			).To(Equal(original.SafetyRules.LimitRequirement.RequiresOnSelect()))
			Expect(roundtrip.DestructiveOps).To(Equal(original.DestructiveOps))

			// Note: RequireExplicitColumns and MaxRowsWithoutLimit are not preserved
			// because they don't exist in the legacy format
		})
	})

	Context("NewTypeSafeSafetyRulesFromLegacy", func() {
		It("should be a convenience wrapper for conversion", func() {
			old := generated.DefaultSafetyRules()
			typeSafe1 := domain.SafetyRulesToTypeSafe(old)
			typeSafe2 := domain.NewTypeSafeSafetyRulesFromLegacy(old)

			Expect(typeSafe1).To(Equal(typeSafe2))
		})
	})
})

var _ = Describe("ParseJSONTagStyle", func() {
	It("should parse valid styles correctly", func() {
		Expect(domain.ParseJSONTagStyle("camel")).To(Equal(domain.JSONTagStyleCamel))
		Expect(domain.ParseJSONTagStyle("snake")).To(Equal(domain.JSONTagStyleSnake))
		Expect(domain.ParseJSONTagStyle("pascal")).To(Equal(domain.JSONTagStylePascal))
		Expect(domain.ParseJSONTagStyle("kebab")).To(Equal(domain.JSONTagStyleKebab))
	})

	It("should return invalid style for unknown strings", func() {
		style := domain.ParseJSONTagStyle("INVALID")
		Expect(style.IsValid()).To(BeFalse())
	})
})
