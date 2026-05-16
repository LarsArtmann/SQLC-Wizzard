package wizard

import (
	"fmt"

	"charm.land/huh/v2"
	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/templates"
	"github.com/LarsArtmann/SQLC-Wizzard/pkg/config"
)

// WizardResult contains the output of running the wizard.
type WizardResult struct {
	Config          *config.SqlcConfig
	TemplateData    generated.TemplateData
	GenerateQueries bool
	GenerateSchema  bool
}

// Wizard manages the interactive configuration flow.
type Wizard struct {
	result    *WizardResult
	themeFunc huh.ThemeFunc
	ui        *UIHelper
	deps      *WizardDependencies // For dependency injection in tests
	context   *FlowContext        // Branching flow context

	// Step handlers
	projectTypeStep *ProjectTypeStep
	databaseStep    *DatabaseStep
	projectDetails  *ProjectDetailsStep
	featuresStep    *FeaturesStep
	outputStep      *OutputStep
}

// NewWizard creates a new wizard instance.
func NewWizard() *Wizard {
	themeFunc := huh.ThemeBase
	ui := NewUIHelper()

	deps := WizardDependencies{
		UI:          ui,
		ProjectType: NewProjectTypeStep(themeFunc, ui),
		Database:    NewDatabaseStep(themeFunc, ui),
		Details:     NewProjectDetailsStep(themeFunc, ui),
		Features:    NewFeaturesStep(themeFunc, ui),
		Output:      NewOutputStep(themeFunc, ui),
		TemplateFunc: func(projectType templates.ProjectType) (templates.Template, error) {
			tmpl, err := templates.GetTemplate(projectType)
			if err != nil {
				return nil, fmt.Errorf("failed to get template for %s: %w", projectType, err)
			}

			return tmpl, nil
		},
	}

	return &Wizard{
		result: &WizardResult{
			GenerateQueries: true,
			GenerateSchema:  true,
		},
		themeFunc: themeFunc,
		ui:        ui,
		deps:      &deps,
		context:   NewFlowContext(),

		// Initialize step handlers
		projectTypeStep: NewProjectTypeStep(themeFunc, ui),
		databaseStep:    NewDatabaseStep(themeFunc, ui),
		projectDetails:  NewProjectDetailsStep(themeFunc, ui),
		featuresStep:    NewFeaturesStep(themeFunc, ui),
		outputStep:      NewOutputStep(themeFunc, ui),
	}
}

// GetResult returns the current wizard result.
func (w *Wizard) GetResult() *WizardResult {
	return w.result
}

// Run executes the interactive wizard with branching flow support.
func (w *Wizard) Run() (*WizardResult, error) {
	// Display welcome banner
	w.showWelcome()

	// Initialize template data with defaults
	data := generated.DefaultTemplateData()

	// Get dynamic steps based on flow context
	steps := w.buildStepList(&data)

	for _, step := range steps {
		w.showStepHeader(step.name)

		err := step.execute(&data)
		if err != nil {
			return nil, fmt.Errorf("step '%s' failed: %w", step.name, err)
		}

		// Update flow context with completed step
		w.context.MarkStepCompleted(step.id)
		w.context.UpdateFromTemplateData(&data)

		w.showStepComplete(step.name, "Completed successfully")
	}

	// Generate config from template
	err := w.generateConfig(&data)
	if err != nil {
		return nil, fmt.Errorf("config generation failed: %w", err)
	}

	// Show final summary
	w.showSummary(&data)

	return w.result, nil
}

// stepDefinition defines a single step in the wizard flow.
type stepDefinition struct {
	name    string
	id      StepID
	execute func(*generated.TemplateData) error
}

// buildStepList builds the dynamic step list based on flow context.
func (w *Wizard) buildStepList(data *generated.TemplateData) []stepDefinition {
	var steps []stepDefinition

	// Project Type step - always first
	steps = append(steps, stepDefinition{
		name:    "Project Type",
		id:      StepProjectType,
		execute: w.getProjectTypeStep().Execute,
	})

	// Update context with project type for dynamic step building
	if data != nil && data.ProjectType != "" {
		w.context.ProjectType = data.ProjectType
	}

	// Database step - always second
	steps = append(steps, stepDefinition{
		name:    "Database",
		id:      StepDatabase,
		execute: w.getDatabaseStep().Execute,
	})

	// Update context with database type
	if data != nil && data.Database.Engine != "" {
		w.context.DatabaseType = data.Database.Engine
	}

	// Project Details step - always third
	steps = append(steps, stepDefinition{
		name:    "Project Details",
		id:      StepProjectDetail,
		execute: w.getProjectDetailsStep().Execute,
	})

	// Features step - conditional based on project type
	shouldShowFeatures := w.context.ProjectType != generated.ProjectTypeHobby &&
		w.context.ProjectType != generated.ProjectTypeTesting
	if shouldShowFeatures {
		steps = append(steps, stepDefinition{
			name:    "Features",
			id:      StepFeatures,
			execute: w.getFeaturesStep().Execute,
		})
	}

	// Output step - always last
	steps = append(steps, stepDefinition{
		name:    "Output Configuration",
		id:      StepOutput,
		execute: w.getOutputStep().Execute,
	})

	return steps
}

// GetFlowContext returns the current flow context.
func (w *Wizard) GetFlowContext() *FlowContext {
	return w.context
}

// Helper methods to get the appropriate step implementation.
func (w *Wizard) getProjectTypeStep() StepInterface {
	if w.deps != nil && w.deps.ProjectType != nil {
		return w.deps.ProjectType
	}

	return w.projectTypeStep
}

func (w *Wizard) getDatabaseStep() StepInterface {
	if w.deps != nil && w.deps.Database != nil {
		return w.deps.Database
	}

	return w.databaseStep
}

func (w *Wizard) getProjectDetailsStep() StepInterface {
	if w.deps != nil && w.deps.Details != nil {
		return w.deps.Details
	}

	return w.projectDetails
}

func (w *Wizard) getFeaturesStep() StepInterface {
	if w.deps != nil && w.deps.Features != nil {
		return w.deps.Features
	}

	return w.featuresStep
}

func (w *Wizard) getOutputStep() StepInterface {
	if w.deps != nil && w.deps.Output != nil {
		return w.deps.Output
	}

	return w.outputStep
}

// Helper methods for UI operations.
func (w *Wizard) showWelcome() {
	if w.deps != nil && w.deps.UI != nil {
		w.deps.UI.ShowWelcome()

		return
	}

	w.ui.ShowWelcome()
}

func (w *Wizard) showStepHeader(title string) {
	if w.deps != nil && w.deps.UI != nil {
		w.deps.UI.ShowStepHeader(title)

		return
	}

	w.ui.ShowStepHeader(title)
}

func (w *Wizard) showStepComplete(title, message string) {
	if w.deps != nil && w.deps.UI != nil {
		w.deps.UI.ShowStepComplete(title, message)

		return
	}

	w.ui.ShowStepComplete(title, message)
}

// generateConfig generates the final sqlc configuration.
func (w *Wizard) generateConfig(data *generated.TemplateData) error {
	// Validate output configuration
	if outputStep := w.getOutputStep(); outputStep != nil {
		if validatableStep, ok := outputStep.(interface {
			ValidateConfiguration(*generated.TemplateData) error
		}); ok {
			err := validatableStep.ValidateConfiguration(data)
			if err != nil {
				return fmt.Errorf("invalid output configuration: %w", err)
			}
		}
	}

	// Get appropriate template
	var (
		tmpl templates.Template
		err  error
	)

	if w.deps != nil && w.deps.TemplateFunc != nil {
		tmpl, err = w.deps.TemplateFunc(data.ProjectType)
	} else {
		tmpl, err = templates.GetTemplate(data.ProjectType)
	}

	if err != nil {
		return fmt.Errorf("failed to get template: %w", err)
	}

	// Generate configuration
	cfg, err := tmpl.Generate(*data)
	if err != nil {
		return fmt.Errorf("failed to generate config: %w", err)
	}

	w.result.Config = cfg
	w.result.TemplateData = *data

	return nil
}

// showSummary displays the final configuration summary.
func (w *Wizard) showSummary(data *generated.TemplateData) {
	var ui UIInterface
	if w.deps != nil && w.deps.UI != nil {
		ui = w.deps.UI
	} else {
		ui = w.ui
	}

	ui.ShowSection("🎉 Configuration Complete")

	summary := fmt.Sprintf(`
Project: %s
Package: %s
Type: %s
Database: %s
Output: %s

Features:
- Interfaces: %t
- Prepared Queries: %t
- JSON Tags: %t

Safety Rules:
- No SELECT *: %t
- Require WHERE: %t
- Require LIMIT: %t

Database Features:
- UUIDs: %t
- JSON: %t
- Arrays: %t
- Full-text: %t
`,
		data.ProjectName,
		data.Package.Name,
		data.ProjectType,
		data.Database.Engine,
		data.Output.BaseDir,
		data.Validation.EmitOptions.EmitInterface,
		data.Validation.EmitOptions.EmitPreparedQueries,
		data.Validation.EmitOptions.EmitJSONTags,
		data.Validation.SafetyRules.NoSelectStar,
		data.Validation.SafetyRules.RequireWhere,
		data.Validation.SafetyRules.RequireLimit,
		data.Database.UseUUIDs,
		data.Database.UseJSON,
		data.Database.UseArrays,
		data.Database.UseFullText,
	)

	ui.ShowInfo(summary)
}
