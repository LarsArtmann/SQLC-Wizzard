package adapters_test

import (
	"context"
	"testing"

	"github.com/LarsArtmann/SQLC-Wizzard/internal/adapters"
)

// BenchmarkRealFileSystemAdapter_ReadWriteFile benchmarks file read/write operations
func BenchmarkRealFileSystemAdapter_ReadWriteFile(b *testing.B) {
	fs := adapters.NewRealFileSystemAdapter()
	content := []byte("test content for benchmarking")
	ctx := context.Background()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Write file
		err := fs.WriteFile(ctx, "/tmp/benchmark-test.txt", content, 0644)
		if err != nil {
			b.Fatal(err)
		}
		
		// Read file
		_, err = fs.ReadFile(ctx, "/tmp/benchmark-test.txt")
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkRealFileSystemAdapter_Exists benchmarks file existence checks
func BenchmarkRealFileSystemAdapter_Exists(b *testing.B) {
	fs := adapters.NewRealFileSystemAdapter()
	ctx := context.Background()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = fs.Exists(ctx, "/tmp")
	}
}

// BenchmarkRealSQLCAdapter_CheckInstallation benchmarks sqlc installation checks
func BenchmarkRealSQLCAdapter_CheckInstallation(b *testing.B) {
	sqlcAdapter := adapters.NewRealSQLCAdapter()
	ctx := context.Background()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = sqlcAdapter.CheckInstallation(ctx)
	}
}

// BenchmarkRealCLIAdapter_RunCommand benchmarks CLI command execution
func BenchmarkRealCLIAdapter_RunCommand(b *testing.B) {
	cliAdapter := adapters.NewRealCLIAdapter()
	ctx := context.Background()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = cliAdapter.RunCommand(ctx, "echo", "benchmark")
	}
}