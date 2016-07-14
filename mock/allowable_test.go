package mock_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/xeger/mongoose/mock"
)

var _ = Describe("Allowed", func() {
	var Receiver Mock

	BeforeEach(func() {
		Receiver = Mock{}
	})

	Context("given a nonsensical call chain", func() {
		Receiver := Mock{}

		It("panics", func() {
			Expect(func() {
				Allow(Receiver).With(true)
			}).To(Panic())
			Expect(func() {
				Allow(Receiver).Return(false)
			}).To(Panic())
		})
	})

	Context("On", func() {
		It("begins call chains", func() {
			Allow(Receiver).On("Foo").Return(true)
			Allow(Receiver).On("Bar").Panic(true)
			Allow(Receiver).On("Baz").With(42).Return(true)
			Allow(Receiver).On("Baz").With(Not(Equal(42))).Panic("not the answer")
		})
	})

	Context("With", func() {
		Context("given basic types", func() {
			PIt("matches equivalency")
		})
		Context("given struct types", func() {
			PIt("matches equivalency")
		})
		Context("given matchers", func() {
			PIt("tests satisfaction")
		})
	})

	Context("Return", func() {
		PIt("programs return values")
	})

	Context("Panic", func() {
		PIt("causes a panic")
	})
})
