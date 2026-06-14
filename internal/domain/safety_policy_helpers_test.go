package domain_test

import (
	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/domain"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/testing"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// safetyRulesExpectations defines the expected values for safety rules testing.
type safetyRulesExpectations struct {
	description            string
	noSelectStar           bool
	requireExplicitColumns bool
	requireWhere           bool
	requireLimit           bool
	maxRowsWithoutLimit    uint
	destructiveOps         domain.DestructiveOperationPolicy
}

// assertSafetyRules validates safety rules against expected values.
func assertSafetyRules(rules any, expectations safetyRulesExpectations) {
	switch r := rules.(type) {
	case domain.TypeSafeSafetyRules:
		Expect(
			r.StyleRules.SelectStarPolicy.ForbidsSelectStar(),
		).To(Equal(expectations.noSelectStar), expectations.description+": SelectStarPolicy")
		Expect(
			r.StyleRules.ColumnExplicitness.RequiresExplicitColumns(),
		).To(Equal(expectations.requireExplicitColumns), expectations.description+": ColumnExplicitness")
		Expect(
			r.SafetyRules.WhereRequirement.RequiresOnDestructive(),
		).To(Equal(expectations.requireWhere), expectations.description+": WhereRequirement")
		Expect(
			r.SafetyRules.LimitRequirement.RequiresOnSelect(),
		).To(Equal(expectations.requireLimit), expectations.description+": LimitRequirement")
		Expect(
			r.SafetyRules.MaxRowsWithoutLimit,
		).To(Equal(expectations.maxRowsWithoutLimit), expectations.description+": MaxRowsWithoutLimit")
		Expect(
			r.DestructiveOps,
		).To(Equal(expectations.destructiveOps), expectations.description+": DestructiveOps")

		err := r.IsValid()
		Expect(err).NotTo(HaveOccurred(), expectations.description+": Should be valid")

	default:
		Fail("Unsupported safety rules type")
	}
}

// DestructiveOperationPolicy validation test suite.
type DestructiveOperationPolicyTestSuite struct{}

func (DestructiveOperationPolicyTestSuite) GetValidValues() []domain.DestructiveOperationPolicy {
	return []domain.DestructiveOperationPolicy{
		domain.DestructiveAllowed,
		domain.DestructiveWithConfirmation,
		domain.DestructiveForbidden,
	}
}

func (DestructiveOperationPolicyTestSuite) GetInvalidValues() []domain.DestructiveOperationPolicy {
	return []domain.DestructiveOperationPolicy{
		"invalid",
		"",
		"sometimes",
	}
}

func (DestructiveOperationPolicyTestSuite) GetTypeName() string {
	return "DestructiveOperationPolicy"
}

// destructivePolicyTestCase defines a test case for DestructiveOperationPolicy behavior.
type destructivePolicyTestCase struct {
	description       string
	allowedPolicy     domain.DestructiveOperationPolicy
	forbiddenPolicies []domain.DestructiveOperationPolicy
	checker           func(domain.DestructiveOperationPolicy) bool
}

// testDestructivePolicyBehavior tests a destructive policy behavior against all policy values.
func testDestructivePolicyBehavior(tc destructivePolicyTestCase) {
	It(tc.description, func() {
		Expect(tc.checker(tc.allowedPolicy)).To(BeTrue())
	})

	It(tc.description+" for other policies", func() {
		for _, policy := range tc.forbiddenPolicies {
			Expect(tc.checker(policy)).To(BeFalse())
		}
	})
}

// DestructivePolicyTestCaseOption configures a destructivePolicyTestCase.
type DestructivePolicyTestCaseOption func(*destructivePolicyTestCase)

// WithAllowedPolicy sets the policy that should be allowed.
func WithAllowedPolicy(policy domain.DestructiveOperationPolicy) DestructivePolicyTestCaseOption {
	return func(tc *destructivePolicyTestCase) {
		tc.allowedPolicy = policy
	}
}

// WithForbiddenPolicies sets the policies that should be forbidden.
func WithForbiddenPolicies(
	policies ...domain.DestructiveOperationPolicy,
) DestructivePolicyTestCaseOption {
	return func(tc *destructivePolicyTestCase) {
		tc.forbiddenPolicies = policies
	}
}

// newDestructivePolicyTestCase creates a test case for DestructiveOperationPolicy behavior.
func newDestructivePolicyTestCase(
	methodName string,
	opts ...DestructivePolicyTestCaseOption,
) destructivePolicyTestCase {
	tc := destructivePolicyTestCase{
		description:   "should return true for method: " + methodName,
		allowedPolicy: domain.DestructiveAllowed,
		forbiddenPolicies: []domain.DestructiveOperationPolicy{
			domain.DestructiveWithConfirmation,
			domain.DestructiveForbidden,
		},
	}

	switch methodName {
	case "AllowsDropTable":
		tc.checker = func(p domain.DestructiveOperationPolicy) bool { return p.AllowsDropTable() }
	case "AllowsTruncate":
		tc.checker = func(p domain.DestructiveOperationPolicy) bool { return p.AllowsTruncate() }
	case "RequiresConfirmation":
		tc.checker = func(p domain.DestructiveOperationPolicy) bool { return p.RequiresConfirmation() }
	}

	for _, opt := range opts {
		opt(&tc)
	}

	return tc
}

// customRuleTestConfig configures a custom rule validation test.
type customRuleTestConfig struct {
	ruleName       string
	ruleExpression string
	expectedError  string
}

// testCustomRuleValidation tests custom rule validation with the given configuration.
func testCustomRuleValidation(cfg customRuleTestConfig) {
	rules := testing.CreateTypeSafeSafetyRules(func(r *domain.TypeSafeSafetyRules) {
		r.CustomRules = []generated.SafetyRule{
			{
				Name:    cfg.ruleName,
				Rule:    cfg.ruleExpression,
				Message: "Test rule",
			},
		}
	})

	err := rules.IsValid()
	Expect(err).To(HaveOccurred())
	Expect(err.Error()).To(ContainSubstring(cfg.expectedError))
}
