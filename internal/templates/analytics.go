package templates

import (
	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/pkg/config"
	"github.com/samber/lo"
)

// AnalyticsTemplate generates sqlc config for analytics and data warehouse projects.
type AnalyticsTemplate struct {
	BaseTemplate
}

// NewAnalyticsTemplate creates a new analytics template.
func NewAnalyticsTemplate() *AnalyticsTemplate {
	return &AnalyticsTemplate{}
}

// Name returns the template name.
func (t *AnalyticsTemplate) Name() string {
	return "analytics"
}

// Description returns a human-readable description.
func (t *AnalyticsTemplate) Description() string {
	return "Optimized for data analytics and reporting with full-text search and array support"
}

// Generate creates a SqlcConfig from template data.
func (t *AnalyticsTemplate) Generate(data generated.TemplateData) (*config.SqlcConfig, error) {
	// Analytics-specific defaults
	if data.Package.Name == "" {
		data.Package.Name = "analytics"
	}
	if data.Package.Path == "" {
		data.Package.Path = "internal/analytics"
	}
	if data.Output.BaseDir == "" {
		data.Output.BaseDir = "internal/analytics"
	}
	if data.Output.QueriesDir == "" {
		data.Output.QueriesDir = "internal/analytics/queries"
	}
	if data.Output.SchemaDir == "" {
		data.Output.SchemaDir = "internal/analytics/schema"
	}
	if data.Database.URL == "" {
		data.Database.URL = "${ANALYTICS_DATABASE_URL}"
	}

	// Build config
	cfg := &config.SqlcConfig{
		Version: "2",
		SQL: []config.SQLConfig{
			{
				Name:                 lo.Ternary(data.ProjectName != "", data.ProjectName, "analytics"),
				Engine:               string(data.Database.Engine),
				Queries:              config.NewPathOrPaths([]string{data.Output.QueriesDir}),
				Schema:               config.NewPathOrPaths([]string{data.Output.SchemaDir}),
				StrictFunctionChecks: lo.ToPtr(true),
				StrictOrderBy:        lo.ToPtr(true),
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

	// Apply validation rules using helper to eliminate duplication
	return t.ApplyValidationRules(cfg, data)
}

// DefaultData returns default TemplateData for analytics template.
func (t *AnalyticsTemplate) DefaultData() TemplateData {
	return generated.TemplateData{
		ProjectName: "",
		ProjectType: MustNewProjectType("analytics"),

		Package: generated.PackageConfig{
			Name: "analytics",
			Path: "internal/analytics",
		},

		Database: generated.DatabaseConfig{
			Engine:      MustNewDatabaseType("postgresql"),
			URL:         "${ANALYTICS_DATABASE_URL}",
			UseManaged:  false,
			UseUUIDs:    false,
			UseJSON:     true,
			UseArrays:   true,
			UseFullText: true,
		},

		Output: generated.OutputConfig{
			BaseDir:    "internal/analytics",
			QueriesDir: "internal/analytics/queries",
			SchemaDir:  "internal/analytics/schema",
		},

		Validation: generated.ValidationConfig{
			StrictFunctions: true,
			StrictOrderBy:   true,
			EmitOptions: generated.EmitOptions{
				EmitJSONTags:             true,
				EmitPreparedQueries:      false,
				EmitInterface:            true,
				EmitEmptySlices:          true,
				EmitResultStructPointers: false,
				EmitParamsStructPointers: false,
				EmitEnumValidMethod:      false,
				EmitAllEnumValues:        true,
				JSONTagsCaseStyle:        "snake",
			},
			SafetyRules: generated.SafetyRules{
				NoSelectStar: false,
				RequireWhere: false,
				NoDropTable:  true,
				NoTruncate:   true,
				RequireLimit: true,
				Rules:        []generated.SafetyRule{},
			},
		},
	}
}

// RequiredFeatures returns which features this template requires.
func (t *AnalyticsTemplate) RequiredFeatures() []string {
	return []string{"emit_interface", "json_tags", "full_text_search", "strict_checks"}
}

// buildGoGenConfig builds the GoGenConfig from template data.
func (t *AnalyticsTemplate) buildGoGenConfig(data generated.TemplateData) *config.GoGenConfig {
	sqlPackage := t.GetSQLPackage(data.Database.Engine)
	cfg := t.BuildGoGenConfig(data, sqlPackage)

	// Analytics uses specific rename rules
	cfg.Rename = map[string]string{
		"id":   "ID",
		"json": "JSON",
		"api":  "API",
		"http": "HTTP",
	}

	return cfg
}
