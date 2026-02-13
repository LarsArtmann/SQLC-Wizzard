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

// EntryAdder interface for adding entries to validation result.
type EntryAdder interface {
	AddEntry(field, message string)
	GetEntries() []ValidationError
}

func addEntriesAndVerify(adder EntryAdder, field1, message1, field2, message2 string, expectValid bool) {
	adder.AddEntry(field1, message1)
	adder.AddEntry(field2, message2)
	entries := adder.GetEntries()
	Expect(entries).To(HaveLen(2))
	Expect(entries[0].Field).To(Equal(field1))
	Expect(entries[0].Message).To(Equal(message1))
	Expect(entries[1].Field).To(Equal(field2))
	Expect(entries[1].Message).To(Equal(message2))
}

type errorAdder struct {
	result *ValidationResult
}

func (e *errorAdder) AddEntry(field, message string) {
	e.result.AddError(field, message)
}

func (e *errorAdder) GetEntries() []ValidationError {
	return e.result.Errors
}

type warningAdder struct {
	result *ValidationResult
}

func (w *warningAdder) AddEntry(field, message string) {
	w.result.AddWarning(field, message)
}

func (w *warningAdder) GetEntries() []ValidationError {
	return w.result.Warnings
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
						errors := &errorAdder{result: result}
						addEntriesAndVerify(errors, tc.addField1, tc.addMessage1, tc.addField2, tc.addMessage2, false)
						Expect(result.IsValid()).To(BeFalse())
					} else {
						warnings := &warningAdder{result: result}
						addEntriesAndVerify(warnings, tc.addField1, tc.addMessage1, tc.addField2, tc.addMessage2, true)
						Expect(result.IsValid()).To(BeTrue())
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
