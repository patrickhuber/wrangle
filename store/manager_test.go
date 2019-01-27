package store_test

import (
	"github.com/patrickhuber/wrangle/config"
	"github.com/patrickhuber/wrangle/store"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type dummyConfigStoreProvider struct {
	name string
}

func (provider *dummyConfigStoreProvider) Name() string {
	return provider.name
}

func (provider *dummyConfigStoreProvider) Create(store *config.Store) (store.Store, error) {
	return &dummyConfigStore{}, nil
}

type dummyConfigStore struct {
}

func (s *dummyConfigStore) Delete(name string) error {
	return nil
}

func (s *dummyConfigStore) Get(name string) (store.Item, error) {
	return store.NewItem(name, nil), nil
}

func (s *dummyConfigStore) Name() string {
	return ""
}

func (s *dummyConfigStore) Type() string {
	return "dummy"
}

func (s *dummyConfigStore) Set(item store.Item) error {
	return nil
}

var _ = Describe("", func() {
	It("can register provider", func() {
		manager := store.NewManager()
		manager.Register(&dummyConfigStoreProvider{name: "test"})
		_, ok := manager.Get("test")
		Expect(ok).To(BeTrue())
	})

	It("can create config store", func() {
		manager := store.NewManager()
		manager.Register(&dummyConfigStoreProvider{name: "dummy"})
		store, err := manager.Create(&config.Store{
			Name:      "test",
			Stores:    []string{"test"},
			StoreType: "dummy",
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
