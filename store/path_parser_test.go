package store_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/patrickhuber/wrangle/store"
)

var _ = Describe("PathParser", func() {
	It("can parse", func() {
		storeAndPath, err := store.ParsePath("hello:kitty")
		Expect(err).To(BeNil())
		Expect(storeAndPath.Store).To(Equal("hello"))
		Expect(storeAndPath.Path).To(Equal("kitty"))
	})
})
