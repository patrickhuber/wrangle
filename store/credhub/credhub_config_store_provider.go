package credhub

import (
	"github.com/patrickhuber/cli-mgr/config"
	"github.com/patrickhuber/cli-mgr/store"
)

type CredHubConfigStoreProvider struct {
}

func (provider *CredHubConfigStoreProvider) GetName() string {
	return "credhub"
}

func (provider *CredHubConfigStoreProvider) Create(configSource *config.ConfigSource) (store.Store, error) {
	configStoreConfig, err := NewCredHubStoreConfig(configSource)
	if err != nil {
		return nil, err
	}
	return NewCredHubStore(configStoreConfig)
}
