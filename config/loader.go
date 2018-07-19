package config

import (
	"bytes"
	"fmt"
	"os/user"
	"path/filepath"

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
	err := loader.ensureExists(configPath)
	if err != nil {
		return nil, err
	}
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

func (loader *loader) ensureExists(configFile string) error {
	fileSystem := loader.FileSystem()
	ok, err := afero.Exists(fileSystem, configFile)
	if err != nil {
		return err
	}
	if ok {
		return nil
	}

	configDirectory := filepath.Dir(configFile)
	err = loader.ensureDirectoryExists(configDirectory)
	if err != nil {
		return err
	}

	data := &bytes.Buffer{}
	fmt.Fprintln(data, "stores:")
	fmt.Fprintln(data, "environments:")
	fmt.Fprintln(data, "packages:")

	return afero.WriteFile(fileSystem, configFile, data.Bytes(), 0644)
}

func (loader *loader) ensureDirectoryExists(directory string) error {
	fileSystem := loader.FileSystem()
	ok, err := afero.DirExists(fileSystem, directory)
	if err != nil {
		return err
	}
	if !ok {
		err = fileSystem.MkdirAll(directory, 0644)
		if err != nil {
			return err
		}
	}
	return nil
}
