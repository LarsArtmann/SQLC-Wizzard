package templates

import (
	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/validation"
	"github.com/LarsArtmann/SQLC-Wizzard/pkg/config"
	"github.com/samber/lo"
)

// APIFirstTemplate generates sqlc config for API-first projects.
type APIFirstTemplate struct {
	BaseTemplate
}

// NewAPIFirstTemplate creates a new API-first template.
func NewAPIFirstTemplate() *APIFirstTemplate {
	return &APIFirstTemplate{}
}

// Name returns the template name.
func (t *APIFirstTemplate) Name() string {
	return "api-first"
}

// Description returns a human-readable description.
func (t *APIFirstTemplate) Description() string {
	return "Optimized for REST/GraphQL API development with JSON support and camelCase naming"
}

// Generate creates a SqlcConfig from template data.
func (t *APIFirstTemplate) Generate(data generated.TemplateData) (*config.SqlcConfig, error) {
	// Set defaults
	cb := ConfigBuilder{
		Data:               data,
		DefaultName:        "api",
		DefaultDatabaseURL: "${DATABASE_URL}",
		Strict:             false,
	}

	// Build config using base builder
	cfg, err := cb.Build()
	if err != nil {
		return nil, err
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

// DefaultData returns default TemplateData for API-first template.
func (t *APIFirstTemplate) DefaultData() TemplateData {
	return generated.TemplateData{
		ProjectName: "",
		ProjectType: MustNewProjectType("api-first"),

		Package: generated.PackageConfig{
			Name: "db",
			Path: "internal/db",
		},

		Database: generated.DatabaseConfig{
			Engine:      MustNewDatabaseType("postgresql"),
			URL:         "${DATABASE_URL}",
			UseManaged:  true,
			UseUUIDs:    true,
			UseJSON:     true,
			UseArrays:   true,
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
				EmitPreparedQueries:      true,
				EmitInterface:            true,
				EmitEmptySlices:          true,
				EmitResultStructPointers: true,
				EmitParamsStructPointers: true,
				EmitEnumValidMethod:      true,
				EmitAllEnumValues:        true,
				JSONTagsCaseStyle:        "camel",
			},
			SafetyRules: generated.SafetyRules{
				NoSelectStar: true,
				RequireWhere: true,
				NoDropTable:  true,
				NoTruncate:   true,
				RequireLimit: false,
				Rules:        []generated.SafetyRule{},
			},
		},
	}
}

// RequiredFeatures returns which features this template requires.
func (t *APIFirstTemplate) RequiredFeatures() []string {
	return []string{"emit_interface", "prepared_queries", "json_tags", "camel_case"}
}
