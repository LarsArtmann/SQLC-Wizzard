package domain_test

import (
	"testing"

	"github.com/LarsArtmann/SQLC-Wizzard/internal/domain"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Domain Services", func() {
	Describe("UUID Generation", func() {
		It("should generate valid UUIDs", func() {
			uuid1 := domain.GenerateUUID()
			uuid2 := domain.GenerateUUID()

			Expect(uuid1).ToNot(BeEmpty())
			Expect(uuid2).ToNot(BeEmpty())
			Expect(uuid1).ToNot(Equal(uuid2))
		})

		It("should generate UUIDs in correct format", func() {
			uuid := domain.GenerateUUID()
			// UUID v4 format: xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx
			Expect(uuid).To(MatchRegexp(`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`))
		})
	})
})

func TestDomain(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Domain Suite")
}
