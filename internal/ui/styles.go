// Package ui provides shared UI styling and helpers.
package ui

import "charm.land/lipgloss/v2"

const (
	// ContentPaddingLeft is the default left padding for content.
	ContentPaddingLeft = 2
	// UIWidth is the default width for UI elements.
	UIWidth = 76
)

// Shared lipgloss styles for consistent UI across the application.
var (
	// Success styles.
	SuccessTitle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#00D084")).
			Padding(1, 0)

	SuccessMessage = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00D084")).
			PaddingLeft(ContentPaddingLeft)

		// Title styles.
	TitlePrimary = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7D56F4")).
			Padding(1, 0).
			MarginBottom(1)

	TitleSection = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7D56F4")).
			MarginBottom(1)

		// Info styles.
	InfoText = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			PaddingLeft(ContentPaddingLeft).
			Width(UIWidth).
			Align(lipgloss.Left)

	InfoBlock = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Padding(0, 0, 1, 0)

		// Error styles.
	ErrorTitle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FF5555")).
			Padding(0, 1)

	ErrorMessage = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF8888")).
			PaddingLeft(ContentPaddingLeft)

		// Command output styles.
	SuccessOutput = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("10")).
			Padding(1, 0)

	NextStepsTitle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("14")).
			Bold(true).
			MarginTop(1)

	CommandText = lipgloss.NewStyle().
			Foreground(lipgloss.Color("8")).
			PaddingLeft(ContentPaddingLeft)

		// Content styles.
	ContentText = lipgloss.NewStyle().
			PaddingLeft(ContentPaddingLeft).
			Width(UIWidth).
			Align(lipgloss.Left)

		// Highlight styles.
	HighlightAccent = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF7E67"))

	HighlightDim = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#99"))

	HighlightBold = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("99")).
			Padding(1, 0)
)

// NewTitleStyle creates a title style with configurable padding.
func NewTitleStyle(vertical, horizontal int) lipgloss.Style {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7D56F4")).
		Padding(vertical, horizontal).
		MarginBottom(1)
}

// NewErrorStyle creates an error style with configurable padding.
func NewErrorStyle(vertical, horizontal int) lipgloss.Style {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FF5555")).
		Padding(vertical, horizontal)
}
