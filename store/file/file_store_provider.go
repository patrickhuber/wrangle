package file

import (
	"github.com/patrickhuber/cli-mgr/config"
	"github.com/patrickhuber/cli-mgr/store"
	"github.com/spf13/afero"
)

type fileStoreProvider struct {
}

// NewFileStoreProvider creates a new file store provider which implements the store.Provider interface
func NewFileStoreProvider() store.Provider {
	return &fileStoreProvider{}
}

func (provider *fileStoreProvider) GetName() string {
	return "file"
}

func (provider *fileStoreProvider) Create(configSource *config.ConfigSource) (store.Store, error) {
	cfg, err := NewFileStoreConfig(configSource)
	if err != nil {
		return nil, err
	}
	store := NewFileStore(cfg.Name, cfg.Path, afero.NewOsFs())
	return store, nil
}
