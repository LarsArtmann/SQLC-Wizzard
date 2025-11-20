package domain_test

import (
	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/domain"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// Test cases for bidirectional conversions between old and new types
// Run via TestDomain in domain_test.go

// emitOptionsTestCase represents a test case for EmitOptions conversion
type emitOptionsTestCase struct {
	description         string
	input              generated.EmitOptions
	expectedNullHandling domain.NullHandlingMode
	expectedPointers    domain.StructPointerMode
	expectedJSONStyle  domain.JSONTagStyle
}

// runEmitOptionsTest runs conversion test with given input and expected values
func runEmitOptionsTest(testCase emitOptionsTestCase) {
	It("should convert "+testCase.description+" correctly", func() {
		typeSafe := domain.EmitOptionsToTypeSafe(testCase.input)

		Expect(typeSafe.NullHandling).To(Equal(testCase.expectedNullHandling))
		Expect(typeSafe.StructPointers).To(Equal(testCase.expectedPointers))
		Expect(typeSafe.JSONTagStyle).To(Equal(testCase.expectedJSONStyle))
	})
}

var _ = Describe("EmitOptions Conversions", func() {
	Context("EmitOptionsToTypeSafe", func() {
		It("should convert empty_slices mode correctly", func() {
			old := generated.EmitOptions{
				EmitJSONTags:             true,
				EmitPreparedQueries:      true,
				EmitInterface:            true,
				EmitEmptySlices:          true,
				EmitResultStructPointers: false,
				EmitParamsStructPointers: false,
				EmitEnumValidMethod:      true,
				EmitAllEnumValues:        true,
				JSONTagsCaseStyle:        "camel",
			}

			typeSafe := domain.EmitOptionsToTypeSafe(old)

			Expect(typeSafe.NullHandling).To(Equal(domain.NullHandlingEmptySlices))
			Expect(typeSafe.EnumMode).To(Equal(domain.EnumGenerationComplete))
			Expect(typeSafe.StructPointers).To(Equal(domain.StructPointerNever))
			Expect(typeSafe.JSONTagStyle).To(Equal(domain.JSONTagStyleCamel))
			Expect(typeSafe.Features.GenerateJSONTags).To(BeTrue())
			Expect(typeSafe.Features.GeneratePreparedQueries).To(BeTrue())
			Expect(typeSafe.Features.GenerateInterface).To(BeTrue())
		})

		It("should convert pointers mode correctly", func() {
			old := generated.EmitOptions{
				EmitJSONTags:             true,
				EmitPreparedQueries:      false,
				EmitInterface:            false,
				EmitEmptySlices:          false,
				EmitResultStructPointers: true,
				EmitParamsStructPointers: true,
				EmitEnumValidMethod:      false,
				EmitAllEnumValues:        false,
				JSONTagsCaseStyle:        "snake",
			}

			typeSafe := domain.EmitOptionsToTypeSafe(old)

			Expect(typeSafe.NullHandling).To(Equal(domain.NullHandlingPointers))
			Expect(typeSafe.EnumMode).To(Equal(domain.EnumGenerationBasic))
			Expect(typeSafe.StructPointers).To(Equal(domain.StructPointerAlways))
			Expect(typeSafe.JSONTagStyle).To(Equal(domain.JSONTagStyleSnake))
		})

		runEmitOptionsTest(emitOptionsTestCase{
		description: "explicit_null mode",
		input: generated.EmitOptions{
			EmitJSONTags:             false,
			EmitPreparedQueries:      false,
			EmitInterface:            false,
			EmitEmptySlices:          false,
			EmitResultStructPointers: false,
			EmitParamsStructPointers: false,
			EmitEnumValidMethod:      true,
			EmitAllEnumValues:        false,
			JSONTagsCaseStyle:        "pascal",
		},
		expectedNullHandling: domain.NullHandlingExplicitNull,
		expectedPointers:    domain.StructPointerNever,
		expectedJSONStyle:  domain.JSONTagStylePascal,
	})

		runEmitOptionsTest(emitOptionsTestCase{
		description: "mixed mode",
		input: generated.EmitOptions{
			EmitJSONTags:             true,
			EmitPreparedQueries:      true,
			EmitInterface:            true,
			EmitEmptySlices:          false,
			EmitResultStructPointers: true,
			EmitParamsStructPointers: false,
			EmitEnumValidMethod:      true,
			EmitAllEnumValues:        true,
			JSONTagsCaseStyle:        "kebab",
		},
		expectedNullHandling: domain.NullHandlingMixed,
		expectedPointers:    domain.StructPointerResults,
		expectedJSONStyle:  domain.JSONTagStyleKebab,
	})

		It("should handle invalid JSON tag style with fallback", func() {
			old := generated.EmitOptions{
				EmitJSONTags:             true,
				EmitPreparedQueries:      true,
				EmitInterface:            true,
				EmitEmptySlices:          true,
				EmitResultStructPointers: false,
				EmitParamsStructPointers: false,
				EmitEnumValidMethod:      true,
				EmitAllEnumValues:        true,
				JSONTagsCaseStyle:        "INVALID",
			}

			typeSafe := domain.EmitOptionsToTypeSafe(old)

			// Should fallback to camel case for invalid style
			Expect(typeSafe.JSONTagStyle).To(Equal(domain.JSONTagStyleCamel))
		})
	})

	Context("TypeSafeEmitOptions.ToLegacy", func() {
		It("should convert back to legacy format correctly", func() {
			typeSafe := domain.TypeSafeEmitOptions{
				NullHandling:   domain.NullHandlingPointers,
				EnumMode:       domain.EnumGenerationComplete,
				StructPointers: domain.StructPointerNever,
				JSONTagStyle:   domain.JSONTagStyleCamel,
				Features: domain.CodeGenerationFeatures{
					GenerateJSONTags:        true,
					GeneratePreparedQueries: true,
					GenerateInterface:       true,
					UseExactTableNames:      false,
				},
			}

			legacy := typeSafe.ToTemplateData()

			Expect(legacy.EmitJSONTags).To(BeTrue())
			Expect(legacy.EmitPreparedQueries).To(BeTrue())
			Expect(legacy.EmitInterface).To(BeTrue())
			Expect(legacy.EmitEmptySlices).To(BeFalse())
			Expect(legacy.EmitResultStructPointers).To(BeFalse())
			Expect(legacy.EmitParamsStructPointers).To(BeFalse())
			Expect(legacy.EmitEnumValidMethod).To(BeTrue())
			Expect(legacy.EmitAllEnumValues).To(BeTrue())
			Expect(legacy.JSONTagsCaseStyle).To(Equal("camel"))
		})

		It("should convert empty_slices mode back correctly", func() {
			typeSafe := domain.TypeSafeEmitOptions{
				NullHandling:   domain.NullHandlingEmptySlices,
				EnumMode:       domain.EnumGenerationBasic,
				StructPointers: domain.StructPointerAlways,
				JSONTagStyle:   domain.JSONTagStyleSnake,
				Features: domain.CodeGenerationFeatures{
					GenerateJSONTags:        false,
					GeneratePreparedQueries: false,
					GenerateInterface:       false,
					UseExactTableNames:      true,
				},
			}

			legacy := typeSafe.ToTemplateData()

			Expect(legacy.EmitEmptySlices).To(BeTrue())
			Expect(legacy.EmitResultStructPointers).To(BeTrue())
			Expect(legacy.EmitParamsStructPointers).To(BeTrue())
			Expect(legacy.EmitEnumValidMethod).To(BeFalse())
			Expect(legacy.EmitAllEnumValues).To(BeFalse())
			Expect(legacy.JSONTagsCaseStyle).To(Equal("snake"))
		})
	})

	Context("Roundtrip Conversions", func() {
		It("should preserve data through old→new→old conversion", func() {
			original := generated.EmitOptions{
				EmitJSONTags:             true,
				EmitPreparedQueries:      true,
				EmitInterface:            true,
				EmitEmptySlices:          true,
				EmitResultStructPointers: false,
				EmitParamsStructPointers: false,
				EmitEnumValidMethod:      true,
				EmitAllEnumValues:        true,
				JSONTagsCaseStyle:        "camel",
			}

			typeSafe := domain.EmitOptionsToTypeSafe(original)
			roundtrip := domain.EmitOptionsToTypeSafe(typeSafe.ToTemplateData())

			Expect(roundtrip).To(Equal(typeSafe))
		})

		It("should preserve data through new→old→new conversion for representable modes", func() {
			// Note: NullHandlingPointers and NullHandlingExplicitNull both map to the same
			// legacy representation (all false), so we use NullHandlingEmptySlices which is
			// distinctly representable in the old format
			original := domain.TypeSafeEmitOptions{
				NullHandling:   domain.NullHandlingEmptySlices,
				EnumMode:       domain.EnumGenerationComplete,
				StructPointers: domain.StructPointerNever,
				JSONTagStyle:   domain.JSONTagStyleCamel,
				Features: domain.CodeGenerationFeatures{
					GenerateJSONTags:        true,
					GeneratePreparedQueries: true,
					GenerateInterface:       true,
					UseExactTableNames:      false,
				},
			}

			legacy := original.ToTemplateData()
			roundtrip := domain.EmitOptionsToTypeSafe(legacy)

			Expect(roundtrip.NullHandling).To(Equal(original.NullHandling))
			Expect(roundtrip.EnumMode).To(Equal(original.EnumMode))
			Expect(roundtrip.StructPointers).To(Equal(original.StructPointers))
			Expect(roundtrip.JSONTagStyle).To(Equal(original.JSONTagStyle))
			// Note: UseExactTableNames is not preserved (doesn't exist in legacy)
			Expect(roundtrip.Features.GenerateJSONTags).To(Equal(original.Features.GenerateJSONTags))
			Expect(roundtrip.Features.GeneratePreparedQueries).To(Equal(original.Features.GeneratePreparedQueries))
			Expect(roundtrip.Features.GenerateInterface).To(Equal(original.Features.GenerateInterface))
		})

		It("should document lossy conversion for NullHandlingPointers vs ExplicitNull", func() {
			// NullHandlingPointers and NullHandlingExplicitNull both map to the same
			// legacy representation, so information is lost
			pointers := domain.TypeSafeEmitOptions{
				NullHandling:   domain.NullHandlingPointers,
				EnumMode:       domain.EnumGenerationComplete,
				StructPointers: domain.StructPointerNever,
				JSONTagStyle:   domain.JSONTagStyleCamel,
				Features:       domain.CodeGenerationFeatures{},
			}

			explicitNull := domain.TypeSafeEmitOptions{
				NullHandling:   domain.NullHandlingExplicitNull,
				EnumMode:       domain.EnumGenerationComplete,
				StructPointers: domain.StructPointerNever,
				JSONTagStyle:   domain.JSONTagStyleCamel,
				Features:       domain.CodeGenerationFeatures{},
			}

			// Both convert to the same legacy representation
			legacyPointers := pointers.ToTemplateData()
			legacyExplicit := explicitNull.ToTemplateData()

			Expect(legacyPointers).To(Equal(legacyExplicit))

			// Both convert back to NullHandlingExplicitNull (information lost)
			roundtripPointers := domain.EmitOptionsToTypeSafe(legacyPointers)
			roundtripExplicit := domain.EmitOptionsToTypeSafe(legacyExplicit)

			Expect(roundtripPointers.NullHandling).To(Equal(domain.NullHandlingExplicitNull))
			Expect(roundtripExplicit.NullHandling).To(Equal(domain.NullHandlingExplicitNull))
		})
	})

	Context("NewTypeSafeEmitOptionsFromLegacy", func() {
		It("should be a convenience wrapper for conversion", func() {
			old := generated.DefaultEmitOptions()
			typeSafe1 := domain.EmitOptionsToTypeSafe(old)
			typeSafe2 := domain.NewTypeSafeEmitOptionsFromLegacy(old)

			Expect(typeSafe1).To(Equal(typeSafe2))
		})
	})
})

var _ = Describe("SafetyRules Conversions", func() {
	Context("SafetyRulesToTypeSafe", func() {
		It("should convert forbidden destructive ops correctly", func() {
			old := generated.SafetyRules{
				NoSelectStar: true,
				RequireWhere: true,
				NoDropTable:  true,
				NoTruncate:   true,
				RequireLimit: false,
				Rules:        []generated.SafetyRule{},
			}

			typeSafe := domain.SafetyRulesToTypeSafe(old)

			Expect(typeSafe.StyleRules.NoSelectStar).To(BeTrue())
			Expect(typeSafe.StyleRules.RequireExplicitColumns).To(BeFalse())
			Expect(typeSafe.SafetyRules.RequireWhere).To(BeTrue())
			Expect(typeSafe.SafetyRules.RequireLimit).To(BeFalse())
			Expect(typeSafe.SafetyRules.MaxRowsWithoutLimit).To(Equal(uint(1000)))
			Expect(typeSafe.DestructiveOps).To(Equal(domain.DestructiveForbidden))
		})

		It("should convert allowed destructive ops correctly", func() {
			old := generated.SafetyRules{
				NoSelectStar: false,
				RequireWhere: false,
				NoDropTable:  false,
				NoTruncate:   false,
				RequireLimit: true,
				Rules:        []generated.SafetyRule{},
			}

			typeSafe := domain.SafetyRulesToTypeSafe(old)

			Expect(typeSafe.StyleRules.NoSelectStar).To(BeFalse())
			Expect(typeSafe.SafetyRules.RequireWhere).To(BeFalse())
			Expect(typeSafe.SafetyRules.RequireLimit).To(BeTrue())
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
			Expect(len(typeSafe.CustomRules)).To(Equal(2))
			Expect(typeSafe.CustomRules[0].Name).To(Equal("no-complex-joins"))
			Expect(typeSafe.CustomRules[1].Name).To(Equal("require-index"))
		})
	})

	Context("TypeSafeSafetyRules.ToLegacy", func() {
		It("should convert back to legacy format correctly", func() {
			typeSafe := domain.TypeSafeSafetyRules{
				StyleRules: domain.QueryStyleRules{
					NoSelectStar:           true,
					RequireExplicitColumns: true,
				},
				SafetyRules: domain.QuerySafetyRules{
					RequireWhere:        true,
					RequireLimit:        true,
					MaxRowsWithoutLimit: 100,
				},
				DestructiveOps: domain.DestructiveForbidden,
				CustomRules:    []generated.SafetyRule{},
			}

			legacy := typeSafe.ToLegacy()

			Expect(legacy.NoSelectStar).To(BeTrue())
			Expect(legacy.RequireWhere).To(BeTrue())
			Expect(legacy.NoDropTable).To(BeTrue())
			Expect(legacy.NoTruncate).To(BeTrue())
			Expect(legacy.RequireLimit).To(BeTrue())
		})

		It("should convert allowed destructive ops correctly", func() {
			typeSafe := domain.TypeSafeSafetyRules{
				StyleRules: domain.QueryStyleRules{
					NoSelectStar:           false,
					RequireExplicitColumns: false,
				},
				SafetyRules: domain.QuerySafetyRules{
					RequireWhere:        false,
					RequireLimit:        false,
					MaxRowsWithoutLimit: 0,
				},
				DestructiveOps: domain.DestructiveAllowed,
				CustomRules:    []generated.SafetyRule{},
			}

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

			typeSafe := domain.TypeSafeSafetyRules{
				StyleRules: domain.QueryStyleRules{
					NoSelectStar:           true,
					RequireExplicitColumns: false,
				},
				SafetyRules: domain.QuerySafetyRules{
					RequireWhere:        true,
					RequireLimit:        false,
					MaxRowsWithoutLimit: 1000,
				},
				DestructiveOps: domain.DestructiveForbidden,
				CustomRules:    customRules,
			}

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
			original := domain.TypeSafeSafetyRules{
				StyleRules: domain.QueryStyleRules{
					NoSelectStar:           true,
					RequireExplicitColumns: false,
				},
				SafetyRules: domain.QuerySafetyRules{
					RequireWhere:        true,
					RequireLimit:        false,
					MaxRowsWithoutLimit: 1000,
				},
				DestructiveOps: domain.DestructiveForbidden,
				CustomRules:    []generated.SafetyRule{},
			}

			legacy := original.ToLegacy()
			roundtrip := domain.SafetyRulesToTypeSafe(legacy)

			Expect(roundtrip.StyleRules.NoSelectStar).To(Equal(original.StyleRules.NoSelectStar))
			Expect(roundtrip.SafetyRules.RequireWhere).To(Equal(original.SafetyRules.RequireWhere))
			Expect(roundtrip.SafetyRules.RequireLimit).To(Equal(original.SafetyRules.RequireLimit))
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
