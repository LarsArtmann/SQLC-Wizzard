package wizard

import (
	"fmt"
	"path/filepath"

	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/charmbracelet/huh"
)

// OutputStep handles output configuration and file paths
type OutputStep struct {
	theme *huh.Theme
	ui    *UIHelper
}

// NewOutputStep creates a new output step
func NewOutputStep(theme *huh.Theme, ui *UIHelper) *OutputStep {
	return &OutputStep{
		theme: theme,
		ui:    ui,
	}
}

// Execute runs output configuration step
func (s *OutputStep) Execute(data *generated.TemplateData) error {
	s.ui.ShowStepHeader("Output Configuration")

	// Output directory
	if err := s.configureOutputDirectory(data); err != nil {
		return err
	}

	// SQL file configuration
	if err := s.configureSQLPaths(data); err != nil {
		return err
	}

	// Generated code configuration
	if err := s.configureGeneratedPaths(data); err != nil {
		return err
	}

	s.ui.ShowStepComplete("Output Configuration", "All file paths configured")
	return nil
}

// configureOutputDirectory configures the main output directory
func (s *OutputStep) configureOutputDirectory(data *generated.TemplateData) error {
	var outputDir string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Where should generated code be output?").
				Description("Relative to project root").
				Placeholder("./internal/db").
				Value(&outputDir).
				Validate(func(str string) error {
					if str == "" {
						return fmt.Errorf("output directory cannot be empty")
					}
					return nil
				}),
		),
	).WithTheme(s.theme)

	if err := form.Run(); err != nil {
		return fmt.Errorf("output directory configuration failed: %w", err)
	}

	// Clean and validate the path
	outputDir = filepath.Clean(outputDir)
	if outputDir == "." {
		outputDir = "./internal/db"
	}

	data.Output.BaseDir = outputDir
	return nil
}

// configureSQLPaths configures SQL file paths
func (s *OutputStep) configureSQLPaths(data *generated.TemplateData) error {
	var (
		schemaDir string
		queriesDir string
	)

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Where are your SQL schema files?").
				Description("Directory containing CREATE TABLE statements").
				Placeholder("./sql/schema").
				Value(&schemaDir),
			huh.NewInput().
				Title("Where are your SQL query files?").
				Description("Directory containing SELECT, INSERT, UPDATE, DELETE queries").
				Placeholder("./sql/queries").
				Value(&queriesDir),
		),
	).WithTheme(s.theme)

	if err := form.Run(); err != nil {
		return fmt.Errorf("SQL paths configuration failed: %w", err)
	}

	// Set defaults if empty
	if schemaDir == "" {
		schemaDir = "./sql/schema"
	}
	if queriesDir == "" {
		queriesDir = "./sql/queries"
	}

	// Clean paths
	schemaDir = filepath.Clean(schemaDir)
	queriesDir = filepath.Clean(queriesDir)

	data.Output.SchemaDir = schemaDir
	data.Output.QueriesDir = queriesDir

	return nil
}

// configureGeneratedPaths configures generated code output paths
func (s *OutputStep) configureGeneratedPaths(data *generated.TemplateData) error {
	var (
		goOutput string
		emitDB   bool
		emitModels bool
		emitQueries bool
	)

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Go package output directory?").
				Description("Where to generate Go code (relative to main output dir)").
				Placeholder("./gen").
				Value(&goOutput),
		),
		huh.NewGroup(
			huh.NewConfirm().
				Title("Generate database code?").
				Description("Generate database connection and migration code").
				Value(&emitDB),
			huh.NewConfirm().
				Title("Generate models?").
				Description("Generate Go struct models from your schema").
				Value(&emitModels),
			huh.NewConfirm().
				Title("Generate query code?").
				Description("Generate Go query functions from your SQL").
				Value(&emitQueries),
		),
	).WithTheme(s.theme)

	if err := form.Run(); err != nil {
		return fmt.Errorf("generated paths configuration failed: %w", err)
	}

	// Set default if empty
	if goOutput == "" {
		goOutput = "./gen"
	}
	goOutput = filepath.Clean(goOutput)

	// Configure output options
	data.Output.Go.Package = goOutput
	data.Output.Go.EmitDB = emitDB
	data.Output.Go.EmitModels = emitModels
	data.Output.Go.EmitQueries = emitQueries

	return nil
}

// ValidateConfiguration validates the complete output configuration
func (s *OutputStep) ValidateConfiguration(data *generated.TemplateData) error {
	// Check for path conflicts
	if data.Output.Directory == "./sql/schema" || data.Output.Directory == "./sql/queries" {
		return fmt.Errorf("output directory cannot be the same as SQL directories")
	}

	// Check for reasonable directory structure
	if len(data.Output.SQL.Schema) == 0 {
		return fmt.Errorf("at least one schema directory is required")
	}

	if len(data.Output.SQL.Queries) == 0 {
		return fmt.Errorf("at least one queries directory is required")
	}

	// Validate Go output path
	if data.Output.Go.Package == "" {
		return fmt.Errorf("Go output package path cannot be empty")
	}

	return nil
}