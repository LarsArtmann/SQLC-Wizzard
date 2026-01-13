package config

import (
	"fmt"
	"slices"
	"strings"
)

// ValidationError represents a configuration validation error.
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// ValidationResult contains validation errors and warnings.
type ValidationResult struct {
	Errors   []ValidationError
	Warnings []ValidationError
}

// IsValid returns true if there are no errors.
func (r *ValidationResult) IsValid() bool {
	return len(r.Errors) == 0
}

// AddError adds a validation error.
func (r *ValidationResult) AddError(field, message string) {
	r.Errors = append(r.Errors, ValidationError{Field: field, Message: message})
}

// AddWarning adds a validation warning.
func (r *ValidationResult) AddWarning(field, message string) {
	r.Warnings = append(r.Warnings, ValidationError{Field: field, Message: message})
}

// Validate performs comprehensive validation on a SqlcConfig.
func Validate(cfg *SqlcConfig) *ValidationResult {
	result := &ValidationResult{}

	// Handle nil configuration
	if cfg == nil {
		result.AddError("config", "configuration cannot be nil")
		return result
	}

	// Validate version
	if cfg.Version == "" {
		result.AddError("version", "version is required")
	} else if cfg.Version != "1" && cfg.Version != "2" {
		result.AddError("version", fmt.Sprintf("unsupported version: %s (expected 1 or 2)", cfg.Version))
	}

	// Validate SQL configurations
	if len(cfg.SQL) == 0 {
		result.AddError("sql", "at least one SQL configuration is required")
		return result // Early return if no SQL configs
	}

	for i, sqlCfg := range cfg.SQL {
		validateSQLConfig(&sqlCfg, i, result)
	}

	return result
}

func validateSQLConfig(cfg *SQLConfig, index int, result *ValidationResult) {
	prefix := fmt.Sprintf("sql[%d]", index)

	// Validate engine
	validEngines := []string{"postgresql", "mysql", "sqlite"}
	if cfg.Engine == "" {
		result.AddError(prefix+".engine", "engine is required")
	} else if !slices.Contains(validEngines, cfg.Engine) {
		result.AddError(prefix+".engine", fmt.Sprintf("invalid engine: %s (must be one of: %s)", cfg.Engine, strings.Join(validEngines, ", ")))
	}

	// Validate queries path
	if cfg.Queries.IsEmpty() {
		result.AddError(prefix+".queries", "queries path is required")
	}

	// Validate schema path
	if cfg.Schema.IsEmpty() {
		result.AddError(prefix+".schema", "schema path is required")
	}

	// Validate gen configuration
	if cfg.Gen.Go == nil && cfg.Gen.Kotlin == nil && cfg.Gen.Python == nil && cfg.Gen.TypeScript == nil {
		result.AddError(prefix+".gen", "at least one language generation config is required")
	}

	// Validate Go gen config if present
	if cfg.Gen.Go != nil {
		validateGoGenConfig(cfg.Gen.Go, prefix+".gen.go", result)
	}
}

func validateGoGenConfig(cfg *GoGenConfig, prefix string, result *ValidationResult) {
	// Validate required fields
	if cfg.Package == "" {
		result.AddError(prefix+".package", "package name is required")
	}

	if cfg.Out == "" {
		result.AddError(prefix+".out", "output directory is required")
	}

	// Validate json_tags_case_style if set
	if cfg.JSONTagsCaseStyle != "" {
		validStyles := []string{"camel", "pascal", "snake"}
		if !slices.Contains(validStyles, cfg.JSONTagsCaseStyle) {
			result.AddError(prefix+".json_tags_case_style", fmt.Sprintf("invalid case style: %s (must be one of: %s)", cfg.JSONTagsCaseStyle, strings.Join(validStyles, ", ")))
		}
	}

	// Add warnings for best practices
	if !cfg.EmitInterface {
		result.AddWarning(prefix+".emit_interface", "consider enabling emit_interface for better testability")
	}

	if !cfg.EmitPreparedQueries {
		result.AddWarning(prefix+".emit_prepared_queries", "consider enabling emit_prepared_queries for better performance")
	}

	if !cfg.EmitJSONTags {
		result.AddWarning(prefix+".emit_json_tags", "consider enabling emit_json_tags for JSON serialization support")
	}
}
