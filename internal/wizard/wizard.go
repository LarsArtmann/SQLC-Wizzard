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

	return &Wizard{
		result: &WizardResult{
			GenerateQueries: true,
			GenerateSchema:  true,
		},
		theme: theme,
		ui:    ui,

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
	w.ui.ShowWelcome()

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
		{"Project Type", w.projectTypeStep.Execute},
		{"Database", w.databaseStep.Execute},
		{"Project Details", w.projectDetails.Execute},
		{"Features", w.featuresStep.Execute},
		{"Output Configuration", w.outputStep.Execute},
	}

	for _, step := range steps {
		w.ui.ShowStepHeader(step.name)

		if err := step.execute(&data); err != nil {
			return nil, fmt.Errorf("step '%s' failed: %w", step.name, err)
		}

		w.ui.ShowStepComplete(step.name, "Completed successfully")
	}

	// Generate config from template
	if err := w.generateConfig(&data); err != nil {
		return nil, fmt.Errorf("config generation failed: %w", err)
	}

	// Show final summary
	w.showSummary(&data)

	return w.result, nil
}

// generateConfig generates the final sqlc configuration
func (w *Wizard) generateConfig(data *generated.TemplateData) error {
	// Validate output configuration
	if err := w.outputStep.ValidateConfiguration(data); err != nil {
		return fmt.Errorf("invalid output configuration: %w", err)
	}

	// Get appropriate template
	tmpl, err := templates.GetTemplate(data.ProjectType)
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
	w.ui.ShowSection("🎉 Configuration Complete")

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

	w.ui.ShowInfo(summary)
}
