package wizard

import (
	"fmt"

	"charm.land/huh/v2"
	"charm.land/lipgloss/v2"
	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/apperrors"
	uistyles "github.com/LarsArtmann/SQLC-Wizzard/internal/ui"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/schema"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/templates"
	"github.com/LarsArtmann/SQLC-Wizzard/pkg/config"
)

// UIHelper manages UI styling and display.
type UIHelper struct {
	themeFunc huh.ThemeFunc
}

// NewUIHelper creates a new UI helper.
func NewUIHelper() *UIHelper {
	return &UIHelper{
		themeFunc: huh.ThemeBase,
	}
}

// ShowStepHeader displays a step header.
func (ui *UIHelper) ShowStepHeader(title string) {
	ui.showTitledSection(title, 1, 0)
}

// ShowStepComplete displays a step completion message.
func (ui *UIHelper) ShowStepComplete(title, message string) {
	fmt.Println(uistyles.SuccessTitle.Render("✅ " + title))
	fmt.Println(uistyles.SuccessMessage.Render(message))
	fmt.Println()
}

// ShowSection displays a section header.
func (ui *UIHelper) ShowSection(title string) {
	ui.showTitledSection(title, 0, 1)
}

// createTitleStyle creates a consistent title style with configurable padding.
func (ui *UIHelper) createTitleStyle(vertical, horizontal int) lipgloss.Style {
	return uistyles.NewTitleStyle(vertical, horizontal)
}

// NewTitleStyle creates a consistent title style with configurable padding.
func (ui *UIHelper) NewTitleStyle(vertical, horizontal int) lipgloss.Style {
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
	fmt.Println(uistyles.InfoText.Render(message))
}

// showWelcome displays welcome banner.
func (ui *UIHelper) ShowWelcome() {
	fmt.Println(uistyles.HighlightBold.Render("🧙‍♂️  SQLC Configuration Wizard"))
	fmt.Println(uistyles.InfoBlock.Render("Let's create a perfect sqlc setup for your project!\n"))
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
		PaddingLeft(uistyles.ContentPaddingLeft).
		Width(uistyles.UIWidth).
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

	err := runConfirmationForm(
		ui.themeFunc,
		"Generate configuration with these settings?",
		"You can edit this later in the generated yaml file",
		&confirmed,
	)
	if err != nil {
		return false, fmt.Errorf("confirmation form failed: %w", err)
	}

	if !confirmed {
		return false, apperrors.NewError(
			apperrors.ErrorCodeValidationError,
			fmt.Sprintf("configuration cancelled by user (confirmed=%t)", confirmed),
		)
	}

	return true, nil
}

// formatConfigurationSummary formats configuration for display.
func (ui *UIHelper) formatConfigurationSummary(
	cfg *schema.Schema,
	data generated.TemplateData,
) string {
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
func (ui *UIHelper) formatCompletionDetails(
	cfg *schema.Schema,
	_ generated.TemplateData,
) string {
	detailStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#99")).
		PaddingLeft(uistyles.ContentPaddingLeft)

	details := detailStyle.Render("Generated Files:")
	details += "\n" + "- sqlc.yaml configuration"
	details += "\n" + fmt.Sprintf("- Database schema (%d tables)", len(cfg.Tables))
	details += "\n" + "- Query files (based on schema)"

	return details
}

// createErrorStyles returns consistent error styling for display.
func (ui *UIHelper) createErrorStyles() (errorStyle, detailStyle lipgloss.Style) {
	errorStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FF5555")).
		Padding(0, 1)

	detailStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF8888")).
		PaddingLeft(uistyles.ContentPaddingLeft)

	return errorStyle, detailStyle
}

// showErrorWithSchemaDetails displays schema errors.
func (ui *UIHelper) showErrorWithSchemaDetails(err *schema.SchemaError) {
	errorStyle, detailStyle := ui.createErrorStyles()

	fmt.Println(errorStyle.Render("❌ Schema Error"))
	fmt.Println(detailStyle.Render("Code: " + err.Code))
	fmt.Println(detailStyle.Render("Message: " + err.Message))
}

// showErrorWithTypedDetails displays typed errors.
func (ui *UIHelper) showErrorWithTypedDetails(err *apperrors.Error) {
	errorStyle, detailStyle := ui.createErrorStyles()

	fmt.Println(errorStyle.Render("❌ Error"))
	fmt.Println(detailStyle.Render("Code: " + string(err.Code)))
	fmt.Println(detailStyle.Render("Message: " + err.Message))

	if err.Description != "" {
		fmt.Println(detailStyle.Render("Description: " + err.Description))
	}
}
