package credhub

import (
	"github.com/patrickhuber/cli-mgr/config"
	"github.com/patrickhuber/cli-mgr/store"
)

type credHubConfigStoreProvider struct {
}

func NewCredHubStoreProvider() store.Provider {
	return &credHubConfigStoreProvider{}
}
func (provider *credHubConfigStoreProvider) GetName() string {
	return "credhub"
}

func (provider *credHubConfigStoreProvider) Create(configSource *config.ConfigSource) (store.Store, error) {
	configStoreConfig, err := NewCredHubStoreConfig(configSource)
	if err != nil {
		return nil, err
	}
	return NewCredHubStore(configStoreConfig)
}
