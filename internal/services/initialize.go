package services

import (
	"github.com/patrickhuber/go-log"
	"github.com/patrickhuber/go-xplat/filepath"
	"github.com/patrickhuber/go-xplat/fs"
	"github.com/patrickhuber/wrangle/internal/config"
	"github.com/patrickhuber/wrangle/internal/global"
)

type initialize struct {
	fs     fs.FS
	logger log.Logger
	path   *filepath.Processor
}

type InitializeRequest struct {
	Directory string
	Force     bool
}

type Initialize interface {
	Execute(r *InitializeRequest) error
}

func NewInitialize(fs fs.FS, path *filepath.Processor, logger log.Logger) Initialize {
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

	err := i.fs.MkdirAll(localDotDir, 0664)
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
		ApiVersion: config.ConfigApiVersion,
		Spec: config.Spec{
			Feeds:       []config.Feed{},
			Stores:      []config.Store{},
			Environment: map[string]string{},
			Packages:    []config.Package{},
		},
	}

	p := config.NewFile(i.fs, localWrangleFile)
	return p.Write(cfg)
}
