package file

import (
	"github.com/patrickhuber/cli-mgr/config"
	"github.com/spf13/afero"
)

type FileConfigStoreProvider struct {
}

func (provider *FileConfigStoreProvider) GetName() string {
	return "file"
}

func (provider *FileConfigStoreProvider) Create(configSource *config.ConfigSource) (config.ConfigStore, error) {
	cfg, err := NewFileConfigStoreConfig(configSource)
	if err != nil {
		return nil, err
	}
	store := NewFileStore(cfg.Name, cfg.Path, afero.NewOsFs())
	return store, nil
}