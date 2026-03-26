package wizard

import (
	"fmt"

	"charm.land/huh/v2"
	"github.com/LarsArtmann/SQLC-Wizzard/generated"
)

// ProjectTypeStep handles project type selection.
type ProjectTypeStep struct {
	themeFunc huh.ThemeFunc
	ui        *UIHelper
}

// NewProjectTypeStep creates a new project type step.
func NewProjectTypeStep(themeFunc huh.ThemeFunc, ui *UIHelper) *ProjectTypeStep {
	return &ProjectTypeStep{
		themeFunc: themeFunc,
		ui:        ui,
	}
}

// Execute runs the project type selection step.
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
	).WithTheme(s.themeFunc)

	err := form.Run()
	if err != nil {
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
