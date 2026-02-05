package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

// newMigrateListCommand creates migrate list command.
func newMigrateListCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List available migration targets",
		Long:  "List available migration targets and supported options.",
		Run: func(cmd *cobra.Command, args []string) {
			runMigrationList()
		},
	}
}

// runMigrationList displays available migration options.
func runMigrationList() {
	fmt.Println("📋 Available Migration Targets:")
	fmt.Println("\n🔧 Configuration Migrations:")
	fmt.Println("  SQLC Versions: 1.x -> 2.0 (latest)")
	fmt.Println("\n🗄️  Database Types:")
	fmt.Println("  postgresql -> mysql")
	fmt.Println("  postgresql -> sqlite")
	fmt.Println("  mysql -> postgresql")
	fmt.Println("  mysql -> sqlite")
	fmt.Println("  sqlite -> postgresql")
	fmt.Println("  sqlite -> mysql")
	fmt.Println("\n📝 Usage Examples:")
	fmt.Println("  sqlc-wizard migrate -s old.sqlc.yaml -d new.sqlc.yaml -b postgresql")
	fmt.Println("  sqlc-wizard migrate -s config.yaml -d config.yaml -v 2.0")
	fmt.Println("  sqlc-wizard migrate status -s ./migrations -d sqlite://test.db")
	fmt.Println("  sqlc-wizard migrate create -n add_users_table")
}
