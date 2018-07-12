package memory

import (
	"github.com/patrickhuber/wrangle/config"
	"github.com/patrickhuber/wrangle/store"
)

type memoryStoreProvider struct {
}

// NewMemoryStoreProvider creates a new memory store provider
func NewMemoryStoreProvider() store.Provider {
	return &memoryStoreProvider{}
}

func (*memoryStoreProvider) GetName() string {
	return "memory"
}

func (*memoryStoreProvider) Create(configSource *config.Store) (store.Store, error) {

	return NewMemoryStore(configSource.Name), nil
}
