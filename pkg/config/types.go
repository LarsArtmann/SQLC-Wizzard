package config

import (
	"github.com/LarsArtmann/SQLC-Wizzard/generated"
)

// SqlcConfig represents the complete sqlc.yaml configuration (v2 schema)
type SqlcConfig struct {
	Version string       `yaml:"version"`
	Cloud   *CloudConfig `yaml:"cloud,omitempty"`
	SQL     []SQLConfig  `yaml:"sql"`
	Rules   []RuleConfig `yaml:"rules,omitempty"`
}

// CloudConfig represents sqlc Cloud integration settings
type CloudConfig struct {
	Organization string `yaml:"organization,omitempty"`
	Project      string `yaml:"project"`
	Token        string `yaml:"token,omitempty"`
	Hostname     string `yaml:"hostname,omitempty"`
}

// SQLConfig represents a single database configuration
type SQLConfig struct {
	Name                 string          `yaml:"name,omitempty"`
	Engine               string          `yaml:"engine"`
	Queries              PathOrPaths     `yaml:"queries"` // Paths to SQL query files
	Schema               PathOrPaths     `yaml:"schema"`  // Paths to SQL schema/migration files
	Gen                  GenConfig       `yaml:"gen"`
	Database             *DatabaseConfig `yaml:"database,omitempty"`
	Rules                []RuleConfig    `yaml:"rules,omitempty"`
	StrictFunctionChecks *bool           `yaml:"strict_function_checks,omitempty"`
	StrictOrderBy        *bool           `yaml:"strict_order_by,omitempty"`
	CodegenPlugins       []CodegenPlugin `yaml:"codegen,omitempty"`
}

// DatabaseConfig represents database connection settings
type DatabaseConfig struct {
	URI     string `yaml:"uri,omitempty"`
	Managed bool   `yaml:"managed,omitempty"`
}

// GenConfig represents code generation settings
type GenConfig struct {
	Go         *GoGenConfig         `yaml:"go,omitempty"`
	Kotlin     *KotlinGenConfig     `yaml:"kotlin,omitempty"`
	Python     *PythonGenConfig     `yaml:"python,omitempty"`
	TypeScript *TypeScriptGenConfig `yaml:"typescript,omitempty"`
}

// GoGenConfig represents Go-specific generation settings
type GoGenConfig struct {
	Package                     string            `yaml:"package"`
	Out                         string            `yaml:"out"`
	SQLPackage                  string            `yaml:"sql_package,omitempty"`
	BuildTags                   string            `yaml:"build_tags,omitempty"`
	EmitInterface               bool              `yaml:"emit_interface,omitempty"`
	EmitJSONTags                bool              `yaml:"emit_json_tags,omitempty"`
	EmitDBTags                  bool              `yaml:"emit_db_tags,omitempty"`
	EmitPreparedQueries         bool              `yaml:"emit_prepared_queries,omitempty"`
	EmitExactTableNames         bool              `yaml:"emit_exact_table_names,omitempty"`
	EmitEmptySlices             bool              `yaml:"emit_empty_slices,omitempty"`
	EmitExportedQueries         bool              `yaml:"emit_exported_queries,omitempty"`
	EmitResultStructPointers    bool              `yaml:"emit_result_struct_pointers,omitempty"`
	EmitParamsStructPointers    bool              `yaml:"emit_params_struct_pointers,omitempty"`
	EmitMethodsWithDBArgument   bool              `yaml:"emit_methods_with_db_argument,omitempty"`
	EmitPointersForNullTypes    bool              `yaml:"emit_pointers_for_null_types,omitempty"`
	EmitEnumValidMethod         bool              `yaml:"emit_enum_valid_method,omitempty"`
	EmitAllEnumValues           bool              `yaml:"emit_all_enum_values,omitempty"`
	JSONTagsCaseStyle           string            `yaml:"json_tags_case_style,omitempty"`
	OmitUnusedStructs           bool              `yaml:"omit_unused_structs,omitempty"`
	OmitSQLCVersion             bool              `yaml:"omit_sqlc_version,omitempty"`
	QueryParameterLimit         int               `yaml:"query_parameter_limit,omitempty"`
	OutputDBFileName            string            `yaml:"output_db_file_name,omitempty"`
	OutputModelsFileName        string            `yaml:"output_models_file_name,omitempty"`
	OutputQuerierFileName       string            `yaml:"output_querier_file_name,omitempty"`
	OutputCopyfromFileName      string            `yaml:"output_copyfrom_file_name,omitempty"`
	OutputBatchFileName         string            `yaml:"output_batch_file_name,omitempty"`
	Overrides                   []Override        `yaml:"overrides,omitempty"`
	Rename                      map[string]string `yaml:"rename,omitempty"`
	InflectionExcludeTableNames []string          `yaml:"inflection_exclude_table_names,omitempty"`
}

// KotlinGenConfig represents Kotlin-specific generation settings
type KotlinGenConfig struct {
	Package string `yaml:"package"`
	Out     string `yaml:"out"`
}

// PythonGenConfig represents Python-specific generation settings
type PythonGenConfig struct {
	Package string `yaml:"package"`
	Out     string `yaml:"out"`
}

// TypeScriptGenConfig represents TypeScript-specific generation settings
type TypeScriptGenConfig struct {
	Package string `yaml:"package"`
	Out     string `yaml:"out"`
}

// Override represents a type override configuration
type Override struct {
	DBType       string `yaml:"db_type,omitempty"`
	GoType       string `yaml:"go_type,omitempty"`
	GoImportPath string `yaml:"go_import_path,omitempty"`
	GoStructTag  string `yaml:"go_struct_tag,omitempty"`
	Nullable     bool   `yaml:"nullable,omitempty"`
	Column       string `yaml:"column,omitempty"`
	Table        string `yaml:"table,omitempty"`
	ColumnName   string `yaml:"column_name,omitempty"`
	GoBasicType  bool   `yaml:"go_basic_type,omitempty"`
	GoPointer    bool   `yaml:"go_pointer,omitempty"`
}

// RuleConfig represents a validation rule (CEL-based)
type RuleConfig struct {
	Name    string `yaml:"name"`
	Rule    string `yaml:"rule"`
	Message string `yaml:"message,omitempty"`
}

// CodegenPlugin represents a WASM plugin configuration
type CodegenPlugin struct {
	WASM    string            `yaml:"wasm,omitempty"`
	SHA256  string            `yaml:"sha256,omitempty"`
	Out     string            `yaml:"out,omitempty"`
	Options map[string]string `yaml:"options,omitempty"`
}

// ApplyEmitOptions applies emit options to a GoGenConfig in a type-safe manner
// This is a type-safe operation that eliminates field-by-field copying (DRY principle)
func ApplyEmitOptions(opts *generated.EmitOptions, cfg *GoGenConfig) {
	if opts.EmitJSONTags {
		cfg.EmitJSONTags = opts.EmitJSONTags
	}
	if opts.EmitPreparedQueries {
		cfg.EmitPreparedQueries = opts.EmitPreparedQueries
	}
	if opts.EmitInterface {
		cfg.EmitInterface = opts.EmitInterface
	}
	if opts.EmitEmptySlices {
		cfg.EmitEmptySlices = opts.EmitEmptySlices
	}
	if opts.EmitResultStructPointers {
		cfg.EmitResultStructPointers = opts.EmitResultStructPointers
	}
	if opts.EmitParamsStructPointers {
		cfg.EmitParamsStructPointers = opts.EmitParamsStructPointers
	}
	if opts.EmitEnumValidMethod {
		cfg.EmitEnumValidMethod = opts.EmitEnumValidMethod
	}
	if opts.EmitAllEnumValues {
		cfg.EmitAllEnumValues = opts.EmitAllEnumValues
	}
	if opts.JSONTagsCaseStyle != "" {
		cfg.JSONTagsCaseStyle = opts.JSONTagsCaseStyle
	}
}
