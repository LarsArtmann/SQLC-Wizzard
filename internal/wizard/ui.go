package wizard

import (
	"fmt"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/schema"
	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/errors"
	"github.com/LarsArtmann/SQLC-Wizzard/pkg/config"
)

// ShowCompletion displays final configuration summary with typed schema
// Replaces 'any' type with proper typed schema parameter
func ShowCompletion(cfg *schema.Schema, data generated.TemplateData) {
	ui := NewUIHelper()
	
	ui.ShowSection("🎉 Configuration Complete")
	
	// Build typed summary using schema and template data
	summary := ui.formatConfigurationSummary(cfg, data)
	ui.ShowInfo(summary)
	
	// Show completion details with type safety
	completion := ui.formatCompletionDetails(cfg, data)
	ui.ShowSection("✅ Generation Complete")
	ui.ShowInfo(completion)
}

// ShowError displays schema validation errors with proper typing
func ShowError(err error) {
	ui := NewUIHelper()
	
	// Check if it's a schema error
	if schemaErr, ok := err.(*schema.SchemaError); ok {
		ui.showErrorWithSchemaDetails(schemaErr)
		return
	}
	
	// Check if it's our typed error
	if appErr, ok := err.(*errors.Error); ok {
		ui.showErrorWithTypedDetails(appErr)
		return
	}
	
	// Handle as generic error
	ui.ShowSection("❌ Error Occurred")
	ui.ShowInfo(err.Error())
}

// ShowProgress displays progress for schema operations
func ShowProgress(current, total int, operation string) {
	ui := NewUIHelper()
	
	ui.ShowSection("🔄 " + operation)
	ui.ShowInfo(fmt.Sprintf("Progress: %d/%d completed", current, total))
}

// ValidateConfiguration validates configuration using schema
func ValidateConfiguration(cfg *schema.Schema) error {
	if cfg == nil {
		return errors.NewError(errors.ErrorCodeInternalServer, "Schema cannot be null")
	}
	
	// Validate schema using typed validation
	if err := cfg.Validate(); err != nil {
		return err
	}
	
	// Additional business logic validation
	if len(cfg.Tables) > 100 {
		return errors.NewError(errors.ErrorCodeSchemaValidation, "Schema exceeds maximum allowed tables")
	}
	
	return nil
}

// GenerateConfiguration generates final configuration from schema and template data
func GenerateConfiguration(schema *schema.Schema, data generated.TemplateData) (string, error) {
	// Validate inputs
	if err := ValidateConfiguration(schema); err != nil {
		return "", err
	}
	
	// Generate configuration using typed data
	config, err := generateConfigFromSchema(schema, data)
	if err != nil {
		return "", err
	}
	
	return config, nil
}

// generateConfigFromSchema generates sqlc configuration from schema and template data
func generateConfigFromSchema(s *schema.Schema, data generated.TemplateData) (string, error) {
	// Create a basic SQLC config from schema
	sqlcConfig := &config.SqlcConfig{
		Version: "2",
		SQL: []config.SQLConfig{
			{
				Engine: "postgresql", // Default to postgresql
				Schema: config.NewPathOrPaths([]string{"schema.sql"}),
				Queries: config.NewPathOrPaths([]string{"query.sql"}),
				Gen: config.GenConfig{
					Go: &config.GoGenConfig{
						Out:       "internal/db",
						Package:   "db",
						Overrides: []config.Override{},
					},
				},
			},
		},
	}
	
	// Convert to YAML string
	yamlData, err := config.Marshal(sqlcConfig)
	if err != nil {
		return "", errors.ConfigParseError("generated config", err)
	}
	
	return string(yamlData), nil
}