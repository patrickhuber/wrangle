package services

import (
	"github.com/patrickhuber/go-log"
	"github.com/patrickhuber/wrangle/pkg/config"
	"github.com/patrickhuber/wrangle/pkg/filesystem"
)

type initialize struct {
	fs             filesystem.FileSystem
	configProvider config.Provider
	logger         log.Logger
}

type InitializeRequest struct {
	ApplicationName string
	Force           bool
}

type Initialize interface {
	Execute(r *InitializeRequest) error
}

func NewInitialize(fs filesystem.FileSystem, provider config.Provider, logger log.Logger) Initialize {
	return &initialize{
		fs:             fs,
		configProvider: provider,
		logger:         logger,
	}
}

func (i *initialize) Execute(r *InitializeRequest) error {
	i.logger.Debugln("initalize")

	i.logger.Debugln("get configuration")
	cfg, err := i.configProvider.Get()
	if err != nil {
		return err
	}

	i.logger.Debugln("set configuration")
	return i.configProvider.Set(cfg)
}
