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
// Deprecated: This is the OLD boolean-heavy version. For new code, use TypeSafeSafetyRules
// from safety_policy.go which provides proper type safety with semantic groupings.
//
// Migration path:
//
//	Old: rules := domain.DefaultSafetyRules()
//	New: rules := domain.NewTypeSafeSafetyRules()
type SafetyRules = generated.SafetyRules
