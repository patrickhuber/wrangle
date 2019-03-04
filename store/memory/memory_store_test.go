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
			item := store.NewItem("key", store.Value, "value")
			err := memoryStore.Set(item)
			Expect(err).To(BeNil())
		})
	})
	Describe("Get", func() {
		It("returns value", func() {
			item := store.NewItem("key", store.Value, "value")
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
			item := store.NewItem(key, store.Value, "value")
			err := memoryStore.Set(item)
			Expect(err).To(BeNil())

			err = memoryStore.Delete(key)
			Expect(err).To(BeNil())
		})
	})
	Describe("List", func() {
		BeforeEach(func() {
			items := []store.Item{
				store.NewItem("/test", store.Value, "test"),
				store.NewItem("/test2", store.Value, "test2"),
				store.NewItem("/parent/child", store.Value, "child"),
				store.NewItem("/parent/child2", store.Value, "child2"),
			}
			for _, i := range items {
				err := memoryStore.Set(i)
				Expect(err).To(BeNil())
			}
		})
		When("Path Rooted", func() {
			It("Lists All Values", func() {
				items, err := memoryStore.List("/")
				Expect(err).To(BeNil())
				Expect(len(items)).To(Equal(4))
			})
		})
		When("Path Blank", func() {
			It("Lists All Values", func() {
				items, err := memoryStore.List("")
				Expect(err).To(BeNil())
				Expect(len(items)).To(Equal(4))
			})
		})
		When("Path Matches Parent", func() {
			It("Returns only children", func() {
				items, err := memoryStore.List("parent")
				Expect(err).To(BeNil())
				Expect(len(items)).To(Equal(2))
			})
		})
		When("Path Matches only one key", func() {
			It("Returns only that key", func() {
				items, err := memoryStore.List("test")
				Expect(err).To(BeNil())
				Expect(len(items)).To(Equal(1))
			})
		})
	})
})
