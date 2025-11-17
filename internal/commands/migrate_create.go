package commands

import (
	"context"
	"fmt"

	"github.com/LarsArtmann/SQLC-Wizzard/internal/adapters"
	"github.com/spf13/cobra"
)

// newMigrateDBCreateCommand creates migration creation command
func newMigrateDBCreateCommand() *cobra.Command {
	var name string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new migration file",
		Long:  "Create a new up/down migration file pair with a specified name.",
		Run: func(cmd *cobra.Command, args []string) {
			migrationConfig := &CreateConfig{
				Name:           name,
				MigrationsPath: cmd.Flag("path").Value.String(),
			}

			if err := runMigrationCreate(migrationConfig); err != nil {
				fmt.Printf("❌ Migration creation failed: %v\n", err)
				return
			}
		},
	}

	cmd.Flags().StringVarP(&name, "name", "n", "", "Migration name")
	cmd.Flags().String("path", ".", "Migrations directory path")
	return cmd
}

// CreateConfig represents migration creation configuration
type CreateConfig struct {
	Name           string
	MigrationsPath string
}

// runMigrationCreate executes migration file creation
func runMigrationCreate(config *CreateConfig) error {
	if config.Name == "" {
		return &MigrationError{
			Code:    "MISSING_NAME",
			Message: "Please specify migration name",
		}
	}

	fmt.Printf("📝 Creating migration: %s\n", config.Name)

	migrationAdapter := adapters.NewRealMigrationAdapter()

	filename, err := migrationAdapter.CreateMigration(context.Background(), config.Name, config.MigrationsPath)
	if err != nil {
		return &MigrationError{
			Code:    "CREATION_FAILED",
			Message: fmt.Sprintf("Migration creation failed: %v", err),
		}
	}

	fmt.Printf("✅ Migration created: %s\n", filename)
	return nil
}
