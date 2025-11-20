package domain_test

import (
	"github.com/LarsArtmann/SQLC-Wizzard/internal/domain"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/testing"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// EnumGenerationMode validation test suite
type EnumGenerationModeTestSuite struct{}

func (EnumGenerationModeTestSuite) GetValidValues() []domain.EnumGenerationMode {
	return []domain.EnumGenerationMode{
		domain.EnumGenerationBasic,
		domain.EnumGenerationWithValidation,
		domain.EnumGenerationComplete,
	}
}

func (EnumGenerationModeTestSuite) GetInvalidValues() []domain.EnumGenerationMode {
	return []domain.EnumGenerationMode{
		"invalid",
		"",
		"advanced",
	}
}

func (EnumGenerationModeTestSuite) GetTypeName() string {
	return "EnumGenerationMode"
}

var _ = Describe("EnumGenerationMode", func() {
	// Use generic validation test suite
	testing.TestValidationSuite(EnumGenerationModeTestSuite{})

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
			testing.RunStringRepresentationTest([]testing.EnumTestCase{
				{EnumValue: domain.EnumGenerationBasic, ExpectedString: "basic"},
				{EnumValue: domain.EnumGenerationWithValidation, ExpectedString: "with_validation"},
				{EnumValue: domain.EnumGenerationComplete, ExpectedString: "complete"},
			})
		})
	})
})