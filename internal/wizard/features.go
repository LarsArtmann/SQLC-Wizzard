package wizard

import (
	"fmt"

	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/charmbracelet/huh"
)

// FeaturesStep handles feature selection and validation configuration
type FeaturesStep struct {
	theme *huh.Theme
	ui    *UIHelper
}

// NewFeaturesStep creates a new features step
func NewFeaturesStep(theme *huh.Theme, ui *UIHelper) *FeaturesStep {
	return &FeaturesStep{
		theme: theme,
		ui:    ui,
	}
}

// Execute runs the feature selection step
func (s *FeaturesStep) Execute(data *generated.TemplateData) error {
	s.ui.ShowStepHeader("Features & Validation")

	// Code generation options
	if err := s.configureCodeGeneration(data); err != nil {
		return err
	}

	// Safety rules
	if err := s.configureSafetyRules(data); err != nil {
		return err
	}

	// Database features
	if err := s.configureDatabaseFeatures(data); err != nil {
		return err
	}

	s.ui.ShowStepComplete("Features", "Code generation, safety rules, and database features configured")
	return nil
}

// fieldAssignment defines how to assign a boolean value to a data structure
type fieldAssignment func(data *generated.TemplateData, value bool)

// confirmationConfig represents a complete confirmation field configuration
type confirmationConfig struct {
	title       string
	description string
	assign      fieldAssignment
}

// Pre-defined configuration sets for different feature types
var (
	// codeGenerationConfigs defines options for code generation features
	codeGenerationConfigs = []confirmationConfig{
		{
			title:       "Generate Go interfaces?",
			description: "Create interfaces for query methods",
			assign: func(data *generated.TemplateData, value bool) {
				data.Validation.EmitOptions.EmitInterface = value
			},
		},
		{
			title:       "Generate prepared queries?",
			description: "Create prepared query methods for better performance",
			assign: func(data *generated.TemplateData, value bool) {
				data.Validation.EmitOptions.EmitPreparedQueries = value
			},
		},
		{
			title:       "Add JSON tags?",
			description: "Add JSON struct tags to generated models",
			assign: func(data *generated.TemplateData, value bool) {
				data.Validation.EmitOptions.EmitJSONTags = value
			},
		},
	}

	// safetyRuleConfigs defines options for safety rule configuration
	safetyRuleConfigs = []confirmationConfig{
		{
			title:       "Forbid SELECT *?",
			description: "Prevent SELECT * queries for better performance and explicitness",
			assign: func(data *generated.TemplateData, value bool) {
				data.Validation.SafetyRules.NoSelectStar = value
			},
		},
		{
			title:       "Require WHERE clause?",
			description: "Force WHERE clauses in UPDATE/DELETE queries to prevent accidental data modification",
			assign: func(data *generated.TemplateData, value bool) {
				data.Validation.SafetyRules.RequireWhere = value
			},
		},
		{
			title:       "Require LIMIT on SELECT?",
			description: "Force LIMIT clauses on SELECT queries to prevent large result sets",
			assign: func(data *generated.TemplateData, value bool) {
				data.Validation.SafetyRules.RequireLimit = value
			},
		},
	}
)

// runConfirmationWithAssignment runs a confirmation form and applies assignments
func (s *FeaturesStep) runConfirmationWithAssignment(data *generated.TemplateData, configs []confirmationConfig, errorContext string) error {
	// Create boolean values for each field
	values := make([]bool, len(configs))
	valuePtrs := make([]*bool, len(configs))
	for i := range values {
		valuePtrs[i] = &values[i]
	}

	// Build form fields
	var formFields []huh.Field
	for i, config := range configs {
		formFields = append(formFields,
			huh.NewConfirm().
				Title(config.title).
				Description(config.description).
				Value(valuePtrs[i]),
		)
	}

	form := huh.NewForm(
		huh.NewGroup(formFields...),
	).WithTheme(s.theme)

	if err := form.Run(); err != nil {
		return fmt.Errorf("%s configuration failed: %w", errorContext, err)
	}

	// Apply assignments
	for i, config := range configs {
		config.assign(data, values[i])
	}

	return nil
}

// configureCodeGeneration configures code generation options
func (s *FeaturesStep) configureCodeGeneration(data *generated.TemplateData) error {
	return s.runConfirmationWithAssignment(data, codeGenerationConfigs, "code generation")
}

// configureSafetyRules configures safety rules
func (s *FeaturesStep) configureSafetyRules(data *generated.TemplateData) error {
	return s.runConfirmationWithAssignment(data, safetyRuleConfigs, "safety rules")
}

// configureDatabaseFeatures configures database-specific features
func (s *FeaturesStep) configureDatabaseFeatures(data *generated.TemplateData) error {
	var useUUIDs, useJSON, useArrays, useFullText bool

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Use UUID primary keys?").
				Description("Generate UUID primary keys instead of auto-increment integers").
				Value(&useUUIDs),
			huh.NewConfirm().
				Title("Use JSON columns?").
				Description("Enable JSON column support for flexible data storage").
				Value(&useJSON),
			huh.NewConfirm().
				Title("Use array columns?").
				Description("Enable array column support (PostgreSQL only)").
				Value(&useArrays),
			huh.NewConfirm().
				Title("Use full-text search?").
				Description("Enable full-text search capabilities").
				Value(&useFullText),
		),
	).WithTheme(s.theme)

	if err := form.Run(); err != nil {
		return fmt.Errorf("database features configuration failed: %w", err)
	}

	data.Database.UseUUIDs = useUUIDs
	data.Database.UseJSON = useJSON
	data.Database.UseArrays = useArrays
	data.Database.UseFullText = useFullText

	return nil
}
