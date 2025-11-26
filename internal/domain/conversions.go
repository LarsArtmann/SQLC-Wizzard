package domain

import "github.com/LarsArtmann/SQLC-Wizzard/generated"

// This file provides bidirectional conversions between OLD boolean-heavy types
// and NEW type-safe enum-based types. These conversions enable gradual migration
// from the legacy generated types to the new type-safe domain types.

// ============================================================================
// EmitOptions Conversions
// ============================================================================

// ToTypeSafe converts old generated.EmitOptions to new TypeSafeEmitOptions
// This is the bridge from legacy boolean-heavy structure to type-safe enums
func EmitOptionsToTypeSafe(old generated.EmitOptions) TypeSafeEmitOptions {
	// Determine NullHandlingMode from combination of boolean flags
	var nullHandling NullHandlingMode
	if old.EmitEmptySlices && !old.EmitResultStructPointers && !old.EmitParamsStructPointers {
		nullHandling = NullHandlingEmptySlices
	} else if !old.EmitEmptySlices && old.EmitResultStructPointers && old.EmitParamsStructPointers {
		nullHandling = NullHandlingPointers
	} else if !old.EmitEmptySlices && !old.EmitResultStructPointers && !old.EmitParamsStructPointers {
		nullHandling = NullHandlingExplicitNull
	} else {
		// Mixed mode: some pointers, some not
		nullHandling = NullHandlingMixed
	}

	// Determine EnumGenerationMode from boolean flags
	var enumMode EnumGenerationMode
	if old.EmitEnumValidMethod && old.EmitAllEnumValues {
		enumMode = EnumGenerationComplete
	} else if old.EmitEnumValidMethod && !old.EmitAllEnumValues {
		enumMode = EnumGenerationWithValidation
	} else {
		enumMode = EnumGenerationBasic
	}

	// Determine StructPointerMode from boolean flags
	var structPointers StructPointerMode
	if old.EmitResultStructPointers && old.EmitParamsStructPointers {
		structPointers = StructPointerAlways
	} else if old.EmitResultStructPointers && !old.EmitParamsStructPointers {
		structPointers = StructPointerResults
	} else if !old.EmitResultStructPointers && old.EmitParamsStructPointers {
		structPointers = StructPointerParams
	} else {
		structPointers = StructPointerNever
	}

	// Parse JSON tag style (with fallback to camel)
	jsonStyle := ParseJSONTagStyle(old.JSONTagsCaseStyle)
	if !jsonStyle.IsValid() {
		jsonStyle = JSONTagStyleCamel
	}

	return TypeSafeEmitOptions{
		NullHandling:   nullHandling,
		EnumMode:       enumMode,
		StructPointers: structPointers,
		JSONTagStyle:   jsonStyle,
		Features: CodeGenerationFeatures{
			GenerateJSONTags:        old.EmitJSONTags,
			GeneratePreparedQueries: old.EmitPreparedQueries,
			GenerateInterface:       old.EmitInterface,
			UseExactTableNames:      false, // Not present in old structure
		},
	}
}

// DEPRECATED: ToLegacy moved to emit_modes.go as ToTemplateData()
// This eliminates split brain and provides single source of truth for conversion
// Use: opts.ToTemplateData() instead of opts.ToLegacy()

// ParseJSONTagStyle converts string to JSONTagStyle enum
func ParseJSONTagStyle(s string) JSONTagStyle {
	switch s {
	case "camel":
		return JSONTagStyleCamel
	case "snake":
		return JSONTagStyleSnake
	case "pascal":
		return JSONTagStylePascal
	case "kebab":
		return JSONTagStyleKebab
	default:
		return JSONTagStyle(s) // Return as-is for validation
	}
}

// ============================================================================
// SafetyRules Conversions
// ============================================================================

// SafetyRulesToTypeSafe converts old generated.SafetyRules to new TypeSafeSafetyRules
// This is the bridge from legacy boolean-heavy structure to type-safe semantic groupings
func SafetyRulesToTypeSafe(old generated.SafetyRules) TypeSafeSafetyRules {
	// Determine DestructiveOperationPolicy from boolean flags
	var destructiveOps DestructiveOperationPolicy
	if !old.NoDropTable && !old.NoTruncate {
		// Both allowed
		destructiveOps = DestructiveAllowed
	} else if old.NoDropTable && old.NoTruncate {
		// Both forbidden
		destructiveOps = DestructiveForbidden
	} else {
		// Mixed state - map to forbidden for safety
		destructiveOps = DestructiveForbidden
	}

	return TypeSafeSafetyRules{
		StyleRules: QueryStyleRules{
			SelectStarPolicy:   convertNoSelectStarPolicy(old.NoSelectStar),
			ColumnExplicitness: ColumnExplicitnessDefault, // Not present in old structure
		},
		SafetyRules: QuerySafetyRules{
			WhereRequirement: convertWhereRequirement(old.RequireWhere),
			LimitRequirement: convertLimitRequirement(old.RequireLimit),
			// MaxRowsWithoutLimit not present in old structure, use safe default
			MaxRowsWithoutLimit: 1000,
		},
		DestructiveOps: destructiveOps,
		CustomRules:    old.Rules,
	}
}

// convertNoSelectStarPolicy converts old boolean to new SelectStarPolicy
func convertNoSelectStarPolicy(noSelectStar bool) SelectStarPolicy {
	if noSelectStar {
		return SelectStarForbidden
	}
	return SelectStarAllowed
}

// convertWhereRequirement converts old boolean to new WhereClauseRequirement
func convertWhereRequirement(requireWhere bool) WhereClauseRequirement {
	if requireWhere {
		return WhereClauseOnDestructive
	}
	return WhereClauseNever
}

// convertLimitRequirement converts old boolean to new LimitClauseRequirement
func convertLimitRequirement(requireLimit bool) LimitClauseRequirement {
	if requireLimit {
		return LimitClauseOnSelect
	}
	return LimitClauseNever
}

// Legacy Conversion - DEPRECATED
// This method exists for backward compatibility but should not be used in new code
// Use: opts.ToTemplateData() instead of opts.ToLegacy()
func (rules TypeSafeSafetyRules) ToLegacy() generated.SafetyRules {
	return generated.SafetyRules{
		NoSelectStar: rules.StyleRules.SelectStarPolicy.ForbidsSelectStar(),
		RequireWhere: rules.SafetyRules.WhereRequirement.RequiresOnDestructive(),
		NoDropTable:  !rules.DestructiveOps.AllowsDropTable(),
		NoTruncate:   !rules.DestructiveOps.AllowsTruncate(),
		RequireLimit: rules.SafetyRules.LimitRequirement.RequiresOnSelect(),
		Rules:        rules.CustomRules,
	}
}

// ============================================================================
// Convenience Constructors with Conversions
// ============================================================================

// Convenience Constructor - DEPRECATED
// This exists for backward compatibility but should not be used in new code
func NewTypeSafeEmitOptionsFromLegacy(old generated.EmitOptions) TypeSafeEmitOptions {
	return EmitOptionsToTypeSafe(old)
}

// Convenience Constructor - DEPRECATED
// This exists for backward compatibility but should not be used in new code
func NewTypeSafeSafetyRulesFromLegacy(old generated.SafetyRules) TypeSafeSafetyRules {
	return SafetyRulesToTypeSafe(old)
}
