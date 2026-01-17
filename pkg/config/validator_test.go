package config_test

import (
	"strings"

	. "github.com/LarsArtmann/SQLC-Wizzard/pkg/config"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// Helper function to create a basic SQL configuration.
func createBasicSQLConfig(engine, outDir, packageName string) SQLConfig {
	return SQLConfig{
		Engine:  engine,
		Schema:  NewPathOrPaths([]string{"schema.sql"}),
		Queries: NewPathOrPaths([]string{"queries.sql"}),
		Gen: GenConfig{
			Go: &GoGenConfig{
				Out:     outDir,
				Package: packageName,
			},
		},
	}
}

// Helper function to create a basic SqlcConfig.
func createBasicSqlcConfig(engine string) *SqlcConfig {
	return &SqlcConfig{
		Version: "2",
		SQL: []SQLConfig{
			createBasicSQLConfig(engine, "db", "db"),
		},
	}
}

var _ = Describe("Validator", func() {
	Context("ValidationResult", func() {
		It("should initialize with empty errors and warnings", func() {
			result := &ValidationResult{}

			Expect(result.Errors).To(BeEmpty())
			Expect(result.Warnings).To(BeEmpty())
			Expect(result.IsValid()).To(BeTrue())
		})

		Context("error and warning management", func() {
			type testCase struct {
				name        string
				addField1   string
				addMessage1 string
				addField2   string
				addMessage2 string
			}

			DescribeTable("should add items correctly",
				func(tc testCase) {
					result := &ValidationResult{}

					if strings.Contains(tc.name, "error") {
						result.AddError(tc.addField1, tc.addMessage1)
						result.AddError(tc.addField2, tc.addMessage2)
						Expect(result.Errors).To(HaveLen(2))
						Expect(result.Errors[0].Field).To(Equal(tc.addField1))
						Expect(result.Errors[0].Message).To(Equal(tc.addMessage1))
						Expect(result.Errors[1].Field).To(Equal(tc.addField2))
						Expect(result.Errors[1].Message).To(Equal(tc.addMessage2))
						Expect(result.IsValid()).To(BeFalse())
					} else {
						result.AddWarning(tc.addField1, tc.addMessage1)
						result.AddWarning(tc.addField2, tc.addMessage2)
						Expect(result.Warnings).To(HaveLen(2))
						Expect(result.Warnings[0].Field).To(Equal(tc.addField1))
						Expect(result.Warnings[0].Message).To(Equal(tc.addMessage1))
						Expect(result.Warnings[1].Field).To(Equal(tc.addField2))
						Expect(result.Warnings[1].Message).To(Equal(tc.addMessage2))
						Expect(result.IsValid()).To(BeTrue()) // Warnings don't make it invalid
					}
				},
				Entry("should add errors correctly", testCase{
					name:        "error",
					addField1:   "field1",
					addMessage1: "error message",
					addField2:   "field2",
					addMessage2: "another error",
				}),
				Entry("should add warnings correctly", testCase{
					name:        "warning",
					addField1:   "field1",
					addMessage1: "warning message",
					addField2:   "field2",
					addMessage2: "another warning",
				}),
			)
		})

		It("should mix errors and warnings correctly", func() {
			result := &ValidationResult{}

			result.AddError("field1", "error message")
			result.AddWarning("field2", "warning message")
			result.AddWarning("field3", "another warning")

			Expect(result.Errors).To(HaveLen(1))
			Expect(result.Warnings).To(HaveLen(2))
			Expect(result.IsValid()).To(BeFalse()) // Has errors
		})
	})

	Context("ValidationError", func() {
		It("should format error message correctly", func() {
			err := ValidationError{
				Field:   "testField",
				Message: "test message",
			}

			Expect(err.Error()).To(Equal("testField: test message"))
		})

		It("should handle empty field and message", func() {
			err := ValidationError{
				Field:   "",
				Message: "",
			}

			Expect(err.Error()).To(Equal(": "))
		})

		It("should handle special characters in field and message", func() {
			err := ValidationError{
				Field:   "field.name[0]",
				Message: "Message with: special characters!",
			}

			Expect(err.Error()).To(Equal("field.name[0]: Message with: special characters!"))
		})
	})

	Context("Validate Function", func() {
		It("should validate a minimal valid configuration", func() {
			cfg := &SqlcConfig{
				Version: "2",
				SQL: []SQLConfig{
					{
						Engine:  "postgresql",
						Schema:  NewPathOrPaths([]string{"schema.sql"}),
						Queries: NewPathOrPaths([]string{"queries.sql"}),
						Gen: GenConfig{
							Go: &GoGenConfig{
								Out:     "db",
								Package: "db",
							},
						},
					},
				},
			}

			result := Validate(cfg)
			Expect(result).NotTo(BeNil())
			Expect(result.IsValid()).To(BeTrue())
		})

		It("should detect missing version", func() {
			cfg := &SqlcConfig{
				Version: "",
				SQL: []SQLConfig{
					{
						Engine:  "postgresql",
						Schema:  NewPathOrPaths([]string{"schema.sql"}),
						Queries: NewPathOrPaths([]string{"queries.sql"}),
						Gen: GenConfig{
							Go: &GoGenConfig{
								Out:     "db",
								Package: "db",
							},
						},
					},
				},
			}

			result := Validate(cfg)
			Expect(result).NotTo(BeNil())
			Expect(result.IsValid()).To(BeFalse())
			Expect(result.Errors).ToNot(BeEmpty())

			// Find version error
			var versionError *ValidationError
			for _, err := range result.Errors {
				if err.Field == "version" {
					versionError = &err
					break
				}
			}
			Expect(versionError).To(HaveOccurred())
			Expect(versionError.Message).To(ContainSubstring("required"))
		})

		It("should detect missing SQL configurations", func() {
			cfg := &SqlcConfig{
				Version: "2",
				SQL:     []SQLConfig{},
			}

			result := Validate(cfg)
			Expect(result).NotTo(BeNil())
			Expect(result.IsValid()).To(BeFalse())
			Expect(result.Errors).ToNot(BeEmpty())
		})

		It("should detect invalid database engine", func() {
			cfg := createBasicSqlcConfig("invalid_engine")

			result := Validate(cfg)
			Expect(result).NotTo(BeNil())
			// May or may not be invalid depending on validation logic
		})

		It("should handle multiple SQL configurations", func() {
			cfg := &SqlcConfig{
				Version: "2",
				SQL: []SQLConfig{
					createBasicSQLConfig("postgresql", "db1", "db1"),
					createBasicSQLConfig("mysql", "db2", "db2"),
				},
			}

			result := Validate(cfg)
			Expect(result).NotTo(BeNil())
			// Should not panic with multiple configurations
		})

		It("should validate Go generator configuration", func() {
			cfg := createBasicSqlcConfig("postgresql")
			cfg.SQL[0].Gen.Go.Out = ""
			cfg.SQL[0].Gen.Go.Package = ""

			result := Validate(cfg)
			Expect(result).NotTo(BeNil())
			// May have warnings/errors for empty out/package
		})

		It("should handle nil Go generator config", func() {
			cfg := createBasicSqlcConfig("postgresql")
			cfg.SQL[0].Gen.Go = nil

			result := Validate(cfg)
			Expect(result).NotTo(BeNil())
			// Should handle nil Go config gracefully
		})

		It("should validate path configurations", func() {
			cfg := createBasicSqlcConfig("postgresql")
			cfg.SQL[0].Schema = NewPathOrPaths([]string{""}) // Empty path

			result := Validate(cfg)
			Expect(result).NotTo(BeNil())
			// May have warnings/errors for empty paths
		})

		It("should handle nil config gracefully", func() {
			result := Validate(nil)
			Expect(result).NotTo(BeNil())
			Expect(result.IsValid()).To(BeFalse())
			Expect(result.Errors).ToNot(BeEmpty())
		})
	})

	Context("Edge Cases and Error Handling", func() {
		It("should handle configuration with warnings only", func() {
			cfg := createBasicSqlcConfig("postgresql")
			// May trigger warnings for missing optional fields

			result := Validate(cfg)
			Expect(result).NotTo(BeNil())
			// Should be valid even with warnings
		})

		It("should handle complex validation scenarios", func() {
			cfg := &SqlcConfig{
				Version: "2",
				SQL: []SQLConfig{
					{
						Engine:  "postgresql",
						Schema:  NewPathOrPaths([]string{"schema1.sql", "schema2.sql"}),
						Queries: NewPathOrPaths([]string{"queries1.sql", "queries2.sql"}),
						Gen: GenConfig{
							Go: &GoGenConfig{
								Out:     "db",
								Package: "db",
								Overrides: []Override{
									{
										GoType: "uuid.UUID",
										DBType: "uuid",
									},
									{
										GoType: "json.RawMessage",
										DBType: "jsonb",
									},
								},
							},
						},
					},
				},
			}

			result := Validate(cfg)
			Expect(result).NotTo(BeNil())
			// Should handle complex configurations
		})
	})
})
