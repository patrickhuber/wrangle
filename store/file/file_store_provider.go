package file

import (
	"github.com/patrickhuber/cli-mgr/config"
	"github.com/patrickhuber/cli-mgr/store"
	"github.com/spf13/afero"
)

type fileStoreProvider struct {
	fileSystem afero.Fs
}

// NewFileStoreProvider creates a new file store provider which implements the store.Provider interface
func NewFileStoreProvider(fileSystem afero.Fs) store.Provider {
	return &fileStoreProvider{fileSystem: fileSystem}
}

func (provider *fileStoreProvider) GetName() string {
	return "file"
}

func (provider *fileStoreProvider) Create(configSource *config.ConfigSource) (store.Store, error) {
	cfg, err := NewFileStoreConfig(configSource)
	if err != nil {
		return nil, err
	}
	store := NewFileStore(cfg.Name, cfg.Path, provider.fileSystem)
	return store, nil
}
