package commands_test

import (
	"github.com/LarsArtmann/SQLC-Wizzard/pkg/config"
)

// Helper function to create test SqlcConfig instances.
func createTestSqlcConfig(schema, out, pkg string, engine ...string) *config.SqlcConfig {
	dbEngine := "postgresql"
	if len(engine) > 0 {
		dbEngine = engine[0]
	}

	return &config.SqlcConfig{
		Version: "2",
		SQL: []config.SQLConfig{
			{
				Engine: dbEngine,
				Schema: config.NewPathOrPaths([]string{schema}),
				Gen: config.GenConfig{
					Go: &config.GoGenConfig{
						Out:     out,
						Package: pkg,
					},
				},
			},
		},
	}
}
