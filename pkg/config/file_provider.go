package config

import (
	"fmt"
	"os"
	"path"

	"github.com/patrickhuber/wrangle/pkg/filesystem"
	"gopkg.in/yaml.v2"
)

type fileProvider struct {
	fs   filesystem.FileSystem
	path string
}

// NewFileProvider creates a new file Provider with the given filesystem and file path
func NewFileProvider(fs filesystem.FileSystem, path string) Provider {
	return &fileProvider{
		fs:   fs,
		path: path,
	}
}

func (p *fileProvider) Get() (*Config, error) {

	// make sure the path exists
	ok, err := p.fs.Exists(p.path)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, os.ErrNotExist
	}

	// make sure the file is a file and not a directory
	ok, err = p.fs.IsDir(p.path)
	if err != nil {
		return nil, err
	}
	if ok {
		return nil, fmt.Errorf("the configuration path (%s) must be a file", p.path)
	}

	// read the data
	data, err := p.fs.Read(p.path)
	if err != nil {
		return nil, err
	}

	// write the config
	config := &Config{}
	err = yaml.UnmarshalStrict(data, config)
	return config, err
}

func (p *fileProvider) Set(config *Config) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	dir := path.Dir(p.path)
	err = p.fs.MkdirAll(dir, 0644)
	if err != nil {
		return err
	}
	return p.fs.Write(p.path, data, 0644)
}
