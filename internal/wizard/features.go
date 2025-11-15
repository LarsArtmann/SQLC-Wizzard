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

// configureCodeGeneration configures code generation options
func (s *FeaturesStep) configureCodeGeneration(data *generated.TemplateData) error {
	var (
		emitInterface       bool
		emitPreparedQueries bool
		emitJSONTags        bool
	)

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Generate Go interfaces?").
				Description("Create interfaces for query methods").
				Value(&emitInterface),
			huh.NewConfirm().
				Title("Generate prepared queries?").
				Description("Create prepared query methods for better performance").
				Value(&emitPreparedQueries),
			huh.NewConfirm().
				Title("Add JSON tags?").
				Description("Add JSON struct tags to generated models").
				Value(&emitJSONTags),
		),
	).WithTheme(s.theme)

	if err := form.Run(); err != nil {
		return fmt.Errorf("code generation configuration failed: %w", err)
	}

	data.Validation.EmitOptions.EmitInterface = emitInterface
	data.Validation.EmitOptions.EmitPreparedQueries = emitPreparedQueries
	data.Validation.EmitOptions.EmitJSONTags = emitJSONTags

	return nil
}

// configureSafetyRules configures safety rules
func (s *FeaturesStep) configureSafetyRules(data *generated.TemplateData) error {
	var (
		noSelectStar  bool
		requireWhere  bool
		requireLimit  bool
	)

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Forbid SELECT *?").
				Description("Prevent SELECT * queries for better performance and explicitness").
				Value(&noSelectStar),
			huh.NewConfirm().
				Title("Require WHERE clause?").
				Description("Force WHERE clauses in UPDATE/DELETE queries to prevent accidental data modification").
				Value(&requireWhere),
			huh.NewConfirm().
				Title("Require LIMIT on SELECT?").
				Description("Force LIMIT clauses on SELECT queries to prevent large result sets").
				Value(&requireLimit),
		),
	).WithTheme(s.theme)

	if err := form.Run(); err != nil {
		return fmt.Errorf("safety rules configuration failed: %w", err)
	}

	data.Validation.SafetyRules.NoSelectStar = noSelectStar
	data.Validation.SafetyRules.RequireWhere = requireWhere
	data.Validation.SafetyRules.RequireLimit = requireLimit

	return nil
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