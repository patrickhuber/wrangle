package credhub

import "github.com/patrickhuber/cli-mgr/config"

type CredHubConfigStoreProvider struct {
}

func (provider *CredHubConfigStoreProvider) GetName() string {
	return "credhub"
}

func (provider *CredHubConfigStoreProvider) Create(configSource *config.ConfigSource) (config.ConfigStore, error) {
	configStoreConfig, err := NewCredHubStoreConfig(configSource)
	if err != nil {
		return nil, err
	}
	return NewCredHubStore(configStoreConfig)
}
