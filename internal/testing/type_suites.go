package testing

import (
	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	. "github.com/onsi/gomega"
)

// ValidProjectTypes contains all valid project types for testing purposes.
var ValidProjectTypes = []generated.ProjectType{
	generated.ProjectTypeHobby,
	generated.ProjectTypeMicroservice,
	generated.ProjectTypeEnterprise,
	generated.ProjectTypeAPIFirst,
	generated.ProjectTypeAnalytics,
	generated.ProjectTypeTesting,
	generated.ProjectTypeMultiTenant,
	generated.ProjectTypeLibrary,
}

// ValidDatabaseTypes contains all valid database types for testing purposes.
var ValidDatabaseTypes = []generated.DatabaseType{
	generated.DatabaseTypePostgreSQL,
	generated.DatabaseTypeMySQL,
	generated.DatabaseTypeSQLite,
}

// ProjectTypeTestSuite implements ValidationTestSuite for ProjectType validation tests.
type ProjectTypeTestSuite struct{}

func (s ProjectTypeTestSuite) GetValidValues() []generated.ProjectType { return ValidProjectTypes }
func (s ProjectTypeTestSuite) GetInvalidValues() []generated.ProjectType {
	return []generated.ProjectType{generated.ProjectType("invalid-type")}
}
func (s ProjectTypeTestSuite) GetTypeName() string { return "ProjectType" }

// DatabaseTypeTestSuite implements ValidationTestSuite for DatabaseType validation tests.
type DatabaseTypeTestSuite struct{}

func (s DatabaseTypeTestSuite) GetValidValues() []generated.DatabaseType { return ValidDatabaseTypes }

func (s DatabaseTypeTestSuite) GetInvalidValues() []generated.DatabaseType {
	return []generated.DatabaseType{generated.DatabaseType("invalid-db")}
}
func (s DatabaseTypeTestSuite) GetTypeName() string { return "DatabaseType" }

// ValidateAllProjectTypes tests that all project types in ValidProjectTypes are valid.
// This helper eliminates duplicate validation code across test files.
func ValidateAllProjectTypes() {
	for _, projectType := range ValidProjectTypes {
		Expect(projectType.IsValid()).To(BeTrue(),
			"Project type %s should be valid", projectType)
	}
}

// ValidateAllDatabaseTypes tests that all database types in ValidDatabaseTypes are valid.
// This helper eliminates duplicate validation code across test files.
func ValidateAllDatabaseTypes() {
	for _, dbType := range ValidDatabaseTypes {
		Expect(dbType.IsValid()).To(BeTrue(),
			"Database type %s should be valid", dbType)
	}
}
