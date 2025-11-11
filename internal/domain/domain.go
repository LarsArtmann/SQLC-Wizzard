// Package domain provides domain models for business logic
package domain

import (
	"github.com/LarsArtmann/SQLC-Wizzard/generated"
)

// EmitOptions is a type alias for the generated EmitOptions
// This ensures we use the single source of truth from TypeSpec generation
type EmitOptions = generated.EmitOptions

// SafetyRules represents CEL-based validation rules
// DEPRECATED: Use generated.SafetyRules directly
// This type is kept for backward compatibility only
type SafetyRules = generated.SafetyRules

// DefaultEmitOptions returns safe defaults for code generation
func DefaultEmitOptions() EmitOptions {
	return EmitOptions{
		EmitJSONTags:             true,
		EmitPreparedQueries:      true,
		EmitInterface:            true,
		EmitEmptySlices:          true,
		EmitResultStructPointers: false,
		EmitParamsStructPointers: false,
		EmitEnumValidMethod:      true,
		EmitAllEnumValues:        true,
		JSONTagsCaseStyle:        "camel",
	}
}

// DefaultSafetyRules returns safe defaults for query validation
func DefaultSafetyRules() SafetyRules {
	return SafetyRules{
		NoSelectStar: true,
		RequireWhere: true,
		RequireLimit: false, // Not too restrictive by default
		Rules:        []generated.SafetyRule{},
	}
}
