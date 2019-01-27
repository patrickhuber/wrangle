package memory

import (
	"github.com/patrickhuber/wrangle/config"
	"github.com/patrickhuber/wrangle/store"
)

type memoryStoreProvider struct {
	stores map[string]store.Store
}

// NewMemoryStoreProviderWithMap creates a new memory store provider with the givin backing map
func NewMemoryStoreProviderWithMap(stores map[string]store.Store) store.Provider {
	return &memoryStoreProvider{
		stores: stores,
	}
}

// NewMemoryStoreProvider creates a new memory store provider
func NewMemoryStoreProvider(stores ...store.Store) store.Provider {
	storeMap := map[string]store.Store{}
	for _, s := range stores {
		storeMap[s.Name()] = s
	}
	return NewMemoryStoreProviderWithMap(storeMap)
}

func (p *memoryStoreProvider) Name() string {
	return "memory"
}

func (p *memoryStoreProvider) Create(configSource *config.Store) (store.Store, error) {
	if s, ok := p.stores[configSource.Name]; ok {
		return s, nil
	}
	s := NewMemoryStore(configSource.Name)
	p.stores[configSource.Name] = s
	return s, nil
}

func (p *memoryStoreProvider) Stores() []store.Store {
	values := []store.Store{}
	for _, v := range p.stores {
		values = append(values, v)
	}
	return values
}
