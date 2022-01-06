package services

import (
	"github.com/patrickhuber/wrangle/pkg/config"
	"github.com/patrickhuber/wrangle/pkg/filesystem"
)

type initialize struct {
	fs     filesystem.FileSystem
	reader config.Reader
}

type InitializeRequest struct {
	ApplicationName  string
	GlobalConfigFile string
	Force            bool
}

type Initialize interface {
	Execute(r *InitializeRequest) error
}

func NewInitialize(fs filesystem.FileSystem, reader config.Reader) Initialize {
	return &initialize{
		fs:     fs,
		reader: reader,
	}
}

func (i *initialize) Execute(r *InitializeRequest) error {
	// does the global config exist?
	exists, err := i.fs.Exists(r.GlobalConfigFile)
	if err != nil {
		return err
	}

	if exists && !r.Force {
		return nil
	}

	cfg, err := i.reader.Get()
	if err != nil {
		return err
	}

	configProvider := config.NewFileProvider(i.fs, r.GlobalConfigFile)
	return configProvider.Set(cfg)
}
