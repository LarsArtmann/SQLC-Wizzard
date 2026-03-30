package wizard

import (
	"github.com/LarsArtmann/SQLC-Wizzard/generated"
)

// BranchingPolicy defines the rules for dynamic wizard flow branching.
// This interface allows different branching strategies to be plugged in.
type BranchingPolicy interface {
	// ShouldShowStep determines if a step should be shown/executed.
	ShouldShowStep(stepID StepID, context *FlowContext) bool

	// ShouldShowFeature determines if a feature should be shown.
	ShouldShowFeature(featureKey string, context *FlowContext) bool

	// GetFeatureDefault returns the default value for a feature.
	GetFeatureDefault(featureKey string, context *FlowContext) bool

	// GetStepDescription returns a description for a step.
	GetStepDescription(stepID StepID) string
}

// DefaultBranchingPolicy implements a sensible default branching strategy.
type DefaultBranchingPolicy struct{}

// NewDefaultBranchingPolicy creates a new default branching policy.
func NewDefaultBranchingPolicy() *DefaultBranchingPolicy {
	return &DefaultBranchingPolicy{}
}

// ShouldShowStep determines if a step should be shown based on the context.
func (p *DefaultBranchingPolicy) ShouldShowStep(stepID StepID, ctx *FlowContext) bool {
	switch stepID {
	case StepProjectType:
		// Always show project type selection
		return true

	case StepDatabase:
		// Always show database selection after project type
		return true

	case StepProjectDetail:
		// Always show project details
		return true

	case StepFeatures:
		// Skip features for simple project types
		return ctx.ProjectType != generated.ProjectTypeHobby &&
			ctx.ProjectType != generated.ProjectTypeTesting

	case StepOutput:
		// Always show output configuration
		return true

	case StepAdvanced:
		// Only show advanced step for complex projects
		return ctx.ProjectType == generated.ProjectTypeEnterprise ||
			ctx.ProjectType == generated.ProjectTypeMultiTenant ||
			ctx.IsAdvancedMode

	case StepReview:
		// Only show review for enterprise projects
		return ctx.ProjectType == generated.ProjectTypeEnterprise ||
			ctx.IsAdvancedMode

	default:
		return true
	}
}

// ShouldShowFeature determines if a feature should be shown based on context.
func (p *DefaultBranchingPolicy) ShouldShowFeature(featureKey string, ctx *FlowContext) bool {
	switch featureKey {
	// UUIDs are only available for databases that support them
	case "uuid":
		return ctx.DatabaseType == generated.DatabaseTypePostgreSQL ||
			ctx.DatabaseType == generated.DatabaseTypeMySQL

	// JSON is available for most databases
	case "json":
		return true

	// Arrays are PostgreSQL-specific
	case "array":
		return ctx.DatabaseType == generated.DatabaseTypePostgreSQL

	// Full-text search is PostgreSQL-specific (MySQL also has it but limited)
	case "fulltext":
		return ctx.DatabaseType == generated.DatabaseTypePostgreSQL ||
			ctx.DatabaseType == generated.DatabaseTypeMySQL

	// Strict mode for enterprise projects
	case "strict_mode":
		return ctx.ProjectType == generated.ProjectTypeEnterprise ||
			ctx.ProjectType == generated.ProjectTypeMultiTenant

	// Prepared queries available for most project types
	case "prepared_queries":
		return ctx.ProjectType != generated.ProjectTypeHobby &&
			ctx.ProjectType != generated.ProjectTypeTesting

	// JSON tags for API-first and enterprise
	case "json_tags":
		return ctx.ProjectType == generated.ProjectTypeAPIFirst ||
			ctx.ProjectType == generated.ProjectTypeEnterprise ||
			ctx.ProjectType == generated.ProjectTypeMultiTenant

	// Interfaces for library and enterprise
	case "interface":
		return ctx.ProjectType == generated.ProjectTypeLibrary ||
			ctx.ProjectType == generated.ProjectTypeEnterprise ||
			ctx.ProjectType == generated.ProjectTypeMultiTenant

	// Strict ORDER BY for analytics
	case "strict_orderby":
		return ctx.ProjectType == generated.ProjectTypeAnalytics ||
			ctx.ProjectType == generated.ProjectTypeEnterprise

	default:
		return true
	}
}

// GetFeatureDefault returns the default value for a feature based on context.
func (p *DefaultBranchingPolicy) GetFeatureDefault(featureKey string, ctx *FlowContext) bool {
	switch featureKey {
	// UUIDs default based on database type
	case "uuid":
		return ctx.ShouldEnableUUIDs()

	// JSON typically enabled for most projects
	case "json":
		return ctx.ShouldEnableJSON()

	// Arrays disabled by default (PostgreSQL only)
	case "array":
		return false

	// Full-text search disabled by default
	case "fulltext":
		return false

	// Strict mode disabled by default
	case "strict_mode":
		return false

	// Prepared queries based on project type
	case "prepared_queries":
		return ctx.ProjectType == generated.ProjectTypeEnterprise ||
			ctx.ProjectType == generated.ProjectTypeAPIFirst

	// JSON tags enabled for API-first and enterprise
	case "json_tags":
		return ctx.ProjectType == generated.ProjectTypeAPIFirst ||
			ctx.ProjectType == generated.ProjectTypeEnterprise

	// Interfaces enabled for library and enterprise
	case "interface":
		return ctx.ProjectType == generated.ProjectTypeLibrary ||
			ctx.ProjectType == generated.ProjectTypeEnterprise

	// Strict ORDER BY for analytics
	case "strict_orderby":
		return ctx.ProjectType == generated.ProjectTypeAnalytics

	default:
		return false
	}
}

// GetStepDescription returns a human-readable description for a step.
func (p *DefaultBranchingPolicy) GetStepDescription(stepID StepID) string {
	switch stepID {
	case StepProjectType:
		return "Select the type of project you're building"
	case StepDatabase:
		return "Choose your database engine"
	case StepProjectDetail:
		return "Configure project name and package"
	case StepFeatures:
		return "Configure code generation and validation options"
	case StepOutput:
		return "Set up output directories"
	case StepAdvanced:
		return "Configure advanced options"
	case StepReview:
		return "Review and confirm your configuration"
	default:
		return ""
	}
}

// SimpleBranchingPolicy is a simplified policy for testing.
type SimpleBranchingPolicy struct{}

// NewSimpleBranchingPolicy creates a new simple branching policy for testing.
func NewSimpleBranchingPolicy() *SimpleBranchingPolicy {
	return &SimpleBranchingPolicy{}
}

// ShouldShowStep always returns true for simple policy.
func (p *SimpleBranchingPolicy) ShouldShowStep(stepID StepID, ctx *FlowContext) bool {
	return true
}

// ShouldShowFeature always returns true for simple policy.
func (p *SimpleBranchingPolicy) ShouldShowFeature(featureKey string, ctx *FlowContext) bool {
	return true
}

// GetFeatureDefault returns false for simple policy.
func (p *SimpleBranchingPolicy) GetFeatureDefault(featureKey string, ctx *FlowContext) bool {
	return false
}

// GetStepDescription returns empty string for simple policy.
func (p *SimpleBranchingPolicy) GetStepDescription(stepID StepID) string {
	return ""
}
