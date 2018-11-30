package memory

import (
	"github.com/patrickhuber/wrangle/config"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MemoryStoreProvider", func() {
	It("can create memory store", func() {
		provider := NewMemoryStoreProvider()
		name := provider.Name()
		Expect(name).To(Equal("memory"))

		store, err := provider.Create(&config.Store{})
		Expect(err).To(BeNil())
		Expect(store).ToNot(BeNil())
	})
})
