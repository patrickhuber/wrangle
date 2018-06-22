package file

import (
	"github.com/patrickhuber/cli-mgr/config"
	"github.com/patrickhuber/cli-mgr/store"
	"github.com/pkg/errors"
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
		return nil, errors.Wrap(err, "unable to create file store config")
	}
	store, err := NewFileStore(cfg.Name, cfg.Path, provider.fileSystem)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create file store")
	}
	return store, nil
}