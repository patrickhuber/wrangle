package config

import (
	"github.com/spf13/afero"
)

type fsProvider struct {
	fs             afero.Fs
	configFilePath string
}

// NewFsProvider creates a new file system provider
func NewFsProvider(fs afero.Fs, configFilePath string) Provider {
	return &fsProvider{
		fs:             fs,
		configFilePath: configFilePath,
	}
}

func (provider *fsProvider) Initialize() (*Config, error) {
	return provider.Get()
}

func (provider *fsProvider) Get() (*Config, error) {

	// validate file exists
	ok, err := afero.Exists(provider.fs, provider.configFilePath)
	if err != nil {
		return nil, err
	}
	if !ok {
		cfg := &Config{
			Stores:    []Store{},
			Processes: []Process{},
			Imports:   []PackageReference{},
		}
		return cfg, provider.Set(cfg)
	}

	// open the file
	file, err := provider.fs.Open(provider.configFilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// read file and return
	reader := NewYamlReader(file)
	return reader.Read()
}

func (provider *fsProvider) Set(c *Config) error {
	// open the file
	file, err := provider.fs.Open(provider.configFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// write the config
	writer := NewYamlWriter(file)
	return writer.Write(c)
}
