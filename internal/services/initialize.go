package services

import (
	"github.com/patrickhuber/wrangle/pkg/config"
	"github.com/patrickhuber/wrangle/pkg/filesystem"
)

type initialize struct {
	fs             filesystem.FileSystem
	configProvider config.Provider
}

type InitializeRequest struct {
	ApplicationName string
	Force           bool
}

type Initialize interface {
	Execute(r *InitializeRequest) error
}

func NewInitialize(fs filesystem.FileSystem, provider config.Provider) Initialize {
	return &initialize{
		fs:             fs,
		configProvider: provider,
	}
}

func (i *initialize) Execute(r *InitializeRequest) error {

	cfg, err := i.configProvider.Get()
	if err != nil {
		return err
	}

	return i.configProvider.Set(cfg)
}
