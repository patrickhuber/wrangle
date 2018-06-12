package store

import (
	"fmt"

	"github.com/patrickhuber/cli-mgr/config"
)

type manager struct {
	providers map[string]Provider
}

// Manager manages providers and provides a factory method for creating Stores
type Manager interface {
	Get(name string) (Provider, bool)
	Register(provider Provider)
	Create(configSource *config.ConfigSource) (Store, error)
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
	manager.providers[provider.GetName()] = provider
}

func (manager *manager) Create(configSource *config.ConfigSource) (Store, error) {
	name := configSource.Name
	value, ok := manager.providers[name]
	if !ok {
		return nil, fmt.Errorf("Unable to find key '%s' in manager.Providers. Did you forget to register it?", name)
	}
	return value.Create(configSource)
}