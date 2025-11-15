package wizard

import (
	"github.com/LarsArtmann/SQLC-Wizzard/internal/schema"
	"github.com/LarsArtmann/SQLC-Wizzard/generated"
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
		return &schema.SchemaError{
			Code:    "NULL_SCHEMA",
			Message: "Schema cannot be null",
		}
	}
	
	// Validate schema using typed validation
	if err := cfg.Validate(); err != nil {
		return err
	}
	
	// Additional business logic validation
	if len(cfg.Tables) > 100 {
		return &schema.SchemaError{
			Code:    "TOO_MANY_TABLES",
			Message: "Schema exceeds maximum allowed tables",
		}
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