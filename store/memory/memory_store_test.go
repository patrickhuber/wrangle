package memory_test

import (
	"github.com/patrickhuber/wrangle/store"
	"github.com/patrickhuber/wrangle/store/memory"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MemoryStore", func() {
	var (
		memoryStoreName string
		memoryStore     store.Store
	)
	BeforeEach(func() {
		memoryStoreName = "test"
		memoryStore = memory.NewMemoryStore(memoryStoreName)
	})
	Describe("Name", func() {
		It("returns name", func() {
			Expect(memoryStore.Name()).To(Equal(memoryStoreName))
		})
	})
	Describe("Type", func() {
		It("returns type", func() {
			Expect(memoryStore.Type()).To(Equal("memory"))
		})
	})
	Describe("Put", func() {
		It("can put value", func() {
			key := "key"
			value := "value"
			_, err := memoryStore.Put(key, value)
			Expect(err).To(BeNil())
		})
	})
	Describe("GetByName", func() {
		It("returns value", func() {
			key := "key"
			value := "value"
			_, err := memoryStore.Put(key, value)
			Expect(err).To(BeNil())

			data, err := memoryStore.GetByName(key)

			Expect(err).To(BeNil())
			Expect(data.Value()).To(Equal(value))
		})
	})
	Describe("Delete", func() {
		It("deletes value", func() {
			key := "key"
			value := "value"

			_, err := memoryStore.Put(key, value)
			Expect(err).To(BeNil())

			count, err := memoryStore.Delete(key)
			Expect(err).To(BeNil())
			Expect(count).To(Equal(1))
		})
	})
})
