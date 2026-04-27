package commands

import (
	"fmt"

	"github.com/LarsArtmann/SQLC-Wizzard/internal/ui"
)

// UI output helpers for consistent command output styling.

// PrintSuccess prints a success message with consistent styling.
func PrintSuccess(message string) {
	fmt.Println(ui.SuccessOutput.Render("✓ " + message))
}

// PrintSuccessf prints a formatted success message.
func PrintSuccessf(format string, args ...any) {
	PrintSuccess(fmt.Sprintf(format, args...))
}

// PrintInfo prints an info message with consistent styling.
func PrintInfo(message string) {
	fmt.Println(ui.InfoBlock.Render(message))
}

// PrintInfoWithSummary prints info message followed by a summary.
func PrintInfoWithSummary(message, summary string) {
	fmt.Println(ui.InfoBlock.Render(message))
	fmt.Println()
	PrintSuccess(summary)
}

// PrintNextSteps prints next steps with consistent styling.
func PrintNextSteps(steps []string) {
	fmt.Println(ui.NextStepsTitle.Render("Next Steps:"))

	for _, step := range steps {
		fmt.Println(ui.CommandText.Render(step))
	}
}

// PrintError prints an error message with consistent styling.
func PrintError(message string) {
	fmt.Println(ui.NewErrorStyle(0, 0).Render("✗ " + message))
}
