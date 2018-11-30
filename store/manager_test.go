package store

import (
	"github.com/patrickhuber/wrangle/config"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type dummyConfigStoreProvider struct {
	name string
}

func (provider *dummyConfigStoreProvider) Name() string {
	return provider.name
}

func (provider *dummyConfigStoreProvider) Create(store *config.Store) (Store, error) {
	return &dummyConfigStore{}, nil
}

type dummyConfigStore struct {
}

func (store *dummyConfigStore) Delete(name string) (int, error) {
	return 0, nil
}

func (store *dummyConfigStore) GetByName(name string) (Data, error) {
	return &data{}, nil
}

func (store *dummyConfigStore) Name() string {
	return ""
}

func (store *dummyConfigStore) Type() string {
	return "dummy"
}

func (store *dummyConfigStore) Put(key string, value string) (string, error) {
	return "", nil
}

var _ = Describe("", func() {
	It("can register provider", func() {
		manager := NewManager()
		manager.Register(&dummyConfigStoreProvider{name: "test"})
		_, ok := manager.Get("test")
		Expect(ok).To(BeTrue())
	})

	It("can create config store", func() {
		manager := NewManager()
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
			manager := NewManager()
			_, err := manager.Create(&config.Store{Name: "test"})
			Expect(err).ToNot(BeNil())
		})
	})
})
