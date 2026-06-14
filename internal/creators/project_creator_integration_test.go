package creators_test

import (
	"context"
	"fmt"
	"strings"

	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/creators"
	"github.com/LarsArtmann/SQLC-Wizzard/pkg/config"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ProjectCreator Integration and Error Handling", func() {
	var (
		mockFS  *MockFileSystemAdapter
		mockCLI *MockCLIAdapter
		creator *creators.ProjectCreator
		ctx     context.Context
	)

	BeforeEach(func() {
		mockFS = &MockFileSystemAdapter{}
		mockCLI = &MockCLIAdapter{}
		creator = creators.NewProjectCreator(mockFS, mockCLI)
		ctx = context.Background()
	})

	Context("Integration", func() {
		It("should create complete project structure in correct order", func() {
			cfg := createBaseConfig("integration-test")

			err := creator.CreateProject(ctx, cfg)

			Expect(err).NotTo(HaveOccurred())

			// Verify order: directories first, then files
			Expect(mockFS.mkdirAllCalls).NotTo(BeEmpty())
			Expect(mockFS.writeFileCalls).To(HaveLen(2))

			// Verify that all mkdir calls occur before any write calls
			// by checking that no "write:" entries appear before any "mkdir:" entries
			foundWrite := false
			for _, entry := range mockFS.callLog {
				if strings.HasPrefix(entry, "write:") {
					foundWrite = true
				} else if strings.HasPrefix(entry, "mkdir:") && foundWrite {
					// Found a mkdir after a write - wrong order!
					Fail(
						fmt.Sprintf(
							"mkdir call %s appears after write call, violates order requirement",
							entry,
						),
					)
				}
			}

			// Ensure we have both mkdir and write operations
			hasMkdir := false
			hasWrite := false

			for _, entry := range mockFS.callLog {
				if strings.HasPrefix(entry, "mkdir:") {
					hasMkdir = true
				} else if strings.HasPrefix(entry, "write:") {
					hasWrite = true
				}
			}

			Expect(hasMkdir).To(BeTrue(), "Expected at least one mkdir operation")
			Expect(hasWrite).To(BeTrue(), "Expected at least one write operation")
		})
	})

	Context("Error Handling", func() {
		createTestConfig := func() *creators.CreateConfig {
			return &creators.CreateConfig{
				ProjectName: "test",
				ProjectType: generated.ProjectTypeMicroservice,
				Database:    generated.DatabaseTypePostgreSQL,
				Config:      &config.SqlcConfig{Version: "2"},
			}
		}

		DescribeTable(
			"should return descriptive error",
			func(setup func(), expectedError string) {
				mockFS.shouldFailMkdir = false
				mockFS.shouldFailWrite = false

				setup()

				err := creator.CreateProject(ctx, createTestConfig())
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring(expectedError))
			},
			Entry(
				"when mkdir fails",
				func() { mockFS.shouldFailMkdir = true },
				"failed to create directory structure",
			),
			Entry(
				"when YAML write fails",
				func() { mockFS.shouldFailWrite = true },
				"failed to generate sqlc.yaml",
			),
		)
	})
})
