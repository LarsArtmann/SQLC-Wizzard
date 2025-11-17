package wizard_test

import (
	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// validateProjectType is a helper function to test project type validation
func validateProjectType(input string, expected generated.ProjectType, valid bool) {
	projectType := generated.ProjectType(input)
	if valid {
		Expect(projectType.IsValid()).To(BeTrue(), "for input: %s", input)
		Expect(projectType).To(Equal(expected), "for input: %s", input)
	} else {
		Expect(projectType.IsValid()).To(BeFalse(), "for input: %s", input)
	}
}

// validateDatabaseType is a helper function to test database type validation
func validateDatabaseType(input string, expected generated.DatabaseType, valid bool) {
	dbType := generated.DatabaseType(input)
	if valid {
		Expect(dbType.IsValid()).To(BeTrue(), "for input: %s", input)
		Expect(dbType).To(Equal(expected), "for input: %s", input)
	} else {
		Expect(dbType.IsValid()).To(BeFalse(), "for input: %s", input)
	}
}

var _ = Describe("Wizard Steps", func() {
	It("should handle project type selection", func() {
		// Test project type validation using helper function
		validateProjectType("microservice", generated.ProjectTypeMicroservice, true)
		validateProjectType("enterprise", generated.ProjectTypeEnterprise, true)
		validateProjectType("hobby", generated.ProjectTypeHobby, true)
		validateProjectType("invalid", generated.ProjectType(""), false)
	})

	It("should handle database type selection", func() {
		// Test database type validation using helper function
		validateDatabaseType("postgresql", generated.DatabaseTypePostgreSQL, true)
		validateDatabaseType("mysql", generated.DatabaseTypeMySQL, true)
		validateDatabaseType("sqlite", generated.DatabaseTypeSQLite, true)
		validateDatabaseType("invalid", generated.DatabaseType(""), false)
	})
})
