package templates

import (
	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/validation"
	"github.com/LarsArtmann/SQLC-Wizzard/pkg/config"
	"github.com/samber/lo"
)

// LibraryTemplate generates sqlc config for reusable Go libraries.
type LibraryTemplate struct {
	BaseTemplate
}

// NewLibraryTemplate creates a new library template.
func NewLibraryTemplate() *LibraryTemplate {
	return &LibraryTemplate{}
}

// Name returns the template name.
func (t *LibraryTemplate) Name() string {
	return "library"
}

// Description returns a human-readable description.
func (t *LibraryTemplate) Description() string {
	return "Library package configuration for reusable Go library development"
}

// Generate creates a SqlcConfig from template data.
func (t *LibraryTemplate) Generate(data generated.TemplateData) (*config.SqlcConfig, error) {
	// Library-specific defaults
	if data.Package.Name == "" {
		data.Package.Name = "db"
	}
	if data.Package.Path == "" {
		data.Package.Path = "internal/db"
	}
	if data.Output.BaseDir == "" {
		data.Output.BaseDir = "internal/db"
	}
	if data.Output.QueriesDir == "" {
		data.Output.QueriesDir = "internal/db/queries"
	}
	if data.Output.SchemaDir == "" {
		data.Output.SchemaDir = "internal/db/schema"
	}
	if data.Database.URL == "" {
		data.Database.URL = "${DATABASE_URL}"
	}

	// Build config
	cfg := &config.SqlcConfig{
		Version: "2",
		SQL: []config.SQLConfig{
			{
				Name:                 lo.Ternary(data.ProjectName != "", data.ProjectName, "library"),
				Engine:               string(data.Database.Engine),
				Queries:              config.NewPathOrPaths([]string{data.Output.QueriesDir}),
				Schema:               config.NewPathOrPaths([]string{data.Output.SchemaDir}),
				StrictFunctionChecks: lo.ToPtr(false),
				StrictOrderBy:        lo.ToPtr(false),
				Database: &config.DatabaseConfig{
					URI:     data.Database.URL,
					Managed: data.Database.UseManaged,
				},
				Gen: config.GenConfig{
					Go: t.buildGoGenConfig(data),
				},
				Rules: []config.RuleConfig{},
			},
		},
	}

	// Apply emit options using type-safe helper function
	config.ApplyEmitOptions(&data.Validation.EmitOptions, cfg.SQL[0].Gen.Go)

	// Convert rule types using the centralized transformer
	transformer := validation.NewRuleTransformer()
	rules := transformer.TransformSafetyRules(&data.Validation.SafetyRules)
	configRules := lo.Map(rules, func(r generated.RuleConfig, _ int) config.RuleConfig {
		return config.RuleConfig{
			Name:    r.Name,
			Rule:    r.Rule,
			Message: r.Message,
		}
	})
	cfg.SQL[0].Rules = configRules

	return cfg, nil
}

// DefaultData returns default TemplateData for library template.
func (t *LibraryTemplate) DefaultData() TemplateData {
	return generated.TemplateData{
		ProjectName: "",
		ProjectType: MustNewProjectType("library"),

		Package: generated.PackageConfig{
			Name: "db",
			Path: "internal/db",
		},

		Database: generated.DatabaseConfig{
			Engine:      MustNewDatabaseType("postgresql"),
			URL:         "${DATABASE_URL}",
			UseManaged:  false,
			UseUUIDs:    false,
			UseJSON:     false,
			UseArrays:   false,
			UseFullText: false,
		},

		Output: generated.OutputConfig{
			BaseDir:    "internal/db",
			QueriesDir: "internal/db/queries",
			SchemaDir:  "internal/db/schema",
		},

		Validation: generated.ValidationConfig{
			StrictFunctions: false,
			StrictOrderBy:   false,
			EmitOptions: generated.EmitOptions{
				EmitJSONTags:             true,
				EmitPreparedQueries:      false,
				EmitInterface:            true,
				EmitEmptySlices:          true,
				EmitResultStructPointers: false,
				EmitParamsStructPointers: false,
				EmitEnumValidMethod:      true,
				EmitAllEnumValues:        true,
				JSONTagsCaseStyle:        "camel",
			},
			SafetyRules: generated.SafetyRules{
				NoSelectStar: false,
				RequireWhere: false,
				NoDropTable:  true,
				NoTruncate:   true,
				RequireLimit: false,
				Rules:        []generated.SafetyRule{},
			},
		},
	}
}

// RequiredFeatures returns which features this template requires.
func (t *LibraryTemplate) RequiredFeatures() []string {
	return []string{"emit_interface", "json_tags", "enum_valid_method"}
}

// buildGoGenConfig builds the GoGenConfig from template data.
func (t *LibraryTemplate) buildGoGenConfig(data generated.TemplateData) *config.GoGenConfig {
	sqlPackage := t.GetSQLPackage(data.Database.Engine)
	cfg := t.BuildGoGenConfig(data, sqlPackage)

	// Libraries use simpler rename rules
	cfg.Rename = map[string]string{
		"id":   "ID",
		"json": "JSON",
		"api":  "API",
	}

	return cfg
}
