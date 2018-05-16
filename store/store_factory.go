package store

import "github.com/patrickhuber/cli-mgr/config"

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
	}
	return nil
}
