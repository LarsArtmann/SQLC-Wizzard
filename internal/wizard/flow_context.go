package wizard

import (
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
	for _, s := range fc.CompletedSteps {
		if s == step {
			return true
		}
	}

	return false
}

// IsStepSkipped checks if a step was skipped.
func (fc *FlowContext) IsStepSkipped(step StepID) bool {
	for _, s := range fc.SkippedSteps {
		if s == step {
			return true
		}
	}

	return false
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
	switch fc.DatabaseType {
	case generated.DatabaseTypePostgreSQL:
		return true
	case generated.DatabaseTypeMySQL:
		return true
	case generated.DatabaseTypeSQLite:
		return false // SQLite doesn't have native UUID support
	default:
		return false
	}
}

// ShouldEnableJSON returns whether JSON support should be enabled by default.
func (fc *FlowContext) ShouldEnableJSON() bool {
	switch fc.DatabaseType {
	case generated.DatabaseTypePostgreSQL:
		return true
	case generated.DatabaseTypeMySQL:
		return true
	case generated.DatabaseTypeSQLite:
		return false // SQLite has limited JSON support
	default:
		return false
	}
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

// GetDatabaseSpecificFeatures returns the database-specific features for the current database.
func (fc *FlowContext) GetDatabaseSpecificFeatures() []FeatureSpec {
	var features []FeatureSpec

	switch fc.DatabaseType {
	case generated.DatabaseTypePostgreSQL:
		features = []FeatureSpec{
			{
				Key:         "uuid",
				Title:       "Use UUID primary keys?",
				Description: "Generate UUID primary keys instead of auto-increment integers",
			},
			{
				Key:         "json",
				Title:       "Use JSON columns?",
				Description: "Enable JSON column support for flexible data storage",
			},
			{Key: "array", Title: "Use array columns?", Description: "Enable array column support"},
			{
				Key:         "fulltext",
				Title:       "Use full-text search?",
				Description: "Enable full-text search capabilities",
			},
		}
	case generated.DatabaseTypeMySQL:
		features = []FeatureSpec{
			{
				Key:         "uuid",
				Title:       "Use UUID primary keys?",
				Description: "Generate UUID primary keys instead of auto-increment integers",
			},
			{
				Key:         "json",
				Title:       "Use JSON columns?",
				Description: "Enable JSON column support for flexible data storage",
			},
			{
				Key:         "fulltext",
				Title:       "Use full-text search?",
				Description: "Enable full-text search capabilities (MySQL 5.7+)",
			},
		}
	case generated.DatabaseTypeSQLite:
		features = []FeatureSpec{
			{
				Key:         "json",
				Title:       "Use JSON columns?",
				Description: "Enable JSON column support for flexible data storage",
			},
		}
	}

	return features
}

// GetProjectTypeFeatures returns project-type-specific features.
func (fc *FlowContext) GetProjectTypeFeatures() []FeatureSpec {
	var features []FeatureSpec

	switch fc.ProjectType {
	case "enterprise", "multi-tenant":
		features = []FeatureSpec{
			{
				Key:         "strict_mode",
				Title:       "Enable strict mode?",
				Description: "Enable strict validation for all queries",
			},
			{
				Key:         "prepared_queries",
				Title:       "Use prepared queries?",
				Description: "Generate prepared query methods for better performance",
			},
		}
	case "api-first":
		features = []FeatureSpec{
			{
				Key:         "json_tags",
				Title:       "Include JSON tags?",
				Description: "Add JSON struct tags to generated models",
			},
			{
				Key:         "interface",
				Title:       "Generate interfaces?",
				Description: "Create interfaces for query methods",
			},
		}
	case "analytics":
		features = []FeatureSpec{
			{
				Key:         "fulltext",
				Title:       "Enable full-text search?",
				Description: "Enable full-text search capabilities",
			},
			{
				Key:         "strict_orderby",
				Title:       "Strict ORDER BY?",
				Description: "Require ORDER BY in SELECT queries",
			},
		}
	case "library":
		features = []FeatureSpec{
			{
				Key:         "interface",
				Title:       "Generate interfaces?",
				Description: "Create interfaces for maximum flexibility",
			},
			{
				Key:         "json_tags",
				Title:       "Include JSON tags?",
				Description: "Add JSON struct tags for serialization",
			},
		}
	}

	return features
}

// FeatureSpec defines a feature specification for branching decisions.
type FeatureSpec struct {
	Key         string
	Title       string
	Description string
	Enabled     bool
}
