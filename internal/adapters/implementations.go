// Package adapters provides concrete implementations of adapter interfaces
// All adapter implementations are now in separate files following SRP
package adapters

// Re-export all adapter constructors for convenience
var (
	NewSQLCAdapterFunc         = NewRealSQLCAdapter
	NewDatabaseAdapterFunc     = NewRealDatabaseAdapter
	NewCLIAdapterFunc         = NewRealCLIAdapter
	NewTemplateAdapterFunc    = NewRealTemplateAdapter
	NewFileSystemAdapterFunc  = NewRealFileSystemAdapter
	NewMigrationAdapterFunc   = NewRealMigrationAdapter
)