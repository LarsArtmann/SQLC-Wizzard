package wizard

import (
	"fmt"

	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/templates"
	"github.com/charmbracelet/huh"
)

// ProjectTypeStep handles project type selection
type ProjectTypeStep struct {
	theme *huh.Theme
	ui    *UIHelper
}

// NewProjectTypeStep creates a new project type step
func NewProjectTypeStep(theme *huh.Theme, ui *UIHelper) *ProjectTypeStep {
	return &ProjectTypeStep{
		theme: theme,
		ui:    ui,
	}
}

// Execute runs the project type selection step
func (s *ProjectTypeStep) Execute(data *generated.TemplateData) error {
	s.ui.ShowStepHeader("Project Type Selection")

	var projectType string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("What type of project are you building?").
				Options(
					huh.NewOption("🏠 Hobby - Simple SQLite setup", string(generated.ProjectTypeHobby)),
					huh.NewOption("⚡ Microservice - Single DB, container-optimized", string(generated.ProjectTypeMicroservice)),
					huh.NewOption("🏢 Enterprise - Multi-DB, comprehensive", string(generated.ProjectTypeEnterprise)),
					huh.NewOption("🔧 API-First - JSON-focused, REST-friendly", string(generated.ProjectTypeAPIFirst)),
					huh.NewOption("📊 Analytics - Read-optimized, warehousing", string(generated.ProjectTypeAnalytics)),
					huh.NewOption("🧪 Testing - Isolated, disposable", string(generated.ProjectTypeTesting)),
					huh.NewOption("🏗️  Multi-Tenant - Shared resources", string(generated.ProjectTypeMultiTenant)),
					huh.NewOption("📦 Library - Embeddable, minimal deps", string(generated.ProjectTypeLibrary)),
				).
				Value(&projectType),
		),
	).WithTheme(s.theme)

	if err := form.Run(); err != nil {
		return fmt.Errorf("project type selection failed: %w", err)
	}

	// Validate project type
	pt := generated.ProjectType(projectType)
	if !pt.IsValid() {
		return fmt.Errorf("invalid project type: %s", projectType)
	}

	data.ProjectType = pt
	s.ui.ShowStepComplete("Project Type", string(data.ProjectType))
	
	return nil
}