// Package wizard provides step definitions for the configuration wizard
package wizard

import (
	"fmt"

	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/apperrors"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/templates"
	"github.com/charmbracelet/huh"
)

// CreateProjectTypeStep creates project type selection step.
func CreateProjectTypeStep(data *generated.TemplateData) *huh.Select[string] {
	projectTypePtr := new(string)
	*projectTypePtr = string(data.ProjectType)

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
	databasePtr := new(string)
	*databasePtr = string(data.Database.Engine)

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
	return huh.NewInput().
		Title("Project Name").
		Description("Enter the name for your project").
		Value(&data.ProjectName).
		Validate(func(name string) error {
			if name == "" {
				return apperrors.NewError(apperrors.ErrorCodeValidationError, "project name cannot be empty")
			}
			if len(name) < 2 {
				return apperrors.NewError(apperrors.ErrorCodeValidationError, "project name must be at least 2 characters")
			}
			if len(name) > 50 {
				return apperrors.NewError(apperrors.ErrorCodeValidationError, "project name must be less than 50 characters")
			}
			// Simple validation for now
			return nil
		})
}

// CreatePackageNameStep creates package name input step.
func CreatePackageNameStep(data *generated.TemplateData) *huh.Input {
	return huh.NewInput().
		Title("Package Name").
		Description("Enter the Go package name for generated code").
		Value(&data.Package.Name).
		Placeholder("db").
		Validate(func(name string) error {
			if name == "" {
				return apperrors.NewError(apperrors.ErrorCodeValidationError, "package name cannot be empty")
			}
			// Simple validation for now
			return nil
		})
}

// CreatePackagePathStep creates package path input step.
func CreatePackagePathStep(data *generated.TemplateData) *huh.Input {
	return huh.NewInput().
		Title("Package Path").
		Description("Enter the Go module path for your project").
		Value(&data.Package.Path).
		Placeholder("github.com/username/project").
		Validate(func(path string) error {
			if path == "" {
				return apperrors.NewError(apperrors.ErrorCodeValidationError, "package path cannot be empty")
			}
			// Simple validation for now
			return nil
		})
}

// CreateOutputDirStep creates output directory input step.
func CreateOutputDirStep(data *generated.TemplateData) *huh.Input {
	return huh.NewInput().
		Title("Output Directory").
		Description("Enter the directory where generated files will be placed").
		Value(&data.Output.BaseDir).
		Placeholder("internal/db").
		Validate(func(dir string) error {
			if dir == "" {
				return apperrors.NewError(apperrors.ErrorCodeValidationError, "output directory cannot be empty")
			}
			return nil
		})
}

// CreateDatabaseURLStep creates database URL input step.
func CreateDatabaseURLStep(data *generated.TemplateData) *huh.Input {
	placeholder := "postgresql://localhost:5432/dbname"
	description := "Enter the database connection URL (use environment variables in production)"

	switch data.Database.Engine {
	case "sqlite":
		placeholder = "./data.db"
		description = "Enter the SQLite database file path"
	case "mysql":
		placeholder = "mysql://localhost:3306/dbname"
		description = "Enter the MySQL connection URL (use environment variables in production)"
	}

	return huh.NewInput().
		Title("Database URL").
		Description(description).
		Value(&data.Database.URL).
		Placeholder(placeholder).
		Validate(func(url string) error {
			if url == "" {
				return apperrors.NewError(apperrors.ErrorCodeValidationError, "database URL cannot be empty")
			}
			return nil
		})
}

// CreateFeatureSteps creates feature selection steps.
func CreateFeatureSteps(data *generated.TemplateData) []huh.Field {
	return []huh.Field{
		huh.NewConfirm().
			Title("Use UUIDs").
			Description("Generate UUID support for primary keys").
			Value(&data.Database.UseUUIDs),

		huh.NewConfirm().
			Title("Use JSON").
			Description("Generate JSON support for flexible columns").
			Value(&data.Database.UseJSON),

		huh.NewConfirm().
			Title("Use Arrays").
			Description("Generate array support (PostgreSQL specific)").
			Value(&data.Database.UseArrays),

		huh.NewConfirm().
			Title("Use Full-Text Search").
			Description("Generate full-text search support (PostgreSQL specific)").
			Value(&data.Database.UseFullText),

		huh.NewConfirm().
			Title("Strict Function Checks").
			Description("Enable strict function validation in SQL").
			Value(&data.Validation.StrictFunctions),

		huh.NewConfirm().
			Title("Strict Order By").
			Description("Enable strict ORDER BY validation").
			Value(&data.Validation.StrictOrderBy),
	}
}
