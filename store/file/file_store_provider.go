package file

import (
	"github.com/patrickhuber/wrangle/filesystem"
	"github.com/patrickhuber/wrangle/config"
	"github.com/patrickhuber/wrangle/crypto"
	"github.com/patrickhuber/wrangle/store"
	"github.com/pkg/errors"	
)

type fileStoreProvider struct {
	fileSystem filesystem.FileSystem
	factory    crypto.PgpFactory
}

// NewFileStoreProvider creates a new file store provider which implements the store.Provider interface
func NewFileStoreProvider(fileSystem filesystem.FileSystem, factory crypto.PgpFactory) store.Provider {
	return &fileStoreProvider{
		fileSystem: fileSystem,
		factory:    factory}
}

func (provider *fileStoreProvider) Name() string {
	return "file"
}

func (provider *fileStoreProvider) Create(configSource *config.Store) (store.Store, error) {
	cfg, err := NewFileStoreConfig(configSource)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create file store config")
	}

	var decrypter crypto.Decrypter
	if provider.factory != nil {
		decrypter, err = provider.factory.CreateDecrypter()
		if err != nil {
			return nil, err
		}
	}

	store, err := NewFileStore(cfg.Name, cfg.Path, provider.fileSystem, decrypter)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create file store")
	}
	return store, nil
}
