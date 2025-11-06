package domain

import "github.com/LarsArtmann/SQLC-Wizzard/pkg/config"

// EmitOptions represents code generation options.
// This is the SINGLE SOURCE OF TRUTH for code generation settings.
//
// Use ToGoGenConfigFields() to get the fields that should be set in config.GoGenConfig.
type EmitOptions struct {
	// Code generation flags
	EmitInterface      bool
	PreparedQueries    bool
	JSONTags           bool
	DBTags             bool
	ExactTableNames    bool
	EmptySlices        bool
	OmitUnusedStructs  bool
	ExportedQueries    bool
	PointersForNull    bool

	// JSON tag style
	JSONTagsCaseStyle string

	// File naming
	DBFileName       string
	ModelsFileName   string
	QuerierFileName  string
	CopyfromFileName string
	BatchFileName    string
}

// ApplyToGoGenConfig applies emit options to a GoGenConfig.
// This eliminates the split brain by having ONE method that sets all fields.
func (eo EmitOptions) ApplyToGoGenConfig(cfg *config.GoGenConfig) {
	cfg.EmitInterface = eo.EmitInterface
	cfg.EmitPreparedQueries = eo.PreparedQueries
	cfg.EmitJSONTags = eo.JSONTags
	cfg.EmitDBTags = eo.DBTags
	cfg.EmitExactTableNames = eo.ExactTableNames
	cfg.EmitEmptySlices = eo.EmptySlices
	cfg.OmitUnusedStructs = eo.OmitUnusedStructs
	cfg.EmitExportedQueries = eo.ExportedQueries
	cfg.EmitPointersForNullTypes = eo.PointersForNull
	cfg.JSONTagsCaseStyle = eo.JSONTagsCaseStyle
	cfg.OmitSQLCVersion = false // Always false for generated code attribution
	cfg.OutputDBFileName = eo.DBFileName
	cfg.OutputModelsFileName = eo.ModelsFileName
	cfg.OutputQuerierFileName = eo.QuerierFileName
	cfg.OutputCopyfromFileName = eo.CopyfromFileName
	cfg.OutputBatchFileName = eo.BatchFileName
}

// DefaultEmitOptions returns recommended emit options for production use.
func DefaultEmitOptions() EmitOptions {
	return EmitOptions{
		EmitInterface:      true,
		PreparedQueries:    true,
		JSONTags:           true,
		DBTags:             false,
		ExactTableNames:    true,
		EmptySlices:        true,
		OmitUnusedStructs:  true,
		ExportedQueries:    true,
		PointersForNull:    false,
		JSONTagsCaseStyle:  "camel",
		DBFileName:         "db.go",
		ModelsFileName:     "models.go",
		QuerierFileName:    "querier.go",
		CopyfromFileName:   "copyfrom.go",
		BatchFileName:      "batch.go",
	}
}
