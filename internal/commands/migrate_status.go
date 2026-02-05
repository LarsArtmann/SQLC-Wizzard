package commands

import (
	"context"
	"fmt"

	"github.com/LarsArtmann/SQLC-Wizzard/internal/adapters"
	"github.com/spf13/cobra"
)

// newMigrateStatusCommand creates migrate status command.
func newMigrateStatusCommand() *cobra.Command {
	var (
		source   string
		database string
	)

	cmd := &cobra.Command{
		Use:   "status",
		Short: "Check migration status",
		Long:  "Check the status of database migrations and show the current version.",
		Run: func(cmd *cobra.Command, args []string) {
			migrationConfig := &StatusConfig{
				Source:   source,
				Database: database,
			}

			if err := runStatusCheck(migrationConfig); err != nil {
				fmt.Printf("❌ Status check failed: %v\n", err)
				return
			}
		},
	}

	cmd.Flags().StringVarP(&source, "source", "s", "", "Migration source path")
	cmd.Flags().StringVarP(&database, "database", "d", "", "Database connection URL")
	return cmd
}

// StatusConfig represents status check configuration.
type StatusConfig struct {
	Source   string
	Database string
}

// runStatusCheck executes migration status check.
func runStatusCheck(config *StatusConfig) error {
	if config.Source == "" {
		return &MigrationError{
			Code:    "MISSING_SOURCE",
			Message: "Please specify migration source path",
		}
	}

	if config.Database == "" {
		return &MigrationError{
			Code:    "MISSING_DATABASE",
			Message: "Please specify database URL",
		}
	}

	// Create migration adapter and check status
	migrationAdapter := adapters.NewRealMigrationAdapter()

	status, err := migrationAdapter.Status(context.Background(), config.Source, config.Database)
	if err != nil {
		return &MigrationError{
			Code:    "STATUS_CHECK_FAILED",
			Message: fmt.Sprintf("Failed to check migration status: %v", err),
		}
	}

	fmt.Printf("📊 Migration Status:\n")
	fmt.Printf("  Current Version: %v\n", status.GetCurrentVersion())
	fmt.Printf("  Dirty State: %v\n", status.IsDirty())
	fmt.Printf("  Total Migrations: %v\n", status.GetMigrationCount())

	return nil
}
