namespace sqlc.wizard;

/// ProjectType represents the type of project template
enum ProjectType {
  hobby = "hobby";
  microservice = "microservice";
  enterprise = "enterprise";
  api_first = "api-first";
  analytics = "analytics";
  testing = "testing";
  multi_tenant = "multi-tenant";
  library = "library";
}

/// DatabaseType represents the supported database engines
enum DatabaseType {
  postgresql = "postgresql";
  mysql = "mysql";
  sqlite = "sqlite";
}

/// EmitOptions defines SQL code generation options
model EmitOptions {
  emit_json_tags: bool;
  emit_prepared_queries: bool;
  emit_interface: bool;
  emit_empty_slices: bool;
  emit_result_struct_pointers: bool;
  emit_params_struct_pointers: bool;
  emit_enum_valid_method: bool;
  emit_all_enum_values: bool;
  json_tags_case_style: string; // TODO: convert to enum
}

/// SafetyRule represents a CEL-based validation rule
model SafetyRule {
  name: string;
  rule: string;
  message: string;
}

/// TemplateData represents the complete data structure for template generation
model TemplateData {
  project_name: string;
  project_type: ProjectType;
  package_path: string;
  
  database: DatabaseType;
  database_url: string?;
  use_managed_db: bool;
  
  output_dir: string;
  queries_dir: string;
  schema_dir: string;
  
  package_name: string;
  sql_package: string;
  build_tags: string?;
  
  use_uuids: bool;
  use_json: bool;
  use_arrays: bool;
  use_full_text_search: bool;
  
  emit_options: EmitOptions;
  safety_rules: SafetyRule[];
  
  strict_functions: bool;
  strict_order_by: bool;
}

/// CreateProjectCommand represents a command to create a new project
model CreateProjectCommand {
  name: string;
  project_type: ProjectType;
  database: DatabaseType;
  output_dir: string;
}

/// ProjectCreated represents a domain event when a project is created
model ProjectCreated {
  id: string;
  name: string;
  project_type: ProjectType;
  database: DatabaseType;
  created_at: string;
}

/// ValidationResult represents the result of validation
model ValidationResult {
  is_valid: bool;
  errors: string[];
  warnings: string[];
}

/// GenerateFilesCommand represents a command to generate files
model GenerateFilesCommand {
  template_data: TemplateData;
  force: bool;
}

/// FilesGenerated represents a domain event when files are generated
model FilesGenerated {
  template_data: TemplateData;
  output_files: string[];
  generated_at: string;
}