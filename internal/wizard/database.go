package wizard

import (
	"fmt"

	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/charmbracelet/huh"
)

// DatabaseStep handles database selection and configuration.
type DatabaseStep struct {
	theme *huh.Theme
	ui    *UIHelper
}

// NewDatabaseStep creates a new database step.
func NewDatabaseStep(theme *huh.Theme, ui *UIHelper) *DatabaseStep {
	return &DatabaseStep{
		theme: theme,
		ui:    ui,
	}
}

// Execute runs the database selection step.
func (s *DatabaseStep) Execute(data *generated.TemplateData) error {
	s.ui.ShowStepHeader("Database Selection")

	var database string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Which database will you use?").
				Options(
					huh.NewOption("🐘 PostgreSQL - Full-featured, recommended", string(generated.DatabaseTypePostgreSQL)),
					huh.NewOption("🗄️  SQLite - Lightweight, embedded", string(generated.DatabaseTypeSQLite)),
					huh.NewOption("🐬 MySQL - Popular, widely supported", string(generated.DatabaseTypeMySQL)),
				).
				Value(&database),
		),
	).WithTheme(s.theme)

	if err := form.Run(); err != nil {
		return fmt.Errorf("database selection failed: %w", err)
	}

	// Validate database type
	dt := generated.DatabaseType(database)
	if !dt.IsValid() {
		return fmt.Errorf("invalid database type: %s", database)
	}

	data.Database.Engine = dt
	s.ui.ShowStepComplete("Database", string(data.Database.Engine))

	// Database-specific configuration
	return s.configureDatabaseOptions(data)
}

// configureDatabaseOptions configures database-specific options.
func (s *DatabaseStep) configureDatabaseOptions(data *generated.TemplateData) error {
	switch data.Database.Engine {
	case generated.DatabaseTypePostgreSQL:
		return s.applyDatabaseConfig(data, "PostgreSQL", true, true, true, true, "UUIDs, JSON, Full-text, Arrays enabled")
	case generated.DatabaseTypeSQLite:
		return s.applyDatabaseConfig(data, "SQLite", false, false, false, false, "Lightweight configuration applied")
	case generated.DatabaseTypeMySQL:
		return s.applyDatabaseConfig(data, "MySQL", true, true, true, false, "UUIDs, JSON, Full-text enabled")
	default:
		return fmt.Errorf("unsupported database engine: %s", data.Database.Engine)
	}
}

// applyDatabaseConfig applies database-specific configuration to the template data.
func (s *DatabaseStep) applyDatabaseConfig(
	data *generated.TemplateData,
	engineName string,
	useUUIDs, useJSON, useFullText, useArrays bool,
	completionMessage string,
) error {
	s.ui.ShowInfo(fmt.Sprintf("Configuring %s options...", engineName))

	data.Database.UseUUIDs = useUUIDs
	data.Database.UseJSON = useJSON
	data.Database.UseFullText = useFullText
	data.Database.UseArrays = useArrays

	s.ui.ShowStepComplete(fmt.Sprintf("%s Config", engineName), completionMessage)
	return nil
}
