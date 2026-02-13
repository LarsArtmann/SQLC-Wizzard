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

var _ = Describe("Schema", func() {
	Context("NewSchema", func() {
		It("should create schema with valid data", func() {
			tables := []schema.Table{
				{
					Name: "users",
					Columns: []schema.Column{
						{Name: "id", Type: schema.ColumnTypeInteger},
					},
				},
			}

			s, err := schema.NewSchema("mydb", tables)

			Expect(err).NotTo(HaveOccurred())
			Expect(s).NotTo(BeNil())
			Expect(s.Name).To(Equal("mydb"))
			Expect(s.Tables).To(HaveLen(1))
			Expect(s.Metadata.DatabaseEngine).To(Equal("unknown"))
			Expect(s.Metadata.Version).To(Equal("1.0.0"))
		})

		It("should reject empty schema name", func() {
			tables := []schema.Table{
				{
					Name: "users",
					Columns: []schema.Column{
						{Name: "id", Type: schema.ColumnTypeInteger},
					},
				},
			}

			s, err := schema.NewSchema("", tables)

			Expect(err).To(HaveOccurred())
			Expect(s).To(BeNil())
			Expect(err.Error()).To(ContainSubstring("empty"))
		})

		It("should reject schema with no tables", func() {
			s, err := schema.NewSchema("mydb", []schema.Table{})

			Expect(err).To(HaveOccurred())
			Expect(s).To(BeNil())
			Expect(err.Error()).To(ContainSubstring("at least one table"))
		})

		It("should reject schema with too many tables", func() {
			tables := make([]schema.Table, 1001)
			for i := range tables {
				tables[i] = schema.Table{
					Name: "table",
					Columns: []schema.Column{
						{Name: "id", Type: schema.ColumnTypeInteger},
					},
				}
			}

			s, err := schema.NewSchema("mydb", tables)

			Expect(err).To(HaveOccurred())
			Expect(s).To(BeNil())
			Expect(err.Error()).To(ContainSubstring("exceeds"))
			Expect(err.Error()).To(ContainSubstring("1000"))
		})

		It("should reject schema with empty table name", func() {
			tables := []schema.Table{
				{
					Name: "",
					Columns: []schema.Column{
						{Name: "id", Type: schema.ColumnTypeInteger},
					},
				},
			}

			s, err := schema.NewSchema("mydb", tables)

			Expect(err).To(HaveOccurred())
			Expect(s).To(BeNil())
		})

		It("should reject schema with table without columns", func() {
			tables := []schema.Table{
				{
					Name:    "empty_table",
					Columns: []schema.Column{},
				},
			}

			s, err := schema.NewSchema("mydb", tables)

			Expect(err).To(HaveOccurred())
			Expect(s).To(BeNil())
		})
	})

	Context("Validate", func() {
		var (
			testTableName   = "users"
			testColumnName  = "id"
			createTestTable = func(tableName, columnName string) schema.Table {
				return schema.Table{
					Name: tableName,
					Columns: []schema.Column{
						{Name: columnName, Type: schema.ColumnTypeInteger},
					},
				}
			}
			createSchemaWithTables = func(name string, tables []schema.Table) *schema.Schema {
				return &schema.Schema{
					Name:   name,
					Tables: tables,
				}
			}
		)

		It("should validate valid schema", func() {
			tables := []schema.Table{
				createTestTable(testTableName, testColumnName),
			}

			s, err := schema.NewSchema("mydb", tables)
			Expect(err).NotTo(HaveOccurred())

			err = s.Validate()
			Expect(err).NotTo(HaveOccurred())
		})

		It("should reject nil schema", func() {
			var s *schema.Schema
			err := s.Validate()

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("null"))
		})

		It("should reject schema with empty name", func() {
			tables := []schema.Table{
				createTestTable(testTableName, testColumnName),
			}
			s := createSchemaWithTables("", tables)

			err := s.Validate()
			Expect(err).To(HaveOccurred())
		})

		It("should reject schema with no tables", func() {
			s := createSchemaWithTables("mydb", []schema.Table{})

			err := s.Validate()
			Expect(err).To(HaveOccurred())
		})

		createSchemaWithInvalidTable := func(name string) *schema.Schema {
			return createSchemaWithTables(name, []schema.Table{
				{Name: "", Columns: []schema.Column{{Name: "id", Type: schema.ColumnTypeInteger}}},
			})
		}

		It("should reject schema with invalid table", func() {
			s := createSchemaWithInvalidTable("mydb")

			err := s.Validate()
			Expect(err).To(HaveOccurred())
		})
	})

	Context("GetTable", func() {
		var s *schema.Schema

		BeforeEach(func() {
			tables := []schema.Table{
				{
					Name: "users",
					Columns: []schema.Column{
						{Name: "id", Type: schema.ColumnTypeInteger},
					},
				},
				{
					Name: "posts",
					Columns: []schema.Column{
						{Name: "id", Type: schema.ColumnTypeInteger},
					},
				},
			}

			var err error
			s, err = schema.NewSchema("mydb", tables)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should find existing table", func() {
			table, found := s.GetTable("users")

			Expect(found).To(BeTrue())
			Expect(table).NotTo(BeNil())
			Expect(table.Name).To(Equal("users"))
		})

		It("should return false for non-existent table", func() {
			table, found := s.GetTable("non_existent")

			Expect(found).To(BeFalse())
			Expect(table).To(BeNil())
		})
	})

	Context("GetColumn", func() {
		var s *schema.Schema

		BeforeEach(func() {
			tables := []schema.Table{
				{
					Name: "users",
					Columns: []schema.Column{
						{Name: "id", Type: schema.ColumnTypeInteger},
						{Name: "email", Type: schema.ColumnTypeString},
						{Name: "created_at", Type: schema.ColumnTypeTimestamp},
					},
				},
			}

			var err error
			s, err = schema.NewSchema("mydb", tables)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should find existing column in existing table", func() {
			col, found := s.GetColumn("users", "email")

			Expect(found).To(BeTrue())
			Expect(col).NotTo(BeNil())
			Expect(col.Name).To(Equal("email"))
			Expect(col.Type).To(Equal(schema.ColumnTypeString))
		})

		It("should return false for non-existent table", func() {
			col, found := s.GetColumn("non_existent", "email")

			Expect(found).To(BeFalse())
			Expect(col).To(BeNil())
		})

		It("should return false for non-existent column", func() {
			col, found := s.GetColumn("users", "non_existent")

			Expect(found).To(BeFalse())
			Expect(col).To(BeNil())
		})
	})
})

var _ = Describe("SchemaError", func() {
	It("should implement error interface", func() {
		err := &schema.SchemaError{
			Code:    "TEST_ERROR",
			Message: "This is a test error",
		}

		Expect(err.Error()).To(Equal("This is a test error"))
	})
})

var _ = Describe("Complex Schema", func() {
	It("should handle schema with multiple tables and complex types", func() {
		tables := []schema.Table{
			{
				Name: "users",
				Columns: []schema.Column{
					{Name: "id", Type: schema.ColumnTypeUUID, PrimaryKey: true},
					{Name: "email", Type: schema.ColumnTypeString, Unique: true},
					{Name: "name", Type: schema.ColumnTypeString, Nullable: true},
					{Name: "age", Type: schema.ColumnTypeInteger, Nullable: true},
					{Name: "created_at", Type: schema.ColumnTypeTimestamp},
					{Name: "metadata", Type: schema.ColumnTypeJSON, Nullable: true},
				},
			},
			{
				Name: "posts",
				Columns: []schema.Column{
					{Name: "id", Type: schema.ColumnTypeInteger, PrimaryKey: true},
					{Name: "user_id", Type: schema.ColumnTypeUUID},
					{Name: "title", Type: schema.ColumnTypeString},
					{Name: "content", Type: schema.ColumnTypeText},
					{Name: "published", Type: schema.ColumnTypeBoolean},
					{Name: "published_at", Type: schema.ColumnTypeDateTime, Nullable: true},
				},
			},
		}

		s, err := schema.NewSchema("blog_db", tables)

		Expect(err).NotTo(HaveOccurred())
		Expect(s).NotTo(BeNil())
		Expect(s.Tables).To(HaveLen(2))

		err = s.Validate()
		Expect(err).NotTo(HaveOccurred())

		// Verify we can find tables and columns
		usersTable, found := s.GetTable("users")
		Expect(found).To(BeTrue())
		Expect(usersTable.Columns).To(HaveLen(6))

		emailCol, found := s.GetColumn("users", "email")
		Expect(found).To(BeTrue())
		Expect(emailCol.Unique).To(BeTrue())

		postsTable, found := s.GetTable("posts")
		Expect(found).To(BeTrue())
		Expect(postsTable.Columns).To(HaveLen(6))
	})
})
