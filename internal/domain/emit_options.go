package domain

import (
	"github.com/LarsArtmann/SQLC-Wizzard/internal/errors"
	"github.com/LarsArtmann/SQLC-Wizzard/pkg/config"
)

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

// Validate checks if the EmitOptions are valid.
// This is a smart constructor pattern - validate after construction.
func (eo EmitOptions) Validate() error {
	// Validate JSON tags case style
	if eo.JSONTags && eo.JSONTagsCaseStyle != "" {
		switch eo.JSONTagsCaseStyle {
		case "camel", "pascal", "snake", "none":
			// Valid
		default:
			return errors.NewValidationErrorf(
				"json_tags_case_style",
				"invalid JSON tags case style: %s (must be one of: camel, pascal, snake, none)",
				eo.JSONTagsCaseStyle,
			).Error
		}
	}

	// Validate file names are not empty
	if eo.DBFileName == "" {
		return errors.NewValidationError("db_file_name", "DB file name cannot be empty").Error
	}
	if eo.ModelsFileName == "" {
		return errors.NewValidationError("models_file_name", "models file name cannot be empty").Error
	}
	if eo.QuerierFileName == "" {
		return errors.NewValidationError("querier_file_name", "querier file name cannot be empty").Error
	}

	// File names should end with .go
	fileNames := map[string]string{
		"db_file_name":       eo.DBFileName,
		"models_file_name":   eo.ModelsFileName,
		"querier_file_name":  eo.QuerierFileName,
		"copyfrom_file_name": eo.CopyfromFileName,
		"batch_file_name":    eo.BatchFileName,
	}

	for field, fileName := range fileNames {
		if fileName != "" && len(fileName) < 3 || !hasGoExtension(fileName) {
			return errors.NewValidationErrorf(
				field,
				"file name must end with .go: %s",
				fileName,
			).Error
		}
	}

	return nil
}

// hasGoExtension checks if a filename ends with .go
func hasGoExtension(filename string) bool {
	return len(filename) >= 3 && filename[len(filename)-3:] == ".go"
}

// NewEmitOptions creates and validates EmitOptions.
// Use this constructor when building EmitOptions from user input.
func NewEmitOptions(opts EmitOptions) (EmitOptions, error) {
	if err := opts.Validate(); err != nil {
		return EmitOptions{}, err
	}
	return opts, nil
}
