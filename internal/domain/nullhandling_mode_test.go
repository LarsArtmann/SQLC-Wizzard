package domain_test

import (
	"github.com/LarsArtmann/SQLC-Wizzard/internal/domain"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/testing"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// NullHandlingMode validation test suite.
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

	Context("Mode-specific methods", func() {
		type modeMethodTest struct {
			methodName    string
			assertTrueFor domain.NullHandlingMode
			methodCaller  func(domain.NullHandlingMode) bool
		}

		DescribeTable("should return true only for specific mode",
			func(tc modeMethodTest) {
				Expect(tc.methodCaller(tc.assertTrueFor)).To(BeTrue())

				allModes := []domain.NullHandlingMode{
					domain.NullHandlingPointers,
					domain.NullHandlingEmptySlices,
					domain.NullHandlingExplicitNull,
					domain.NullHandlingMixed,
				}

				for _, mode := range allModes {
					if mode != tc.assertTrueFor {
						Expect(tc.methodCaller(mode)).To(BeFalse())
					}
				}
			},
			Entry("UseEmptySlices", modeMethodTest{
				methodName:    "UseEmptySlices",
				assertTrueFor: domain.NullHandlingEmptySlices,
				methodCaller:  func(m domain.NullHandlingMode) bool { return m.UseEmptySlices() },
			}),
			Entry("UseExplicitNull", modeMethodTest{
				methodName:    "UseExplicitNull",
				assertTrueFor: domain.NullHandlingExplicitNull,
				methodCaller:  func(m domain.NullHandlingMode) bool { return m.UseExplicitNull() },
			}),
		)
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
