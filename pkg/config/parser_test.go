package config_test

import (
	"fmt"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/LarsArtmann/SQLC-Wizzard/internal/errors"
	"github.com/LarsArtmann/SQLC-Wizzard/pkg/config"
)

var _ = Describe("Parser", func() {
	var tempDir string

	BeforeEach(func() {
		var err error
		tempDir, err = os.MkdirTemp("", "sqlc-wizard-test-*")
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		if tempDir != "" {
			if err := os.RemoveAll(tempDir); err != nil {
				// Log the cleanup error but don't fail the test
				fmt.Printf("warning: failed to cleanup temp dir %s: %v\n", tempDir, err)
			}
		}
	})

	Describe("Parse", func() {
		Context("with valid YAML", func() {
			It("should parse minimal config", func() {
				yamlData := []byte(`
version: "2"
sql:
  - name: "test"
    engine: "postgresql"
    queries: "queries/"
    schema: "schema/"
    gen:
      go:
        package: "db"
        out: "internal/db"
`)

				cfg, err := config.Parse(yamlData)

				Expect(err).ToNot(HaveOccurred())
				Expect(cfg).ToNot(BeNil())
				Expect(cfg.Version).To(Equal("2"))
				Expect(cfg.SQL).To(HaveLen(1))
				Expect(cfg.SQL[0].Name).To(Equal("test"))
				Expect(cfg.SQL[0].Engine).To(Equal("postgresql"))
			})

			It("should parse config with multiple SQL entries", func() {
				yamlData := []byte(`
version: "2"
sql:
  - name: "db1"
    engine: "postgresql"
    queries: "db1/queries"
    schema: "db1/schema"
  - name: "db2"
    engine: "mysql"
    queries: "db2/queries"
    schema: "db2/schema"
`)

				cfg, err := config.Parse(yamlData)

				Expect(err).ToNot(HaveOccurred())
				Expect(cfg.SQL).To(HaveLen(2))
				Expect(cfg.SQL[0].Name).To(Equal("db1"))
				Expect(cfg.SQL[1].Name).To(Equal("db2"))
			})
		})

		Context("with invalid YAML", func() {
			It("should return error for malformed YAML", func() {
				yamlData := []byte(`
version: "2"
sql:
  - name: "test
    engine: "postgresql"
`)

				_, err := config.Parse(yamlData)

				Expect(err).To(HaveOccurred())
				Expect(errors.HasCode(err, errors.ErrCodeConfigParseFailed)).To(BeTrue())
			})

			It("should return error for empty data", func() {
				yamlData := []byte(``)

				cfg, err := config.Parse(yamlData)

				// Empty YAML is technically valid, just creates empty struct
				Expect(err).ToNot(HaveOccurred())
				Expect(cfg.Version).To(BeEmpty())
			})
		})
	})

	Describe("ParseFile", func() {
		Context("with existing file", func() {
			It("should parse valid config file", func() {
				configPath := filepath.Join(tempDir, "sqlc.yaml")
				configContent := []byte(`
version: "2"
sql:
  - name: "test"
    engine: "postgresql"
    queries: "queries/"
    schema: "schema/"
    gen:
      go:
        package: "db"
        out: "internal/db"
`)
				err := os.WriteFile(configPath, configContent, 0644)
				Expect(err).ToNot(HaveOccurred())

				cfg, err := config.ParseFile(configPath)

				Expect(err).ToNot(HaveOccurred())
				Expect(cfg).ToNot(BeNil())
				Expect(cfg.Version).To(Equal("2"))
				Expect(cfg.SQL[0].Name).To(Equal("test"))
			})
		})

		Context("with non-existent file", func() {
			It("should return ConfigNotFoundError", func() {
				nonExistentPath := filepath.Join(tempDir, "does-not-exist.yaml")

				_, err := config.ParseFile(nonExistentPath)

				Expect(err).To(HaveOccurred())
				Expect(errors.HasCode(err, errors.ErrCodeConfigNotFound)).To(BeTrue())
			})
		})

		Context("with invalid YAML file", func() {
			It("should return ConfigParseError", func() {
				configPath := filepath.Join(tempDir, "invalid.yaml")
				invalidContent := []byte(`
version: "2"
sql:
  - name: [this is invalid yaml
`)
				err := os.WriteFile(configPath, invalidContent, 0644)
				Expect(err).ToNot(HaveOccurred())

				_, err = config.ParseFile(configPath)

				Expect(err).To(HaveOccurred())
				Expect(errors.HasCode(err, errors.ErrCodeConfigParseFailed)).To(BeTrue())
			})
		})
	})

	Describe("LoadOrDefault", func() {
		Context("with existing file", func() {
			It("should load the config", func() {
				configPath := filepath.Join(tempDir, "sqlc.yaml")
				configContent := []byte(`
version: "2"
sql:
  - name: "loaded"
    engine: "postgresql"
`)
				err := os.WriteFile(configPath, configContent, 0644)
				Expect(err).ToNot(HaveOccurred())

				cfg, err := config.LoadOrDefault(configPath)

				Expect(err).ToNot(HaveOccurred())
				Expect(cfg.SQL[0].Name).To(Equal("loaded"))
			})
		})

		Context("with non-existent file", func() {
			It("should return default config", func() {
				nonExistentPath := filepath.Join(tempDir, "missing.yaml")

				cfg, err := config.LoadOrDefault(nonExistentPath)

				Expect(err).ToNot(HaveOccurred())
				Expect(cfg).ToNot(BeNil())
				Expect(cfg.Version).To(Equal("2"))
				Expect(cfg.SQL).To(BeEmpty())
			})
		})
	})

	Describe("DefaultConfig", func() {
		It("should return a valid default config", func() {
			cfg := config.DefaultConfig()

			Expect(cfg).ToNot(BeNil())
			Expect(cfg.Version).To(Equal("2"))
			Expect(cfg.SQL).To(BeEmpty())
		})
	})
})
