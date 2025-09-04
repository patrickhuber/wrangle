package initialize

import (
	"fmt"

	"github.com/patrickhuber/go-cross/filepath"
	"github.com/patrickhuber/go-cross/fs"
	"github.com/patrickhuber/go-log"
	"github.com/patrickhuber/wrangle/internal/config"
	"github.com/patrickhuber/wrangle/internal/global"
)

type service struct {
	fs     fs.FS
	logger log.Logger
	path   filepath.Provider
}

type Request struct {
	Directory string // defaults to ./.wrangle
	Force     bool
}

type Service interface {
	Execute(r *Request) error
}

func NewService(fs fs.FS, path filepath.Provider, logger log.Logger) Service {
	return &service{
		fs:     fs,
		logger: logger,
		path:   path,
	}
}

func (i *service) Execute(r *Request) error {
	i.logger.Debugln("initalize")

	localWrangleFile := i.path.Join(r.Directory, global.LocalConfigurationFileName)

	i.logger.Infof("creating '%s'", localWrangleFile)
	exists, err := i.fs.Exists(localWrangleFile)
	if err != nil {
		return err
	}

	if exists {
		if !r.Force {
			return fmt.Errorf("file '%s' exists, force == false, skipping create", localWrangleFile)
		}
		i.logger.Infof("force == true, overwriting '%s'")
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
