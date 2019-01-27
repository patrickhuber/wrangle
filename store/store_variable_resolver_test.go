package store_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/wrangle/store"
	"github.com/patrickhuber/wrangle/store/memory"
)

var _ = Describe("StoreVariableResolver", func() {
	It("can get value from resolver", func() {

		memoryStore := memory.NewMemoryStore("test")
		item := store.NewItem("key", "value")
		err := memoryStore.Set(item)
		Expect(err).To(BeNil())

		resolver := store.NewStoreVariableResolver(memoryStore)
		value, err := resolver.Get("key")
		Expect(err).To(BeNil())
		Expect(value).To(Equal("value"))
	})
})
