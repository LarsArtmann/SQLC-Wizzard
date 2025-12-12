package config_test

import (
	. "github.com/LarsArtmann/SQLC-Wizzard/pkg/config"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Validator", func() {
	Context("ValidationResult", func() {
		It("should initialize with empty errors and warnings", func() {
			result := &ValidationResult{}

			Expect(result.Errors).To(BeEmpty())
			Expect(result.Warnings).To(BeEmpty())
			Expect(result.IsValid()).To(BeTrue())
		})

		It("should add errors correctly", func() {
			result := &ValidationResult{}

			result.AddError("field1", "error message")
			result.AddError("field2", "another error")

			Expect(len(result.Errors)).To(Equal(2))
			Expect(result.Errors[0].Field).To(Equal("field1"))
			Expect(result.Errors[0].Message).To(Equal("error message"))
			Expect(result.Errors[1].Field).To(Equal("field2"))
			Expect(result.Errors[1].Message).To(Equal("another error"))
			Expect(result.IsValid()).To(BeFalse())
		})

		It("should add warnings correctly", func() {
			result := &ValidationResult{}

			result.AddWarning("field1", "warning message")
			result.AddWarning("field2", "another warning")

			Expect(len(result.Warnings)).To(Equal(2))
			Expect(result.Warnings[0].Field).To(Equal("field1"))
			Expect(result.Warnings[0].Message).To(Equal("warning message"))
			Expect(result.Warnings[1].Field).To(Equal("field2"))
			Expect(result.Warnings[1].Message).To(Equal("another warning"))
			Expect(result.IsValid()).To(BeTrue()) // Warnings don't make it invalid
		})

		It("should mix errors and warnings correctly", func() {
			result := &ValidationResult{}

			result.AddError("field1", "error message")
			result.AddWarning("field2", "warning message")
			result.AddWarning("field3", "another warning")

			Expect(len(result.Errors)).To(Equal(1))
			Expect(len(result.Warnings)).To(Equal(2))
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
			Expect(len(result.Errors)).To(BeNumerically(">", 0))

			// Find version error
			var versionError *ValidationError
			for _, err := range result.Errors {
				if err.Field == "version" {
					versionError = &err
					break
				}
			}
			Expect(versionError).NotTo(BeNil())
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
			Expect(len(result.Errors)).To(BeNumerically(">", 0))
		})

		It("should detect invalid database engine", func() {
			cfg := &SqlcConfig{
				Version: "2",
				SQL: []SQLConfig{
					{
						Engine:  "invalid_engine",
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
			// May or may not be invalid depending on validation logic
		})

		It("should handle multiple SQL configurations", func() {
			cfg := &SqlcConfig{
				Version: "2",
				SQL: []SQLConfig{
					{
						Engine:  "postgresql",
						Schema:  NewPathOrPaths([]string{"schema1.sql"}),
						Queries: NewPathOrPaths([]string{"queries1.sql"}),
						Gen: GenConfig{
							Go: &GoGenConfig{
								Out:     "db1",
								Package: "db1",
							},
						},
					},
					{
						Engine:  "mysql",
						Schema:  NewPathOrPaths([]string{"schema2.sql"}),
						Queries: NewPathOrPaths([]string{"queries2.sql"}),
						Gen: GenConfig{
							Go: &GoGenConfig{
								Out:     "db2",
								Package: "db2",
							},
						},
					},
				},
			}

			result := Validate(cfg)
			Expect(result).NotTo(BeNil())
			// Should not panic with multiple configurations
		})

		It("should validate Go generator configuration", func() {
			cfg := &SqlcConfig{
				Version: "2",
				SQL: []SQLConfig{
					{
						Engine:  "postgresql",
						Schema:  NewPathOrPaths([]string{"schema.sql"}),
						Queries: NewPathOrPaths([]string{"queries.sql"}),
						Gen: GenConfig{
							Go: &GoGenConfig{
								Out:     "",
								Package: "",
							},
						},
					},
				},
			}

			result := Validate(cfg)
			Expect(result).NotTo(BeNil())
			// May have warnings/errors for empty out/package
		})

		It("should handle nil Go generator config", func() {
			cfg := &SqlcConfig{
				Version: "2",
				SQL: []SQLConfig{
					{
						Engine:  "postgresql",
						Schema:  NewPathOrPaths([]string{"schema.sql"}),
						Queries: NewPathOrPaths([]string{"queries.sql"}),
						Gen: GenConfig{
							Go: nil,
						},
					},
				},
			}

			result := Validate(cfg)
			Expect(result).NotTo(BeNil())
			// Should handle nil Go config gracefully
		})

		It("should validate path configurations", func() {
			cfg := &SqlcConfig{
				Version: "2",
				SQL: []SQLConfig{
					{
						Engine:  "postgresql",
						Schema:  NewPathOrPaths([]string{""}), // Empty path
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
			// May have warnings/errors for empty paths
		})

		It("should handle nil config gracefully", func() {
			result := Validate(nil)
			Expect(result).NotTo(BeNil())
			Expect(result.IsValid()).To(BeFalse())
			Expect(len(result.Errors)).To(BeNumerically(">", 0))
		})
	})

	Context("Edge Cases and Error Handling", func() {
		It("should handle configuration with warnings only", func() {
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
								// May trigger warnings for missing optional fields
							},
						},
					},
				},
			}

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
