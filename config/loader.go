package config

import (
	"fmt"
	"os/user"

	"github.com/patrickhuber/wrangle/filepath"

	"github.com/spf13/afero"
)

type loader struct {
	fileSystem afero.Fs
}

// Loader - loads a config
type Loader interface {
	FileSystem() afero.Fs
	Load(configPath string) (*Config, error)
}

// NewLoader creates a new config loader
func NewLoader(fileSystem afero.Fs) Loader {
	return &loader{fileSystem: fileSystem}
}

func (loader *loader) FileSystem() afero.Fs {
	return loader.fileSystem
}

func (loader *loader) Load(configPath string) (*Config, error) {
	// load the config file
	ok, err := afero.Exists(loader.fileSystem, configPath)

	// if failure finding file, return the error
	if err != nil {
		return nil, err
	}

	// if not found, return error
	if !ok {
		return nil, fmt.Errorf(
			fmt.Sprintf("file %s does not exist", configPath))
	}

	// red the file contents and return a serialized Config struct
	data, err := afero.ReadFile(loader.fileSystem, configPath)
	if err != nil {
		return nil, err
	}
	return Serialize(data)
}

// GetDefaultConfigPath returns the default config path
func GetDefaultConfigPath() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	configDir := filepath.Join(usr.HomeDir, ".wrangle", "config.yml")
	configDir = filepath.ToSlash(configDir)
	return configDir, nil
}
