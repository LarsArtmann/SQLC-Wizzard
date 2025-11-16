package schema

import (
	"fmt"
	"strings"
)

// Schema represents a database schema with proper type safety
// Eliminates interface{} returns from adapter methods
type Schema struct {
	Name        string             `json:"name"`
	Tables      []Table           `json:"tables"`
	Views       []View            `json:"views,omitempty"`
	Indexes     []Index           `json:"indexes,omitempty"`
	Constraints []Constraint      `json:"constraints,omitempty"`
	Enums       []Enum            `json:"enums,omitempty"`
	Metadata    SchemaMetadata    `json:"metadata"`
}

// Table represents a database table with typed columns
type Table struct {
	Name       string    `json:"name"`
	Columns    []Column  `json:"columns"`
	PrimaryKey *Index    `json:"primary_key,omitempty"`
	Indexes    []Index   `json:"indexes,omitempty"`
	ForeignKeys []ForeignKey `json:"foreign_keys,omitempty"`
}

// Column represents a database column with type safety
type Column struct {
	Name       string      `json:"name"`
	Type       ColumnType  `json:"type"`
	Nullable   bool        `json:"nullable"`
	Default    *string     `json:"default,omitempty"`
	PrimaryKey bool        `json:"primary_key"`
	Unique     bool        `json:"unique"`
}

// ColumnType represents strongly-typed column types
type ColumnType string

const (
	ColumnTypeString    ColumnType = "string"
	ColumnTypeInteger   ColumnType = "integer"
	ColumnTypeBigInt    ColumnType = "bigint"
	ColumnTypeFloat     ColumnType = "float"
	ColumnTypeDouble    ColumnType = "double"
	ColumnTypeBoolean   ColumnType = "boolean"
	ColumnTypeDate      ColumnType = "date"
	ColumnTypeDateTime  ColumnType = "datetime"
	ColumnTypeTimestamp ColumnType = "timestamp"
	ColumnTypeJSON      ColumnType = "json"
	ColumnTypeUUID      ColumnType = "uuid"
	ColumnTypeText      ColumnType = "text"
	ColumnTypeBlob      ColumnType = "blob"
)

// IsValid returns true if ColumnType is valid
func (ct ColumnType) IsValid() bool {
	switch ct {
	case ColumnTypeString, ColumnTypeInteger, ColumnTypeBigInt, ColumnTypeFloat,
		 ColumnTypeDouble, ColumnTypeBoolean, ColumnTypeDate, ColumnTypeDateTime,
		 ColumnTypeTimestamp, ColumnTypeJSON, ColumnTypeUUID, ColumnTypeText, ColumnTypeBlob:
		return true
	default:
		return false
	}
}

// View represents a database view
type View struct {
	Name       string `json:"name"`
	Definition string `json:"definition"`
	Columns    []Column `json:"columns"`
}

// Index represents a database index
type Index struct {
	Name    string   `json:"name"`
	Table   string   `json:"table"`
	Columns []string `json:"columns"`
	Unique  bool     `json:"unique"`
	Type    IndexType `json:"type"`
}

// IndexType represents index types
type IndexType string

const (
	IndexTypeBTree    IndexType = "btree"
	IndexTypeHash     IndexType = "hash"
	IndexTypeGin      IndexType = "gin"
	IndexTypeGiST     IndexType = "gist"
	IndexTypeBRIN     IndexType = "brin"
)

// ForeignKey represents a foreign key constraint
type ForeignKey struct {
	Name           string `json:"name"`
	SourceTable    string `json:"source_table"`
	SourceColumn   string `json:"source_column"`
	TargetTable    string `json:"target_table"`
	TargetColumn   string `json:"target_column"`
	OnDeleteAction string `json:"on_delete_action"`
	OnUpdateAction string `json:"on_update_action"`
}

// Constraint represents a database constraint
type Constraint struct {
	Name        string           `json:"name"`
	Type        ConstraintType   `json:"type"`
	Table       string          `json:"table"`
	Columns     []string        `json:"columns"`
	Definition  string          `json:"definition,omitempty"`
	Check       *CheckConstraint `json:"check,omitempty"`
}

// ConstraintType represents constraint types
type ConstraintType string

const (
	ConstraintTypePrimaryKey ConstraintType = "primary_key"
	ConstraintTypeForeignKey ConstraintType = "foreign_key"
	ConstraintTypeUnique     ConstraintType = "unique"
	ConstraintTypeCheck      ConstraintType = "check"
	ConstraintTypeNotNull    ConstraintType = "not_null"
)

// CheckConstraint represents check constraint details
type CheckConstraint struct {
	Condition string `json:"condition"`
}

// Enum represents an enum type
type Enum struct {
	Name   string   `json:"name"`
	Values []string `json:"values"`
}

// SchemaMetadata contains metadata about the schema
type SchemaMetadata struct {
	DatabaseEngine string    `json:"database_engine"`
	Version       string    `json:"version"`
	CreatedAt     string    `json:"created_at"`
	UpdatedAt     string    `json:"updated_at"`
}

// NewSchema creates a new Schema with validation
func NewSchema(name string, tables []Table) (*Schema, error) {
	if strings.TrimSpace(name) == "" {
		return nil, &SchemaError{
			Code:    "EMPTY_SCHEMA_NAME",
			Message: "Schema name cannot be empty",
		}
	}
	
	if len(tables) == 0 {
		return nil, &SchemaError{
			Code:    "NO_TABLES",
			Message: "Schema must contain at least one table",
		}
	}
	
	if len(tables) > 1000 {
		return nil, &SchemaError{
			Code:    "TOO_MANY_TABLES",
			Message: "Schema exceeds reasonable limit of 1000 tables",
		}
	}
	
	// Validate each table
	for i, table := range tables {
		if strings.TrimSpace(table.Name) == "" {
			return nil, &SchemaError{
				Code:    fmt.Sprintf("EMPTY_TABLE_NAME_%d", i),
				Message: fmt.Sprintf("Table at index %d has empty name", i),
			}
		}
		
		if len(table.Columns) == 0 {
			return nil, &SchemaError{
				Code:    fmt.Sprintf("NO_COLUMNS_%s", table.Name),
				Message: fmt.Sprintf("Table %s must contain at least one column", table.Name),
			}
		}
	}
	
	return &Schema{
		Name:     name,
		Tables:   tables,
		Metadata: SchemaMetadata{
			DatabaseEngine: "unknown",
			Version:       "1.0.0",
		},
	}, nil
}

// GetTable returns a table by name
func (s *Schema) GetTable(name string) (*Table, bool) {
	for _, table := range s.Tables {
		if table.Name == name {
			return &table, true
		}
	}
	return nil, false
}

// GetColumn returns a column from a specific table
func (s *Schema) GetColumn(tableName, columnName string) (*Column, bool) {
	table, found := s.GetTable(tableName)
	if !found {
		return nil, false
	}
	
	for _, column := range table.Columns {
		if column.Name == columnName {
			return &column, true
		}
	}
	return nil, false
}

// SchemaError represents schema-specific errors
type SchemaError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *SchemaError) Error() string {
	return e.Message
}

// Validate validates the entire schema
func (s *Schema) Validate() error {
	if s == nil {
		return &SchemaError{
			Code:    "NULL_SCHEMA",
			Message: "Schema cannot be null",
		}
	}
	
	if strings.TrimSpace(s.Name) == "" {
		return &SchemaError{
			Code:    "EMPTY_SCHEMA_NAME",
			Message: "Schema name cannot be empty",
		}
	}
	
	if len(s.Tables) == 0 {
		return &SchemaError{
			Code:    "NO_TABLES",
			Message: "Schema must contain at least one table",
		}
	}
	
	if len(s.Tables) > 1000 {
		return &SchemaError{
			Code:    "TOO_MANY_TABLES",
			Message: "Schema exceeds reasonable limit of 1000 tables",
		}
	}
	
	// Validate each table
	for i, table := range s.Tables {
		if err := table.Validate(); err != nil {
			return &SchemaError{
				Code:    fmt.Sprintf("TABLE_VALIDATION_%d_%s", i, table.Name),
				Message: fmt.Sprintf("Table %s validation failed: %s", table.Name, err.Error()),
			}
		}
	}
	
	return nil
}

// Validate validates a table
func (t *Table) Validate() error {
	if strings.TrimSpace(t.Name) == "" {
		return &SchemaError{
			Code:    "EMPTY_TABLE_NAME",
			Message: "Table name cannot be empty",
		}
	}
	
	if len(t.Columns) == 0 {
		return &SchemaError{
			Code:    "NO_COLUMNS",
			Message: fmt.Sprintf("Table %s must contain at least one column", t.Name),
		}
	}
	
	// Validate each column
	for i, column := range t.Columns {
		if err := column.Validate(); err != nil {
			return &SchemaError{
				Code:    fmt.Sprintf("COLUMN_VALIDATION_%d_%s", i, column.Name),
				Message: fmt.Sprintf("Column %s validation failed: %s", column.Name, err.Error()),
			}
		}
	}
	
	return nil
}

// Validate validates a column
func (c *Column) Validate() error {
	if strings.TrimSpace(c.Name) == "" {
		return &SchemaError{
			Code:    "EMPTY_COLUMN_NAME",
			Message: "Column name cannot be empty",
		}
	}
	
	if !c.Type.IsValid() {
		return &SchemaError{
			Code:    "INVALID_COLUMN_TYPE",
			Message: fmt.Sprintf("Column %s has invalid type: %s", c.Name, string(c.Type)),
		}
	}
	
	return nil
}