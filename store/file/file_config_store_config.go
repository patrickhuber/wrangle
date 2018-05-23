package file

import "github.com/patrickhuber/cli-mgr/config"

type FileConfigStoreConfig struct {
	Name string
	Path string
}

func NewFileConfigStoreConfig(configSource *config.ConfigSource) (*FileConfigStoreConfig, error) {
	cfg := &FileConfigStoreConfig{}
	if value, ok := configSource.Params["path"]; ok {
		cfg.Path = value
	}
	cfg.Name = configSource.Name
	return cfg, nil
}
