// Package domain provides domain models for business logic
package domain

import (
	"github.com/LarsArtmann/SQLC-Wizzard/generated"
)

// EmitOptions is a type alias for the generated EmitOptions
// This ensures we use the single source of truth from the generated types
//
// DEPRECATED: This is the OLD boolean-heavy version. For new code, use TypeSafeEmitOptions
// from emit_modes.go which provides proper type safety with semantic enums.
//
// Migration path:
//
//	Old: opts := domain.DefaultEmitOptions()
//	New: opts := domain.NewTypeSafeEmitOptions()
type EmitOptions = generated.EmitOptions

// SafetyRules represents CEL-based validation rules
//
// DEPRECATED: This is the OLD boolean-heavy version. For new code, use TypeSafeSafetyRules
// from safety_policy.go which provides proper type safety with semantic groupings.
//
// Migration path:
//
//	Old: rules := domain.DefaultSafetyRules()
//	New: rules := domain.NewTypeSafeSafetyRules()
type SafetyRules = generated.SafetyRules

// DefaultEmitOptions returns safe defaults for code generation
//
// DEPRECATED: Use NewTypeSafeEmitOptions() instead for new code.
// This function is kept for backward compatibility with existing code that uses
// the boolean-heavy EmitOptions structure. The new TypeSafeEmitOptions provides
// proper type safety with semantic enums that prevent invalid state combinations.
//
// NOTE: This is also defined in generated/types.go (line 139). That version is
// the canonical source for the OLD boolean-based structure. This wrapper exists
// to provide a domain-layer access point for backward compatibility.
//
// For new code, use:
//
//	opts := domain.NewTypeSafeEmitOptions()
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
//
// DEPRECATED: Use NewTypeSafeSafetyRules() instead for new code.
// This function is kept for backward compatibility with existing code that uses
// the boolean-heavy SafetyRules structure. The new TypeSafeSafetyRules provides
// proper type safety with semantic groupings and enums.
//
// NOTE: This is also defined in generated/types.go (line 154). That version is
// generated and should be considered the canonical source for the
// OLD boolean-based structure. This wrapper exists to provide a domain-layer
// access point for backward compatibility.
//
// For new code, use:
//
//	rules := domain.NewTypeSafeSafetyRules()           // Production defaults
//	rules := domain.NewDevelopmentSafetyRules()       // Dev-friendly
//	rules := domain.NewProductionSafetyRules()        // Extra strict
func DefaultSafetyRules() SafetyRules {
	return SafetyRules{
		NoSelectStar: true,
		RequireWhere: true,
		RequireLimit: false, // Not too restrictive by default
		Rules:        []generated.SafetyRule{},
	}
}
