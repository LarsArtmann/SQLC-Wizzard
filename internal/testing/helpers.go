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

func (s ProjectTypeTestSuite) GetValidValues() []generated.ProjectType { return ValidProjectTypes }
func (s ProjectTypeTestSuite) GetInvalidValues() []generated.ProjectType {
	return []generated.ProjectType{generated.ProjectType("invalid-type")}
}
func (s ProjectTypeTestSuite) GetTypeName() string { return "ProjectType" }

// DatabaseTypeTestSuite implements ValidationTestSuite for DatabaseType validation tests.
type DatabaseTypeTestSuite struct{}

func (s DatabaseTypeTestSuite) GetValidValues() []generated.DatabaseType { return ValidDatabaseTypes }
func (s DatabaseTypeTestSuite) GetInvalidValues() []generated.DatabaseType {
	return []generated.DatabaseType{generated.DatabaseType("invalid-db")}
}
func (s DatabaseTypeTestSuite) GetTypeName() string { return "DatabaseType" }

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

// CreateBaseTypeSafeSafetyRules creates a default TypeSafeSafetyRules structure for test reuse.
// This helper eliminates duplicate fixture code across test files.
func CreateBaseTypeSafeSafetyRules() *domain.TypeSafeSafetyRules {
	return &domain.TypeSafeSafetyRules{
		StyleRules: domain.QueryStyleRules{
			SelectStarPolicy:   domain.SelectStarAllowed,
			ColumnExplicitness: domain.ColumnExplicitnessDefault,
		},
		SafetyRules: domain.QuerySafetyRules{
			WhereRequirement:    domain.WhereClauseNever,
			LimitRequirement:    domain.LimitClauseNever,
			MaxRowsWithoutLimit: 0,
		},
		DestructiveOps: domain.DestructiveAllowed,
		CustomRules:    []generated.SafetyRule{},
	}
}

// CreateTypeSafeSafetyRules creates a TypeSafeSafetyRules with a common test configuration.
// This helper provides a default for most test scenarios and accepts an optional configuration callback.
//
// Example usage:
//
//	rules := testing.CreateTypeSafeSafetyRules(func(r *domain.TypeSafeSafetyRules) {
//	    r.SafetyRules.WhereRequirement = domain.WhereClauseAlways
//	})
func CreateTypeSafeSafetyRules(configure func(*domain.TypeSafeSafetyRules)) *domain.TypeSafeSafetyRules {
	rules := &domain.TypeSafeSafetyRules{
		StyleRules: domain.QueryStyleRules{
			SelectStarPolicy:   domain.SelectStarForbidden,
			ColumnExplicitness: domain.ColumnExplicitnessDefault,
		},
		SafetyRules: domain.QuerySafetyRules{
			WhereRequirement:    domain.WhereClauseOnDestructive,
			LimitRequirement:    domain.LimitClauseOnSelect,
			MaxRowsWithoutLimit: 100,
		},
		DestructiveOps: domain.DestructiveForbidden,
		CustomRules:    []generated.SafetyRule{},
	}
	if configure != nil {
		configure(rules)
	}
	return rules
}

// CreateRestrictiveTypeSafeSafetyRulesWithCustomRules creates a TypeSafeSafetyRules with
// restrictive safety settings and accepts custom rules for testing custom rule preservation.
// This helper eliminates duplicate fixture code for testing custom safety rules.
func CreateRestrictiveTypeSafeSafetyRulesWithCustomRules(customRules []generated.SafetyRule) *domain.TypeSafeSafetyRules {
	return CreateTypeSafeSafetyRules(func(r *domain.TypeSafeSafetyRules) {
		r.SafetyRules.WhereRequirement = domain.WhereClauseAlways
		r.SafetyRules.LimitRequirement = domain.LimitClauseNever
		r.SafetyRules.MaxRowsWithoutLimit = 1000
		r.CustomRules = customRules
	})
}

// CreateStrictTypeSafeSafetyRules creates a strictly configured TypeSafeSafetyRules for testing.
// All safety rules are set to their most restrictive values.
func CreateStrictTypeSafeSafetyRules() *domain.TypeSafeSafetyRules {
	return &domain.TypeSafeSafetyRules{
		StyleRules: domain.QueryStyleRules{
			SelectStarPolicy:   domain.SelectStarForbidden,
			ColumnExplicitness: domain.ColumnExplicitnessRequired,
		},
		SafetyRules: domain.QuerySafetyRules{
			WhereRequirement:    domain.WhereClauseAlways,
			LimitRequirement:    domain.LimitClauseAlways,
			MaxRowsWithoutLimit: 100,
		},
		DestructiveOps: domain.DestructiveForbidden,
		CustomRules:    []generated.SafetyRule{},
	}
}

// CreateQueryStyleRulesForbiddenSelectStar creates QueryStyleRules that forbids SELECT *.
// This helper eliminates duplicate fixture code for testing select star policies.
func CreateQueryStyleRulesForbiddenSelectStar() domain.QueryStyleRules {
	return domain.QueryStyleRules{
		SelectStarPolicy:   domain.SelectStarForbidden,
		ColumnExplicitness: domain.ColumnExplicitnessDefault,
	}
}

// CreateQuerySafetyRulesStrict creates QuerySafetyRules with strict safety settings.
// All safety rules are set to their most restrictive values.
func CreateQuerySafetyRulesStrict() domain.QuerySafetyRules {
	return domain.QuerySafetyRules{
		WhereRequirement:    domain.WhereClauseAlways,
		LimitRequirement:    domain.LimitClauseNever,
		MaxRowsWithoutLimit: 1000,
	}
}

// CreateGeneratedSafetyRulesForbidden creates a SafetyRules with all forbidden flags set.
// This helper eliminates duplicate fixture code across test files.
func CreateGeneratedSafetyRulesForbidden() generated.SafetyRules {
	return generated.SafetyRules{
		NoSelectStar: true,
		RequireWhere: true,
		NoDropTable:  true,
		NoTruncate:   true,
		RequireLimit: false,
		Rules:        []generated.SafetyRule{},
	}
}

// CreateGeneratedSafetyRulesAllowed creates a SafetyRules with all forbidden flags cleared.
// This helper eliminates duplicate fixture code across test files.
func CreateGeneratedSafetyRulesAllowed() generated.SafetyRules {
	return generated.SafetyRules{
		NoSelectStar: false,
		RequireWhere: false,
		NoDropTable:  false,
		NoTruncate:   false,
		RequireLimit: true,
		Rules:        []generated.SafetyRule{},
	}
}

// CreateGeneratedSafetyRulesAllowedWithCustomRules creates a SafetyRules with all forbidden flags
// cleared and accepts custom rules for testing custom rule preservation.
// This helper eliminates duplicate fixture code for testing custom safety rules.
func CreateGeneratedSafetyRulesAllowedWithCustomRules(customRules []generated.SafetyRule) *domain.SafetyRules {
	return &domain.SafetyRules{
		NoSelectStar: false,
		RequireWhere: false,
		NoDropTable:  false,
		NoTruncate:   false,
		RequireLimit: false,
		Rules:        customRules,
	}
}

// GetNullHandlingModeTestCases returns test cases for NullHandlingMode string representation tests.
// This helper eliminates duplicate test case definitions across test files.
func GetNullHandlingModeTestCases() []EnumTestCase {
	return []EnumTestCase{
		{EnumValue: domain.NullHandlingPointers, ExpectedString: "pointers"},
		{EnumValue: domain.NullHandlingEmptySlices, ExpectedString: "empty_slices"},
		{EnumValue: domain.NullHandlingExplicitNull, ExpectedString: "explicit_null"},
		{EnumValue: domain.NullHandlingMixed, ExpectedString: "mixed"},
	}
}

// GetStructPointerModeTestCases returns test cases for StructPointerMode string representation tests.
// This helper eliminates duplicate test case definitions across test files.
func GetStructPointerModeTestCases() []EnumTestCase {
	return []EnumTestCase{
		{EnumValue: domain.StructPointerNever, ExpectedString: "never"},
		{EnumValue: domain.StructPointerResults, ExpectedString: "results"},
		{EnumValue: domain.StructPointerParams, ExpectedString: "params"},
		{EnumValue: domain.StructPointerAlways, ExpectedString: "always"},
	}
}

// GetJSONTagStyleTestCases returns test cases for JSONTagStyle string representation tests.
// This helper eliminates duplicate test case definitions across test files.
func GetJSONTagStyleTestCases() []EnumTestCase {
	return []EnumTestCase{
		{EnumValue: domain.JSONTagStyleCamel, ExpectedString: "camel"},
		{EnumValue: domain.JSONTagStyleSnake, ExpectedString: "snake"},
		{EnumValue: domain.JSONTagStylePascal, ExpectedString: "pascal"},
		{EnumValue: domain.JSONTagStyleKebab, ExpectedString: "kebab"},
	}
}
