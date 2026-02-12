package testing

import (
	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/domain"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// ValidationTestSuite defines interface for types that need validation testing.
type ValidationTestSuite[T interface {
	IsValid() bool
	String() string
}] interface {
	GetValidValues() []T
	GetInvalidValues() []T
	GetTypeName() string
}

// TestValidationSuite runs generic validation tests for any type implementing ValidationTestSuite.
func TestValidationSuite[T interface {
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

// RunBooleanMethodTest runs generic tests for boolean methods.
func RunBooleanMethodTest(context string, trueModes, falseModes []string, method func(string) bool, methodDisplay string) {
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

// RunStringRepresentationTest runs generic tests for String() method of enums.
func RunStringRepresentationTest(enumTestCases []EnumTestCase) {
	for _, testCase := range enumTestCases {
		Expect(testCase.EnumValue.String()).To(Equal(testCase.ExpectedString))
	}
}

// EnumTestCase represents a test case for enum string representation.
type EnumTestCase struct {
	EnumValue      interface{ String() string }
	ExpectedString string
}

// AssertProductionSafetyRules validates that safety rules have production-safe defaults.
// This helper ensures consistent validation across integration and domain tests.
func AssertProductionSafetyRules(rules domain.TypeSafeSafetyRules, description string) {
	By(description + " - verifying code quality rules")
	Expect(rules.StyleRules.SelectStarPolicy.ForbidsSelectStar()).
		To(BeTrue(), description+": should forbid SELECT * for code quality")

	By(description + " - verifying bug prevention rules")
	Expect(rules.SafetyRules.WhereRequirement.RequiresOnDestructive()).
		To(BeTrue(), description+": should require WHERE on destructive operations")

	By(description + " - verifying security policy rules")
	Expect(rules.DestructiveOps).
		To(Equal(domain.DestructiveForbidden), description+": should forbid destructive operations")
}

// ValidProjectTypes contains all valid project types for testing purposes.
var ValidProjectTypes = []generated.ProjectType{
	generated.ProjectTypeHobby,
	generated.ProjectTypeMicroservice,
	generated.ProjectTypeEnterprise,
	generated.ProjectTypeAPIFirst,
	generated.ProjectTypeAnalytics,
	generated.ProjectTypeTesting,
	generated.ProjectTypeMultiTenant,
	generated.ProjectTypeLibrary,
}

// ValidDatabaseTypes contains all valid database types for testing purposes.
var ValidDatabaseTypes = []generated.DatabaseType{
	generated.DatabaseTypePostgreSQL,
	generated.DatabaseTypeMySQL,
	generated.DatabaseTypeSQLite,
}

// ProjectTypeTestSuite implements ValidationTestSuite for ProjectType validation tests.
type ProjectTypeTestSuite struct{}

func (s ProjectTypeTestSuite) GetValidValues() []generated.ProjectType   { return ValidProjectTypes }
func (s ProjectTypeTestSuite) GetInvalidValues() []generated.ProjectType { return []generated.ProjectType{generated.ProjectType("invalid-type")} }
func (s ProjectTypeTestSuite) GetTypeName() string                       { return "ProjectType" }

// DatabaseTypeTestSuite implements ValidationTestSuite for DatabaseType validation tests.
type DatabaseTypeTestSuite struct{}

func (s DatabaseTypeTestSuite) GetValidValues() []generated.DatabaseType   { return ValidDatabaseTypes }
func (s DatabaseTypeTestSuite) GetInvalidValues() []generated.DatabaseType { return []generated.DatabaseType{generated.DatabaseType("invalid-db")} }
func (s DatabaseTypeTestSuite) GetTypeName() string                        { return "DatabaseType" }

// ValidateAllProjectTypes tests that all project types in ValidProjectTypes are valid.
// This helper eliminates duplicate validation code across test files.
func ValidateAllProjectTypes() {
	for _, projectType := range ValidProjectTypes {
		Expect(projectType.IsValid()).To(BeTrue(),
			"Project type %s should be valid", projectType)
	}
}

// ValidateAllDatabaseTypes tests that all database types in ValidDatabaseTypes are valid.
// This helper eliminates duplicate validation code across test files.
func ValidateAllDatabaseTypes() {
	for _, dbType := range ValidDatabaseTypes {
		Expect(dbType.IsValid()).To(BeTrue(),
			"Database type %s should be valid", dbType)
	}
}
