package meta

import (
	"github.com/patrickhuber/wrangle/config"
	"github.com/patrickhuber/wrangle/store"
)

type metaStoreProvider struct {
	configFile string
}

// NewMetaStoreProvider creates a new meta store provider
func NewMetaStoreProvider(configFile string) store.Provider {
	return &metaStoreProvider{
		configFile: configFile,
	}
}

func (*metaStoreProvider) Name() string {
	return "meta"
}

func (p *metaStoreProvider) Create(configSource *config.Store) (store.Store, error) {
	return NewMetaStore(configSource.Name, p.configFile), nil
}
