package store

import (
	"github.com/patrickhuber/cli-mgr/config"
	"github.com/spf13/afero"
)

type StoreFactory struct {
}

func NewStoreFactory() *StoreFactory {
	return &StoreFactory{}
}

func (store *StoreFactory) Create(configSource *config.ConfigSource) (ConfigStore, error) {
	switch configSource.ConfigSourceType {
	case "memory":
		memoryStore := NewMemoryStore(configSource.Name)
		return memoryStore, nil
	case "file":
		fileStore := NewFileStore(configSource.Name, configSource.Params["path"], afero.OsFs{})
		return fileStore, nil
	case "credhub":
		credHubStoreConfig, err := NewCredHubStoreConfig(configSource)
		if err != nil {
			return nil, err
		}
		credHubStore, err := NewCredHubStore(
			credHubStoreConfig,
		)
		if err != nil {
			return nil, err
		}
		return credHubStore, nil
	}
	return nil, nil
}
