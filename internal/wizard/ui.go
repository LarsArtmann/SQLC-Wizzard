// Package wizard provides UI components and theme management
package wizard

import (
	"fmt"
	"strings"

	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/charmbracelet/lipgloss"
)

// ShowWelcome displays welcome message
func ShowWelcome() {
	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color("207")).
		Bold(true).
		Width(80).
		Align(lipgloss.Center).
		MarginTop(1).
		MarginBottom(1)

	title := style.Render("🧙 SQLC-Wizard")
	
	subtitleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("251")).
		Width(80).
		Align(lipgloss.Center).
		MarginBottom(2)
		
	subtitle := subtitleStyle.Render("Perfect sqlc configurations in minutes, not hours")

	description := lipgloss.NewStyle().
		Foreground(lipgloss.Color("242")).
		Width(80).
		Align(lipgloss.Center).
		MarginBottom(3).
		Render("Type-safe SQL accessible to everyone through an intuitive wizard that guides you through creating production-ready sqlc setups with smart defaults and comprehensive validation.")

	fmt.Println(title)
	fmt.Println(subtitle)
	fmt.Println(description)
}

// ShowCompletion displays completion message
func ShowCompletion(config interface{}, data generated.TemplateData) {
	successStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("48")).
		Bold(true).
		Width(80).
		Align(lipgloss.Center).
		MarginTop(1).
		MarginBottom(1)

	success := successStyle.Render("✅ Configuration Complete!")

	detailsStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("251")).
		Width(80).
		MarginBottom(1)

	details := []string{
		fmt.Sprintf("Project Type: %s", strings.Title(string(data.ProjectType))),
		fmt.Sprintf("Database: %s", strings.Title(string(data.Database.Engine))),
		fmt.Sprintf("Output: %s", data.Output.BaseDir),
		fmt.Sprintf("Package: %s", data.Package.Name),
	}

	fmt.Println(success)
	fmt.Println(detailsStyle.Render(strings.Join(details, "\n")))

	nextStepsStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("242")).
		Width(80).
		Align(lipgloss.Center).
		MarginTop(1)

	nextSteps := []string{
		"📝 Write your SQL queries in the queries directory",
		"🗄️  Add your schema files to the schema directory", 
		"🚀 Run 'sqlc generate' to create type-safe Go code",
		"🧪 Import and use the generated database types",
	}

	fmt.Println(nextStepsStyle.Render(strings.Join(nextSteps, "\n")))
}

// ShowError displays error messages with styling
func ShowError(err error, step string) {
	errorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("9")).
		Bold(true).
		Width(80).
		Align(lipgloss.Center).
		MarginTop(1).
		MarginBottom(1)

	errorMsg := errorStyle.Render("❌ Error")

	detailsStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("251")).
		Width(80).
		MarginBottom(1)

	details := fmt.Sprintf("Step: %s\nError: %s", step, err.Error())

	fmt.Println(errorMsg)
	fmt.Println(detailsStyle.Render(details))
}

// CreatePromptStyle returns styled prompt configuration
func CreatePromptStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("251")).
		Width(80).
		Align(lipgloss.Left)
}

// CreateTitleStyle returns styled title configuration
func CreateTitleStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true).
		Width(80).
		Align(lipgloss.Center).
		MarginBottom(1)
}