package wizard

import (
	"fmt"

	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/templates"
	"github.com/LarsArtmann/SQLC-Wizzard/pkg/config"
	"github.com/charmbracelet/huh"
)

// WizardResult contains the output of running the wizard
type WizardResult struct {
	Config          *config.SqlcConfig
	TemplateData    generated.TemplateData
	GenerateQueries bool
	GenerateSchema  bool
}

// Wizard manages the interactive configuration flow
type Wizard struct {
	result *WizardResult
	theme  *huh.Theme
	ui     *UIHelper
	deps   *WizardDependencies // For dependency injection in tests

	// Step handlers
	projectTypeStep *ProjectTypeStep
	databaseStep    *DatabaseStep
	projectDetails  *ProjectDetailsStep
	featuresStep    *FeaturesStep
	outputStep      *OutputStep
}

// NewWizard creates a new wizard instance
func NewWizard() *Wizard {
	theme := huh.ThemeBase()
	ui := NewUIHelper()

	deps := WizardDependencies{
		UI:           ui,
		ProjectType:  NewProjectTypeStep(theme, ui),
		Database:     NewDatabaseStep(theme, ui),
		Details:      NewProjectDetailsStep(theme, ui),
		Features:     NewFeaturesStep(theme, ui),
		Output:       NewOutputStep(theme, ui),
		TemplateFunc: func(projectType templates.ProjectType) (TemplateInterface, error) {
			tmpl, err := templates.GetTemplate(projectType)
			if err != nil {
				return nil, err
			}
			return tmpl, nil
		},
	}

	return &Wizard{
		result: &WizardResult{
			GenerateQueries: true,
			GenerateSchema:  true,
		},
		theme: theme,
		ui:    ui,
		deps:  &deps,

		// Initialize step handlers
		projectTypeStep: NewProjectTypeStep(theme, ui),
		databaseStep:    NewDatabaseStep(theme, ui),
		projectDetails:  NewProjectDetailsStep(theme, ui),
		featuresStep:    NewFeaturesStep(theme, ui),
		outputStep:      NewOutputStep(theme, ui),
	}
}

// GetResult returns the current wizard result
func (w *Wizard) GetResult() *WizardResult {
	return w.result
}

// Run executes the interactive wizard
func (w *Wizard) Run() (*WizardResult, error) {
	// Display welcome banner
	w.showWelcome()

	// Initialize template data with defaults
	data := generated.TemplateData{
		Package: generated.PackageConfig{
			Name: "myproject",
			Path: "github.com/myorg/myproject",
		},
		Database: generated.DatabaseConfig{
			UseUUIDs:    true,
			UseJSON:     true,
			UseArrays:   false,
			UseFullText: false,
		},
		Output: generated.OutputConfig{
			BaseDir:    "./internal/db",
			QueriesDir: "./sql/queries",
			SchemaDir:  "./sql/schema",
		},
		Validation: generated.ValidationConfig{
			EmitOptions: generated.DefaultEmitOptions(),
			SafetyRules: generated.DefaultSafetyRules(),
		},
	}

	// Execute wizard steps in order
	steps := []struct {
		name    string
		execute func(*generated.TemplateData) error
	}{
		{"Project Type", w.getProjectTypeStep().Execute},
		{"Database", w.getDatabaseStep().Execute},
		{"Project Details", w.getProjectDetailsStep().Execute},
		{"Features", w.getFeaturesStep().Execute},
		{"Output Configuration", w.getOutputStep().Execute},
	}

	for _, step := range steps {
		w.showStepHeader(step.name)

		if err := step.execute(&data); err != nil {
			return nil, fmt.Errorf("step '%s' failed: %w", step.name, err)
		}

		w.showStepComplete(step.name, "Completed successfully")
	}

	// Generate config from template
	if err := w.generateConfig(&data); err != nil {
		return nil, fmt.Errorf("config generation failed: %w", err)
	}

	// Show final summary
	w.showSummary(&data)

	return w.result, nil
}

// Helper methods to get the appropriate step implementation
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

// Helper methods for UI operations
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

// generateConfig generates the final sqlc configuration
func (w *Wizard) generateConfig(data *generated.TemplateData) error {
	// Validate output configuration
	if outputStep := w.getOutputStep(); outputStep != nil {
		if validatableStep, ok := outputStep.(interface{ ValidateConfiguration(*generated.TemplateData) error }); ok {
			if err := validatableStep.ValidateConfiguration(data); err != nil {
				return fmt.Errorf("invalid output configuration: %w", err)
			}
		}
	}

	// Get appropriate template
	var tmpl TemplateInterface
	var err error
	
	if w.deps != nil && w.deps.TemplateFunc != nil {
		tmpl, err = w.deps.TemplateFunc(templates.ProjectType(data.ProjectType))
	} else {
		tmpl, err = templates.GetTemplate(templates.ProjectType(data.ProjectType))
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

// showSummary displays the final configuration summary
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
