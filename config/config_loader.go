package config

import (
	"os/user"
	"path/filepath"

	"github.com/patrickhuber/cli-mgr/option"
	"github.com/spf13/afero"
	yaml "gopkg.in/yaml.v2"
)

// ConfigLoader - loads a config
type configLoader struct {
	fileSystem afero.Fs
}

type ConfigLoader interface {
	FileSystem() afero.Fs
	Load(configPath string) (*Config, error)
}

// NewConfigLoader creates a new config loader
func NewConfigLoader(fileSystem afero.Fs) ConfigLoader {
	return &configLoader{fileSystem: fileSystem}
}

func (loader *configLoader) FileSystem() afero.Fs {
	return loader.fileSystem
}

func (loader *configLoader) Load(configPath string) (*Config, error) {
	loader.ensureExists(configPath)
	data, err := afero.ReadFile(loader.fileSystem, configPath)
	if err != nil {
		return nil, err
	}
	config := &Config{}
	err = yaml.Unmarshal([]byte(data), config)
	return config, err
}

// GetConfigPath gets the config path from the options
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

func (loader *configLoader) ensureExists(configFile string) error {
	fileSystem := loader.FileSystem()
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
