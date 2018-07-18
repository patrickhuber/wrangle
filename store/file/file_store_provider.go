package file

import (
	"github.com/patrickhuber/wrangle/config"
	"github.com/patrickhuber/wrangle/crypto"
	"github.com/patrickhuber/wrangle/store"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

type fileStoreProvider struct {
	fileSystem afero.Fs
	factory    crypto.PgpFactory
}

// NewFileStoreProvider creates a new file store provider which implements the store.Provider interface
func NewFileStoreProvider(fileSystem afero.Fs, factory crypto.PgpFactory) store.Provider {
	return &fileStoreProvider{fileSystem: fileSystem, factory: factory}
}

func (provider *fileStoreProvider) GetName() string {
	return "file"
}

func (provider *fileStoreProvider) Create(configSource *config.Store) (store.Store, error) {
	cfg, err := NewFileStoreConfig(configSource)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create file store config")
	}

	var decryptor crypto.Decryptor
	if provider.factory != nil {
		decryptor, err = provider.factory.CreateDecryptor()
		if err != nil {
			return nil, err
		}
	}

	store, err := NewFileStore(cfg.Name, cfg.Path, provider.fileSystem, decryptor)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create file store")
	}
	return store, nil
}
