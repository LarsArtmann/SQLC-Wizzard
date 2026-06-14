package testing

import (
	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/domain"
)

// GetNullHandlingModeTestCases returns test cases for NullHandlingMode string representation tests.
// This helper eliminates duplicate test case definitions across test files.
func GetNullHandlingModeTestCases() []EnumTestCase {
	return []EnumTestCase{
		{EnumValue: domain.NullHandlingPointers, ExpectedString: "pointers"},
		{EnumValue: domain.NullHandlingEmptySlices, ExpectedString: "empty_slices"},
		{EnumValue: domain.NullHandlingExplicitNull, ExpectedString: "explicit_null"},
		{EnumValue: domain.NullHandlingMixed, ExpectedString: "mixed"},
	}
}

// GetStructPointerModeTestCases returns test cases for StructPointerMode string representation tests.
// This helper eliminates duplicate test case definitions across test files.
func GetStructPointerModeTestCases() []EnumTestCase {
	return []EnumTestCase{
		{EnumValue: domain.StructPointerNever, ExpectedString: "never"},
		{EnumValue: domain.StructPointerResults, ExpectedString: "results"},
		{EnumValue: domain.StructPointerParams, ExpectedString: "params"},
		{EnumValue: domain.StructPointerAlways, ExpectedString: "always"},
	}
}

// GetJSONTagStyleTestCases returns test cases for JSONTagStyle string representation tests.
// This helper eliminates duplicate test case definitions across test files.
func GetJSONTagStyleTestCases() []EnumTestCase {
	return []EnumTestCase{
		{EnumValue: domain.JSONTagStyleCamel, ExpectedString: "camel"},
		{EnumValue: domain.JSONTagStyleSnake, ExpectedString: "snake"},
		{EnumValue: domain.JSONTagStylePascal, ExpectedString: "pascal"},
		{EnumValue: domain.JSONTagStyleKebab, ExpectedString: "kebab"},
	}
}

// CreateDefaultValidationConfig creates a ValidationConfig with production-ready defaults.
// This helper eliminates duplicate fixture code across test files.
func CreateDefaultValidationConfig() generated.ValidationConfig {
	return generated.ValidationConfig{
		StrictFunctions: true,
		StrictOrderBy:   true,
		EmitOptions: generated.EmitOptions{
			EmitJSONTags:             true,
			EmitPreparedQueries:      true,
			EmitInterface:            true,
			EmitEmptySlices:          true,
			EmitResultStructPointers: false,
			EmitParamsStructPointers: false,
			EmitEnumValidMethod:      true,
			EmitAllEnumValues:        false,
		},
		SafetyRules: generated.SafetyRules{
			NoSelectStar: true,
			RequireWhere: true,
			RequireLimit: false,
			NoDropTable:  true,
		},
	}
}
