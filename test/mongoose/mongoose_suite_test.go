package mongoose_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/xeger/mongoose/mock"
)

var _ = Describe("mongoose dialect", func() {
	It("has an RSpec-like DSL", func() {
		Allow(Mock{}).ToReceive("Foo", 2)
	})

	It("has a gomega-like DSL", func() {
		Â(Mock{}).On("Foo", 1)
	})

	It("generates code", func() {
		v := &mongoose.MockVehicle{}
		w := &mongoose.MockWheel{}

		v.On("Attach", []mongoose.Wheel{w})
		w.On("Diameter").Return(17.0)

		v.Attach(w)
	})
})

func TestMongoose(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Mongoose Suite")
}
