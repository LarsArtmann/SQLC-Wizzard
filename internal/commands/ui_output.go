package commands

import (
	"fmt"

	"charm.land/lipgloss/v2"
)

const (
	paddingLeft = 2
)

// UI output helpers for consistent command output styling.
var (
	successStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("10")).
			Padding(1, 0)

	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Padding(0, 0, 1, 0)

	nextStepsStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("14")).
			Bold(true).
			MarginTop(1)

	commandStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("8")).
			PaddingLeft(paddingLeft)
)

// PrintSuccess prints a success message with consistent styling.
func PrintSuccess(message string) {
	fmt.Println(successStyle.Render("✓ " + message))
}

// PrintSuccessf prints a formatted success message.
func PrintSuccessf(format string, args ...interface{}) {
	PrintSuccess(fmt.Sprintf(format, args...))
}

// PrintInfo prints an info message with consistent styling.
func PrintInfo(message string) {
	fmt.Println(infoStyle.Render(message))
}

// PrintInfoWithSummary prints info message followed by a summary.
func PrintInfoWithSummary(message, summary string) {
	fmt.Println(infoStyle.Render(message))
	fmt.Println()
	PrintSuccess(summary)
}

// PrintNextSteps prints next steps with consistent styling.
func PrintNextSteps(steps []string) {
	fmt.Println(nextStepsStyle.Render("Next Steps:"))
	for _, step := range steps {
		fmt.Println(commandStyle.Render(step))
	}
}

// PrintError prints an error message with consistent styling.
func PrintError(message string) {
	errorStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("9")).
		Padding(0, 0, 1, 0)

	fmt.Println(errorStyle.Render("✗ " + message))
}
