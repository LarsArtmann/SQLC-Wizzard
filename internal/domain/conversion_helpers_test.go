package domain_test

import (
	"github.com/LarsArtmann/SQLC-Wizzard/generated"
)

// baseEmitOptions creates base EmitOptions with minimal settings.
func baseEmitOptions() generated.EmitOptions {
	return generated.EmitOptions{
		EmitJSONTags:             false,
		EmitPreparedQueries:      false,
		EmitInterface:            false,
		EmitEmptySlices:          false,
		EmitResultStructPointers: false,
		EmitParamsStructPointers: false,
		EmitEnumValidMethod:      false,
		EmitAllEnumValues:        false,
		JSONTagsCaseStyle:        "camel",
	}
}

// commonEmitOptions creates EmitOptions with common features for complete enum mode.
func commonEmitOptions() generated.EmitOptions {
	return generated.EmitOptions{
		EmitJSONTags:             true,
		EmitPreparedQueries:      true,
		EmitInterface:            true,
		EmitEmptySlices:          false,
		EmitResultStructPointers: false,
		EmitParamsStructPointers: false,
		EmitEnumValidMethod:      true,
		EmitAllEnumValues:        true,
		JSONTagsCaseStyle:        "camel",
	}
}

// emitOptionsEmptySlices creates EmitOptions for empty_slices mode.
func emitOptionsEmptySlices() generated.EmitOptions {
	opts := commonEmitOptions()
	opts.EmitEmptySlices = true

	return opts
}

// emitOptionsPointers creates EmitOptions for pointers mode.
func emitOptionsPointers() generated.EmitOptions {
	opts := baseEmitOptions()
	opts.EmitJSONTags = true
	opts.EmitResultStructPointers = true
	opts.EmitParamsStructPointers = true
	opts.JSONTagsCaseStyle = "snake"

	return opts
}
