package gomuti_test

import (
	"net/url"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/xeger/gomuti"
	"github.com/xeger/mongoose/test/fixtures"
)

var _ = Describe("gomuti dialect", func() {
	Context("mocking", func() {
		var v fixtures.Vehicle
		var w fixtures.Wheel

		URL, err := url.Parse("http://null-island.com")
		if err != nil {
			panic("bad test fixture")
		}

		BeforeEach(func() {
			v = &fixtures.MockVehicle{}
			w = &fixtures.MockWheel{}
			Â(v).Call("Attach").With("hello", []fixtures.Wheel{w})
			Â(v).Call("Drive").With("north", 42.0).Return(*URL)
			Â(w).Call("Diameter").Panic("big wheel keep on turnin'")
		})

		It("matches calls", func() {
			Ω(v.Drive("north", 42.0)).Should(Equal(*URL))
			Ω(func() {
				w.Diameter()
			}).Should(Panic())
		})

		It("panics on unmatched calls", func() {
			Ω(func() {
				v.Drive("southeast", 12)
			}).Should(Panic())
		})
	})

	Context("stubbing", func() {
		It("allows stubbing", func() {
			v := fixtures.MockVehicle{Stub: true}
			w := fixtures.MockWheel{Stub: true}
			Ω(v.Range()).Should(Equal(0))
			v.Attach("getting in")
			Ω(v.Wheels()).Should(BeNil())
			Ω(v.Drive("east", -5)).Should(BeEquivalentTo(url.URL{}))
			Ω(v.Refuel(&fixtures.MockFuelCan{Stub: true})).Should(BeNil())
			Ω(v.Occupants()).Should(BeNil())
			Ω(w.Diameter()).Should(BeEquivalentTo(0.0))
		})
	})

	Context("return values", func() {
		It("detects type mismatch", func() {
			v := &fixtures.MockVehicle{}

			Ω(func() (result interface{}) {
				defer func() {
					if r := recover(); r != nil {
						result = r
					}
				}()

				Â(v).Call("Range").Return("fourty two")
				result = v.Range()

				return result
			}()).Should(ContainSubstring("return type mismatch"))
		})

		It("detects return value mismatch", func() {
			v := &fixtures.MockVehicle{}

			Ω(func() (result interface{}) {
				defer func() {
					if r := recover(); r != nil {
						result = r
					}
				}()

				Â(v).Call("Range").Return(42, 422)
				result = v.Range()

				return result
			}()).Should(ContainSubstring("return value mismatch"))
		})
	})
})

func TestMongoose(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gomuti Suite")
}
