package store_test

import (
	"github.com/patrickhuber/wrangle/store"
	"github.com/patrickhuber/wrangle/store/memory"
)

var _ = Describe("", func() {
	It("can get value from resolver", func() {

		memoryStore := memory.NewMemoryStore("test")
		_, err := memoryStore.Put("key", "value")
		Expect(err).To(BeNil())

		resolver := store.NewStoreVariableResolver(memoryStore)
		value, err := resolver.Get("key")
		Expect(err).To(BeNil())
	})
})
