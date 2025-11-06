package templates_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/LarsArtmann/SQLC-Wizzard/internal/errors"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/templates"
)

var _ = Describe("Registry", func() {
	var registry *templates.Registry

	BeforeEach(func() {
		registry = templates.NewRegistry()
	})

	Describe("NewRegistry", func() {
		It("should create a registry with microservice template", func() {
			Expect(registry).ToNot(BeNil())
			Expect(registry.HasTemplate(templates.ProjectTypeMicroservice)).To(BeTrue())
		})

		It("should have microservice template registered", func() {
			tmpl, err := registry.Get(templates.ProjectTypeMicroservice)

			Expect(err).ToNot(HaveOccurred())
			Expect(tmpl).ToNot(BeNil())
			Expect(tmpl.Name()).To(Equal(string(templates.ProjectTypeMicroservice)))
		})
	})

	Describe("Register", func() {
		It("should register a new template", func() {
			mockTemplate := templates.NewMicroserviceTemplate()

			registry.Register(mockTemplate)

			Expect(registry.HasTemplate(templates.ProjectTypeMicroservice)).To(BeTrue())
		})
	})

	Describe("Get", func() {
		Context("with registered template", func() {
			It("should return the template", func() {
				tmpl, err := registry.Get(templates.ProjectTypeMicroservice)

				Expect(err).ToNot(HaveOccurred())
				Expect(tmpl).ToNot(BeNil())
				Expect(tmpl.Name()).To(Equal("microservice"))
			})
		})

		Context("with unregistered template", func() {
			It("should return TemplateNotFoundError", func() {
				_, err := registry.Get(templates.ProjectTypeHobby)

				Expect(err).To(HaveOccurred())
				Expect(errors.HasCode(err, errors.ErrCodeTemplateNotFound)).To(BeTrue())
			})
		})
	})

	Describe("List", func() {
		It("should return all registered templates", func() {
			templates := registry.List()

			Expect(templates).ToNot(BeEmpty())
			Expect(templates).To(HaveLen(1)) // Only microservice is registered currently
		})

		It("should return templates with correct names", func() {
			tmpls := registry.List()

			names := make([]string, len(tmpls))
			for i, t := range tmpls {
				names[i] = t.Name()
			}

			Expect(names).To(ContainElement("microservice"))
		})
	})

	Describe("HasTemplate", func() {
		It("should return true for registered template", func() {
			Expect(registry.HasTemplate(templates.ProjectTypeMicroservice)).To(BeTrue())
		})

		It("should return false for unregistered template", func() {
			Expect(registry.HasTemplate(templates.ProjectTypeHobby)).To(BeFalse())
			Expect(registry.HasTemplate(templates.ProjectTypeEnterprise)).To(BeFalse())
		})
	})

	Describe("GetTemplate (package function)", func() {
		It("should use the default registry", func() {
			tmpl, err := templates.GetTemplate(templates.ProjectTypeMicroservice)

			Expect(err).ToNot(HaveOccurred())
			Expect(tmpl).ToNot(BeNil())
		})

		It("should return error for unregistered template", func() {
			_, err := templates.GetTemplate(templates.ProjectTypeHobby)

			Expect(err).To(HaveOccurred())
			Expect(errors.HasCode(err, errors.ErrCodeTemplateNotFound)).To(BeTrue())
		})
	})

	Describe("ListTemplates (package function)", func() {
		It("should return templates from default registry", func() {
			templates := templates.ListTemplates()

			Expect(templates).ToNot(BeEmpty())
		})
	})
})
