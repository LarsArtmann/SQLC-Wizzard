package wizard

import (
	"fmt"

	"github.com/LarsArtmann/SQLC-Wizzard/pkg/config"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/templates"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

// UIHelper manages UI styling and display
type UIHelper struct {
	theme *huh.Theme
}

// NewUIHelper creates a new UI helper
func NewUIHelper() *UIHelper {
	return &UIHelper{
		theme: huh.ThemeBase(),
	}
}

// showWelcome displays welcome banner
func (ui *UIHelper) ShowWelcome() {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("99")).
		Padding(1, 0)

	descStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Padding(0, 0, 1, 0)

	fmt.Println(titleStyle.Render("🧙‍♂️  SQLC Configuration Wizard"))
	fmt.Println(descStyle.Render("Let's create a perfect sqlc setup for your project!\n"))
}

// showPreview displays configuration preview
func (ui *UIHelper) ShowPreview(data *templates.TemplateData, cfg *config.SqlcConfig) string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7D56F4")).
		Padding(0, 1)

	sectionStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FF7E67")).
		MarginBottom(1)

	contentStyle := lipgloss.NewStyle().
		PaddingLeft(2).
		Width(76).
		Align(lipgloss.Left)

	preview := titleStyle.Render("Configuration Preview")
	
	preview += "\n" + sectionStyle.Render("Project") + "\n" + 
		contentStyle.Render(fmt.Sprintf(`
  Name: %s
  Type: %s
  Database: %s
  Output: %s
  Package: %s
`, data.ProjectName, data.ProjectType, data.Database.Engine, data.Output.BaseDir, data.Package.Path))

	preview += "\n" + sectionStyle.Render("Generation") + "\n" +
		contentStyle.Render(fmt.Sprintf(`
  Queries: %t
  Schema: %t
  UUIDs: %t
  JSON Tags: %t
  Strict Mode: %t
`, true, true, data.Database.UseUUIDs, true, data.Validation.StrictFunctions))

	return preview
}

// getConfirmation shows final confirmation
func (ui *UIHelper) GetConfirmation() (bool, error) {
	var confirmed bool
	
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Generate configuration with these settings?").
				Description("You can edit this later in the generated yaml file").
				Value(&confirmed),
		),
	).WithTheme(ui.theme).Run()

	if err != nil {
		return false, err
	}

	if !confirmed {
		return false, fmt.Errorf("configuration cancelled by user")
	}

	return true, nil
}

// ShowSuccess displays success message
func (ui *UIHelper) ShowSuccess(outputPath string) {
	successStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#00FF00")).
		Padding(1, 0)

	pathStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Padding(0, 0, 1, 0)

	fmt.Println(successStyle.Render("✅ Configuration generated successfully!"))
	fmt.Println(pathStyle.Render("📁 Output: " + outputPath))
}

// ShowError displays error message
func (ui *UIHelper) ShowError(err error) {
	errorStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FF0000")).
		Padding(1, 0)

	fmt.Println(errorStyle.Render("❌ Error: " + err.Error()))
}