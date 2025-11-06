package wizard

import (
	"fmt"
	"strings"

	"github.com/LarsArtmann/SQLC-Wizzard/internal/domain"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/templates"
	"github.com/LarsArtmann/SQLC-Wizzard/pkg/config"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

// WizardResult contains the output of running the wizard
type WizardResult struct {
	Config          *config.SqlcConfig
	TemplateData    templates.TemplateData
	GenerateQueries bool
	GenerateSchema  bool
}

// Wizard manages the interactive configuration flow
type Wizard struct {
	result *WizardResult
	theme  *huh.Theme
}

// NewWizard creates a new wizard instance
func NewWizard() *Wizard {
	return &Wizard{
		result: &WizardResult{
			GenerateQueries: true,
			GenerateSchema:  true,
		},
		theme: huh.ThemeBase(),
	}
}

// Run executes the interactive wizard
func (w *Wizard) Run() (*WizardResult, error) {
	// Display welcome banner
	w.showWelcome()

	// Initialize template data with defaults
	data := templates.TemplateData{
		Features:    templates.DefaultFeatures(),
		SafetyRules: domain.DefaultSafetyRules(),
	}

	// Step 1: Project Type
	if err := w.selectProjectType(&data); err != nil {
		return nil, err
	}

	// Step 2: Database
	if err := w.selectDatabase(&data); err != nil {
		return nil, err
	}

	// Step 3: Project Details
	if err := w.projectDetails(&data); err != nil {
		return nil, err
	}

	// Step 4: Features
	if err := w.selectFeatures(&data); err != nil {
		return nil, err
	}

	// Step 5: Output Configuration
	if err := w.outputConfiguration(&data); err != nil {
		return nil, err
	}

	// Generate config from template
	tmpl, err := templates.GetTemplate(data.ProjectType)
	if err != nil {
		return nil, fmt.Errorf("failed to get template: %w", err)
	}

	cfg, err := tmpl.Generate(data)
	if err != nil {
		return nil, fmt.Errorf("failed to generate config: %w", err)
	}

	w.result.Config = cfg
	w.result.TemplateData = data

	return w.result, nil
}

func (w *Wizard) showWelcome() {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("99")).
		Padding(1, 0)

	descStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Padding(0, 0, 1, 0)

	fmt.Println(titleStyle.Render("🧙‍♂️  SQLC Configuration Wizard"))
	fmt.Println(descStyle.Render("Let's create the perfect sqlc setup for your project!\n"))
}

func (w *Wizard) selectProjectType(data *templates.TemplateData) error {
	var projectType string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("What type of project are you building?").
				Options(
					huh.NewOption("🏠 Hobby - Simple SQLite setup", string(templates.ProjectTypeHobby)),
					huh.NewOption("⚡ Microservice - Single DB, container-optimized", string(templates.ProjectTypeMicroservice)),
					huh.NewOption("🏢 Enterprise - Multi-DB, comprehensive", string(templates.ProjectTypeEnterprise)),
					huh.NewOption("🔧 API-First - JSON-focused, REST-friendly", string(templates.ProjectTypeAPIFirst)),
					huh.NewOption("📦 Library - Embeddable, minimal deps", string(templates.ProjectTypeLibrary)),
				).
				Value(&projectType),
		),
	)

	if err := form.Run(); err != nil {
		return err
	}

	data.ProjectType = templates.ProjectType(projectType)
	return nil
}

func (w *Wizard) selectDatabase(data *templates.TemplateData) error {
	var database string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Which database will you use?").
				Options(
					huh.NewOption("🐘 PostgreSQL - Full-featured, recommended", string(templates.DatabaseTypePostgreSQL)),
					huh.NewOption("🗄️  SQLite - Lightweight, embedded", string(templates.DatabaseTypeSQLite)),
					huh.NewOption("🐬 MySQL - Popular, widely supported", string(templates.DatabaseTypeMySQL)),
				).
				Value(&database),
		),
	)

	if err := form.Run(); err != nil {
		return err
	}

	data.Database = templates.DatabaseType(database)
	return nil
}

func (w *Wizard) projectDetails(data *templates.TemplateData) error {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Project Name").
				Description("Used for naming the SQL configuration").
				Placeholder("my-awesome-service").
				Value(&data.ProjectName).
				Validate(func(s string) error {
					if strings.TrimSpace(s) == "" {
						return fmt.Errorf("project name is required")
					}
					return nil
				}),

			huh.NewInput().
				Title("Go Package Path").
				Description("The full import path for your project").
				Placeholder("github.com/user/project").
				Value(&data.PackagePath).
				Validate(func(s string) error {
					if strings.TrimSpace(s) == "" {
						return fmt.Errorf("package path is required")
					}
					if !strings.Contains(s, "/") {
						return fmt.Errorf("package path should be a full import path (e.g., github.com/user/project)")
					}
					return nil
				}),
		),
	)

	return form.Run()
}

func (w *Wizard) selectFeatures(data *templates.TemplateData) error {
	var features []string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Select database features to enable").
				Description("Choose features your project will use").
				Options(
					huh.NewOption("UUID support", "uuid"),
					huh.NewOption("JSON/JSONB columns", "json"),
					huh.NewOption("Array types", "arrays"),
					huh.NewOption("Full-text search", "fts"),
				).
				Value(&features).
				Limit(4),
		),
	)

	if err := form.Run(); err != nil {
		return err
	}

	// Update features based on selection
	data.Features.UUIDs = contains(features, "uuid")
	data.Features.JSON = contains(features, "json")
	data.Features.Arrays = contains(features, "arrays")
	data.Features.FullTextSearch = contains(features, "fts")

	return nil
}

func (w *Wizard) outputConfiguration(data *templates.TemplateData) error {
	// Set defaults if not already set
	if data.OutputDir == "" {
		data.OutputDir = "internal/db"
	}
	if data.PackageName == "" {
		data.PackageName = "db"
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Output Directory").
				Description("Where generated Go code will be placed").
				Value(&data.OutputDir).
				Validate(func(s string) error {
					if strings.TrimSpace(s) == "" {
						return fmt.Errorf("output directory is required")
					}
					return nil
				}),

			huh.NewInput().
				Title("Go Package Name").
				Description("Package name for generated code").
				Value(&data.PackageName).
				Validate(func(s string) error {
					if strings.TrimSpace(s) == "" {
						return fmt.Errorf("package name is required")
					}
					return nil
				}),

			huh.NewConfirm().
				Title("Generate example queries?").
				Description("Create starter CRUD queries").
				Value(&w.result.GenerateQueries),

			huh.NewConfirm().
				Title("Generate example schema?").
				Description("Create a sample database schema").
				Value(&w.result.GenerateSchema),
		),
	)

	return form.Run()
}

// contains checks if a string slice contains a value
func contains(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}
