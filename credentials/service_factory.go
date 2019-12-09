package credentials

import (
	"github.com/patrickhuber/wrangle/config"
	"github.com/patrickhuber/wrangle/filesystem"
	"github.com/patrickhuber/wrangle/store"
)

// ServiceFactory defines a factory for creating credential services
type ServiceFactory interface {
	Create(configFile string) (Service, error)
}

type serviceFactory struct {
	manager store.Manager
	fs      filesystem.FileSystem
}

// NewServiceFactory creates a new credential service factory
func NewServiceFactory(manager store.Manager, fs filesystem.FileSystem) ServiceFactory {
	return &serviceFactory{
		manager: manager,
		fs:      fs,
	}
}

func (factory *serviceFactory) Create(configFile string) (Service, error) {

	provider := config.NewFsProvider(factory.fs, configFile)
	cfg, err := provider.Get()
	if err != nil {
		return nil, err
	}

	graph, err := config.NewConfigurationGraph(cfg)
	if err != nil {
		return nil, err
	}

	return NewService(cfg, graph, factory.manager)
}
