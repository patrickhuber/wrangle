package store

import (
	"fmt"

	"github.com/patrickhuber/wrangle/config"
)

type manager struct {
	providers map[string]Provider
}

// Manager manages providers and provides a factory method for creating Stores
type Manager interface {
	Get(name string) (Provider, bool)
	Register(provider Provider)
	Create(configSource *config.Store) (Store, error)
}

// NewManager creates a new manager
func NewManager() Manager {
	return &manager{
		providers: make(map[string]Provider),
	}

}

func (manager *manager) Get(name string) (Provider, bool) {
	value, ok := manager.providers[name]
	return value, ok
}

func (manager *manager) Register(provider Provider) {
	manager.providers[provider.Name()] = provider
}

func (manager *manager) Create(configSource *config.Store) (Store, error) {
	configSourceType := configSource.StoreType
	provider, ok := manager.Get(configSourceType)
	if !ok {
		return nil, fmt.Errorf("Unable to find key '%s' in manager.Providers. Did you forget to register it?", configSourceType)
	}
	return provider.Create(configSource)
}
