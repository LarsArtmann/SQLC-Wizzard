package schema_test

import (
	"testing"

	"github.com/LarsArtmann/SQLC-Wizzard/internal/schema"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestSchema(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Schema Suite")
}

var singleTableWithIntID = []schema.Table{
	{
		Name: "users",
		Columns: []schema.Column{
			{Name: "id", Type: schema.ColumnTypeInteger},
		},
	},
}

func makeTable(name string) []schema.Table {
	return []schema.Table{
		{
			Name: name,
			Columns: []schema.Column{
				{Name: "id", Type: schema.ColumnTypeInteger},
			},
		},
	}
}

var _ = Describe("ColumnType", func() {
	Context("IsValid", func() {
		It("should return true for valid column types", func() {
			validTypes := []schema.ColumnType{
				schema.ColumnTypeString,
				schema.ColumnTypeInteger,
				schema.ColumnTypeBigInt,
				schema.ColumnTypeFloat,
				schema.ColumnTypeDouble,
				schema.ColumnTypeBoolean,
				schema.ColumnTypeDate,
				schema.ColumnTypeDateTime,
				schema.ColumnTypeTimestamp,
				schema.ColumnTypeJSON,
				schema.ColumnTypeUUID,
				schema.ColumnTypeText,
				schema.ColumnTypeBlob,
			}

			for _, ct := range validTypes {
				Expect(ct.IsValid()).To(BeTrue(), "ColumnType %s should be valid", ct)
			}
		})

		It("should return false for invalid column types", func() {
			invalidTypes := []schema.ColumnType{
				"invalid",
				"INVALID",
				"",
				"random_type",
			}

			for _, ct := range invalidTypes {
				Expect(ct.IsValid()).To(BeFalse(), "ColumnType %s should be invalid", ct)
			}
		})
	})
})

var _ = Describe("Column", func() {
	Context("Validate", func() {
		It("should validate column with valid data", func() {
			col := schema.Column{
				Name:       "user_id",
				Type:       schema.ColumnTypeInteger,
				Nullable:   false,
				PrimaryKey: true,
				Unique:     true,
			}

			err := col.Validate()
			Expect(err).NotTo(HaveOccurred())
		})

		It("should reject empty column name", func() {
			col := schema.Column{
				Name: "",
				Type: schema.ColumnTypeString,
			}

			err := col.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("empty"))
		})

		It("should reject whitespace-only column name", func() {
			col := schema.Column{
				Name: "   ",
				Type: schema.ColumnTypeString,
			}

			err := col.Validate()
			Expect(err).To(HaveOccurred())
		})

		It("should reject invalid column type", func() {
			col := schema.Column{
				Name: "test_col",
				Type: "invalid_type",
			}

			err := col.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("invalid type"))
		})
	})
})

var _ = Describe("Table", func() {
	Context("Validate", func() {
		It("should validate table with valid data", func() {
			table := schema.Table{
				Name: "users",
				Columns: []schema.Column{
					{
						Name:       "id",
						Type:       schema.ColumnTypeInteger,
						PrimaryKey: true,
					},
					{
						Name: "email",
						Type: schema.ColumnTypeString,
					},
				},
			}

			err := table.Validate()
			Expect(err).NotTo(HaveOccurred())
		})

		It("should reject empty table name", func() {
			table := schema.Table{
				Name: "",
				Columns: []schema.Column{
					{Name: "id", Type: schema.ColumnTypeInteger},
				},
			}

			err := table.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("empty"))
		})

		It("should reject table with no columns", func() {
			table := schema.Table{
				Name:    "users",
				Columns: []schema.Column{},
			}

			err := table.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("at least one column"))
		})

		It("should reject table with invalid column", func() {
			table := schema.Table{
				Name: "users",
				Columns: []schema.Column{
					{Name: "", Type: schema.ColumnTypeString},
				},
			}

			err := table.Validate()
			Expect(err).To(HaveOccurred())
		})
	})
})
