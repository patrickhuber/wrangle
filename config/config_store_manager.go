package config

import "fmt"

type ConfigStoreManager struct {
	Providers map[string]ConfigStoreProvider
}

func NewConfigStoreManager() *ConfigStoreManager {
	return &ConfigStoreManager{
		Providers: make(map[string]ConfigStoreProvider),
	}
}

func (manager *ConfigStoreManager) Register(provider ConfigStoreProvider) {
	manager.Providers[provider.GetName()] = provider
}

func (manager *ConfigStoreManager) Create(configSource *ConfigSource) (ConfigStore, error) {
	name := configSource.Name
	value, ok := manager.Providers[name]
	if !ok {
		return nil, fmt.Errorf("Unable to find key '%s' in manager.Providers. Did you forget to register it?", name)
	}
	return value.Create(configSource)
}
