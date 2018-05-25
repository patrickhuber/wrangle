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

func (loader *ConfigLoader) Load(configPath string) (*Config, error) {
	loader.ensureExists(configPath)
	data, err := afero.ReadFile(loader.FileSystem, configPath)
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

func (loader *ConfigLoader) ensureExists(configFile string) error {
	fileSystem := loader.FileSystem
	ok, err := afero.Exists(fileSystem, configFile)
	if err != nil {
		return err
	}
	if ok {
		return nil
	}
	data := "config-sources:\nprocesses:\n"
	return afero.WriteFile(fileSystem, configFile, []byte(data), 0644)
}
