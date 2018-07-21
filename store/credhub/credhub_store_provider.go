package credhub

import (
	"github.com/patrickhuber/wrangle/config"
	"github.com/patrickhuber/wrangle/store"
)

type credHubConfigStoreProvider struct {
}

func NewCredHubStoreProvider() store.Provider {
	return &credHubConfigStoreProvider{}
}
func (provider *credHubConfigStoreProvider) Name() string {
	return "credhub"
}

func (provider *credHubConfigStoreProvider) Create(configSource *config.Store) (store.Store, error) {
	configStoreConfig, err := NewCredHubStoreConfig(configSource)
	if err != nil {
		return nil, err
	}
	return NewCredHubStore(configStoreConfig)
}
