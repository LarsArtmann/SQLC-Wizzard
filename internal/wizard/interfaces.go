package wizard

import (
	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/templates"
)

// UIInterface defines the interface for UI operations
// This allows for mocking in tests.
type UIInterface interface {
	ShowStepHeader(title string)
	ShowStepComplete(title, message string)
	ShowSection(title string)
	ShowInfo(message string)
	ShowWelcome()
}

// StepInterface defines the interface for wizard steps.
type StepInterface interface {
	Execute(data *generated.TemplateData) error
}

// ValidatableStepInterface extends StepInterface with validation capability.
type ValidatableStepInterface interface {
	StepInterface
	ValidateConfiguration(data *generated.TemplateData) error
}

// WizardDependencies contains all wizard dependencies for dependency injection.
type WizardDependencies struct {
	UI           UIInterface
	ProjectType  StepInterface
	Database     StepInterface
	Details      StepInterface
	Features     StepInterface
	Output       StepInterface
	TemplateFunc func(projectType templates.ProjectType) (templates.Template, error)
}

// NewTestableWizard creates a wizard with injected dependencies for testing.
func NewTestableWizard(deps WizardDependencies) *Wizard {
	return &Wizard{
		result: &WizardResult{
			GenerateQueries: true,
			GenerateSchema:  true,
		},
		deps: &deps,
	}
}
