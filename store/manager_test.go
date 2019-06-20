package store_test

import (
	"github.com/patrickhuber/wrangle/config"
	"github.com/patrickhuber/wrangle/store"
	"github.com/patrickhuber/wrangle/store/memory"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("", func() {
	It("can register provider", func() {
		manager := store.NewManager()
		manager.Register(memory.NewMemoryStoreProvider())
		_, ok := manager.Get("memory")
		Expect(ok).To(BeTrue())
	})

	It("can create config store", func() {
		manager := store.NewManager()
		manager.Register(memory.NewMemoryStoreProvider())
		store, err := manager.Create(&config.Store{
			Name:      "test",
			Stores:    []string{"test"},
			StoreType: "memory",
		})
		Expect(err).To(BeNil())
		Expect(store).ToNot(BeNil())
	})

	Context("missing config store provider", func() {

		It("throws error", func() {
			manager := store.NewManager()
			_, err := manager.Create(&config.Store{Name: "test"})
			Expect(err).ToNot(BeNil())
		})
	})
})
