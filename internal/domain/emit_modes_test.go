package domain_test

import (
	"github.com/LarsArtmann/SQLC-Wizzard/internal/domain"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// Test cases for TypeSafeEmitOptions and related enums
// Run via TestDomain in domain_test.go

// ValidationTestSuite defines interface for types that need validation testing
type ValidationTestSuite[T interface {
	IsValid() bool
	String() string
}] interface {
	GetValidValues() []T
	GetInvalidValues() []T
	GetTypeName() string
}

// testValidationSuite runs generic validation tests for any type implementing ValidationTestSuite
func testValidationSuite[T interface {
	IsValid() bool
	String() string
}](suite ValidationTestSuite[T]) {
	Context("IsValid", func() {
		It("should validate all defined "+suite.GetTypeName(), func() {
			validValues := suite.GetValidValues()
			for _, value := range validValues {
				Expect(value.IsValid()).To(BeTrue(), "%s %s should be valid", suite.GetTypeName(), value)
			}
		})

		It("should reject invalid "+suite.GetTypeName(), func() {
			invalidValues := suite.GetInvalidValues()
			for _, value := range invalidValues {
				Expect(value.IsValid()).To(BeFalse(), "%s %s should be invalid", suite.GetTypeName(), value)
			}
		})
	})
}

// NullHandlingMode validation test suite
type NullHandlingModeTestSuite struct{}

func (NullHandlingModeTestSuite) GetValidValues() []domain.NullHandlingMode {
	return []domain.NullHandlingMode{
		domain.NullHandlingPointers,
		domain.NullHandlingEmptySlices,
		domain.NullHandlingExplicitNull,
		domain.NullHandlingMixed,
	}
}

func (NullHandlingModeTestSuite) GetInvalidValues() []domain.NullHandlingMode {
	return []domain.NullHandlingMode{
		"invalid",
		"",
		"pointer",
	}
}

func (NullHandlingModeTestSuite) GetTypeName() string {
	return "NullHandlingMode"
}

// StructPointerMode validation test suite
type StructPointerModeTestSuite struct{}

func (StructPointerModeTestSuite) GetValidValues() []domain.StructPointerMode {
	return []domain.StructPointerMode{
		domain.StructPointerNever,
		domain.StructPointerResults,
		domain.StructPointerParams,
		domain.StructPointerAlways,
	}
}

func (StructPointerModeTestSuite) GetInvalidValues() []domain.StructPointerMode {
	return []domain.StructPointerMode{
		"invalid",
		"",
		"sometimes",
	}
}

func (StructPointerModeTestSuite) GetTypeName() string {
	return "StructPointerMode"
}

// JSONTagStyle validation test suite
type JSONTagStyleTestSuite struct{}

func (JSONTagStyleTestSuite) GetValidValues() []domain.JSONTagStyle {
	return []domain.JSONTagStyle{
		domain.JSONTagStyleCamel,
		domain.JSONTagStyleSnake,
		domain.JSONTagStylePascal,
		domain.JSONTagStyleKebab,
	}
}

func (JSONTagStyleTestSuite) GetInvalidValues() []domain.JSONTagStyle {
	return []domain.JSONTagStyle{
		"invalid",
		"",
		"UPPER",
	}
}

func (JSONTagStyleTestSuite) GetTypeName() string {
	return "JSONTagStyle"
}

// runBooleanMethodTest runs generic tests for boolean methods
func runBooleanMethodTest(context string, trueModes []string, falseModes []string, method func(string) bool, methodDisplay string) {
	It("should return true for "+context, func() {
		for _, mode := range trueModes {
			Expect(method(mode)).To(BeTrue(), "Mode %s should return true for "+context, mode)
		}
	})

	It("should return false for modes without "+context, func() {
		for _, mode := range falseModes {
			Expect(method(mode)).To(BeFalse(), "Mode %s should return false for "+context, mode)
		}
	})
}

var _ = Describe("NullHandlingMode", func() {
	// Use generic validation test suite
testValidationSuite(NullHandlingModeTestSuite{})

	runBooleanMethodTest("pointer modes", 
	[]string{"pointers", "mixed"}, 
	[]string{"empty_slices", "explicit_null"}, 
	func(mode string) bool { return domain.NullHandlingMode(mode).UsePointers() }, 
	"UsePointers")

	Context("UseEmptySlices", func() {
		It("should return true only for empty slices mode", func() {
			Expect(domain.NullHandlingEmptySlices.UseEmptySlices()).To(BeTrue())
		})

		It("should return false for other modes", func() {
			Expect(domain.NullHandlingPointers.UseEmptySlices()).To(BeFalse())
			Expect(domain.NullHandlingExplicitNull.UseEmptySlices()).To(BeFalse())
			Expect(domain.NullHandlingMixed.UseEmptySlices()).To(BeFalse())
		})
	})

	Context("String", func() {
		It("should return correct string representation", func() {
			Expect(domain.NullHandlingPointers.String()).To(Equal("pointers"))
			Expect(domain.NullHandlingEmptySlices.String()).To(Equal("empty_slices"))
			Expect(domain.NullHandlingExplicitNull.String()).To(Equal("explicit_null"))
			Expect(domain.NullHandlingMixed.String()).To(Equal("mixed"))
		})
	})
})

var _ = Describe("EnumGenerationMode", func() {
	Context("IsValid", func() {
		It("should validate all defined modes", func() {
			validModes := []domain.EnumGenerationMode{
				domain.EnumGenerationBasic,
				domain.EnumGenerationWithValidation,
				domain.EnumGenerationComplete,
			}

			for _, mode := range validModes {
				Expect(mode.IsValid()).To(BeTrue(), "Mode %s should be valid", mode)
			}
		})

		It("should reject invalid modes", func() {
			invalidModes := []domain.EnumGenerationMode{
				"invalid",
				"",
				"advanced",
			}

			for _, mode := range invalidModes {
				Expect(mode.IsValid()).To(BeFalse(), "Mode %s should be invalid", mode)
			}
		})
	})

	Context("IncludesValidation", func() {
		It("should return true for modes with validation", func() {
			Expect(domain.EnumGenerationWithValidation.IncludesValidation()).To(BeTrue())
			Expect(domain.EnumGenerationComplete.IncludesValidation()).To(BeTrue())
		})

		It("should return false for basic mode", func() {
			Expect(domain.EnumGenerationBasic.IncludesValidation()).To(BeFalse())
		})
	})

	Context("IncludesAllValues", func() {
		It("should return true only for complete mode", func() {
			Expect(domain.EnumGenerationComplete.IncludesAllValues()).To(BeTrue())
		})

		It("should return false for other modes", func() {
			Expect(domain.EnumGenerationBasic.IncludesAllValues()).To(BeFalse())
			Expect(domain.EnumGenerationWithValidation.IncludesAllValues()).To(BeFalse())
		})
	})

	Context("String", func() {
		It("should return correct string representation", func() {
			Expect(domain.EnumGenerationBasic.String()).To(Equal("basic"))
			Expect(domain.EnumGenerationWithValidation.String()).To(Equal("with_validation"))
			Expect(domain.EnumGenerationComplete.String()).To(Equal("complete"))
		})
	})
})

var _ = Describe("StructPointerMode", func() {
	// Use generic validation test suite
testValidationSuite(StructPointerModeTestSuite{})

	runBooleanMethodTest("modes with result pointers", 
	[]string{"results", "always"}, 
	[]string{"never", "params"}, 
	func(mode string) bool { return domain.StructPointerMode(mode).UseResultPointers() }, 
	"UseResultPointers")

	runBooleanMethodTest("param pointers",
	[]string{"params", "always"},
	[]string{"never", "results"},
	func(mode string) bool { return domain.StructPointerMode(mode).UseParamPointers() },
	"UseParamPointers")

	Context("String", func() {
		It("should return correct string representation", func() {
			Expect(domain.StructPointerNever.String()).To(Equal("never"))
			Expect(domain.StructPointerResults.String()).To(Equal("results"))
			Expect(domain.StructPointerParams.String()).To(Equal("params"))
			Expect(domain.StructPointerAlways.String()).To(Equal("always"))
		})
	})
})

var _ = Describe("JSONTagStyle", func() {
	// Use generic validation test suite
testValidationSuite(JSONTagStyleTestSuite{})

	Context("String", func() {
		It("should return correct string representation", func() {
			Expect(domain.JSONTagStyleCamel.String()).To(Equal("camel"))
			Expect(domain.JSONTagStyleSnake.String()).To(Equal("snake"))
			Expect(domain.JSONTagStylePascal.String()).To(Equal("pascal"))
			Expect(domain.JSONTagStyleKebab.String()).To(Equal("kebab"))
		})
	})
})

var _ = Describe("TypeSafeEmitOptions", func() {
	Context("IsValid", func() {
		It("should validate options with all valid modes", func() {
			opts := domain.TypeSafeEmitOptions{
				NullHandling:   domain.NullHandlingPointers,
				EnumMode:       domain.EnumGenerationComplete,
				StructPointers: domain.StructPointerNever,
				JSONTagStyle:   domain.JSONTagStyleCamel,
			}

			err := opts.IsValid()
			Expect(err).NotTo(HaveOccurred())
		})

		It("should reject invalid null handling mode", func() {
			opts := domain.TypeSafeEmitOptions{
				NullHandling:   "invalid",
				EnumMode:       domain.EnumGenerationComplete,
				StructPointers: domain.StructPointerNever,
				JSONTagStyle:   domain.JSONTagStyleCamel,
			}

			err := opts.IsValid()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("NullHandling"))
		})

		It("should reject invalid enum mode", func() {
			opts := domain.TypeSafeEmitOptions{
				NullHandling:   domain.NullHandlingPointers,
				EnumMode:       "invalid",
				StructPointers: domain.StructPointerNever,
				JSONTagStyle:   domain.JSONTagStyleCamel,
			}

			err := opts.IsValid()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("EnumMode"))
		})

		It("should reject invalid struct pointer mode", func() {
			opts := domain.TypeSafeEmitOptions{
				NullHandling:   domain.NullHandlingPointers,
				EnumMode:       domain.EnumGenerationComplete,
				StructPointers: "invalid",
				JSONTagStyle:   domain.JSONTagStyleCamel,
			}

			err := opts.IsValid()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("StructPointers"))
		})

		It("should reject invalid JSON tag style", func() {
			opts := domain.TypeSafeEmitOptions{
				NullHandling:   domain.NullHandlingPointers,
				EnumMode:       domain.EnumGenerationComplete,
				StructPointers: domain.StructPointerNever,
				JSONTagStyle:   "invalid",
			}

			err := opts.IsValid()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("JSONTagStyle"))
		})
	})

	Context("NewTypeSafeEmitOptions", func() {
		It("should return valid default options", func() {
			defaults := domain.NewTypeSafeEmitOptions()

			Expect(defaults.NullHandling).To(Equal(domain.NullHandlingPointers))
			Expect(defaults.EnumMode).To(Equal(domain.EnumGenerationComplete))
			Expect(defaults.StructPointers).To(Equal(domain.StructPointerNever))
			Expect(defaults.JSONTagStyle).To(Equal(domain.JSONTagStyleCamel))

			err := defaults.IsValid()
			Expect(err).NotTo(HaveOccurred())
		})

		It("should enable recommended features by default", func() {
			defaults := domain.NewTypeSafeEmitOptions()

			Expect(defaults.Features.GenerateJSONTags).To(BeTrue())
			Expect(defaults.Features.GeneratePreparedQueries).To(BeTrue())
			Expect(defaults.Features.GenerateInterface).To(BeTrue())
			Expect(defaults.Features.UseExactTableNames).To(BeFalse())
		})
	})
})

var _ = Describe("CodeGenerationFeatures", func() {
	It("should allow independent feature configuration", func() {
		features := domain.CodeGenerationFeatures{
			GenerateJSONTags:        true,
			GeneratePreparedQueries: false,
			GenerateInterface:       true,
			UseExactTableNames:      false,
		}

		Expect(features.GenerateJSONTags).To(BeTrue())
		Expect(features.GeneratePreparedQueries).To(BeFalse())
		Expect(features.GenerateInterface).To(BeTrue())
		Expect(features.UseExactTableNames).To(BeFalse())
	})
})

var _ = Describe("Type Safety Benefits", func() {
	Context("Invalid State Prevention", func() {
		It("should prevent incompatible null handling combinations", func() {
			// Before: Could have EmitEmptySlices=true AND EmitPointersForNull=true (contradictory)
			// After: Can only choose ONE NullHandlingMode

			opts := domain.TypeSafeEmitOptions{
				NullHandling: domain.NullHandlingPointers, // Clear, unambiguous choice
			}

			Expect(opts.NullHandling).To(Equal(domain.NullHandlingPointers))
			Expect(opts.NullHandling.UsePointers()).To(BeTrue())
			Expect(opts.NullHandling.UseEmptySlices()).To(BeFalse())
		})

		It("should provide semantic meaning to options", func() {
			// Before: EmitEnumValidMethod=true, EmitAllEnumValues=true (what does this mean?)
			// After: EnumGenerationComplete (clear semantic meaning)

			opts := domain.TypeSafeEmitOptions{
				EnumMode: domain.EnumGenerationComplete,
			}

			Expect(opts.EnumMode.IncludesValidation()).To(BeTrue())
			Expect(opts.EnumMode.IncludesAllValues()).To(BeTrue())
		})

		It("should group related options semantically", func() {
			// Before: EmitResultStructPointers and EmitParamsStructPointers separate
			// After: Single StructPointerMode with clear options

			optsNever := domain.TypeSafeEmitOptions{
				StructPointers: domain.StructPointerNever,
			}

			optsAlways := domain.TypeSafeEmitOptions{
				StructPointers: domain.StructPointerAlways,
			}

			Expect(optsNever.StructPointers.UseResultPointers()).To(BeFalse())
			Expect(optsNever.StructPointers.UseParamPointers()).To(BeFalse())

			Expect(optsAlways.StructPointers.UseResultPointers()).To(BeTrue())
			Expect(optsAlways.StructPointers.UseParamPointers()).To(BeTrue())
		})
	})

	Context("State Space Reduction", func() {
		It("should reduce valid state space dramatically", func() {
			// Before: 8 booleans = 256 possible states (most invalid)
			// After: 4 enums + 4 bool flags = ~80 states (all valid)

			// All combinations of our new enums are semantically valid
			nullModes := []domain.NullHandlingMode{
				domain.NullHandlingPointers,
				domain.NullHandlingEmptySlices,
				domain.NullHandlingExplicitNull,
				domain.NullHandlingMixed,
			}

			enumModes := []domain.EnumGenerationMode{
				domain.EnumGenerationBasic,
				domain.EnumGenerationWithValidation,
				domain.EnumGenerationComplete,
			}

			// Every combination is valid and has clear semantics
			for _, nullMode := range nullModes {
				for _, enumMode := range enumModes {
					opts := domain.TypeSafeEmitOptions{
						NullHandling:   nullMode,
						EnumMode:       enumMode,
						StructPointers: domain.StructPointerNever,
						JSONTagStyle:   domain.JSONTagStyleCamel,
					}

					Expect(opts.IsValid()).NotTo(HaveOccurred())
				}
			}
		})
	})
})
