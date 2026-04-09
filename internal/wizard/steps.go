// Package wizard provides step definitions for the configuration wizard
package wizard

import (
	"fmt"

	"charm.land/huh/v2"
	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/apperrors"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/templates"
)

// Validation constants for wizard steps.
const (
	minProjectNameLength = 2
	maxProjectNameLength = 50
)

// createValidatedInput creates a validated input field with common error checking.
func createValidatedInput(
	title, description, placeholder, fieldName string,
	value *string,
) *huh.Input {
	var localValue string
	if value != nil {
		localValue = *value
	}

	return huh.NewInput().
		Title(title).
		Description(description).
		Value(&localValue).
		Placeholder(placeholder).
		Validate(func(val string) error {
			if val == "" {
				return apperrors.NewError(
					apperrors.ErrorCodeValidationError,
					fieldName+" cannot be empty",
				)
			}

			return nil
		})
}

// CreateProjectTypeStep creates project type selection step.
func CreateProjectTypeStep(data *generated.TemplateData) *huh.Select[string] {
	projectType := generated.ProjectTypeMicroservice
	if data != nil {
		projectType = data.ProjectType
	}

	projectTypePtr := new(string)
	*projectTypePtr = string(projectType)

	return huh.NewSelect[string]().
		Title("Select Project Type").
		Description("Choose the type of project you want to create").
		Options(
			huh.NewOption("🛠️  Hobby Project - Simple personal projects", "hobby"),
			huh.NewOption("🏢  Microservice - API service with single database", "microservice"),
			huh.NewOption("🏭  Enterprise - Complex business applications", "enterprise"),
			huh.NewOption("🚀  API-First - REST/GraphQL API development", "api-first"),
			huh.NewOption("📊  Analytics - Data warehouse and reporting", "analytics"),
			huh.NewOption("🧪  Testing - Test projects and utilities", "testing"),
			huh.NewOption("🏢  Multi-Tenant - SaaS applications", "multi-tenant"),
			huh.NewOption("📚  Library - Reusable code packages", "library"),
		).
		Value(projectTypePtr).
		Validate(func(projectType string) error {
			if !templates.IsValidProjectType(projectType) {
				return fmt.Errorf("invalid project type: %s", projectType)
			}

			return nil
		})
}

// CreateDatabaseStep creates database selection step.
func CreateDatabaseStep(data *generated.TemplateData) *huh.Select[string] {
	engine := generated.DatabaseTypePostgreSQL
	if data != nil {
		engine = data.Database.Engine
	}

	databasePtr := new(string)
	*databasePtr = string(engine)

	return huh.NewSelect[string]().
		Title("Select Database Engine").
		Description("Choose the database engine for your project").
		Options(
			huh.NewOption("🐘  PostgreSQL - Production-ready relational database", "postgresql"),
			huh.NewOption("🐬  MySQL - Popular relational database", "mysql"),
			huh.NewOption("📁  SQLite - Lightweight file-based database", "sqlite"),
		).
		Value(databasePtr).
		Validate(func(database string) error {
			if !templates.IsValidDatabaseType(database) {
				return fmt.Errorf("invalid database type: %s", database)
			}

			return nil
		})
}

// CreateProjectNameStep creates project name input step.
func CreateProjectNameStep(data *generated.TemplateData) *huh.Input {
	projectName := ""
	if data != nil {
		projectName = data.ProjectName
	}

	return huh.NewInput().
		Title("Project Name").
		Description("Enter the name for your project").
		Value(&projectName).
		Validate(func(name string) error {
			if name == "" {
				return apperrors.NewError(
					apperrors.ErrorCodeValidationError,
					"project name cannot be empty",
				)
			}

			if len(name) < minProjectNameLength {
				return apperrors.NewError(
					apperrors.ErrorCodeValidationError,
					fmt.Sprintf("project name must be at least %d characters", minProjectNameLength),
				)
			}

			if len(name) > maxProjectNameLength {
				return apperrors.NewError(
					apperrors.ErrorCodeValidationError,
					fmt.Sprintf("project name must be less than %d characters", maxProjectNameLength),
				)
			}

			return nil
		})
}

// CreatePackageNameStep creates package name input step.
func CreatePackageNameStep(data *generated.TemplateData) *huh.Input {
	var packageName string
	if data != nil {
		packageName = data.Package.Name
	}

	return createValidatedInput(
		"Package Name",
		"Enter the Go package name for generated code",
		"db",
		"package name",
		&packageName,
	)
}

// CreatePackagePathStep creates package path input step.
func CreatePackagePathStep(data *generated.TemplateData) *huh.Input {
	var packagePath string
	if data != nil {
		packagePath = data.Package.Path
	}

	return createValidatedInput(
		"Package Path",
		"Enter the Go module path for your project",
		"github.com/username/project",
		"package path",
		&packagePath,
	)
}

// CreateOutputDirStep creates output directory input step.
func CreateOutputDirStep(data *generated.TemplateData) *huh.Input {
	var outputDir string
	if data != nil {
		outputDir = data.Output.BaseDir
	}

	return createValidatedInput(
		"Output Directory",
		"Enter the directory where generated files will be placed",
		"internal/db",
		"output directory",
		&outputDir,
	)
}

// CreateDatabaseURLStep creates database URL input step.
func CreateDatabaseURLStep(data *generated.TemplateData) *huh.Input {
	engine := generated.DatabaseTypePostgreSQL
	if data != nil {
		engine = data.Database.Engine
	}

	placeholder := "postgresql://localhost:5432/dbname"
	description := "Enter the database connection URL (use environment variables in production)"

	switch engine {
	case generated.DatabaseTypePostgreSQL:
		placeholder = "postgresql://localhost:5432/dbname"
		description = "Enter the PostgreSQL connection URL (use environment variables in production)"
	case generated.DatabaseTypeSQLite:
		placeholder = "./data.db"
		description = "Enter the SQLite database file path"
	case generated.DatabaseTypeMySQL:
		placeholder = "mysql://localhost:3306/dbname"
		description = "Enter the MySQL connection URL (use environment variables in production)"
	}

	var databaseURL string
	if data != nil {
		databaseURL = data.Database.URL
	}

	return huh.NewInput().
		Title("Database URL").
		Description(description).
		Value(&databaseURL).
		Placeholder(placeholder).
		Validate(func(url string) error {
			if url == "" {
				return apperrors.NewError(
					apperrors.ErrorCodeValidationError,
					"database URL cannot be empty",
				)
			}

			return nil
		})
}

// CreateFeatureSteps creates feature selection steps.
func CreateFeatureSteps(data *generated.TemplateData) []huh.Field {
	useUUIDs := true
	useJSON := true
	useArrays := false
	useFullText := false
	strictFunctions := false
	strictOrderBy := false

	if data != nil {
		useUUIDs = data.Database.UseUUIDs
		useJSON = data.Database.UseJSON
		useArrays = data.Database.UseArrays
		useFullText = data.Database.UseFullText
		strictFunctions = data.Validation.StrictFunctions
		strictOrderBy = data.Validation.StrictOrderBy
	}

	return []huh.Field{
		huh.NewConfirm().
			Title("Use UUIDs").
			Description("Generate UUID support for primary keys").
			Value(&useUUIDs),

		huh.NewConfirm().
			Title("Use JSON").
			Description("Generate JSON support for flexible columns").
			Value(&useJSON),

		huh.NewConfirm().
			Title("Use Arrays").
			Description("Generate array support (PostgreSQL specific)").
			Value(&useArrays),

		huh.NewConfirm().
			Title("Use Full-Text Search").
			Description("Generate full-text search support (PostgreSQL specific)").
			Value(&useFullText),

		huh.NewConfirm().
			Title("Strict Function Checks").
			Description("Enable strict function validation in SQL").
			Value(&strictFunctions),

		huh.NewConfirm().
			Title("Strict Order By").
			Description("Enable strict ORDER BY validation").
			Value(&strictOrderBy),
	}
}
