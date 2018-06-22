package file

import (
	"github.com/patrickhuber/cli-mgr/config"
	"github.com/pkg/errors"
)

type FileStoreConfig struct {
	Name string
	Path string
}

func NewFileStoreConfig(configSource *config.ConfigSource) (*FileStoreConfig, error) {
	cfg := &FileStoreConfig{}
	value, ok := configSource.Params["path"]
	if !ok {
		return nil, errors.New("unable to find required parameter 'path' in configuration source.")
	}
	cfg.Path = value
	cfg.Name = configSource.Name
	return cfg, nil
}
