package wizard

import (
	"fmt"
	"strings"

	"charm.land/huh/v2"
	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/apperrors"
)

// Package name validation constants.
const (
	minPackageNameLength = 2
	maxPackageNameLength = 50
)

// ASCII case conversion constants.
const (
	// lowercaseOffsetASCII is the difference between uppercase and lowercase ASCII letters.
	lowercaseOffsetASCII = 32
)

// ProjectDetailsStep handles project name and configuration details.
type ProjectDetailsStep struct {
	themeFunc huh.ThemeFunc
	ui        *UIHelper
}

// NewProjectDetailsStep creates a new project details step.
func NewProjectDetailsStep(themeFunc huh.ThemeFunc, ui *UIHelper) *ProjectDetailsStep {
	return &ProjectDetailsStep{
		themeFunc: themeFunc,
		ui:        ui,
	}
}

// Execute runs project details configuration step.
func (s *ProjectDetailsStep) Execute(data *generated.TemplateData) error {
	s.ui.ShowStepHeader("Project Details")

	// Project name
	var projectName string

	nameForm := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("What is your project name?").
				Placeholder("my-awesome-project").
				Value(&projectName).
				Validate(func(str string) error {
					if len(str) < MinProjectNameLength {
						return apperrors.NewError(
							apperrors.ErrorCodeValidationError,
							fmt.Sprintf(
								"project name must be at least %d characters (got %d)",
								MinProjectNameLength,
								len(str),
							),
						)
					}

					if len(str) > MaxProjectNameLength {
						return apperrors.NewError(
							apperrors.ErrorCodeValidationError,
							fmt.Sprintf(
								"project name must be less than %d characters (got %d)",
								MaxProjectNameLength,
								len(str),
							),
						)
					}

					return nil
				}),
		),
	).WithTheme(s.themeFunc)

	err := nameForm.Run()
	if err != nil {
		return fmt.Errorf("project name input failed for %q: %w", projectName, err)
	}

	data.ProjectName = projectName

	// Package name (auto-generated from project name)
	var packageName string

	packageForm := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("What is your Go package name?").
				Placeholder(s.generatePackageName(projectName)).
				Value(&packageName).
				Validate(func(str string) error {
					if len(str) < minPackageNameLength {
						return apperrors.NewError(
							apperrors.ErrorCodeValidationError,
							fmt.Sprintf(
								"package name must be at least %d characters (got %d)",
								minPackageNameLength,
								len(str),
							),
						)
					}

					return nil
				}),
		),
	).WithTheme(s.themeFunc)

	err = packageForm.Run()
	if err != nil {
		return fmt.Errorf(
			"package name input failed for %q (projectName=%q): %w",
			packageName,
			projectName,
			err,
		)
	}

	if packageName == "" {
		packageName = s.generatePackageName(projectName)
	}

	data.Package.Name = packageName
	data.Package.Path = "github.com/yourorg/" + packageName

	s.ui.ShowStepComplete(
		"Project Details",
		fmt.Sprintf("Name: %s, Package: %s", data.ProjectName, data.Package.Name),
	)

	return nil
}

// generatePackageName converts project name to valid Go package name.
func (s *ProjectDetailsStep) generatePackageName(projectName string) string {
	// Simple conversion: replace spaces and hyphens with underscores, remove invalid characters
	packageName := projectName
	packageName = s.replaceInvalidChars(packageName)
	packageName = s.lowerCaseFirst(packageName)

	// Ensure it's a valid Go identifier
	if s.isGoKeyword(packageName) {
		packageName += "pkg"
	}

	return packageName
}

// replaceInvalidChars replaces characters invalid in Go package names.
func (s *ProjectDetailsStep) replaceInvalidChars(input string) string {
	var result strings.Builder

	for _, char := range input {
		if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || char == '_' ||
			(char >= '0' && char <= '9') {
			result.WriteRune(char)
		} else {
			result.WriteString("_")
		}
	}

	return result.String()
}

// lowerCaseFirst converts first character to lowercase if it's uppercase.
func (s *ProjectDetailsStep) lowerCaseFirst(input string) string {
	if len(input) == 0 {
		return input
	}

	if input[0] >= 'A' && input[0] <= 'Z' {
		return string(input[0]+lowercaseOffsetASCII) + input[1:]
	}

	return input
}

// isGoKeyword checks if string is a Go keyword.
func (s *ProjectDetailsStep) isGoKeyword(keyword string) bool {
	keywords := map[string]bool{
		"break": true, "case": true, "chan": true, "const": true, "continue": true,
		"default": true, "defer": true, "else": true, "fallthrough": true, "for": true,
		"func": true, "go": true, "goto": true, "if": true, "import": true,
		"interface": true, "map": true, "package": true, "range": true, "return": true,
		"select": true, "struct": true, "switch": true, "type": true, "var": true,
	}

	return keywords[keyword]
}
