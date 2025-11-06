package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// Version information (set via ldflags during build)
	Version   = "dev"
	Commit    = "unknown"
	BuildDate = "unknown"
)

func main() {
	if err := newRootCmd().Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func newRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "sqlc-wizard",
		Short: "🧙‍♂️ An interactive CLI wizard for generating sqlc configurations",
		Long: `SQLC-Wizard makes type-safe SQL accessible to everyone by providing
an intuitive wizard that guides developers through creating production-ready
sqlc setups with smart defaults and comprehensive validation.

Generate perfect sqlc.yaml configurations in minutes, not hours!`,
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	// Add version command
	rootCmd.AddCommand(newVersionCmd())

	// TODO: Add init command
	// TODO: Add generate command
	// TODO: Add validate command
	// TODO: Add doctor command
	// TODO: Add plugins command
	// TODO: Add migrate command

	return rootCmd
}

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("sqlc-wizard %s\n", Version)
			fmt.Printf("  commit: %s\n", Commit)
			fmt.Printf("  built:  %s\n", BuildDate)
		},
	}
}
