package domain_test

import (
	"github.com/LarsArtmann/SQLC-Wizzard/internal/domain"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/testing"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

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
		"unknown",
	}
}

func (NullHandlingModeTestSuite) GetTypeName() string {
	return "NullHandlingMode"
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

	Context("UseExplicitNull", func() {
		It("should return true only for explicit null mode", func() {
			Expect(domain.NullHandlingExplicitNull.UseExplicitNull()).To(BeTrue())
		})

		It("should return false for other modes", func() {
			Expect(domain.NullHandlingPointers.UseExplicitNull()).To(BeFalse())
			Expect(domain.NullHandlingEmptySlices.UseExplicitNull()).To(BeFalse())
			Expect(domain.NullHandlingMixed.UseExplicitNull()).To(BeFalse())
		})
	})

	Context("String", func() {
		It("should return correct string representation", func() {
			testing.RunStringRepresentationTest([]testing.EnumTestCase{
				{EnumValue: domain.NullHandlingPointers, ExpectedString: "pointers"},
				{EnumValue: domain.NullHandlingEmptySlices, ExpectedString: "empty_slices"},
				{EnumValue: domain.NullHandlingExplicitNull, ExpectedString: "explicit_null"},
				{EnumValue: domain.NullHandlingMixed, ExpectedString: "mixed"},
			})
		})
	})
})