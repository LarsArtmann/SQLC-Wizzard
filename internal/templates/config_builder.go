package templates

import (
	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/pkg/config"
	"github.com/samber/lo"
)

// ConfigBuilder helps construct sqlc configurations with common patterns.
// This eliminates duplication across template implementations.
type ConfigBuilder struct {
	// Data holds the template configuration data.
	Data generated.TemplateData
	// DefaultName is used when ProjectName is empty.
	DefaultName string
	// DefaultDatabaseURL is used when Database.URL is empty.
	DefaultDatabaseURL string
	// Strict controls strict mode settings.
	Strict bool
}

// Build creates a SqlcConfig from the configured values.
func (cb *ConfigBuilder) Build() (*config.SqlcConfig, error) {
	base := &BaseTemplate{}

	return &config.SqlcConfig{
		Version: "2",
		SQL: []config.SQLConfig{
			{
				Name: lo.Ternary(
					cb.Data.ProjectName != "",
					cb.Data.ProjectName,
					cb.DefaultName,
				),
				Engine:               string(cb.Data.Database.Engine),
				Queries:              config.NewPathOrPaths([]string{cb.Data.Output.QueriesDir}),
				Schema:               config.NewPathOrPaths([]string{cb.Data.Output.SchemaDir}),
				StrictFunctionChecks: new(cb.Strict),
				StrictOrderBy:        new(cb.Strict),
				Database: &config.DatabaseConfig{
					URI: lo.Ternary(
						cb.Data.Database.URL != "",
						cb.Data.Database.URL,
						cb.DefaultDatabaseURL,
					),
					Managed: cb.Data.Database.UseManaged,
				},
				Gen: config.GenConfig{
					Go: base.BuildGoGenConfig(cb.Data, base.GetSQLPackage(cb.Data.Database.Engine)),
				},
				Rules: []config.RuleConfig{},
			},
		},
	}, nil
}
