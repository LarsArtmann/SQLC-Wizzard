package wizard_test

import (
	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Wizard Steps", func() {
	It("should handle project type selection", func() {
		// Test project type validation
		testCases := []struct {
			input    string
			expected generated.ProjectType
			valid    bool
		}{
			{"microservice", generated.ProjectTypeMicroservice, true},
			{"enterprise", generated.ProjectTypeEnterprise, true},
			{"hobby", generated.ProjectTypeHobby, true},
			{"invalid", generated.ProjectType(""), false},
		}

		for _, tc := range testCases {
			projectType := generated.ProjectType(tc.input)
			if tc.valid {
				Expect(projectType.IsValid()).To(BeTrue(), "for input: %s", tc.input)
				Expect(projectType).To(Equal(tc.expected), "for input: %s", tc.input)
			} else {
				Expect(projectType.IsValid()).To(BeFalse(), "for input: %s", tc.input)
			}
		}
	})

	It("should handle database type selection", func() {
		// Test database type validation
		testCases := []struct {
			input    string
			expected generated.DatabaseType
			valid    bool
		}{
			{"postgresql", generated.DatabaseTypePostgreSQL, true},
			{"mysql", generated.DatabaseTypeMySQL, true},
			{"sqlite", generated.DatabaseTypeSQLite, true},
			{"invalid", generated.DatabaseType(""), false},
		}

		for _, tc := range testCases {
			dbType := generated.DatabaseType(tc.input)
			if tc.valid {
				Expect(dbType.IsValid()).To(BeTrue(), "for input: %s", tc.input)
				Expect(dbType).To(Equal(tc.expected), "for input: %s", tc.input)
			} else {
				Expect(dbType.IsValid()).To(BeFalse(), "for input: %s", tc.input)
			}
		}
	})
})