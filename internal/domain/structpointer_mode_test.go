package domain_test

import (
	"github.com/LarsArtmann/SQLC-Wizzard/internal/domain"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/testing"
	. "github.com/onsi/ginkgo/v2"
)

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

var _ = Describe("StructPointerMode", func() {
	// Use generic validation test suite
	testing.TestValidationSuite(StructPointerModeTestSuite{})

	testing.RunBooleanMethodTest("modes with result pointers",
		[]string{"results", "always"},
		[]string{"never", "params"},
		func(mode string) bool { return domain.StructPointerMode(mode).UseResultPointers() },
		"UseResultPointers")

	testing.RunBooleanMethodTest("param pointers",
		[]string{"params", "always"},
		[]string{"never", "results"},
		func(mode string) bool { return domain.StructPointerMode(mode).UseParamPointers() },
		"UseParamPointers")

	Context("String", func() {
		It("should return correct string representation", func() {
			testing.RunStringRepresentationTest([]testing.EnumTestCase{
				{EnumValue: domain.StructPointerNever, ExpectedString: "never"},
				{EnumValue: domain.StructPointerResults, ExpectedString: "results"},
				{EnumValue: domain.StructPointerParams, ExpectedString: "params"},
				{EnumValue: domain.StructPointerAlways, ExpectedString: "always"},
			})
		})
	})
})
