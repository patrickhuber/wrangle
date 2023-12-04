package config

import (
	"regexp"

	"github.com/patrickhuber/go-xplat/filepath"
	"github.com/patrickhuber/go-xplat/fs"
	"github.com/patrickhuber/go-xplat/os"
)

// LocalProvider config provider provides the local default configuration as well as fetches the existing configurations
type LocalProvider struct {
	os           os.OS
	fs           fs.FS
	path         *filepath.Processor
	localDefault Config
}

func NewLocalProvider(localDefault Config, os os.OS, fs fs.FS, path *filepath.Processor) LocalProvider {
	return LocalProvider{
		os:           os,
		fs:           fs,
		path:         path,
		localDefault: localDefault,
	}
}

func (l LocalProvider) Default() Config {
	return l.localDefault
}

const (
	LocalConfigFilePattern = "[.]wrangle[.](yml|yaml|json)"
)

func (l LocalProvider) Get() ([]Config, error) {

	// look in the current directory
	pwd, err := l.os.WorkingDirectory()
	if err != nil {
		return nil, err
	}

	// work up the directory hierarchy to find every path to the root
	current := pwd
	dirs := []string{}
	for {
		dirs = append(dirs, current)

		parent := l.path.Dir(current)
		if parent == current {
			break
		}

		current = parent
	}

	// loop through all the directories looking for configuration files
	var cfgs []Config
	for _, dir := range dirs {

		files, err := l.fs.ReadDir(dir)
		if err != nil {
			return nil, err
		}

		reg := regexp.MustCompile(LocalConfigFilePattern)

		// for each path match only files that match the config file pattern
		for _, file := range files {

			if file.IsDir() {
				continue
			}

			if !reg.MatchString(file.Name()) {
				continue
			}

			filePath := l.path.Join(dir, file.Name())

			// load the configuration file
			cfg, err := NewFile(l.fs, filePath).Read()
			if err != nil {
				return nil, err
			}

			// set the path metadata
			cfg.Metadata["file"] = filePath
			cfgs = append(cfgs, cfg)
		}
	}
	return cfgs, nil
}

func NewLocalDefault() Config {
	return Config{
		ApiVersion: ConfigApiVersion,
		Metadata:   map[string]string{},
		Spec: Spec{
			Feeds:       []Feed{},
			Environment: map[string]string{},
			Stores:      []Store{},
			Packages:    []Package{},
		},
	}
}
