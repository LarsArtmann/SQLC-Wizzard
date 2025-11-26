package domain_test

import (
	"github.com/LarsArtmann/SQLC-Wizzard/internal/domain"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/testing"
	. "github.com/onsi/ginkgo/v2"
)

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

var _ = Describe("JSONTagStyle", func() {
	// Use generic validation test suite
	testing.TestValidationSuite(JSONTagStyleTestSuite{})

	Context("String", func() {
		It("should return correct string representation", func() {
			testing.RunStringRepresentationTest([]testing.EnumTestCase{
				{EnumValue: domain.JSONTagStyleCamel, ExpectedString: "camel"},
				{EnumValue: domain.JSONTagStyleSnake, ExpectedString: "snake"},
				{EnumValue: domain.JSONTagStylePascal, ExpectedString: "pascal"},
				{EnumValue: domain.JSONTagStyleKebab, ExpectedString: "kebab"},
			})
		})
	})
})
