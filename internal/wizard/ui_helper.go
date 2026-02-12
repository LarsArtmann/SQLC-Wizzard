package wizard

import (
	"fmt"

	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/apperrors"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/schema"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/templates"
	"github.com/LarsArtmann/SQLC-Wizzard/pkg/config"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

// UIHelper manages UI styling and display.
type UIHelper struct {
	theme *huh.Theme
}

// NewUIHelper creates a new UI helper.
func NewUIHelper() *UIHelper {
	return &UIHelper{
		theme: huh.ThemeBase(),
	}
}

// ShowStepHeader displays a step header.
func (ui *UIHelper) ShowStepHeader(title string) {
	ui.showTitledSection(title, 1, 0)
}

// ShowStepComplete displays a step completion message.
func (ui *UIHelper) ShowStepComplete(title, message string) {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#00D084")).
		Padding(1, 0)

	messageStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00D084")).
		PaddingLeft(2)

	fmt.Println(titleStyle.Render("✅ " + title))
	fmt.Println(messageStyle.Render(message))
	fmt.Println()
}

// ShowSection displays a section header.
func (ui *UIHelper) ShowSection(title string) {
	ui.showTitledSection(title, 0, 1)
}

// createTitleStyle creates a consistent title style with configurable padding.
func (ui *UIHelper) createTitleStyle(vertical, horizontal int) lipgloss.Style {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7D56F4")).
		Padding(vertical, horizontal).
		MarginBottom(1)
}

// showTitledSection renders a titled section with consistent styling.
func (ui *UIHelper) showTitledSection(title string, vertical, horizontal int) {
	titleStyle := ui.createTitleStyle(vertical, horizontal)
	fmt.Println(titleStyle.Render("📍 " + title))
}

// ShowInfo displays information.
func (ui *UIHelper) ShowInfo(message string) {
	infoStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		PaddingLeft(2).
		Width(76).
		Align(lipgloss.Left)

	fmt.Println(infoStyle.Render(message))
}

// showWelcome displays welcome banner.
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

// showPreview displays configuration preview.
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

// getConfirmation shows final confirmation.
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
		return false, apperrors.NewError(apperrors.ErrorCodeValidationError, "configuration cancelled by user")
	}

	return true, nil
}

// formatConfigurationSummary formats configuration for display.
func (ui *UIHelper) formatConfigurationSummary(cfg *schema.Schema, data generated.TemplateData) string {
	summaryStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF7E67")).
		Padding(0, 1)

	summary := summaryStyle.Render("Configuration Summary")
	summary += "\n" + fmt.Sprintf("Schema: %s (Tables: %d)", cfg.Name, len(cfg.Tables))
	summary += "\n" + fmt.Sprintf("Project: %s (%s)", data.ProjectName, data.ProjectType)
	summary += "\n" + fmt.Sprintf("Database: %s", data.Database.Engine)
	summary += "\n" + "Output: " + data.Output.BaseDir

	return summary
}

// formatCompletionDetails formats completion details for display.
func (ui *UIHelper) formatCompletionDetails(cfg *schema.Schema, data generated.TemplateData) string {
	detailStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#99")).
		PaddingLeft(2)

	details := detailStyle.Render("Generated Files:")
	details += "\n" + "- sqlc.yaml configuration"
	details += "\n" + fmt.Sprintf("- Database schema (%d tables)", len(cfg.Tables))
	details += "\n" + "- Query files (based on schema)"

	return details
}

// showErrorWithSchemaDetails displays schema apperrors.
func (ui *UIHelper) showErrorWithSchemaDetails(err *schema.SchemaError) {
	errorStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FF5555")).
		Padding(0, 1)

	detailStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF8888")).
		PaddingLeft(2)

	fmt.Println(errorStyle.Render("❌ Schema Error"))
	fmt.Println(detailStyle.Render("Code: " + err.Code))
	fmt.Println(detailStyle.Render("Message: " + err.Message))
}

// showErrorWithTypedDetails displays typed apperrors.
func (ui *UIHelper) showErrorWithTypedDetails(err *apperrors.Error) {
	errorStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FF5555")).
		Padding(0, 1)

	detailStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF8888")).
		PaddingLeft(2)

	fmt.Println(errorStyle.Render("❌ Error"))
	fmt.Println(detailStyle.Render("Code: " + string(err.Code)))
	fmt.Println(detailStyle.Render("Message: " + err.Message))

	if err.Description != "" {
		fmt.Println(detailStyle.Render("Description: " + err.Description))
	}
}
