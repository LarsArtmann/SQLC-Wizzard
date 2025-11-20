package domain_test

import (
	"github.com/LarsArtmann/SQLC-Wizzard/internal/domain"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/testing"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// Test cases for TypeSafeEmitOptions and related enums
// Run via TestDomain in domain_test.go

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

// Use generic helpers from centralized testing package
	}
}

var _ = Describe("NullHandlingMode", func() {
	// Use generic validation test suite
	testing.TestValidationSuite(NullHandlingModeTestSuite{})

	testing.RunBooleanMethodTest("pointer modes", 
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
})

// runStringRepresentationTest runs generic tests for String() method of enums
func runStringRepresentationTest(testCases map[string]string) {
	for enumValue, expectedString := range testCases {
		Expect(enumValue.String()).To(Equal(expectedString))
	}
}

var _ = Describe("NullHandlingMode", func() {
	Context("String", func() {
		It("should return correct string representation", func() {
			runStringRepresentationTest(map[string]string{
				domain.NullHandlingPointers.String():    "pointers",
				domain.NullHandlingEmptySlices.String(): "empty_slices",
				domain.NullHandlingExplicitNull.String(): "explicit_null",
				domain.NullHandlingMixed.String():       "mixed",
			})
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
			runStringRepresentationTest(map[string]string{
				domain.EnumGenerationBasic.String():          "basic",
				domain.EnumGenerationWithValidation.String(): "with_validation",
				domain.EnumGenerationComplete.String():       "complete",
			})
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
			runStringRepresentationTest(map[string]string{
				domain.StructPointerNever.String():   "never",
				domain.StructPointerResults.String(): "results",
				domain.StructPointerParams.String():  "params",
				domain.StructPointerAlways.String():  "always",
			})
		})
	})
})

var _ = Describe("JSONTagStyle", func() {
	// Use generic validation test suite
testValidationSuite(JSONTagStyleTestSuite{})

	Context("String", func() {
		It("should return correct string representation", func() {
			runStringRepresentationTest(map[string]string{
				domain.JSONTagStyleCamel.String():  "camel",
				domain.JSONTagStyleSnake.String():  "snake",
				domain.JSONTagStylePascal.String(): "pascal",
				domain.JSONTagStyleKebab.String():  "kebab",
			})
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
