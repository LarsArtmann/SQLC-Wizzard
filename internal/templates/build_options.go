package templates

// BuildOptions holds all configuration options for BuildDefaultData.
// This reduces the function signature from 21 parameters to a single struct,
// making it easier to extend and maintain.
type BuildOptions struct {
	// Required fields
	ProjectType   string
	DbEngine      string
	DatabaseURL   string
	PackagePath   string
	BaseOutputDir string

	// Database features - all default to true
	UseManaged  bool
	UseUUIDs    bool
	UseJSON     bool
	UseArrays   bool
	UseFullText bool

	// Emit options
	EmitJSONTags             bool
	EmitPreparedQueries      bool
	EmitInterface            bool
	EmitEmptySlices          bool
	EmitResultStructPointers bool
	EmitParamsStructPointers bool
	EmitEnumValidMethod      bool
	EmitAllEnumValues        bool
	JSONTagsCaseStyle        string

	// Emit options - extended
	StrictFunctions bool
	StrictOrderBy   bool

	// Safety rules
	NoSelectStar bool
	RequireWhere bool
	NoDropTable  bool
	NoTruncate   bool
	RequireLimit bool
}

// NewBuildOptions creates BuildOptions with sensible defaults.
// This is the preferred way to create options for BuildDefaultData.
func NewBuildOptions(projectType, dbEngine string) BuildOptions {
	return BuildOptions{
		// Required fields - caller must override
		ProjectType:   projectType,
		DbEngine:      dbEngine,
		DatabaseURL:   DefaultDatabaseURL,
		PackagePath:   DefaultPackagePath,
		BaseOutputDir: DefaultPackagePath,

		// Database features - default to true
		UseManaged:  true,
		UseUUIDs:    true,
		UseJSON:     true,
		UseArrays:   true,
		UseFullText: false,

		// Emit options
		EmitPreparedQueries:      true,
		EmitResultStructPointers: true,
		EmitParamsStructPointers: true,
		EmitJSONTags:             false,
		EmitInterface:            false,
		EmitEmptySlices:          false,
		EmitEnumValidMethod:      false,
		EmitAllEnumValues:        false,
		JSONTagsCaseStyle:        DefaultJSONStyle,

		// Emit options - extended
		StrictFunctions: false,
		StrictOrderBy:   false,

		// Safety rules - conservative defaults
		NoSelectStar: true,
		RequireWhere: true,
		NoDropTable:  false,
		NoTruncate:   false,
		RequireLimit: false,
	}
}
