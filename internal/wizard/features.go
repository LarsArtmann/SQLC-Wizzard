package wizard

import (
	"fmt"

	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/charmbracelet/huh"
)

// FeaturesStep handles feature selection and validation configuration.
type FeaturesStep struct {
	theme *huh.Theme
	ui    *UIHelper
}

// NewFeaturesStep creates a new features step.
func NewFeaturesStep(theme *huh.Theme, ui *UIHelper) *FeaturesStep {
	return &FeaturesStep{
		theme: theme,
		ui:    ui,
	}
}

// Execute runs the feature selection step.
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

// fieldAssignment defines how to assign a boolean value to a data structure.
type fieldAssignment func(data *generated.TemplateData, value bool)

// Feature configuration interface to handle different config types generically.
type FeatureConfig interface {
	GetTitle() string
	GetDescription() string
	Assign(data *generated.TemplateData, value bool)
}

// featureConfig implements FeatureConfig interface.
type featureConfig struct {
	title       string
	description string
	assign      fieldAssignment
}

func (c featureConfig) GetTitle() string                                { return c.title }
func (c featureConfig) GetDescription() string                          { return c.description }
func (c featureConfig) Assign(data *generated.TemplateData, value bool) { c.assign(data, value) }

// createFeatureConfig creates a new feature configuration.
func createFeatureConfig(title, description string, assign fieldAssignment) FeatureConfig {
	return &featureConfig{
		title:       title,
		description: description,
		assign:      assign,
	}
}

// rawFeatureConfig holds the static data for building a FeatureConfig.
type rawFeatureConfig struct {
	title       string
	description string
	assign      fieldAssignment
}

// buildFeatureConfigs transforms raw configs into FeatureConfig slice.
func buildFeatureConfigs(rawConfigs []rawFeatureConfig) []FeatureConfig {
	configs := make([]FeatureConfig, len(rawConfigs))
	for i, raw := range rawConfigs {
		configs[i] = createFeatureConfig(raw.title, raw.description, raw.assign)
	}
	return configs
}

// Pre-defined configuration sets

// Code generation configs.
var codeGenerationConfigs = buildFeatureConfigs([]rawFeatureConfig{
	{
		title:       "Generate Go interfaces?",
		description: "Create interfaces for query methods",
		assign:      func(data *generated.TemplateData, val bool) { data.Validation.EmitOptions.EmitInterface = val },
	},
	{
		title:       "Generate prepared queries?",
		description: "Create prepared query methods for better performance",
		assign:      func(data *generated.TemplateData, val bool) { data.Validation.EmitOptions.EmitPreparedQueries = val },
	},
	{
		title:       "Add JSON tags?",
		description: "Add JSON struct tags to generated models",
		assign:      func(data *generated.TemplateData, val bool) { data.Validation.EmitOptions.EmitJSONTags = val },
	},
})

// Safety rule configs.
var safetyRuleConfigs = buildFeatureConfigs([]rawFeatureConfig{
	{
		title:       "Forbid SELECT *?",
		description: "Prevent SELECT * queries for better performance and explicitness",
		assign:      func(data *generated.TemplateData, val bool) { data.Validation.SafetyRules.NoSelectStar = val },
	},
	{
		title:       "Require WHERE clause?",
		description: "Force WHERE clauses in UPDATE/DELETE queries to prevent accidental data modification",
		assign:      func(data *generated.TemplateData, val bool) { data.Validation.SafetyRules.RequireWhere = val },
	},
	{
		title:       "Require LIMIT on SELECT?",
		description: "Force LIMIT clauses on SELECT queries to prevent large result sets",
		assign:      func(data *generated.TemplateData, val bool) { data.Validation.SafetyRules.RequireLimit = val },
	},
})

// runFeatureConfigForm runs confirmation form for any feature configuration.
func (s *FeaturesStep) runFeatureConfigForm(data *generated.TemplateData, configs []FeatureConfig, errorContext string) error {
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
				Title(config.GetTitle()).
				Description(config.GetDescription()).
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
		config.Assign(data, values[i])
	}

	return nil
}

// configureCodeGeneration configures code generation options.
func (s *FeaturesStep) configureCodeGeneration(data *generated.TemplateData) error {
	return s.runFeatureConfigForm(data, codeGenerationConfigs, "code generation")
}

// configureSafetyRules configures safety rules.
func (s *FeaturesStep) configureSafetyRules(data *generated.TemplateData) error {
	return s.runFeatureConfigForm(data, safetyRuleConfigs, "safety rules")
}

// configureDatabaseFeatures configures database-specific features.
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
