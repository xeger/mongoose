package testify_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/xeger/mongoose/test/fixtures"

	"testing"
)

var _ = Describe("testify dialect", func() {
	It("generates code", func() {
		v := &fixtures.MockVehicle{}
		w := &fixtures.MockWheel{}

		v.On("Attach", "hello", []fixtures.Wheel{w})
		w.On("Diameter").Return(17.0)

		v.Attach("hello", w)
	})
})

func TestTestify(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Testify Suite")
}
