package memory

import (
	"github.com/patrickhuber/cli-mgr/config"
	"github.com/patrickhuber/cli-mgr/store"
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

func (*memoryStoreProvider) Create(configSource *config.ConfigSource) (store.Store, error) {

	return NewMemoryStore(configSource.Name), nil
}
