package file

import "github.com/patrickhuber/cli-mgr/config"

type FileStoreConfig struct {
	Name string
	Path string
}

func NewFileStoreConfig(configSource *config.ConfigSource) (*FileStoreConfig, error) {
	cfg := &FileStoreConfig{}
	if value, ok := configSource.Params["path"]; ok {
		cfg.Path = value
	}
	cfg.Name = configSource.Name
	return cfg, nil
}
