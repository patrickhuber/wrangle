package store

import (
	"testing"

	"github.com/patrickhuber/wrangle/config"
	"github.com/stretchr/testify/require"
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

func TestManager(t *testing.T) {

	t.Run("CanRegisterProvider", func(t *testing.T) {
		r := require.New(t)
		manager := NewManager()
		manager.Register(&dummyConfigStoreProvider{name: "test"})
		_, ok := manager.Get("test")
		r.True(ok)
	})

	t.Run("CanCreateConfigStore", func(t *testing.T) {
		r := require.New(t)
		manager := NewManager()
		manager.Register(&dummyConfigStoreProvider{name: "dummy"})
		store, err := manager.Create(&config.Store{
			Name:      "test",
			Stores:    []string{"test"},
			StoreType: "dummy",
		})
		r.Nil(err)
		r.NotNil(store)
	})

	t.Run("MissingConfigStoreProviderThrowsError", func(t *testing.T) {
		r := require.New(t)
		manager := NewManager()
		_, err := manager.Create(&config.Store{Name: "test"})
		r.NotNil(err)
	})
}
