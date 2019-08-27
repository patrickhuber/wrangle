package config

import (
	"fmt"

	"github.com/patrickhuber/wrangle/filepath"

	"github.com/patrickhuber/wrangle/filesystem"
)

type loader struct {
	fileSystem filesystem.FileSystem
}

// Loader - loads a config
type Loader interface {
	FileSystem() filesystem.FileSystem
	LoadConfig(configPath string) (*Config, error)
}

// NewLoader creates a new config loader
func NewLoader(fileSystem filesystem.FileSystem) Loader {
	return &loader{fileSystem: fileSystem}
}

func (loader *loader) FileSystem() filesystem.FileSystem {
	return loader.fileSystem
}

func (loader *loader) LoadConfig(configPath string) (*Config, error) {
	data, err := loader.loadFileData(configPath)
	if err != nil {
		return nil, err
	}
	return DeserializeConfig(data)
}

func (loader *loader) loadFileData(path string) ([]byte, error) {
	// load the package file
	ok, err := loader.fileSystem.Exists(path)

	// if failure finding file, return the error
	if err != nil {
		return nil, err
	}

	// if not found, return error
	if !ok {
		return nil, fmt.Errorf(
			fmt.Sprintf("file %s does not exist", path))
	}

	// red the file contents and return a serialized Config struct
	data, err := loader.fileSystem.Read(path)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// GetDefaultConfigPath returns the default config path
func GetDefaultConfigPath(workingDirectory string) (string, error) {
	configDir := filepath.Join(workingDirectory, "config.yml")
	configDir = filepath.ToSlash(configDir)
	return configDir, nil
}
