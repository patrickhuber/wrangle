package services

import (
	"github.com/patrickhuber/wrangle/pkg/config"
	"github.com/patrickhuber/wrangle/pkg/filesystem"
	"github.com/patrickhuber/wrangle/pkg/ilog"
)

type initialize struct {
	fs             filesystem.FileSystem
	configProvider config.Provider
	logger         ilog.Logger
}

type InitializeRequest struct {
	ApplicationName string
	Force           bool
}

type Initialize interface {
	Execute(r *InitializeRequest) error
}

func NewInitialize(fs filesystem.FileSystem, provider config.Provider, logger ilog.Logger) Initialize {
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
