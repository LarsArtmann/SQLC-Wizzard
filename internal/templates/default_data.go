package templates

import (
	"fmt"

	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/pkg/config"
)

// BuildValidationConfig creates a ValidationConfig with the provided parameters.
// This eliminates duplication of validation configuration across templates.
func (t *BaseTemplate) BuildValidationConfig(
	strictFunctions, strictOrderBy bool,
	emitJSONTags, emitPreparedQueries, emitInterface, emitEmptySlices bool,
	emitResultStructPointers, emitParamsStructPointers, emitEnumValidMethod, emitAllEnumValues bool,
	jsonTagsCaseStyle string,
	noSelectStar, requireWhere, noDropTable, noTruncate, requireLimit bool,
) generated.ValidationConfig {
	return generated.ValidationConfig{
		StrictFunctions: strictFunctions,
		StrictOrderBy:   strictOrderBy,
		EmitOptions: generated.EmitOptions{
			EmitJSONTags:             emitJSONTags,
			EmitPreparedQueries:      emitPreparedQueries,
			EmitInterface:            emitInterface,
			EmitEmptySlices:          emitEmptySlices,
			EmitResultStructPointers: emitResultStructPointers,
			EmitParamsStructPointers: emitParamsStructPointers,
			EmitEnumValidMethod:      emitEnumValidMethod,
			EmitAllEnumValues:        emitAllEnumValues,
			JSONTagsCaseStyle:        jsonTagsCaseStyle,
		},
		SafetyRules: generated.SafetyRules{
			NoSelectStar: noSelectStar,
			RequireWhere: requireWhere,
			NoDropTable:  noDropTable,
			NoTruncate:   noTruncate,
			RequireLimit: requireLimit,
			Rules:        []generated.SafetyRule{},
		},
	}
}

// BuildDefaultData creates default TemplateData with the provided parameters.
// This eliminates duplication in template DefaultData() methods by providing
// a template method that accepts the variable configuration values.
//
// Deprecated: Use BuildDefaultDataFromOptions with BuildOptions struct instead.
// This method will be removed in a future version.
func (t *BaseTemplate) BuildDefaultData(
	projectType string,
	dbEngine string,
	databaseURL string,
	packagePath string,
	baseOutputDir string,
	useManaged, useUUIDs, useJSON, useArrays, useFullText bool,
	emitJSONTags, emitPreparedQueries, emitInterface, emitEmptySlices bool,
	emitResultStructPointers, emitParamsStructPointers, emitEnumValidMethod, emitAllEnumValues bool,
	jsonTagsCaseStyle string,
	strictFunctions, strictOrderBy bool,
	noSelectStar, requireWhere, noDropTable, noTruncate, requireLimit bool,
) generated.TemplateData {
	return generated.TemplateData{
		ProjectName: "",
		ProjectType: MustNewProjectType(projectType),

		Package: generated.PackageConfig{
			Name: "db",
			Path: packagePath,
		},

		Database: generated.DatabaseConfig{
			Engine:      MustNewDatabaseType(dbEngine),
			URL:         databaseURL,
			UseManaged:  useManaged,
			UseUUIDs:    useUUIDs,
			UseJSON:     useJSON,
			UseArrays:   useArrays,
			UseFullText: useFullText,
		},

		Output: generated.OutputConfig{
			BaseDir:    baseOutputDir,
			QueriesDir: baseOutputDir + "/queries",
			SchemaDir:  baseOutputDir + "/schema",
		},

		Validation: t.BuildValidationConfig(
			strictFunctions,
			strictOrderBy,
			emitJSONTags,
			emitPreparedQueries,
			emitInterface,
			emitEmptySlices,
			emitResultStructPointers,
			emitParamsStructPointers,
			emitEnumValidMethod,
			emitAllEnumValues,
			jsonTagsCaseStyle,
			noSelectStar,
			requireWhere,
			noDropTable,
			noTruncate,
			requireLimit,
		),
	}
}

// BuildDefaultDataFromOptions creates default TemplateData using BuildOptions.
// This is the preferred API for new templates as it reduces the 21-parameter
// signature to a single struct parameter.
//
// Example usage:
//
//	options := NewBuildOptions("microservice", "postgresql")
//	options.UseArrays = false
//	options.EmitJSONTags = true
//	data := t.BuildDefaultDataFromOptions(options)
func (t *BaseTemplate) BuildDefaultDataFromOptions(opts BuildOptions) generated.TemplateData {
	return generated.TemplateData{
		ProjectName: "",
		ProjectType: MustNewProjectType(opts.ProjectType),

		Package: generated.PackageConfig{
			Name: "db",
			Path: opts.PackagePath,
		},

		Database: generated.DatabaseConfig{
			Engine:      MustNewDatabaseType(opts.DbEngine),
			URL:         opts.DatabaseURL,
			UseManaged:  opts.UseManaged,
			UseUUIDs:    opts.UseUUIDs,
			UseJSON:     opts.UseJSON,
			UseArrays:   opts.UseArrays,
			UseFullText: opts.UseFullText,
		},

		Output: generated.OutputConfig{
			BaseDir:    opts.BaseOutputDir,
			QueriesDir: opts.BaseOutputDir + "/queries",
			SchemaDir:  opts.BaseOutputDir + "/schema",
		},

		Validation: t.BuildValidationConfig(
			opts.StrictFunctions,
			opts.StrictOrderBy,
			opts.EmitJSONTags,
			opts.EmitPreparedQueries,
			opts.EmitInterface,
			opts.EmitEmptySlices,
			opts.EmitResultStructPointers,
			opts.EmitParamsStructPointers,
			opts.EmitEnumValidMethod,
			opts.EmitAllEnumValues,
			opts.JSONTagsCaseStyle,
			opts.NoSelectStar,
			opts.RequireWhere,
			opts.NoDropTable,
			opts.NoTruncate,
			opts.RequireLimit,
		),
	}
}

// GenerateWithDefaults is a template method that eliminates duplicated code in template implementations.
// It applies template-specific defaults to the data and builds a SqlcConfig using the shared ConfigBuilder.
// Template implementations should call this method with their specific values.
//
// Parameters:
//   - data: The template data to apply defaults to
//   - packageName: Default package name to use when data.Package.Name is empty
//   - packagePath: Default package path to use when data.Package.Path is empty
//   - baseDir: Default base directory to use when data.Output.BaseDir is empty
//   - queriesDir: Default queries directory to use when data.Output.QueriesDir is empty
//   - schemaDir: Default schema directory to use when data.Output.SchemaDir is empty
//   - databaseURL: Default database URL to use when data.Database.URL is empty
//   - projectName: Default project name to use when data.ProjectName is empty
//   - strict: Whether to enable strict mode settings
//
// Returns: A SqlcConfig configured with the provided values, or an error if building fails.
func (t *BaseTemplate) GenerateWithDefaults(
	data generated.TemplateData,
	packageName string,
	packagePath string,
	baseDir string,
	queriesDir string,
	schemaDir string,
	databaseURL string,
	projectName string,
	strict bool,
) (*config.SqlcConfig, error) {
	applyGeneratedDefaults(&data, packageName, packagePath, baseDir, queriesDir, schemaDir, databaseURL)

	builder := &ConfigBuilder{
		Data:               data,
		DefaultName:        projectName,
		DefaultDatabaseURL: databaseURL,
		Strict:             strict,
	}

	cfg, err := builder.Build()
	if err != nil {
		return nil, fmt.Errorf(
			"config builder failed for project %q (package=%s, path=%s, baseDir=%s, queriesDir=%s, schemaDir=%s, strict=%v): %w",
			projectName,
			packageName,
			packagePath,
			baseDir,
			queriesDir,
			schemaDir,
			strict,
			err,
		)
	}

	cfg.SQL[0].Gen.Go = t.BuildGoConfigWithOverrides(data)

	return t.ApplyValidationRules(cfg, data)
}

// applyGeneratedDefaults fills empty fields in TemplateData with provided defaults.
// Centralizes the zero-value handling used by GenerateWithDefaults.
func applyGeneratedDefaults(
	data *generated.TemplateData,
	packageName, packagePath, baseDir, queriesDir, schemaDir, databaseURL string,
) {
	if data.Package.Name == "" {
		data.Package.Name = packageName
	}

	if data.Package.Path == "" {
		data.Package.Path = packagePath
	}

	if data.Output.BaseDir == "" {
		data.Output.BaseDir = baseDir
	}

	if data.Output.QueriesDir == "" {
		data.Output.QueriesDir = queriesDir
	}

	if data.Output.SchemaDir == "" {
		data.Output.SchemaDir = schemaDir
	}

	if data.Database.URL == "" {
		data.Database.URL = databaseURL
	}
}
