package config

import (
	"fmt"

	"github.com/patrickhuber/wrangle/filepath"

	"github.com/spf13/afero"
)

type loader struct {
	fileSystem afero.Fs
}

// Loader - loads a config
type Loader interface {
	FileSystem() afero.Fs
	LoadConfig(configPath string) (*Config, error)
	LoadPackage(packagePath string) (*Package, error)
	LoadPackageAsInterface(packagePath string) (interface{}, error)
}

// NewLoader creates a new config loader
func NewLoader(fileSystem afero.Fs) Loader {
	return &loader{fileSystem: fileSystem}
}

func (loader *loader) FileSystem() afero.Fs {
	return loader.fileSystem
}

func (loader *loader) LoadConfig(configPath string) (*Config, error) {
	data, err := loader.loadFileData(configPath)
	if err != nil {
		return nil, err
	}
	return DeserializeConfig(data)
}

func (loader *loader) LoadPackage(packagePath string) (*Package, error) {
	data, err := loader.loadFileData(packagePath)
	if err != nil {
		return nil, err
	}
	return DeserializePackage(data)
}

func (loader *loader) LoadPackageAsInterface(packagePath string) (interface{}, error) {
	data, err := loader.loadFileData(packagePath)
	if err != nil {
		return nil, err
	}
	return DeserializePackageAsInterface(data)
}

func (loader *loader) loadFileData(path string) ([]byte, error) {
	// load the package file
	ok, err := afero.Exists(loader.fileSystem, path)

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
	data, err := afero.ReadFile(loader.fileSystem, path)
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
