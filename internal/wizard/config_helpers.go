package wizard

import (
	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/pkg/config"
)

// ApplyEmitOptions applies emit options to a GoGenConfig in a type-safe manner
// This is a type-safe operation that eliminates field-by-field copying (DRY principle)
func ApplyEmitOptions(opts *generated.EmitOptions, cfg *config.GoGenConfig) {
	if opts.EmitJSONTags {
		cfg.EmitJSONTags = opts.EmitJSONTags
	}
	if opts.EmitPreparedQueries {
		cfg.EmitPreparedQueries = opts.EmitPreparedQueries
	}
	if opts.EmitInterface {
		cfg.EmitInterface = opts.EmitInterface
	}
	if opts.EmitEmptySlices {
		cfg.EmitEmptySlices = opts.EmitEmptySlices
	}
	if opts.EmitResultStructPointers {
		cfg.EmitResultStructPointers = opts.EmitResultStructPointers
	}
	if opts.EmitParamsStructPointers {
		cfg.EmitParamsStructPointers = opts.EmitParamsStructPointers
	}
	if opts.EmitEnumValidMethod {
		cfg.EmitEnumValidMethod = opts.EmitEnumValidMethod
	}
	if opts.EmitAllEnumValues {
		cfg.EmitAllEnumValues = opts.EmitAllEnumValues
	}
	if opts.JSONTagsCaseStyle != "" {
		cfg.JSONTagsCaseStyle = opts.JSONTagsCaseStyle
	}
}