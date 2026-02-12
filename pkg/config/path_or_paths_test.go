package config_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v3"

	"github.com/LarsArtmann/SQLC-Wizzard/pkg/config"
)

var _ = Describe("PathOrPaths", func() {
	testSinglePath := func(path string, expectEmpty bool) {
		pop := config.NewSinglePath(path)

		Expect(pop.Strings()).To(Equal([]string{path}))
		Expect(pop.First()).To(Equal(path))
		Expect(pop.IsEmpty()).To(Equal(expectEmpty))
	}

	Describe("NewPathOrPaths", func() {
		It("should create from a slice of paths", func() {
			paths := []string{"path1", "path2", "path3"}
			pop := config.NewPathOrPaths(paths)

			Expect(pop.Strings()).To(Equal(paths))
			Expect(pop.IsEmpty()).To(BeFalse())
		})

		It("should handle empty slice", func() {
			pop := config.NewPathOrPaths([]string{})

			Expect(pop.Strings()).To(BeEmpty())
			Expect(pop.IsEmpty()).To(BeTrue())
		})

		It("should handle nil slice", func() {
			pop := config.NewPathOrPaths(nil)

			Expect(pop.Strings()).To(BeEmpty())
			Expect(pop.IsEmpty()).To(BeTrue())
		})
	})

	Describe("NewSinglePath", func() {
		It("should create from a single path", func() {
			testSinglePath("path/to/queries", false)
		})

		It("should handle empty string", func() {
			testSinglePath("", false)
		})
	})

	Describe("First", func() {
		It("should return first path when multiple paths exist", func() {
			pop := config.NewPathOrPaths([]string{"first", "second", "third"})

			Expect(pop.First()).To(Equal("first"))
		})

		It("should return empty string when no paths exist", func() {
			pop := config.NewPathOrPaths([]string{})

			Expect(pop.First()).To(Equal(""))
		})
	})

	Describe("UnmarshalYAML", func() {
		Context("with a single string", func() {
			It("should unmarshal to a single-element slice", func() {
				yamlData := `path: "internal/db/queries"`
				var result struct {
					Path config.PathOrPaths `yaml:"path"`
				}

				err := yaml.Unmarshal([]byte(yamlData), &result)

				Expect(err).ToNot(HaveOccurred())
				Expect(result.Path.Strings()).To(Equal([]string{"internal/db/queries"}))
				Expect(result.Path.First()).To(Equal("internal/db/queries"))
			})
		})

		Context("with an array of strings", func() {
			It("should unmarshal to a multi-element slice", func() {
				yamlData := `
path:
  - "internal/db/queries"
  - "internal/db/migrations"
  - "pkg/queries"`
				var result struct {
					Path config.PathOrPaths `yaml:"path"`
				}

				err := yaml.Unmarshal([]byte(yamlData), &result)

				Expect(err).ToNot(HaveOccurred())
				Expect(result.Path.Strings()).To(Equal([]string{
					"internal/db/queries",
					"internal/db/migrations",
					"pkg/queries",
				}))
			})
		})

		Context("with invalid data", func() {
			It("should return error for object/map value", func() {
				yamlData := `
path:
  foo: bar`
				var result struct {
					Path config.PathOrPaths `yaml:"path"`
				}

				err := yaml.Unmarshal([]byte(yamlData), &result)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("path_or_paths must be either a string or array of strings"))
			})
		})
	})

	Describe("MarshalYAML", func() {
		It("should marshal as array", func() {
			pop := config.NewPathOrPaths([]string{"path1", "path2"})
			data := struct {
				Path config.PathOrPaths `yaml:"path"`
			}{
				Path: pop,
			}

			yamlBytes, err := yaml.Marshal(data)

			Expect(err).ToNot(HaveOccurred())
			Expect(string(yamlBytes)).To(ContainSubstring("- path1"))
			Expect(string(yamlBytes)).To(ContainSubstring("- path2"))
		})

		It("should marshal single path as array", func() {
			pop := config.NewSinglePath("single/path")
			data := struct {
				Path config.PathOrPaths `yaml:"path"`
			}{
				Path: pop,
			}

			yamlBytes, err := yaml.Marshal(data)

			Expect(err).ToNot(HaveOccurred())
			Expect(string(yamlBytes)).To(ContainSubstring("- single/path"))
		})

		It("should marshal empty as empty array", func() {
			pop := config.NewPathOrPaths([]string{})
			data := struct {
				Path config.PathOrPaths `yaml:"path"`
			}{
				Path: pop,
			}

			yamlBytes, err := yaml.Marshal(data)

			Expect(err).ToNot(HaveOccurred())
			Expect(string(yamlBytes)).To(ContainSubstring("path: []"))
		})
	})

	Describe("Round-trip marshal/unmarshal", func() {
		It("should preserve single path", func() {
			original := struct {
				Path config.PathOrPaths `yaml:"path"`
			}{
				Path: config.NewSinglePath("test/path"),
			}

			yamlBytes, err := yaml.Marshal(original)
			Expect(err).ToNot(HaveOccurred())

			var result struct {
				Path config.PathOrPaths `yaml:"path"`
			}
			err = yaml.Unmarshal(yamlBytes, &result)
			Expect(err).ToNot(HaveOccurred())

			Expect(result.Path.Strings()).To(Equal(original.Path.Strings()))
		})

		It("should preserve multiple paths", func() {
			original := struct {
				Path config.PathOrPaths `yaml:"path"`
			}{
				Path: config.NewPathOrPaths([]string{"path1", "path2", "path3"}),
			}

			yamlBytes, err := yaml.Marshal(original)
			Expect(err).ToNot(HaveOccurred())

			var result struct {
				Path config.PathOrPaths `yaml:"path"`
			}
			err = yaml.Unmarshal(yamlBytes, &result)
			Expect(err).ToNot(HaveOccurred())

			Expect(result.Path.Strings()).To(Equal(original.Path.Strings()))
		})
	})
})
