package templates

// SQL package paths for generated sqlc code.
const (
	SQLPackagePostgreSQL = "pgx/v5"
	SQLPackageStdlib     = "database/sql"
)

// Build tags for generated sqlc code per database engine.
const (
	BuildTagPostgreSQL = "postgres"
	BuildTagMySQL      = "mysql"
	BuildTagSQLite     = "sqlite"
)

// Default values used by templates when caller does not provide one.
const (
	DefaultDatabaseURL = "${DATABASE_URL}"
	DefaultPackagePath = "internal/db"
	DefaultJSONStyle   = "snake"
)
