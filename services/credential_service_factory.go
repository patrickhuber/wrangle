package services

import (
	"github.com/patrickhuber/wrangle/config"
	"github.com/patrickhuber/wrangle/store"
)

type CredentialServiceFactory interface {
	Create(configFile string) (CredentialService, error)
}

type credentialServiceFactory struct {
	manager store.Manager
	loader  config.Loader
}

func NewCredentialServiceFactory(manager store.Manager, loader config.Loader) CredentialServiceFactory {
	return &credentialServiceFactory{
		manager: manager,
		loader:  loader,
	}
}

func (f *credentialServiceFactory) Create(configFile string) (CredentialService, error) {
	cfg, err := f.loader.LoadConfig(configFile)
	if err != nil {
		return nil, err
	}
	graph, err := config.NewConfigurationGraph(cfg)
	if err != nil {
		return nil, err
	}
	return NewCredentialService(cfg, graph, f.manager)
}
