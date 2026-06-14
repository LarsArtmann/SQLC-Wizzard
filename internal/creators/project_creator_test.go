package creators_test

import (
	"fmt"
	"testing"

	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/creators"
	"github.com/LarsArtmann/SQLC-Wizzard/pkg/config"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCreators(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Creators Suite")
}

// Enum test constants - single source of truth for all enum tests.
var allProjectTypes = []generated.ProjectType{
	generated.ProjectTypeMicroservice,
	generated.ProjectTypeHobby,
	generated.ProjectTypeEnterprise,
}

var allDatabaseTypes = []generated.DatabaseType{
	generated.DatabaseTypePostgreSQL,
	generated.DatabaseTypeMySQL,
	generated.DatabaseTypeSQLite,
}

// testEnumAssignment tests that an enum field can be correctly assigned and retrieved.
// Eliminated duplicate test structure by centralizing the pattern.
func testEnumAssignment[C, E any](
	fieldName string,
	values []E,
	setField func(*C, E),
	getField func(*C) E,
) {
	It(fmt.Sprintf("should support all %s types", fieldName), func() {
		for _, v := range values {
			cfg := new(C)
			setField(cfg, v)
			Expect(getField(cfg)).To(Equal(v), "%s should be correctly assigned", fieldName)
		}
	})
}

// createBaseConfig generates a base project configuration with standard defaults.
func createBaseConfig(projectName string) *creators.CreateConfig {
	return &creators.CreateConfig{
		ProjectName: projectName,
		ProjectType: generated.ProjectTypeMicroservice,
		Database:    generated.DatabaseTypePostgreSQL,
		Config: &config.SqlcConfig{
			Version: "2",
			SQL: []config.SQLConfig{
				{
					Engine:  "postgresql",
					Queries: config.NewPathOrPaths([]string{"queries/"}),
					Schema:  config.NewPathOrPaths([]string{"schema/"}),
					Gen: config.GenConfig{
						Go: &config.GoGenConfig{
							Package: "db",
							Out:     "internal/db",
						},
					},
				},
			},
		},
	}
}
