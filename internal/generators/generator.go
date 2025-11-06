package generators

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/LarsArtmann/SQLC-Wizzard/internal/templates"
	"github.com/LarsArtmann/SQLC-Wizzard/pkg/config"
)

// Generator handles file generation
type Generator struct {
	outputDir string
}

// NewGenerator creates a new generator
func NewGenerator(outputDir string) *Generator {
	return &Generator{
		outputDir: outputDir,
	}
}

// GenerateAll generates all files (config, queries, schema)
func (g *Generator) GenerateAll(cfg *config.SqlcConfig, data templates.TemplateData, includeQueries, includeSchema bool) error {
	// Generate sqlc.yaml
	if err := g.GenerateSqlcConfig(cfg); err != nil {
		return fmt.Errorf("failed to generate sqlc.yaml: %w", err)
	}

	// Generate example queries if requested
	if includeQueries {
		if err := g.GenerateExampleQueries(data); err != nil {
			return fmt.Errorf("failed to generate queries: %w", err)
		}
	}

	// Generate example schema if requested
	if includeSchema {
		if err := g.GenerateExampleSchema(data); err != nil {
			return fmt.Errorf("failed to generate schema: %w", err)
		}
	}

	return nil
}

// GenerateSqlcConfig writes the sqlc.yaml file
func (g *Generator) GenerateSqlcConfig(cfg *config.SqlcConfig) error {
	path := filepath.Join(g.outputDir, "sqlc.yaml")

	// Ensure directory exists
	if err := os.MkdirAll(g.outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Write config file
	if err := config.WriteFileFormatted(cfg, path); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

// generateFileWithTemplate is a helper to generate files with templates
func (g *Generator) generateFileWithTemplate(data templates.TemplateData, dirKey, defaultDir, templateType, filename string) error {
	// Determine directory
	var dir string
	var templateContent func(templates.DatabaseType) string
	
	switch templateType {
	case "queries":
		dir = data.QueriesDir
		if dir == "" {
			dir = defaultDir
		}
		templateContent = getQueryTemplate
	case "schema":
		dir = data.SchemaDir
		if dir == "" {
			dir = defaultDir
		}
		templateContent = getSchemaTemplate
	default:
		return fmt.Errorf("unsupported template type: %s", templateType)
	}

	// Make it absolute path if relative
	if !filepath.IsAbs(dir) {
		dir = filepath.Join(g.outputDir, dir)
	}

	// Ensure directory exists
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create %s directory: %w", templateType, err)
	}

	// Get template content based on database type
	content := templateContent(data.Database)
	if content == "" {
		return fmt.Errorf("no %s template for database: %s", templateType, data.Database)
	}

	// Write to output
	outputPath := filepath.Join(dir, filename)
	if err := os.WriteFile(outputPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write %s file: %w", templateType, err)
	}

	return nil
}

// GenerateExampleQueries copies example query files
func (g *Generator) GenerateExampleQueries(data templates.TemplateData) error {
	return g.generateFileWithTemplate(data, "queries", "internal/db/queries", "queries", "users.sql")
}

// GenerateExampleSchema copies example schema files
func (g *Generator) GenerateExampleSchema(data templates.TemplateData) error {
	return g.generateFileWithTemplate(data, "schema", "internal/db/schema", "schema", "001_users_table.sql")
}

// GenerateSummary creates a summary of what was generated
func (g *Generator) GenerateSummary(cfg *config.SqlcConfig, includeQueries, includeSchema bool) string {
	summary := "✓ Generated files:\n"
	summary += fmt.Sprintf("  • sqlc.yaml (%d SQL configuration(s))\n", len(cfg.SQL))

	if includeQueries {
		summary += "  • Example queries (CRUD operations)\n"
	}

	if includeSchema {
		summary += "  • Example schema (users table)\n"
	}

	summary += "\nNext steps:\n"
	summary += "  1. Review and customize sqlc.yaml\n"
	summary += "  2. Add your schema files\n"
	summary += "  3. Write SQL queries\n"
	summary += "  4. Run: sqlc generate\n"

	return summary
}
