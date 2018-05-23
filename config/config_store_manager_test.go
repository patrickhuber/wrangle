package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type dummyConfigStoreProvider struct {
	Name string
}

func (provider *dummyConfigStoreProvider) GetName() string {
	return provider.Name
}

func (provider *dummyConfigStoreProvider) Create(configSource *ConfigSource) (ConfigStore, error) {
	return &MemoryStore{}, nil
}

func TestConfigStoreManager(t *testing.T) {

	t.Run("CanRegisterProvider", func(t *testing.T) {
		r := require.New(t)
		manager := NewConfigStoreManager()
		manager.Register(&dummyConfigStoreProvider{Name: "test"})
		r.Equal(1, len(manager.Providers))
	})

	t.Run("CanCreateConfigStore", func(t *testing.T) {
		r := require.New(t)
		manager := NewConfigStoreManager()
		manager.Register(&dummyConfigStoreProvider{Name: "test"})
		store, err := manager.Create(&ConfigSource{
			Name: "test",
		})
		r.Nil(err)
		r.NotNil(store)
	})

	t.Run("MissingConfigStoreProviderThrowsError", func(t *testing.T) {
		r := require.New(t)
		manager := NewConfigStoreManager()
		_, err := manager.Create(&ConfigSource{Name: "test"})
		r.NotNil(err)
	})
}
