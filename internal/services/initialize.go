package services

import (
	"github.com/patrickhuber/go-cross/filepath"
	"github.com/patrickhuber/go-cross/fs"
	"github.com/patrickhuber/go-log"
	"github.com/patrickhuber/wrangle/internal/config"
	"github.com/patrickhuber/wrangle/internal/global"
)

type initialize struct {
	fs     fs.FS
	logger log.Logger
	path   filepath.Provider
}

type InitializeRequest struct {
	Directory string
	Force     bool
}

type Initialize interface {
	Execute(r *InitializeRequest) error
}

func NewInitialize(fs fs.FS, path filepath.Provider, logger log.Logger) Initialize {
	return &initialize{
		fs:     fs,
		logger: logger,
		path:   path,
	}
}

func (i *initialize) Execute(r *InitializeRequest) error {
	i.logger.Debugln("initalize")
	i.logger.Infof("initializing in local directory '%s'", r.Directory)

	localDotDir := i.path.Join(r.Directory, global.LocalConfigurationDirectoryName)
	i.logger.Infof("creating '%s'", localDotDir)

	err := i.fs.MkdirAll(localDotDir, 0775)
	if err != nil {
		return err
	}

	localWrangleFile := i.path.Join(r.Directory, global.LocalConfigurationFileName)

	i.logger.Infof("creating '%s'", localWrangleFile)
	exists, err := i.fs.Exists(localWrangleFile)
	if err != nil {
		return err
	}

	if exists {
		if !r.Force {
			i.logger.Infof("file '%s' exists, force=false, skipping create")
			return nil
		}
		i.logger.Infof("force = true, overwriting '%s'")
	}
	cfg := config.Config{
		ApiVersion: config.ApiVersion,
		Kind:       config.Kind,
		Spec: config.Spec{
			Feeds:       []config.Feed{},
			Stores:      []config.Store{},
			Environment: map[string]string{},
			Packages:    []config.Package{},
		},
	}
	return config.WriteFile(i.fs, localWrangleFile, cfg)
}
