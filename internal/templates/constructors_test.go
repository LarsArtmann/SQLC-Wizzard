package templates_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/LarsArtmann/SQLC-Wizzard/internal/errors"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/templates"
)

var _ = Describe("Constructors", func() {
	Describe("NewProjectType", func() {
		Context("with valid project types", func() {
			It("should accept hobby", func() {
				pt, err := templates.NewProjectType("hobby")

				Expect(err).ToNot(HaveOccurred())
				Expect(pt).To(Equal(templates.ProjectTypeHobby))
			})

			It("should accept microservice", func() {
				pt, err := templates.NewProjectType("microservice")

				Expect(err).ToNot(HaveOccurred())
				Expect(pt).To(Equal(templates.ProjectTypeMicroservice))
			})

			It("should accept enterprise", func() {
				pt, err := templates.NewProjectType("enterprise")

				Expect(err).ToNot(HaveOccurred())
				Expect(pt).To(Equal(templates.ProjectTypeEnterprise))
			})

			It("should accept api-first", func() {
				pt, err := templates.NewProjectType("api-first")

				Expect(err).ToNot(HaveOccurred())
				Expect(pt).To(Equal(templates.ProjectTypeAPIFirst))
			})

			It("should accept library", func() {
				pt, err := templates.NewProjectType("library")

				Expect(err).ToNot(HaveOccurred())
				Expect(pt).To(Equal(templates.ProjectTypeLibrary))
			})
		})

		Context("with invalid project type", func() {
			It("should return validation error", func() {
				_, err := templates.NewProjectType("invalid-type")

				Expect(err).To(HaveOccurred())
				Expect(errors.HasCode(err, errors.ErrCodeValidationFailed)).To(BeTrue())
			})

			It("should provide helpful error message", func() {
				_, err := templates.NewProjectType("invalid")

				Expect(err.Error()).To(ContainSubstring("invalid project type"))
				Expect(err.Error()).To(ContainSubstring("hobby"))
				Expect(err.Error()).To(ContainSubstring("microservice"))
			})
		})
	})

	Describe("MustNewProjectType", func() {
		It("should return project type for valid input", func() {
			pt := templates.MustNewProjectType("microservice")

			Expect(pt).To(Equal(templates.ProjectTypeMicroservice))
		})

		It("should panic for invalid input", func() {
			Expect(func() {
				templates.MustNewProjectType("invalid")
			}).To(Panic())
		})
	})

	Describe("NewDatabaseType", func() {
		Context("with valid database types", func() {
			It("should accept postgresql", func() {
				dt, err := templates.NewDatabaseType("postgresql")

				Expect(err).ToNot(HaveOccurred())
				Expect(dt).To(Equal(templates.MustNewDatabaseType("postgresql")))
			})

			It("should accept mysql", func() {
				dt, err := templates.NewDatabaseType("mysql")

				Expect(err).ToNot(HaveOccurred())
				Expect(dt).To(Equal(templates.MustNewDatabaseType("mysql")))
			})

			It("should accept sqlite", func() {
				dt, err := templates.NewDatabaseType("sqlite")

				Expect(err).ToNot(HaveOccurred())
				Expect(dt).To(Equal(templates.MustNewDatabaseType("sqlite")))
			})
		})

		Context("with invalid database type", func() {
			It("should return validation error", func() {
				_, err := templates.NewDatabaseType("mongodb")

				Expect(err).To(HaveOccurred())
				Expect(errors.HasCode(err, errors.ErrCodeValidationFailed)).To(BeTrue())
			})

			It("should provide helpful error message", func() {
				_, err := templates.NewDatabaseType("invalid")

				Expect(err.Error()).To(ContainSubstring("invalid database type"))
				Expect(err.Error()).To(ContainSubstring("postgresql"))
				Expect(err.Error()).To(ContainSubstring("mysql"))
				Expect(err.Error()).To(ContainSubstring("sqlite"))
			})
		})
	})

	Describe("MustNewDatabaseType", func() {
		It("should return database type for valid input", func() {
			dt := templates.MustNewDatabaseType("postgresql")

			Expect(dt).To(Equal(templates.MustNewDatabaseType("postgresql")))
		})

		It("should panic for invalid input", func() {
			Expect(func() {
				templates.MustNewDatabaseType("invalid")
			}).To(Panic())
		})
	})

	Describe("Validation Functions", func() {
		Context("ProjectType validation", func() {
			It("should return true for valid types", func() {
				Expect(templates.IsValidProjectType("hobby")).To(BeTrue())
			Expect(templates.IsValidProjectType("microservice")).To(BeTrue())
			Expect(templates.IsValidProjectType("enterprise")).To(BeTrue())
			})

			It("should return false for invalid types", func() {
				Expect(templates.IsValidProjectType("invalid")).To(BeFalse())
			Expect(templates.IsValidProjectType("")).To(BeFalse())
			Expect(templates.IsValidProjectType("MICROSERVICE")).To(BeFalse()) // Case sensitive
			})
		})

		Context("DatabaseType validation", func() {
			It("should return true for valid types", func() {
				Expect(templates.IsValidDatabaseType("postgresql")).To(BeTrue())
			Expect(templates.IsValidDatabaseType("mysql")).To(BeTrue())
			Expect(templates.IsValidDatabaseType("sqlite")).To(BeTrue())
			})

			It("should return false for invalid types", func() {
				Expect(templates.IsValidDatabaseType("mongodb")).To(BeFalse())
			Expect(templates.IsValidDatabaseType("")).To(BeFalse())
			Expect(templates.IsValidDatabaseType("PostgreSQL")).To(BeFalse()) // Case sensitive
			})
		})
	})
})
