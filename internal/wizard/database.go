package wizard

import (
	"fmt"

	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/charmbracelet/huh"
)

// DatabaseStep handles database selection and configuration
type DatabaseStep struct {
	theme *huh.Theme
	ui    *UIHelper
}

// NewDatabaseStep creates a new database step
func NewDatabaseStep(theme *huh.Theme, ui *UIHelper) *DatabaseStep {
	return &DatabaseStep{
		theme: theme,
		ui:    ui,
	}
}

// Execute runs the database selection step
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

// configureDatabaseOptions configures database-specific options
func (s *DatabaseStep) configureDatabaseOptions(data *generated.TemplateData) error {
	switch data.Database.Engine {
	case generated.DatabaseTypePostgreSQL:
		return s.configurePostgreSQL(data)
	case generated.DatabaseTypeSQLite:
		return s.configureSQLite(data)
	case generated.DatabaseTypeMySQL:
		return s.configureMySQL(data)
	default:
		return fmt.Errorf("unsupported database engine: %s", data.Database.Engine)
	}
}

// configurePostgreSQL configures PostgreSQL-specific options
func (s *DatabaseStep) configurePostgreSQL(data *generated.TemplateData) error {
	s.ui.ShowInfo("Configuring PostgreSQL options...")

	// Default PostgreSQL configuration
	data.Database.UseUUIDs = true
	data.Database.UseJSON = true
	data.Database.UseFullText = true
	data.Database.UseArrays = true

	s.ui.ShowStepComplete("PostgreSQL Config", "UUIDs, JSON, Full-text, Arrays enabled")
	return nil
}

// configureSQLite configures SQLite-specific options
func (s *DatabaseStep) configureSQLite(data *generated.TemplateData) error {
	s.ui.ShowInfo("Configuring SQLite options...")

	// Default SQLite configuration - minimal features
	data.Database.UseUUIDs = false
	data.Database.UseJSON = false
	data.Database.UseFullText = false
	data.Database.UseArrays = false

	s.ui.ShowStepComplete("SQLite Config", "Lightweight configuration applied")
	return nil
}

// configureMySQL configures MySQL-specific options
func (s *DatabaseStep) configureMySQL(data *generated.TemplateData) error {
	s.ui.ShowInfo("Configuring MySQL options...")

	// Default MySQL configuration
	data.Database.UseUUIDs = true
	data.Database.UseJSON = true
	data.Database.UseFullText = true
	data.Database.UseArrays = false // MySQL arrays are not standard

	s.ui.ShowStepComplete("MySQL Config", "UUIDs, JSON, Full-text enabled")
	return nil
}
