package config

import (
	"fmt"
	"os"
	"path"

	"github.com/patrickhuber/wrangle/pkg/filesystem"
	"gopkg.in/yaml.v2"
)

type fileProvider struct {
	fs         filesystem.FileSystem
	properties Properties
}

// NewFileProvider creates a new file Provider with the given filesystem and file path
func NewFileProvider(fs filesystem.FileSystem, properties Properties) Provider {
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

	// make sure the path exists
	ok, err := p.fs.Exists(globalConfigFilePath)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, os.ErrNotExist
	}

	// make sure the file is a file and not a directory
	ok, err = p.fs.IsDir(globalConfigFilePath)
	if err != nil {
		return nil, err
	}
	if ok {
		return nil, fmt.Errorf("the configuration path (%s) must be a file", globalConfigFilePath)
	}

	// read the data
	data, err := p.fs.Read(globalConfigFilePath)
	if err != nil {
		return nil, err
	}

	// write the config
	config := &Config{}
	err = yaml.UnmarshalStrict(data, config)
	return config, err
}

func (p *fileProvider) Lookup() (*Config, bool, error) {

	globalConfigFilePath, err := p.globalConfigFilePath()
	if err != nil {
		return nil, false, err
	}

	// make sure the path exists
	ok, err := p.fs.Exists(globalConfigFilePath)
	if err != nil {
		return nil, false, err
	}
	if !ok {
		return nil, false, nil
	}

	// make sure the file is a file and not a directory
	ok, err = p.fs.IsDir(globalConfigFilePath)
	if err != nil {
		return nil, false, err
	}
	if ok {
		return nil, false, fmt.Errorf("the configuration path (%s) must be a file", globalConfigFilePath)
	}

	// read the data
	data, err := p.fs.Read(globalConfigFilePath)
	if err != nil {
		return nil, false, err
	}

	// write the config
	config := &Config{}
	err = yaml.UnmarshalStrict(data, config)
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
	return p.fs.Write(globalConfigFilePath, data, 0644)
}

func (p *fileProvider) globalConfigFilePath() (string, error) {

	path, ok := p.properties.Lookup(GlobalConfigFilePathProperty)
	if !ok {
		return "", fmt.Errorf("missing global config file property")
	}
	return path, nil
}
