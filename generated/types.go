// Package generated contains type-safe enums and models for the SQLC-Wizard
// Generated types ensure compile-time safety and prevent invalid states

package generated

import "time"

// ProjectType represents the type of project template
type ProjectType string

const (
	ProjectTypeHobby        ProjectType = "hobby"
	ProjectTypeMicroservice ProjectType = "microservice"
	ProjectTypeEnterprise   ProjectType = "enterprise"
	ProjectTypeAPIFirst     ProjectType = "api-first"
	ProjectTypeAnalytics    ProjectType = "analytics"
	ProjectTypeTesting      ProjectType = "testing"
	ProjectTypeMultiTenant  ProjectType = "multi-tenant"
	ProjectTypeLibrary      ProjectType = "library"
)

// IsValid returns true if ProjectType is valid
// This prevents invalid states at runtime
func (p ProjectType) IsValid() bool {
	switch p {
	case ProjectTypeHobby, ProjectTypeMicroservice, ProjectTypeEnterprise,
		ProjectTypeAPIFirst, ProjectTypeAnalytics, ProjectTypeTesting,
		ProjectTypeMultiTenant, ProjectTypeLibrary:
		return true
	default:
		return false
	}
}

// DatabaseType represents the supported database engines
type DatabaseType string

const (
	DatabaseTypePostgreSQL DatabaseType = "postgresql"
	DatabaseTypeMySQL      DatabaseType = "mysql"
	DatabaseTypeSQLite     DatabaseType = "sqlite"
)

// IsValid returns true if DatabaseType is valid
// This prevents invalid states at runtime
func (d DatabaseType) IsValid() bool {
	switch d {
	case DatabaseTypePostgreSQL, DatabaseTypeMySQL, DatabaseTypeSQLite:
		return true
	default:
		return false
	}
}

// PackageConfig represents complete package configuration
// TypeSpec: model PackageConfig { ... }
type PackageConfig struct {
	Name      string `json:"name"`
	Path      string `json:"path"`
	BuildTags string `json:"build_tags,omitempty"`
}

// DatabaseConfig represents database-specific configuration
// TypeSpec: model DatabaseConfig { ... }
type DatabaseConfig struct {
	Engine      DatabaseType `json:"engine"`
	URL         string       `json:"url,omitempty"`
	UseManaged  bool         `json:"use_managed"`
	UseUUIDs    bool         `json:"use_uuids"`
	UseJSON     bool         `json:"use_json"`
	UseArrays   bool         `json:"use_arrays"`
	UseFullText bool         `json:"use_full_text"`
}

// OutputConfig represents output directory configuration
type OutputConfig struct {
	BaseDir    string `json:"base_dir"`
	QueriesDir string `json:"queries_dir"`
	SchemaDir  string `json:"schema_dir"`
}

// ValidationConfig represents validation settings
// TypeSpec: model ValidationConfig { ... }
type ValidationConfig struct {
	StrictFunctions bool        `json:"strict_functions"`
	StrictOrderBy   bool        `json:"strict_order_by"`
	EmitOptions     EmitOptions `json:"emit_options"`
	SafetyRules     SafetyRules `json:"safety_rules"`
}

// EmitOptions defines SQL code generation options
type EmitOptions struct {
	EmitJSONTags             bool   `json:"emit_json_tags"`
	EmitPreparedQueries      bool   `json:"emit_prepared_queries"`
	EmitInterface            bool   `json:"emit_interface"`
	EmitEmptySlices          bool   `json:"emit_empty_slices"`
	EmitResultStructPointers bool   `json:"emit_result_struct_pointers"`
	EmitParamsStructPointers bool   `json:"emit_params_struct_pointers"`
	EmitEnumValidMethod      bool   `json:"emit_enum_valid_method"`
	EmitAllEnumValues        bool   `json:"emit_all_enum_values"`
	JSONTagsCaseStyle        string `json:"json_tags_case_style"`
}

// SafetyRule represents a CEL-based validation rule
// TypeSpec: model SafetyRule { ... }
type SafetyRule struct {
	Name    string `json:"name"`
	Rule    string `json:"rule"`
	Message string `json:"message"`
}

// SafetyRules represents validation rules configuration
// TypeSpec: model SafetyRules { ... }
type SafetyRules struct {
	NoSelectStar bool         `json:"no_select_star"`
	RequireWhere bool         `json:"require_where"`
	NoDropTable  bool         `json:"no_drop_table"`
	NoTruncate   bool         `json:"no_truncate"`
	RequireLimit bool         `json:"require_limit"`
	Rules        []SafetyRule `json:"rules"`
}

// NOTE: ToRuleConfigs method removed - use internal/validation/rule_transformer.TransformSafetyRules instead
// This eliminates the split brain and provides a single source of truth for rule transformation

// RuleConfig represents a validation rule configuration
// TypeSpec: model RuleConfig { ... }
type RuleConfig struct {
	Name    string `json:"name"`
	Rule    string `json:"rule"`
	Message string `json:"message"`
}

// DefaultEmitOptions returns safe defaults for code generation
func DefaultEmitOptions() EmitOptions {
	return EmitOptions{
		EmitJSONTags:             true,
		EmitPreparedQueries:      true,
		EmitInterface:            true,
		EmitEmptySlices:          true,
		EmitResultStructPointers: false,
		EmitParamsStructPointers: false,
		EmitEnumValidMethod:      true,
		EmitAllEnumValues:        true,
		JSONTagsCaseStyle:        "camel",
	}
}

// DefaultSafetyRules returns safe defaults for query validation
func DefaultSafetyRules() SafetyRules {
	return SafetyRules{
		NoSelectStar: true,
		RequireWhere: true,
		NoDropTable:  true,
		NoTruncate:   true,
		RequireLimit: false,
		Rules:        []SafetyRule{},
	}
}

// DefaultTemplateData returns a TemplateData struct with default values
func DefaultTemplateData() TemplateData {
	return TemplateData{
		Package: PackageConfig{
			Name: "myproject",
			Path: "github.com/myorg/myproject",
		},
		Database: DatabaseConfig{
			UseUUIDs:    true,
			UseJSON:     true,
			UseArrays:   false,
			UseFullText: false,
		},
		Output: OutputConfig{
			BaseDir:    "./internal/db",
			QueriesDir: "./sql/queries",
			SchemaDir:  "./sql/schema",
		},
		Validation: ValidationConfig{
			EmitOptions: DefaultEmitOptions(),
			SafetyRules: DefaultSafetyRules(),
		},
	}
}

// TemplateData represents the complete data structure for template generation
type TemplateData struct {
	ProjectName string      `json:"project_name"`
	ProjectType ProjectType `json:"project_type"`

	Package    PackageConfig    `json:"package"`
	Database   DatabaseConfig   `json:"database"`
	Output     OutputConfig     `json:"output"`
	Validation ValidationConfig `json:"validation"`
}

// CreateProjectCommand represents a command to create a new project
type CreateProjectCommand struct {
	Name        string       `json:"name"`
	ProjectType ProjectType  `json:"project_type"`
	Database    DatabaseType `json:"database"`
	OutputDir   string       `json:"output_dir"`
}

// ProjectCreated represents a domain event when a project is created
type ProjectCreated struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	ProjectType ProjectType  `json:"project_type"`
	Database    DatabaseType `json:"database"`
	CreatedAt   time.Time    `json:"created_at"`
}

// ValidationResult represents the result of validation
type ValidationResult struct {
	IsValid  bool     `json:"is_valid"`
	Errors   []string `json:"errors"`
	Warnings []string `json:"warnings"`
}

// GenerateFilesCommand represents a command to generate files
type GenerateFilesCommand struct {
	TemplateData TemplateData `json:"template_data"`
	Force        bool         `json:"force"`
}

// FilesGenerated represents a domain event when files are generated
type FilesGenerated struct {
	TemplateData TemplateData `json:"template_data"`
	OutputFiles  []string     `json:"output_files"`
	GeneratedAt  time.Time    `json:"generated_at"`
}
