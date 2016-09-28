package generate

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("go:generate support", func() {
	It("generates mocks", func() {
		m := MockGenerated{Stub: true}
		Expect(m.Foo("bar", 7)).NotTo(HaveOccurred())
	})
})

func TestMongoose(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gomuti Generation Suite")
}
