package validation_test

import (
	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/domain"
	testingHelper "github.com/LarsArtmann/SQLC-Wizzard/internal/testing"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/validation"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// runRuleTransformationTest runs a generic test for safety rule transformation.
func runRuleTransformationTest(
	transformer *validation.RuleTransformer,
	ruleName, expectedRule, expectedMessage string,
	setupRules func() *generated.SafetyRules,
) {
	rules := setupRules()
	configRules := transformer.TransformSafetyRules(rules)

	Expect(configRules).To(HaveLen(1))
	Expect(configRules[0].Name).To(Equal(ruleName))
	Expect(configRules[0].Rule).To(Equal(expectedRule))
	Expect(configRules[0].Message).To(Equal(expectedMessage))
}

// expectSingleRule verifies that a single rule with the expected name and rule expression is generated.
func expectSingleRule(
	transformer *validation.RuleTransformer,
	rules *domain.TypeSafeSafetyRules,
	expectedName, expectedRule string,
) {
	configRules := transformer.TransformTypeSafeSafetyRules(rules)

	Expect(configRules).To(HaveLen(1))
	Expect(configRules[0].Name).To(Equal(expectedName))
	Expect(configRules[0].Rule).To(Equal(expectedRule))
}

// assertIdenticalRules checks that boolean-based and type-safe rule transformations produce identical results.
func assertIdenticalRules(boolResult, typeSafeResult []generated.RuleConfig, expectedRule string) {
	Expect(boolResult).To(HaveLen(1))
	Expect(typeSafeResult).To(HaveLen(1))
	Expect(boolResult[0].Rule).To(Equal(typeSafeResult[0].Rule))
	Expect(boolResult[0].Name).To(Equal(typeSafeResult[0].Name))
	Expect(boolResult[0].Rule).To(Equal(expectedRule))
}

// toGeneratedSafetyRules converts domain.SafetyRules to generated.SafetyRules for testing.
//
// Deprecated: domain.SafetyRules is the legacy boolean version. New tests should use
// the TypeSafeSafetyRules and the TypeSafeRuleTransformer API instead.
func toGeneratedSafetyRules(rules *domain.SafetyRules) *generated.SafetyRules {
	if rules == nil {
		return nil
	}

	return &generated.SafetyRules{
		NoSelectStar: rules.NoSelectStar,
		RequireWhere: rules.RequireWhere,
		RequireLimit: rules.RequireLimit,
		Rules:        rules.Rules,
	}
}

// runParityTest runs a parity test between boolean and type-safe safety rules.
func runParityTest(
	transformer *validation.RuleTransformer,
	boolRules *generated.SafetyRules,
	setupTypeSafeRules func() *domain.TypeSafeSafetyRules,
	expectedExpression string,
) {
	boolResult := transformer.TransformSafetyRules(boolRules)
	typeSafeRules := setupTypeSafeRules()
	typeSafeResult := transformer.TransformTypeSafeSafetyRules(typeSafeRules)
	assertIdenticalRules(boolResult, typeSafeResult, expectedExpression)
}

// runFieldParityTest is a generic helper for testing parity between boolean flags and type-safe enums.
// It reduces duplication for tests that verify boolean flags produce identical expressions to type-safe enum values.
func runFieldParityTest(
	transformer *validation.RuleTransformer,
	boolRules *generated.SafetyRules,
	setTypeSafeField func(*domain.TypeSafeSafetyRules),
	expectedExpression string,
) {
	runParityTest(
		transformer,
		boolRules,
		func() *domain.TypeSafeSafetyRules {
			rules := testingHelper.CreateBaseTypeSafeSafetyRules()
			setTypeSafeField(rules)

			return rules
		},
		expectedExpression,
	)
}

// runBooleanFlagParityTest is a parameterised helper for testing that boolean safety rule flags
// produce identical expressions to their type-safe enum equivalents.
// This eliminates repetitive test boilerplate for NoSelectStar, RequireWhere, and RequireLimit.
func runBooleanFlagParityTest(
	transformer *validation.RuleTransformer,
	testName string,
	isSelectStar, isWhere, isLimit bool,
	typeSafeFieldSetter func(*domain.TypeSafeSafetyRules),
	expectedExpression string,
) {
	It(testName, func() {
		runFieldParityTest(
			transformer,
			&generated.SafetyRules{
				NoSelectStar: isSelectStar,
				RequireWhere: isWhere,
				RequireLimit: isLimit,
			},
			typeSafeFieldSetter,
			expectedExpression,
		)
	})
}

// newRuleTransformerForParity creates a fresh RuleTransformer for parity tests.
// Parity tests run in their own Describe block, so they need their own
// transformer instance independent of the BeforeEach hook in the main suite.
func newRuleTransformerForParity() *validation.RuleTransformer {
	return validation.NewRuleTransformer()
}

// runMultiRuleTransformationTest is a parameterized helper for testing transformation of multiple
// safety rules. It eliminates duplication for tests that verify boolean flags produce correct
// rule sets with specific expected names in a defined order.
func runMultiRuleTransformationTest(
	transformer *validation.RuleTransformer,
	setupRules func() *generated.SafetyRules,
	transformFunc func(*validation.RuleTransformer, *generated.SafetyRules) []generated.RuleConfig,
	expectedRuleNames []string,
) {
	configRules := transformFunc(transformer, setupRules())

	Expect(configRules).To(HaveLen(len(expectedRuleNames)))

	for i, expectedName := range expectedRuleNames {
		Expect(configRules[i].Name).To(Equal(expectedName))
	}
}
