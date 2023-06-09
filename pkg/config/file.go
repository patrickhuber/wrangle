package config

import (
	"errors"
	"fmt"
	iofs "io/fs"
	"path"

	"github.com/patrickhuber/go-xplat/fs"
	"gopkg.in/yaml.v3"
)

type fileProvider struct {
	fs         fs.FS
	properties Properties
}

// NewFileProvider creates a new file Provider with the given filesystem and file path
func NewFileProvider(fs fs.FS, properties Properties) Provider {
	return &fileProvider{
		fs:         fs,
		properties: properties,
	}
}

func (p *fileProvider) Get() (*Config, error) {

	globalConfigFilePath, err := p.globalConfigFilePath()
	if err != nil {
		return nil, err
	}

	cfg, ok, err := p.Lookup()
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("unable to find file %s", globalConfigFilePath)
	}
	return cfg, nil
}

func (p *fileProvider) Lookup() (*Config, bool, error) {

	globalConfigFilePath, err := p.globalConfigFilePath()
	if err != nil {
		return nil, false, err
	}

	// check if the file exists
	fi, err := p.fs.Stat(globalConfigFilePath)
	if err != nil {
		if errors.Is(err, iofs.ErrNotExist) {
			return nil, false, nil
		}
		return nil, false, err
	}

	if fi.IsDir() {
		return nil, false, fmt.Errorf("the configuration path (%s) must be a file", globalConfigFilePath)
	}

	// read the data
	data, err := p.fs.ReadFile(globalConfigFilePath)
	if err != nil {
		return nil, false, err
	}

	// write the config
	config := &Config{}
	err = yaml.Unmarshal(data, config)
	return config, false, err
}

func (p *fileProvider) Set(config *Config) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	globalConfigFilePath, err := p.globalConfigFilePath()
	if err != nil {
		return err
	}
	dir := path.Dir(globalConfigFilePath)
	err = p.fs.MkdirAll(dir, 0644)
	if err != nil {
		return err
	}
	return p.fs.WriteFile(globalConfigFilePath, data, 0644)
}

func (p *fileProvider) globalConfigFilePath() (string, error) {

	path, ok := p.properties.Lookup(GlobalConfigFilePathProperty)
	if !ok {
		return "", fmt.Errorf("missing global config file property")
	}
	return path, nil
}
