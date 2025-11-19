package domain

// This file contains type-safe enums for code generation options
// Replaces the boolean-heavy EmitOptions from generated/types.go with
// semantic groupings that prevent invalid state combinations

// NullHandlingMode defines how nullable database values are represented in generated code
type NullHandlingMode string

const (
	// NullHandlingPointers uses pointers for all nullable fields (*string, *int64, etc.)
	NullHandlingPointers NullHandlingMode = "pointers"

	// NullHandlingEmptySlices treats nil slices as empty slices (never nil)
	NullHandlingEmptySlices NullHandlingMode = "empty_slices"

	// NullHandlingExplicitNull uses sql.Null* types (sql.NullString, sql.NullInt64, etc.)
	NullHandlingExplicitNull NullHandlingMode = "explicit_null"

	// NullHandlingMixed allows pointers for some types, sql.Null* for others (advanced)
	NullHandlingMixed NullHandlingMode = "mixed"
)

// IsValid returns true if the null handling mode is recognized
func (n NullHandlingMode) IsValid() bool {
	switch n {
	case NullHandlingPointers, NullHandlingEmptySlices, NullHandlingExplicitNull, NullHandlingMixed:
		return true
	default:
		return false
	}
}

// String returns the string representation of the mode
func (n NullHandlingMode) String() string {
	return string(n)
}

// UsePointers returns true if this mode uses pointer types for nullability
func (n NullHandlingMode) UsePointers() bool {
	return n == NullHandlingPointers || n == NullHandlingMixed
}

// UseEmptySlices returns true if this mode prefers empty slices over nil
func (n NullHandlingMode) UseEmptySlices() bool {
	return n == NullHandlingEmptySlices
}

// EnumGenerationMode defines how database enums are generated in code
type EnumGenerationMode string

const (
	// EnumGenerationBasic generates simple enum constants
	EnumGenerationBasic EnumGenerationMode = "basic"

	// EnumGenerationWithValidation adds IsValid() methods to enum types
	EnumGenerationWithValidation EnumGenerationMode = "with_validation"

	// EnumGenerationComplete adds IsValid() + All() methods (all enum values)
	EnumGenerationComplete EnumGenerationMode = "complete"
)

// IsValid returns true if the enum generation mode is recognized
func (e EnumGenerationMode) IsValid() bool {
	switch e {
	case EnumGenerationBasic, EnumGenerationWithValidation, EnumGenerationComplete:
		return true
	default:
		return false
	}
}

// String returns the string representation of the mode
func (e EnumGenerationMode) String() string {
	return string(e)
}

// IncludesValidation returns true if this mode generates IsValid() methods
func (e EnumGenerationMode) IncludesValidation() bool {
	return e == EnumGenerationWithValidation || e == EnumGenerationComplete
}

// IncludesAllValues returns true if this mode generates All() methods
func (e EnumGenerationMode) IncludesAllValues() bool {
	return e == EnumGenerationComplete
}

// StructPointerMode defines when to use pointers for generated structs
type StructPointerMode string

const (
	// StructPointerNever never uses pointers for result/param structs
	StructPointerNever StructPointerMode = "never"

	// StructPointerResults uses pointers only for result structs
	StructPointerResults StructPointerMode = "results"

	// StructPointerParams uses pointers only for parameter structs
	StructPointerParams StructPointerMode = "params"

	// StructPointerAlways uses pointers for both result and parameter structs
	StructPointerAlways StructPointerMode = "always"
)

// IsValid returns true if the struct pointer mode is recognized
func (s StructPointerMode) IsValid() bool {
	switch s {
	case StructPointerNever, StructPointerResults, StructPointerParams, StructPointerAlways:
		return true
	default:
		return false
	}
}

// String returns the string representation of the mode
func (s StructPointerMode) String() string {
	return string(s)
}

// UseResultPointers returns true if result structs should be pointers
func (s StructPointerMode) UseResultPointers() bool {
	return s == StructPointerResults || s == StructPointerAlways
}

// UseParamPointers returns true if parameter structs should be pointers
func (s StructPointerMode) UseParamPointers() bool {
	return s == StructPointerParams || s == StructPointerAlways
}

// JSONTagStyle defines the case style for JSON tags
type JSONTagStyle string

const (
	// JSONTagStyleCamel uses camelCase (e.g., userId)
	JSONTagStyleCamel JSONTagStyle = "camel"

	// JSONTagStyleSnake uses snake_case (e.g., user_id)
	JSONTagStyleSnake JSONTagStyle = "snake"

	// JSONTagStylePascal uses PascalCase (e.g., UserId)
	JSONTagStylePascal JSONTagStyle = "pascal"

	// JSONTagStyleKebab uses kebab-case (e.g., user-id)
	JSONTagStyleKebab JSONTagStyle = "kebab"
)

// IsValid returns true if the JSON tag style is recognized
func (j JSONTagStyle) IsValid() bool {
	switch j {
	case JSONTagStyleCamel, JSONTagStyleSnake, JSONTagStylePascal, JSONTagStyleKebab:
		return true
	default:
		return false
	}
}

// String returns the string representation of the style
func (j JSONTagStyle) String() string {
	return string(j)
}

// CodeGenerationFeatures represents optional code generation features
// These are independent boolean flags that can be enabled/disabled
type CodeGenerationFeatures struct {
	// GenerateJSONTags adds json:"..." tags to struct fields
	GenerateJSONTags bool

	// GeneratePreparedQueries generates prepared statement methods
	GeneratePreparedQueries bool

	// GenerateInterface generates an interface for the Queries struct
	GenerateInterface bool

	// UseExactTableNames uses exact database table names (no pluralization)
	UseExactTableNames bool
}

// TypeSafeEmitOptions represents type-safe code generation configuration
// This is the NEW type-safe version that replaces the boolean-heavy generated.EmitOptions
// with semantic groupings that prevent invalid state combinations
//
// Migration path: Use this for new code, gradually migrate existing code from EmitOptions
type TypeSafeEmitOptions struct {
	// NullHandling defines how nullable values are represented
	NullHandling NullHandlingMode

	// EnumMode defines how enums are generated
	EnumMode EnumGenerationMode

	// StructPointers defines when to use pointers for structs
	StructPointers StructPointerMode

	// JSONTagStyle defines the case style for JSON tags
	JSONTagStyle JSONTagStyle

	// Features contains independent feature flags
	Features CodeGenerationFeatures
}

// DomainValidationError represents a validation error in the domain layer
type DomainValidationError struct {
	Field   string
	Message string
}

func (e *DomainValidationError) Error() string {
	return e.Field + ": " + e.Message
}

// IsValid validates that all configuration options are valid
func (e *TypeSafeEmitOptions) IsValid() error {
	if !e.NullHandling.IsValid() {
		return &DomainValidationError{
			Field:   "NullHandling",
			Message: "Invalid null handling mode: " + string(e.NullHandling),
		}
	}

	if !e.EnumMode.IsValid() {
		return &DomainValidationError{
			Field:   "EnumMode",
			Message: "Invalid enum generation mode: " + string(e.EnumMode),
		}
	}

	if !e.StructPointers.IsValid() {
		return &DomainValidationError{
			Field:   "StructPointers",
			Message: "Invalid struct pointer mode: " + string(e.StructPointers),
		}
	}

	if !e.JSONTagStyle.IsValid() {
		return &DomainValidationError{
			Field:   "JSONTagStyle",
			Message: "Invalid JSON tag style: " + string(e.JSONTagStyle),
		}
	}

	return nil
}

// NewTypeSafeEmitOptions returns production-ready defaults for code generation
// This creates the NEW type-safe version with semantic enums
func NewTypeSafeEmitOptions() TypeSafeEmitOptions {
	return TypeSafeEmitOptions{
		NullHandling:   NullHandlingPointers,
		EnumMode:       EnumGenerationComplete,
		StructPointers: StructPointerNever,
		JSONTagStyle:   JSONTagStyleCamel,
		Features: CodeGenerationFeatures{
			GenerateJSONTags:        true,
			GeneratePreparedQueries: true,
			GenerateInterface:       true,
			UseExactTableNames:      false,
		},
	}
}
