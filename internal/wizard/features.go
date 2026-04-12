package wizard

import (
	"fmt"

	"charm.land/huh/v2"
	"github.com/LarsArtmann/SQLC-Wizzard/generated"
)

// runConfirmationForm creates and runs a confirmation form, returning the result in the provided value pointer.
func runConfirmationForm(themeFunc huh.ThemeFunc, title, description string, result *bool) error {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title(title).
				Description(description).
				Value(result),
		),
	).WithTheme(themeFunc)

	return form.Run()
}

// FeaturesStep handles feature selection and validation configuration.
type FeaturesStep struct {
	themeFunc huh.ThemeFunc
	ui        *UIHelper
}

// NewFeaturesStep creates a new features step.
func NewFeaturesStep(themeFunc huh.ThemeFunc, ui *UIHelper) *FeaturesStep {
	return &FeaturesStep{
		themeFunc: themeFunc,
		ui:        ui,
	}
}

// Execute runs the feature selection step with branching support.
func (s *FeaturesStep) Execute(data *generated.TemplateData) error {
	s.ui.ShowStepHeader("Features & Validation")

	// Code generation options - conditional based on project type
	err := s.configureCodeGeneration(data)
	if err != nil {
		return err
	}

	// Safety rules - conditional based on project type
	err = s.configureSafetyRules(data)
	if err != nil {
		return err
	}

	// Database features - fully dynamic based on database engine
	err = s.configureDatabaseFeatures(data)
	if err != nil {
		return err
	}

	// Project-type specific features - conditional based on project type
	err = s.configureProjectTypeFeatures(data)
	if err != nil {
		return err
	}

	s.ui.ShowStepComplete(
		"Features",
		"Code generation, safety rules, and database features configured",
	)

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

// buildFeatures creates a FeatureConfig slice from variadic arguments, eliminating
// the need for rawFeatureConfig structs and buildFeatureConfigs intermediate step.
func buildFeatures(configs ...FeatureConfig) []FeatureConfig {
	return configs
}

// configSpec defines a feature configuration specification.
type configSpec struct {
	title       string
	description string
	fieldPath   string // Dot-notation path to the field: "EmitOptions.EmitInterface", "SafetyRules.NoSelectStar"
}

// configFieldMapper maps a configSpec to a field assignment function.
var configFieldMapper = map[string]fieldAssignment{
	"EmitOptions.EmitInterface":       func(data *generated.TemplateData, val bool) { data.Validation.EmitOptions.EmitInterface = val },
	"EmitOptions.EmitPreparedQueries": func(data *generated.TemplateData, val bool) { data.Validation.EmitOptions.EmitPreparedQueries = val },
	"EmitOptions.EmitJSONTags":        func(data *generated.TemplateData, val bool) { data.Validation.EmitOptions.EmitJSONTags = val },
	"SafetyRules.NoSelectStar":        func(data *generated.TemplateData, val bool) { data.Validation.SafetyRules.NoSelectStar = val },
	"SafetyRules.RequireWhere":        func(data *generated.TemplateData, val bool) { data.Validation.SafetyRules.RequireWhere = val },
	"SafetyRules.RequireLimit":        func(data *generated.TemplateData, val bool) { data.Validation.SafetyRules.RequireLimit = val },
}

// buildConfigs creates FeatureConfig slice from config specifications.
func buildConfigs(specs []configSpec) []FeatureConfig {
	configs := make([]FeatureConfig, len(specs))
	for i, spec := range specs {
		assign := configFieldMapper[spec.fieldPath]
		configs[i] = createFeatureConfig(spec.title, spec.description, assign)
	}

	return configs
}

// Pre-defined configuration sets

// Code generation configs.
var codeGenerationConfigs = buildConfigs([]configSpec{
	{"Generate Go interfaces?", "Create interfaces for query methods", "EmitOptions.EmitInterface"},
	{
		"Generate prepared queries?",
		"Create prepared query methods for better performance",
		"EmitOptions.EmitPreparedQueries",
	},
	{"Add JSON tags?", "Add JSON struct tags to generated models", "EmitOptions.EmitJSONTags"},
})

// Safety rule configs.
var safetyRuleConfigs = buildConfigs([]configSpec{
	{
		"Forbid SELECT *?",
		"Prevent SELECT * queries for better performance and explicitness",
		"SafetyRules.NoSelectStar",
	},
	{
		"Require WHERE clause?",
		"Force WHERE clauses in UPDATE/DELETE queries to prevent accidental data modification",
		"SafetyRules.RequireWhere",
	},
	{
		"Require LIMIT on SELECT?",
		"Force LIMIT clauses on SELECT queries to prevent large result sets",
		"SafetyRules.RequireLimit",
	},
})

// runFeatureConfigForm runs confirmation form for any feature configuration.
func (s *FeaturesStep) runFeatureConfigForm(
	data *generated.TemplateData,
	configs []FeatureConfig,
	errorContext string,
) error {
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
	).WithTheme(s.themeFunc)

	err := form.Run()
	if err != nil {
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

// configureProjectTypeFeatures configures project-type specific features.
// This method is called after database features and adds additional configuration
// based on the project type (e.g., strict mode for enterprise, prepared queries for API-first).
func (s *FeaturesStep) configureProjectTypeFeatures(data *generated.TemplateData) error {
	// Only show additional project-type specific features for complex project types
	switch data.ProjectType {
	case generated.ProjectTypeEnterprise, generated.ProjectTypeMultiTenant:
		// Enterprise projects benefit from strict mode and prepared queries
		return s.configureEnterpriseFeatures(data)

	case generated.ProjectTypeAPIFirst:
		// API-first projects benefit from interface generation
		return s.configureAPIFirstFeatures(data)

	case generated.ProjectTypeAnalytics:
		// Analytics projects benefit from strict ORDER BY
		return s.configureAnalyticsFeatures(data)

	default:
		// For hobby, microservice, testing, and library - skip additional features
		return nil
	}
}

// configureEnterpriseFeatures adds enterprise-specific feature configuration.
func (s *FeaturesStep) configureEnterpriseFeatures(data *generated.TemplateData) error {
	var enableStrictMode bool

	err := runConfirmationForm(
		s.themeFunc,
		"Enable strict mode?",
		"Enable strict validation for all queries to catch potential issues early",
		&enableStrictMode,
	)
	if err != nil {
		return fmt.Errorf("enterprise features configuration failed: %w", err)
	}

	data.Validation.StrictFunctions = enableStrictMode
	data.Validation.StrictOrderBy = enableStrictMode

	return nil
}

// configureAPIFirstFeatures adds API-first specific feature configuration.
func (s *FeaturesStep) configureAPIFirstFeatures(data *generated.TemplateData) error {
	// Ensure JSON tags are enabled for API-first projects
	data.Validation.EmitOptions.EmitJSONTags = true

	return nil
}

// configureAnalyticsFeatures adds analytics-specific feature configuration.
func (s *FeaturesStep) configureAnalyticsFeatures(data *generated.TemplateData) error {
	var enableStrictOrderBy bool

	err := runConfirmationForm(
		s.themeFunc,
		"Enable strict ORDER BY?",
		"Require ORDER BY in all SELECT queries to ensure predictable results",
		&enableStrictOrderBy,
	)
	if err != nil {
		return fmt.Errorf("analytics features configuration failed: %w", err)
	}

	data.Validation.StrictOrderBy = enableStrictOrderBy

	return nil
}

// configureDatabaseFeatures configures database-specific features using branching context.
func (s *FeaturesStep) configureDatabaseFeatures(data *generated.TemplateData) error {
	// Build dynamic database features based on engine type
	var fields []huh.Field

	switch data.Database.Engine {
	case generated.DatabaseTypePostgreSQL:
		// PostgreSQL supports all features
		fields = []huh.Field{
			huh.NewConfirm().
				Title("Use UUID primary keys?").
				Description("Generate UUID primary keys instead of auto-increment integers").
				Value(&data.Database.UseUUIDs),
			huh.NewConfirm().
				Title("Use JSON columns?").
				Description("Enable JSON column support for flexible data storage").
				Value(&data.Database.UseJSON),
			huh.NewConfirm().
				Title("Use array columns?").
				Description("Enable array column support").
				Value(&data.Database.UseArrays),
			huh.NewConfirm().
				Title("Use full-text search?").
				Description("Enable full-text search capabilities").
				Value(&data.Database.UseFullText),
		}

	case generated.DatabaseTypeMySQL:
		// MySQL supports UUIDs, JSON, and full-text (but not arrays)
		fields = []huh.Field{
			huh.NewConfirm().
				Title("Use UUID primary keys?").
				Description("Generate UUID primary keys instead of auto-increment integers").
				Value(&data.Database.UseUUIDs),
			huh.NewConfirm().
				Title("Use JSON columns?").
				Description("Enable JSON column support for flexible data storage").
				Value(&data.Database.UseJSON),
			huh.NewConfirm().
				Title("Use full-text search?").
				Description("Enable full-text search capabilities (MySQL 5.7+)").
				Value(&data.Database.UseFullText),
		}

	case generated.DatabaseTypeSQLite:
		// SQLite has limited features
		fields = []huh.Field{
			huh.NewConfirm().
				Title("Use JSON columns?").
				Description("Enable JSON column support for flexible data storage").
				Value(&data.Database.UseJSON),
		}

	default:
		// For unknown database types, show no specific features
		s.ui.ShowInfo("No database-specific features available for this engine type")

		return nil
	}

	form := huh.NewForm(
		huh.NewGroup(fields...),
	).WithTheme(s.themeFunc)

	err := form.Run()
	if err != nil {
		return fmt.Errorf(
			"database features form input failed (UUIDs=%v, JSON=%v, Arrays=%v, FullText=%v): %w",
			data.Database.UseUUIDs,
			data.Database.UseJSON,
			data.Database.UseArrays,
			data.Database.UseFullText,
			err,
		)
	}

	return nil
}
