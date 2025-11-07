package wizard

import (
	"fmt"
	"strings"

	"github.com/LarsArtmann/SQLC-Wizzard/internal/templates"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

// selectProject handles project type selection
func (w *Wizard) selectProject(data *templates.TemplateData) error {
	return huh.NewForm(
		huh.NewGroup("Project Configuration",
			huh.NewInput().
				Title("Project name").
				Placeholder("my-sqlc-project").
				Value(&data.ProjectName).
				Validate(func(s string) error {
					if len(strings.TrimSpace(s)) == 0 {
						return fmt.Errorf("project name is required")
					}
					return nil
				}),

			huh.NewSelect[templates.ProjectType]().
				Title("Project type").
				Options(
					huh.NewOption("🏠 Hobby", "hobby"),
					huh.NewOption("🔧 Microservice", "microservice"),
					huh.NewOption("🏢 Enterprise", "enterprise"),
					huh.NewOption("🚀 API-First", "api-first"),
					huh.NewOption("📊 Analytics", "analytics"),
					huh.NewOption("🧪 Testing", "testing"),
					huh.NewOption("🏗️  Multi-Tenant", "multi-tenant"),
					huh.NewOption("📚 Library", "library"),
				).
				Value(&data.ProjectType),
		),
	).WithTheme(w.theme).Run()
}

// selectDatabase handles database selection
func (w *Wizard) selectDatabase(data *generated.TemplateData) error {
	return huh.NewForm(
		huh.NewGroup("Database Configuration",
			huh.NewSelect[generated.DatabaseType]().
				Title("Database").
				Description("Select the primary database for your project").
				Options(
					huh.NewOption("🐘 PostgreSQL", "postgresql"),
					huh.NewOption("🐬 MySQL", "mysql"),
					huh.NewOption("📁 SQLite", "sqlite"),
				).
				Value(&data.Database),

			huh.NewConfirm().
				Title("Use managed database?").
				Description("Let SQLC-Wizard handle database setup and migrations").
				Value(&data.UseManagedDB),
		),
	).WithTheme(w.theme).Run()
}

// projectDetails handles project detail collection
func (w *Wizard) projectDetails(data *generated.TemplateData) error {
	return huh.NewForm(
		huh.NewGroup("Project Details",
			huh.NewInput().
				Title("Package path").
				Placeholder("github.com/user/project").
				Value(&data.PackagePath).
				Validate(func(s string) error {
					if len(strings.TrimSpace(s)) == 0 {
						return fmt.Errorf("package path is required")
					}
					return nil
				}),

			huh.NewInput().
				Title("Output directory").
				Placeholder("./generated").
				Value(&data.OutputDir).
				Validate(func(s string) error {
					if len(strings.TrimSpace(s)) == 0 {
						return fmt.Errorf("output directory is required")
					}
					return nil
				}),

			huh.NewInput().
				Title("Package name").
				Placeholder("models").
				Value(&data.PackageName).
				Validate(func(s string) error {
					if len(strings.TrimSpace(s)) == 0 {
						return fmt.Errorf("package name is required")
					}
					return nil
				}),
		),
	).WithTheme(w.theme).Run()
}