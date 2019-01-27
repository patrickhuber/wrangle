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
	Describe("Set", func() {
		It("can set value", func() {
			item := store.NewItem("key", "value")
			err := memoryStore.Set(item)
			Expect(err).To(BeNil())
		})
	})
	Describe("GetByName", func() {
		It("returns value", func() {
			item := store.NewItem("key", "value")
			err := memoryStore.Set(item)
			Expect(err).To(BeNil())

			data, err := memoryStore.Get(item.Name())

			Expect(err).To(BeNil())
			Expect(data.Value()).To(Equal(item.Value()))
		})
		Context("when doesn't exist", func() {
			It("returns an error", func() {
				_, err := memoryStore.Get("key")
				Expect(err).ToNot(BeNil())
			})
		})
	})
	Describe("Delete", func() {
		It("deletes value", func() {
			key := "key"
			item := store.NewItem(key, "value")
			err := memoryStore.Set(item)
			Expect(err).To(BeNil())

			err = memoryStore.Delete(key)
			Expect(err).To(BeNil())
		})
	})
})
