package memory

import (
	"github.com/patrickhuber/cli-mgr/config"
	"github.com/patrickhuber/cli-mgr/store"
)

type MemoryStoreProvider struct {
}

func (*MemoryStoreProvider) GetName() string {
	return ""
}

func (*MemoryStoreProvider) Create(configSource *config.ConfigSource) (store.Store, error) {
	return nil, nil
}
