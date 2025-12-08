package validation_test

import (
	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/domain"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/validation"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// Integration tests for end-to-end flows using new type-safe types
// These tests prove the ghost system is no longer a ghost - it's integrated!

var _ = Describe("Type-Safe Integration Tests", func() {
	var transformer *validation.RuleTransformer

	BeforeEach(func() {
		transformer = validation.NewRuleTransformer()
	})

	Context("End-to-End Flow: New TypeSafe → RuleConfigs", func() {
		It("should transform production-safe rules end-to-end", func() {
			// Step 1: Create type-safe rules using production defaults
			typeSafeRules := domain.NewTypeSafeSafetyRules()

			// Step 2: Transform to configuration format
			configRules := transformer.TransformTypeSafeSafetyRules(&typeSafeRules)

			// Step 3: Verify correct CEL rules are generated
			Expect(configRules).To(HaveLen(5))

			// Verify style rules
			var noSelectStarFound bool
			for _, rule := range configRules {
				if rule.Name == "no-select-star" {
					noSelectStarFound = true
					Expect(rule.Rule).To(ContainSubstring("SELECT *"))
					Expect(rule.Message).To(ContainSubstring("explicit column names"))
				}
			}
			Expect(noSelectStarFound).To(BeTrue(), "no-select-star rule should be present")

			// Verify safety rules
			var requireWhereFound bool
			for _, rule := range configRules {
				if rule.Name == "require-where" {
					requireWhereFound = true
					Expect(rule.Rule).To(ContainSubstring("hasWhereClause"))
				}
			}
			Expect(requireWhereFound).To(BeTrue(), "require-where rule should be present")

			// Verify destructive policy
			var noDropTableFound, noTruncateFound bool
			for _, rule := range configRules {
				if rule.Name == "no-drop-table" {
					noDropTableFound = true
				}
				if rule.Name == "no-truncate" {
					noTruncateFound = true
				}
			}
			Expect(noDropTableFound).To(BeTrue(), "no-drop-table rule should be present")
			Expect(noTruncateFound).To(BeTrue(), "no-truncate rule should be present")
		})

		It("should handle development environment end-to-end", func() {
			// Step 1: Create dev-friendly rules
			devRules := domain.NewDevelopmentSafetyRules()

			// Step 2: Transform to configuration format
			configRules := transformer.TransformTypeSafeSafetyRules(&devRules)

			// Step 3: Verify no restrictive rules in dev mode
			Expect(configRules).To(BeEmpty(), "Development mode should have no safety rules")
		})

		It("should handle strict production environment end-to-end", func() {
			// Step 1: Create strict production rules
			prodRules := domain.NewProductionSafetyRules()

			// Step 2: Transform to configuration format
			configRules := transformer.TransformTypeSafeSafetyRules(&prodRules)

			// Step 3: Verify all strict rules are present
			Expect(configRules).To(HaveLen(7))

			ruleNames := make(map[string]bool)
			for _, rule := range configRules {
				ruleNames[rule.Name] = true
			}

			// Verify comprehensive rule coverage
			Expect(ruleNames["no-select-star"]).To(BeTrue())
			Expect(ruleNames["require-explicit-columns"]).To(BeTrue())
			Expect(ruleNames["require-where"]).To(BeTrue())
			Expect(ruleNames["require-limit"]).To(BeTrue())
			Expect(ruleNames["max-rows-without-limit"]).To(BeTrue())
			Expect(ruleNames["no-drop-table"]).To(BeTrue())
			Expect(ruleNames["no-truncate"]).To(BeTrue())
		})
	})

	Context("Migration Path: Old → TypeSafe → RuleConfigs", func() {
		It("should convert legacy rules to type-safe and transform", func() {
			// Step 1: Start with old legacy rules (simulating existing code)
			oldRules := generated.DefaultSafetyRules()

			// Step 2: Convert to type-safe format
			typeSafeRules := domain.SafetyRulesToTypeSafe(oldRules)

			// Step 3: Transform to configuration format
			configRules := transformer.TransformTypeSafeSafetyRules(&typeSafeRules)

			// Step 4: Verify rules are correctly transformed
			Expect(configRules).ToNot(BeEmpty())

			// Should have NoSelectStar and RequireWhere from defaults
			ruleNames := make(map[string]bool)
			for _, rule := range configRules {
				ruleNames[rule.Name] = true
			}

			Expect(ruleNames["no-select-star"]).To(BeTrue())
			Expect(ruleNames["require-where"]).To(BeTrue())
		})

		It("should maintain compatibility with legacy transformation", func() {
			// Step 1: Create identical rules in both formats
			oldRules := &generated.SafetyRules{
				NoSelectStar: true,
				RequireWhere: true,
				NoDropTable:  true,
				NoTruncate:   true,
				RequireLimit: false,
			}

			// Step 2: Transform using legacy method
			legacyResult := transformer.TransformSafetyRules(oldRules)

			// Step 3: Convert old → new and transform using new method
			typeSafeRules := domain.SafetyRulesToTypeSafe(*oldRules)
			typeSafeResult := transformer.TransformTypeSafeSafetyRules(&typeSafeRules)

			// Step 4: Verify both produce same core rules
			// Note: Type-safe version may have additional rules (MaxRowsWithoutLimit)
			// but should have all the same base rules

			legacyRuleNames := make(map[string]bool)
			for _, rule := range legacyResult {
				legacyRuleNames[rule.Name] = true
			}

			typeSafeRuleNames := make(map[string]bool)
			for _, rule := range typeSafeResult {
				typeSafeRuleNames[rule.Name] = true
			}

			// All legacy rules should be in type-safe result
			for name := range legacyRuleNames {
				Expect(typeSafeRuleNames[name]).To(BeTrue(), "Rule %s should be present in type-safe result", name)
			}
		})
	})

	Context("Custom Rules Integration", func() {
		It("should preserve custom rules through full migration path", func() {
			// Step 1: Start with old rules including custom CEL rules
			oldRules := generated.SafetyRules{
				NoSelectStar: true,
				RequireWhere: true,
				Rules: []generated.SafetyRule{
					{
						Name:    "no-cross-db-joins",
						Rule:    "query.databases().count() <= 1",
						Message: "Cross-database joins not allowed",
					},
					{
						Name:    "max-subqueries",
						Rule:    "query.subqueries().count() <= 2",
						Message: "Maximum 2 subqueries allowed",
					},
				},
			}

			// Step 2: Convert to type-safe
			typeSafeRules := domain.SafetyRulesToTypeSafe(oldRules)

			// Step 3: Transform to configuration
			configRules := transformer.TransformTypeSafeSafetyRules(&typeSafeRules)

			// Step 4: Verify custom rules are preserved
			customRulesFound := 0
			for _, rule := range configRules {
				if rule.Name == "no-cross-db-joins" {
					customRulesFound++
					Expect(rule.Rule).To(Equal("query.databases().count() <= 1"))
					Expect(rule.Message).To(Equal("Cross-database joins not allowed"))
				}
				if rule.Name == "max-subqueries" {
					customRulesFound++
					Expect(rule.Rule).To(Equal("query.subqueries().count() <= 2"))
					Expect(rule.Message).To(Equal("Maximum 2 subqueries allowed"))
				}
			}

			Expect(customRulesFound).To(Equal(2), "Both custom rules should be preserved")
		})
	})

	Context("Type Safety Benefits in Practice", func() {
		It("should prevent negative row limits at compile time", func() {
			// This test demonstrates the uint type safety benefit
			rules := &domain.TypeSafeSafetyRules{
				StyleRules: domain.QueryStyleRules{},
				SafetyRules: domain.QuerySafetyRules{
					MaxRowsWithoutLimit: 100, // uint can't be negative
				},
				DestructiveOps: domain.DestructiveAllowed,
			}

			configRules := transformer.TransformTypeSafeSafetyRules(rules)

			// Verify the rule was created with positive limit
			var maxRowsRule *generated.RuleConfig
			for i := range configRules {
				if configRules[i].Name == "max-rows-without-limit" {
					maxRowsRule = &configRules[i]
					break
				}
			}

			Expect(maxRowsRule).NotTo(BeNil())
			Expect(maxRowsRule.Rule).To(ContainSubstring("100"))

			// The compiler would reject: MaxRowsWithoutLimit: -100
			// This is a compile-time safety guarantee!
		})

		It("should enforce consistent destructive operation policy", func() {
			// OLD WAY: Could have NoDropTable=true but NoTruncate=false (inconsistent!)
			// NEW WAY: Single policy applies to ALL destructive operations

			rules := &domain.TypeSafeSafetyRules{
				StyleRules:     domain.QueryStyleRules{},
				SafetyRules:    domain.QuerySafetyRules{},
				DestructiveOps: domain.DestructiveForbidden,
			}

			configRules := transformer.TransformTypeSafeSafetyRules(rules)

			// Both DROP and TRUNCATE should be forbidden (consistent policy)
			ruleNames := make(map[string]bool)
			for _, rule := range configRules {
				ruleNames[rule.Name] = true
			}

			Expect(ruleNames["no-drop-table"]).To(BeTrue())
			Expect(ruleNames["no-truncate"]).To(BeTrue())

			// Impossible to have split brain: both are governed by same policy
		})

		It("should provide semantic clarity over boolean flags", func() {
			// OLD WAY: EmitEmptySlices=true, EmitPointers=false (what does this mean?)
			// NEW WAY: NullHandling=NullHandlingEmptySlices (clear semantic meaning)

			// This test demonstrates type-safe safety rules have same benefit
			prodRules := domain.NewProductionSafetyRules()

			// Clear semantic groupings:
			Expect(prodRules.StyleRules.SelectStarPolicy.ForbidsSelectStar()).To(BeTrue())      // Code quality
			Expect(prodRules.SafetyRules.WhereRequirement.RequiresOnDestructive()).To(BeTrue()) // Prevent bugs
			Expect(prodRules.DestructiveOps).To(Equal(domain.DestructiveForbidden))             // Security policy

			// vs old way: NoSelectStar, RequireWhere, NoDropTable, NoTruncate
			// (all flat booleans, no semantic grouping)
		})
	})

	Context("Roundtrip Integration: Old → TypeSafe → Old → RuleConfigs", func() {
		It("should maintain functional equivalence through conversion roundtrip", func() {
			// Step 1: Start with old rules
			original := generated.SafetyRules{
				NoSelectStar: true,
				RequireWhere: true,
				NoDropTable:  true,
				NoTruncate:   true,
				RequireLimit: false,
			}

			// Step 2: Convert to type-safe
			typeSafe := domain.SafetyRulesToTypeSafe(original)

			// Step 3: Convert back to legacy
			roundtrip := typeSafe.ToLegacy()

			// Step 4: Transform both original and roundtrip
			originalRules := transformer.TransformSafetyRules(&original)
			roundtripRules := transformer.TransformSafetyRules(&roundtrip)

			// Step 5: Verify functional equivalence
			Expect(len(originalRules)).To(Equal(len(roundtripRules)))

			for i := range originalRules {
				Expect(roundtripRules[i].Name).To(Equal(originalRules[i].Name))
				Expect(roundtripRules[i].Rule).To(Equal(originalRules[i].Rule))
			}
		})
	})
})
