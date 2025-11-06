package commands

import (
	"context"
	"fmt"
	"os"
	"runtime"

	"github.com/LarsArtmann/SQLC-Wizzard/internal/adapters"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/errors"
	"github.com/LarsArtmann/SQLC-Wizzard/pkg/config"
	"github.com/spf13/cobra"
)

// DoctorCommand creates a new doctor command
func NewDoctorCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "doctor",
		Short: "🩺 Diagnose potential issues with your environment",
		Long: `The doctor command performs a comprehensive health check of your development
environment to identify potential issues before they cause problems.`,
		RunE: runDoctor,
	}
}

// DoctorCheck represents a single diagnostic check
type DoctorCheck struct {
	Name        string
	Description string
	Checker     func(context.Context) *DoctorResult
}

// DoctorResult represents the result of a diagnostic check
type DoctorResult struct {
	Status    string // "PASS", "FAIL", "WARN"
	Message   string
	Solution  string
	Error     error
}

// runDoctor executes the diagnostic checks
func runDoctor(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	
	fmt.Println("🩺 SQLC-Wizard Health Check")
	fmt.Println("============================")
	
	checks := []DoctorCheck{
		{
			Name:        "go-version",
			Description: "Check Go version compatibility",
			Checker:     checkGoVersion,
		},
		{
			Name:        "sqlc-installation",
			Description: "Check sqlc installation",
			Checker:     checkSQLCInstallation,
		},
		{
			Name:        "database-drivers",
			Description: "Check database driver availability",
			Checker:     checkDatabaseDrivers,
		},
		{
			Name:        "filesystem-permissions",
			Description: "Check filesystem permissions",
			Checker:     checkFileSystemPermissions,
		},
		{
			Name:        "memory-availability",
			Description: "Check available memory",
			Checker:     checkMemoryAvailability,
		},
	}
	
	var failed, warned int
	
	for _, check := range checks {
		fmt.Printf("\n🔍 %s\n", check.Description)
		result := check.Checker(ctx)
		
		switch result.Status {
		case "PASS":
			fmt.Printf("   ✅ %s\n", result.Message)
		case "WARN":
			fmt.Printf("   ⚠️  %s\n", result.Message)
			if result.Solution != "" {
				fmt.Printf("   💡 Solution: %s\n", result.Solution)
			}
			warned++
		case "FAIL":
			fmt.Printf("   ❌ %s\n", result.Message)
			if result.Solution != "" {
				fmt.Printf("   💡 Solution: %s\n", result.Solution)
			}
			failed++
		}
		
		if result.Error != nil {
			fmt.Printf("   🐛 Error: %v\n", result.Error)
		}
	}
	
	// Summary
	fmt.Println("\n" + "============================")
	fmt.Println("🏁 Health Check Summary")
	
	switch {
	case failed == 0 && warned == 0:
		fmt.Println("🎉 All checks passed! Your environment is ready for SQLC-Wizard.")
		return nil
	case failed == 0 && warned > 0:
		fmt.Printf("⚠️  %d warning(s) found. Consider addressing them for optimal experience.\n", warned)
		return nil
	default:
		fmt.Printf("❌ %d error(s) and %d warning(s) found. Please fix errors before continuing.\n", failed, warned)
		return errors.New("health_check_failed", "environment issues detected")
	}
}

// checkGoVersion checks Go version compatibility
func checkGoVersion(ctx context.Context) *DoctorResult {
	goVersion := runtime.Version()
	
	// Check minimum Go version (simplified check)
	minVersion := "go1.21"
	if goVersion < minVersion {
		return &DoctorResult{
			Status:   "FAIL",
			Message:  fmt.Sprintf("Go version %s is too old", goVersion),
			Solution: fmt.Sprintf("Please upgrade to Go %s or later", minVersion),
		}
	}
	
	return &DoctorResult{
		Status:  "PASS",
		Message: fmt.Sprintf("Go version %s is compatible", goVersion),
	}
}

// checkSQLCInstallation checks sqlc installation
func checkSQLCInstallation(ctx context.Context) *DoctorResult {
	sqlcAdapter := adapters.NewRealSQLCAdapter()
	err := sqlcAdapter.CheckInstallation(ctx)
	
	if err != nil {
		return &DoctorResult{
			Status:   "FAIL",
			Message:  "sqlc is not installed or not in PATH",
			Solution: "Install sqlc following instructions at https://docs.sqlc.dev",
			Error:    err,
		}
	}
	
	// Try to get version
	_, err = sqlcAdapter.Version(ctx)
	if err != nil {
		return &DoctorResult{
			Status:   "WARN",
			Message:  "sqlc is installed but version check failed",
			Solution: "Ensure sqlc is properly configured",
			Error:    err,
		}
	}
	
	return &DoctorResult{
		Status:  "PASS",
		Message: "sqlc is installed and working",
	}
}

// checkDatabaseDrivers checks database driver availability
func checkDatabaseDrivers(ctx context.Context) *DoctorResult {
	dbAdapter := adapters.NewRealDatabaseAdapter()
	
	// Check connection to SQLite (most common for development)
	sqliteConfig := &config.DatabaseConfig{
		URI:     ":memory:",
		Managed: false,
	}
	
	err := dbAdapter.TestConnection(ctx, sqliteConfig)
	if err != nil {
		return &DoctorResult{
			Status:   "WARN",
			Message:  "SQLite driver may not be available",
			Solution: "Install SQLite3: brew install sqlite3 (macOS) or apt-get install sqlite3 (Ubuntu)",
			Error:    err,
		}
	}
	
	return &DoctorResult{
		Status:  "PASS",
		Message: "SQLite driver is available",
	}
}

// checkFileSystemPermissions checks filesystem permissions
func checkFileSystemPermissions(ctx context.Context) *DoctorResult {
	fsAdapter := adapters.NewRealFileSystemAdapter()
	
	// Try to create a temporary file
	testContent := []byte("test")
	testFile := "/tmp/sqlc-wizard-test"
	
	err := fsAdapter.WriteFile(ctx, testFile, testContent, 0644)
	if err != nil {
		return &DoctorResult{
			Status:   "FAIL",
			Message:  "Cannot write to filesystem",
			Solution: "Check directory permissions and disk space",
			Error:    err,
		}
	}
	
	// Try to read it back
	_, err = fsAdapter.ReadFile(ctx, testFile)
	if err != nil {
		return &DoctorResult{
			Status:   "FAIL",
			Message:  "Cannot read from filesystem",
			Solution: "Check file permissions",
			Error:    err,
		}
	}
	
	// Clean up
	_ = os.Remove(testFile)
	
	return &DoctorResult{
		Status:  "PASS",
		Message: "Filesystem permissions are OK",
	}
}

// checkMemoryAvailability checks available memory
func checkMemoryAvailability(ctx context.Context) *DoctorResult {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	
	// Convert bytes to MB
	availableMB := int(m.Sys / 1024 / 1024)
	minMemoryMB := 512 // 512MB minimum recommended
	
	if availableMB < minMemoryMB {
		return &DoctorResult{
			Status:   "WARN",
			Message:  fmt.Sprintf("Only %d MB memory available", availableMB),
			Solution: "Consider closing other applications or increasing available memory",
		}
	}
	
	return &DoctorResult{
		Status:  "PASS",
		Message: fmt.Sprintf("Memory is sufficient (%d MB available)", availableMB),
	}
}