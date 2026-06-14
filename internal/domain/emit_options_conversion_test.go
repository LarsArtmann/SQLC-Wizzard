package domain_test

import (
	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/domain"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/testing"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("EmitOptions Conversions", func() {
	Context("EmitOptionsToTypeSafe", func() {
		It("should convert empty_slices mode correctly", func() {
			old := emitOptionsEmptySlices()

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
			old := emitOptionsPointers()

			typeSafe := domain.EmitOptionsToTypeSafe(old)

			Expect(typeSafe.NullHandling).To(Equal(domain.NullHandlingPointers))
			Expect(typeSafe.EnumMode).To(Equal(domain.EnumGenerationBasic))
			Expect(typeSafe.StructPointers).To(Equal(domain.StructPointerAlways))
			Expect(typeSafe.JSONTagStyle).To(Equal(domain.JSONTagStyleSnake))
		})

		It("should handle invalid JSON tag style with fallback", func() {
			old := commonEmitOptions()
			old.JSONTagsCaseStyle = "INVALID"

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
			original := testing.CreateFullEmitOptions()

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
			Expect(
				roundtrip.Features.GenerateJSONTags,
			).To(Equal(original.Features.GenerateJSONTags))
			Expect(
				roundtrip.Features.GeneratePreparedQueries,
			).To(Equal(original.Features.GeneratePreparedQueries))
			Expect(
				roundtrip.Features.GenerateInterface,
			).To(Equal(original.Features.GenerateInterface))
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
