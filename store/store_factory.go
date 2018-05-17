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

func (store *StoreFactory) Create(configSource *config.ConfigSource) ConfigStore {
	switch configSource.ConfigSourceType {
	case "memory":
		memoryStore := NewMemoryStore(configSource.Name)
		return memoryStore
	case "file":
		fileStore := NewFileStore(configSource.Name, configSource.Params["path"], afero.OsFs{})
		return fileStore
	}
	return nil
}
