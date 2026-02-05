package wizard

import (
	"fmt"

	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/apperrors"
	"github.com/charmbracelet/huh"
)

// OutputStep handles output configuration and file paths.
type OutputStep struct {
	theme *huh.Theme
	ui    *UIHelper
}

// NewOutputStep creates a new output step.
func NewOutputStep(theme *huh.Theme, ui *UIHelper) *OutputStep {
	return &OutputStep{
		theme: theme,
		ui:    ui,
	}
}

// Execute runs output configuration step.
func (s *OutputStep) Execute(data *generated.TemplateData) error {
	s.ui.ShowStepHeader("Output Configuration")

	var baseDir, queriesDir, schemaDir string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Base output directory").
				Description("Where generated Go code will be placed").
				Placeholder("./internal/db").
				Value(&baseDir),
			huh.NewInput().
				Title("SQL queries directory").
				Description("Directory containing your SQL query files").
				Placeholder("./sql/queries").
				Value(&queriesDir),
			huh.NewInput().
				Title("SQL schema directory").
				Description("Directory containing your SQL schema files").
				Placeholder("./sql/schema").
				Value(&schemaDir),
		),
	).WithTheme(s.theme)

	if err := form.Run(); err != nil {
		return fmt.Errorf("output configuration failed: %w", err)
	}

	// Set defaults if empty
	if baseDir == "" {
		baseDir = "./internal/db"
	}
	if queriesDir == "" {
		queriesDir = "./sql/queries"
	}
	if schemaDir == "" {
		schemaDir = "./sql/schema"
	}

	// Configure output
	data.Output.BaseDir = baseDir
	data.Output.QueriesDir = queriesDir
	data.Output.SchemaDir = schemaDir

	s.ui.ShowInfo("Base directory: " + baseDir)
	s.ui.ShowInfo("Queries directory: " + queriesDir)
	s.ui.ShowInfo("Schema directory: " + schemaDir)

	return nil
}

// ValidateConfiguration validates the complete output configuration.
func (s *OutputStep) ValidateConfiguration(data *generated.TemplateData) error {
	// Check for path conflicts
	if data.Output.BaseDir == data.Output.QueriesDir {
		return apperrors.NewError(apperrors.ErrorCodeValidationError, "base output directory cannot be the same as queries directory")
	}

	if data.Output.BaseDir == data.Output.SchemaDir {
		return apperrors.NewError(apperrors.ErrorCodeValidationError, "base output directory cannot be the same as schema directory")
	}

	if data.Output.QueriesDir == data.Output.SchemaDir {
		return apperrors.NewError(apperrors.ErrorCodeValidationError, "queries directory cannot be the same as schema directory")
	}

	// Check for reasonable directory structure
	if data.Output.BaseDir == "" {
		return apperrors.NewError(apperrors.ErrorCodeValidationError, "base output directory is required")
	}

	if data.Output.QueriesDir == "" {
		return apperrors.NewError(apperrors.ErrorCodeValidationError, "queries directory is required")
	}

	if data.Output.SchemaDir == "" {
		return apperrors.NewError(apperrors.ErrorCodeValidationError, "schema directory is required")
	}

	return nil
}
