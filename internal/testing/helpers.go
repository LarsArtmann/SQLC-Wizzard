package testing

import (
	. "github.com/onsi/gomega"
)

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

// runStringRepresentationTest runs generic tests for String() method of enums
func runStringRepresentationTest(testCases map[string]string) {
	for enumValue, expectedString := range testCases {
		Expect(enumValue.String()).To(Equal(expectedString))
	}
}