package domain_test

import (
	"github.com/LarsArtmann/SQLC-Wizzard/internal/domain"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("TypeSafeEmitOptions", func() {
	Context("IsValid", func() {
		It("should validate options with all valid modes", func() {
			opts := domain.TypeSafeEmitOptions{
				NullHandling:   domain.NullHandlingPointers,
				EnumMode:       domain.EnumGenerationComplete,
				StructPointers: domain.StructPointerNever,
				JSONTagStyle:   domain.JSONTagStyleCamel,
			}

			err := opts.IsValid()
			Expect(err).NotTo(HaveOccurred())
		})

		It("should reject invalid null handling mode", func() {
			opts := domain.TypeSafeEmitOptions{
				NullHandling:   "invalid",
				EnumMode:       domain.EnumGenerationComplete,
				StructPointers: domain.StructPointerNever,
				JSONTagStyle:   domain.JSONTagStyleCamel,
			}

			err := opts.IsValid()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("null handling"))
		})

		It("should reject invalid enum mode", func() {
			opts := domain.TypeSafeEmitOptions{
				NullHandling:   domain.NullHandlingPointers,
				EnumMode:       "invalid",
				StructPointers: domain.StructPointerNever,
				JSONTagStyle:   domain.JSONTagStyleCamel,
			}

			err := opts.IsValid()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("enum generation mode"))
		})

		It("should reject invalid struct pointer mode", func() {
			opts := domain.TypeSafeEmitOptions{
				NullHandling:   domain.NullHandlingPointers,
				EnumMode:       domain.EnumGenerationComplete,
				StructPointers: "invalid",
				JSONTagStyle:   domain.JSONTagStyleCamel,
			}

			err := opts.IsValid()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("struct pointer"))
		})

		It("should reject invalid JSON tag style", func() {
			opts := domain.TypeSafeEmitOptions{
				NullHandling:   domain.NullHandlingPointers,
				EnumMode:       domain.EnumGenerationComplete,
				StructPointers: domain.StructPointerNever,
				JSONTagStyle:   "invalid",
			}

			err := opts.IsValid()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("JSON tag style"))
		})
	})

	Context("ToEmitOptions", func() {
		It("should convert to emit options correctly", func() {
			opts := domain.TypeSafeEmitOptions{
				NullHandling:   domain.NullHandlingPointers,
				EnumMode:       domain.EnumGenerationWithValidation,
				StructPointers: domain.StructPointerResults,
				JSONTagStyle:   domain.JSONTagStyleSnake,
			}

			data := opts.ToTemplateData()

			Expect(data.EmitEmptySlices).To(BeFalse())
			Expect(data.EmitEnumValidMethod).To(BeTrue())
			Expect(data.EmitResultStructPointers).To(BeTrue())
			Expect(data.EmitParamsStructPointers).To(BeFalse())
			Expect(data.JSONTagsCaseStyle).To(Equal("snake"))
		})
	})

	Context("ApplyDefaults", func() {
		It("should apply default values when empty", func() {
			opts := domain.TypeSafeEmitOptions{}

			opts.ApplyDefaults()

			Expect(opts.NullHandling).To(Equal(domain.NullHandlingPointers))
			Expect(opts.EnumMode).To(Equal(domain.EnumGenerationBasic))
			Expect(opts.StructPointers).To(Equal(domain.StructPointerNever))
			Expect(opts.JSONTagStyle).To(Equal(domain.JSONTagStyleCamel))
		})

		It("should not override existing values", func() {
			opts := domain.TypeSafeEmitOptions{
				NullHandling:   domain.NullHandlingEmptySlices,
				EnumMode:       domain.EnumGenerationComplete,
				StructPointers: domain.StructPointerAlways,
				JSONTagStyle:   domain.JSONTagStyleKebab,
			}

			opts.ApplyDefaults()

			Expect(opts.NullHandling).To(Equal(domain.NullHandlingEmptySlices))
			Expect(opts.EnumMode).To(Equal(domain.EnumGenerationComplete))
			Expect(opts.StructPointers).To(Equal(domain.StructPointerAlways))
			Expect(opts.JSONTagStyle).To(Equal(domain.JSONTagStyleKebab))
		})
	})
})
