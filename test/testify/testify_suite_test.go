package testify_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/xeger/mongoose/test/testify"

	"testing"
)

var _ = Describe("testify dialect", func() {
	It("generates code", func() {
		v := &testify.MockVehicle{}
		w := &testify.MockWheel{}

		v.On("Attach", []testify.Wheel{w})
		w.On("Diameter").Return(17.0)

		v.Attach(w)
	})
})

func TestTestify(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Testify Suite")
}
