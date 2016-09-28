package generate

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("go:generate support", func() {
	It("generates mocks", func() {
		m := MockGenerated{}
		m.On("Foo", "a", 7).Return(nil)
		Expect(m.Foo("a", 7)).NotTo(HaveOccurred())
	})
})

func TestMongoose(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Testify Generation Suite")
}
