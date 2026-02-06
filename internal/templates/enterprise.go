package templates

import (
	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/validation"
	"github.com/LarsArtmann/SQLC-Wizzard/pkg/config"
	"github.com/samber/lo"
)

// EnterpriseTemplate generates sqlc config for enterprise-scale projects.
type EnterpriseTemplate struct {
	BaseTemplate
}

// NewEnterpriseTemplate creates a new enterprise template.
func NewEnterpriseTemplate() *EnterpriseTemplate {
	return &EnterpriseTemplate{}
}

// Name returns the template name.
func (t *EnterpriseTemplate) Name() string {
	return "enterprise"
}

// Description returns a human-readable description.
func (t *EnterpriseTemplate) Description() string {
	return "Production-ready configuration with strict safety rules for enterprise applications"
}

// Generate creates a SqlcConfig from template data.
func (t *EnterpriseTemplate) Generate(data generated.TemplateData) (*config.SqlcConfig, error) {
	// Set defaults
	cb := ConfigBuilder{
		Data:               data,
		DefaultName:        "enterprise",
		DefaultDatabaseURL: "${DATABASE_URL}",
		Strict:             true,
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

// DefaultData returns default TemplateData for enterprise template.
func (t *EnterpriseTemplate) DefaultData() TemplateData {
	return generated.TemplateData{
		ProjectName: "",
		ProjectType: MustNewProjectType("enterprise"),

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
			UseFullText: true,
		},

		Output: generated.OutputConfig{
			BaseDir:    "internal/db",
			QueriesDir: "internal/db/queries",
			SchemaDir:  "internal/db/schema",
		},

		Validation: generated.ValidationConfig{
			StrictFunctions: true,
			StrictOrderBy:   true,
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
				RequireLimit: true,
				Rules:        []generated.SafetyRule{},
			},
		},
	}
}

// RequiredFeatures returns which features this template requires.
func (t *EnterpriseTemplate) RequiredFeatures() []string {
	return []string{"emit_interface", "prepared_queries", "json_tags", "strict_checks"}
}