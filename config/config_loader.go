package config

import (
	"os/user"
	"path/filepath"

	"github.com/patrickhuber/cli-mgr/option"
	"github.com/spf13/afero"
	yaml "gopkg.in/yaml.v2"
)

// ConfigLoader - loads a config
type ConfigLoader struct {
	FileSystem afero.Fs
}

func (loader *ConfigLoader) Load(op *option.Options) (*Config, error) {
	path, err := GetConfigPath(op)
	if err != nil {
		return nil, err
	}
	data, err := afero.ReadFile(loader.FileSystem, path)
	if err != nil {
		return nil, err
	}
	config := &Config{}
	err = yaml.Unmarshal([]byte(data), config)
	return config, err
}

func GetConfigPath(op *option.Options) (string, error) {
	if op.ConfigPath != "" {
		return op.ConfigPath, nil
	}
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	configDir := filepath.Join(usr.HomeDir, ".cli-mgr", "config.yml")
	return configDir, nil
}
