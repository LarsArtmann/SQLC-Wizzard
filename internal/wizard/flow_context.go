package wizard

import (
	"slices"

	"github.com/LarsArtmann/SQLC-Wizzard/generated"
)

// FlowContext stores the current state and branching decisions for the wizard flow.
// This enables dynamic step execution and feature filtering based on user selections.
type FlowContext struct {
	// Current selections
	ProjectType  generated.ProjectType
	DatabaseType generated.DatabaseType
	ProjectName  string
	PackageName  string

	// Step execution tracking
	CompletedSteps []StepID
	SkippedSteps   []StepID

	// Feature flags
	EnableUUIDs      bool
	EnableJSON       bool
	EnableArrays     bool
	EnableFullText   bool
	EnableStrictMode bool

	// Configuration flags
	IsAdvancedMode    bool
	SkipOptionalSteps bool

	// Step metadata
	CurrentStep StepID
}

// StepID represents a unique identifier for wizard steps.
type StepID string

const (
	StepProjectType   StepID = "project_type"
	StepDatabase      StepID = "database"
	StepProjectDetail StepID = "project_details"
	StepFeatures      StepID = "features"
	StepOutput        StepID = "output"
	StepAdvanced      StepID = "advanced"
	StepReview        StepID = "review"
)

// NewFlowContext creates a new flow context with default values.
func NewFlowContext() *FlowContext {
	return &FlowContext{
		CompletedSteps:    make([]StepID, 0),
		SkippedSteps:      make([]StepID, 0),
		EnableUUIDs:       true,
		EnableJSON:        true,
		EnableArrays:      false,
		EnableFullText:    false,
		EnableStrictMode:  false,
		SkipOptionalSteps: false,
	}
}

// UpdateFromTemplateData updates the flow context from template data.
func (fc *FlowContext) UpdateFromTemplateData(data *generated.TemplateData) {
	if data == nil {
		return
	}

	fc.ProjectType = data.ProjectType
	fc.DatabaseType = data.Database.Engine
	fc.ProjectName = data.ProjectName
	fc.PackageName = data.Package.Name
	fc.EnableUUIDs = data.Database.UseUUIDs
	fc.EnableJSON = data.Database.UseJSON
	fc.EnableArrays = data.Database.UseArrays
	fc.EnableFullText = data.Database.UseFullText
}

// MarkStepCompleted marks a step as completed.
func (fc *FlowContext) MarkStepCompleted(step StepID) {
	fc.CompletedSteps = append(fc.CompletedSteps, step)
}

// MarkStepSkipped marks a step as skipped.
func (fc *FlowContext) MarkStepSkipped(step StepID) {
	fc.SkippedSteps = append(fc.SkippedSteps, step)
}

// IsStepCompleted checks if a step was completed.
func (fc *FlowContext) IsStepCompleted(step StepID) bool {
	return slices.Contains(fc.CompletedSteps, step)
}

// IsStepSkipped checks if a step was skipped.
func (fc *FlowContext) IsStepSkipped(step StepID) bool {
	return slices.Contains(fc.SkippedSteps, step)
}

// GetRequiredSteps returns the list of steps that should be executed based on context.
func (fc *FlowContext) GetRequiredSteps() []StepID {
	steps := []StepID{StepProjectType}

	// Add database step after project type
	steps = append(steps, StepDatabase)

	// Project details always required
	steps = append(steps, StepProjectDetail)

	// Features step - may be simplified for certain project types
	if fc.ProjectType == "hobby" || fc.ProjectType == "testing" {
		// Simplified feature selection for hobby/testing projects
		if !fc.SkipOptionalSteps {
			steps = append(steps, StepFeatures)
		}
	} else {
		// Full feature selection for other project types
		steps = append(steps, StepFeatures)
	}

	// Output step always required
	steps = append(steps, StepOutput)

	// Add advanced step for enterprise/complex projects
	if fc.ProjectType == "enterprise" || fc.ProjectType == "multi-tenant" || fc.IsAdvancedMode {
		steps = append(steps, StepAdvanced)
	}

	// Review step for complex projects
	if fc.ProjectType == "enterprise" || fc.IsAdvancedMode {
		steps = append(steps, StepReview)
	}

	return steps
}

// ShouldShowDatabaseFeatures returns whether database-specific features should be shown.
func (fc *FlowContext) ShouldShowDatabaseFeatures() bool {
	// Database-specific features are available for all non-trivial project types
	return fc.ProjectType != "hobby" && fc.ProjectType != "testing"
}

// ShouldShowAdvancedFeatures returns whether advanced features should be shown.
func (fc *FlowContext) ShouldShowAdvancedFeatures() bool {
	return fc.ProjectType == "enterprise" ||
		fc.ProjectType == "multi-tenant" ||
		fc.ProjectType == "api-first" ||
		fc.IsAdvancedMode
}

// ShouldEnableUUIDs returns whether UUIDs should be enabled by default.
func (fc *FlowContext) ShouldEnableUUIDs() bool {
	return fc.DatabaseType == generated.DatabaseTypePostgreSQL || fc.DatabaseType == generated.DatabaseTypeMySQL
}

// ShouldEnableJSON returns whether JSON support should be enabled by default.
func (fc *FlowContext) ShouldEnableJSON() bool {
	return fc.ShouldEnableUUIDs() // Same databases support JSON
}

// ShouldEnableArrays returns whether array support should be enabled by default.
func (fc *FlowContext) ShouldEnableArrays() bool {
	// Arrays are PostgreSQL-specific
	return fc.DatabaseType == generated.DatabaseTypePostgreSQL
}

// ShouldEnableFullText returns whether full-text search should be enabled by default.
func (fc *FlowContext) ShouldEnableFullText() bool {
	// Full-text search is PostgreSQL-specific
	return fc.DatabaseType == generated.DatabaseTypePostgreSQL
}

// featureSpec creates a FeatureSpec with the given key, title, and description.
func featureSpec(key, title, description string) FeatureSpec {
	return FeatureSpec{Key: key, Title: title, Description: description}
}

// GetDatabaseSpecificFeatures returns the database-specific features for the current database.
func (fc *FlowContext) GetDatabaseSpecificFeatures() []FeatureSpec {
	switch fc.DatabaseType {
	case generated.DatabaseTypePostgreSQL:
		return []FeatureSpec{
			featureSpec("uuid", "Use UUID primary keys?", "Generate UUID primary keys instead of auto-increment integers"),
			featureSpec("json", "Use JSON columns?", "Enable JSON column support for flexible data storage"),
			featureSpec("array", "Use array columns?", "Enable array column support"),
			featureSpec("fulltext", "Use full-text search?", "Enable full-text search capabilities"),
		}
	case generated.DatabaseTypeMySQL:
		return []FeatureSpec{
			featureSpec("uuid", "Use UUID primary keys?", "Generate UUID primary keys instead of auto-increment integers"),
			featureSpec("json", "Use JSON columns?", "Enable JSON column support for flexible data storage"),
			featureSpec("fulltext", "Use full-text search?", "Enable full-text search capabilities (MySQL 5.7+)"),
		}
	case generated.DatabaseTypeSQLite:
		return []FeatureSpec{
			featureSpec("json", "Use JSON columns?", "Enable JSON column support for flexible data storage"),
		}
	}

	return nil
}

// GetProjectTypeFeatures returns project-type-specific features.
func (fc *FlowContext) GetProjectTypeFeatures() []FeatureSpec {
	switch fc.ProjectType {
	case "enterprise", "multi-tenant":
		return []FeatureSpec{
			featureSpec("strict_mode", "Enable strict mode?", "Enable strict validation for all queries"),
			featureSpec("prepared_queries", "Use prepared queries?", "Generate prepared query methods for better performance"),
		}
	case "api-first":
		return []FeatureSpec{
			featureSpec("json_tags", "Include JSON tags?", "Add JSON struct tags to generated models"),
			featureSpec("interface", "Generate interfaces?", "Create interfaces for query methods"),
		}
	case "analytics":
		return []FeatureSpec{
			featureSpec("fulltext", "Enable full-text search?", "Enable full-text search capabilities"),
			featureSpec("strict_orderby", "Strict ORDER BY?", "Require ORDER BY in SELECT queries"),
		}
	case "library":
		return []FeatureSpec{
			featureSpec("interface", "Generate interfaces?", "Create interfaces for maximum flexibility"),
			featureSpec("json_tags", "Include JSON tags?", "Add JSON struct tags for serialization"),
		}
	}

	return nil
}

// FeatureSpec defines a feature specification for branching decisions.
type FeatureSpec struct {
	Key         string
	Title       string
	Description string
	Enabled     bool
}
