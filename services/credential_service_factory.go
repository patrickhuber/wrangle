package services

import (
	"github.com/patrickhuber/wrangle/config"
	"github.com/patrickhuber/wrangle/store"
	"github.com/spf13/afero"
)

// CredentialServiceFactory defines a factory for creating credential services
type CredentialServiceFactory interface {
	Create(configFile string) (CredentialService, error)
}

type credentialServiceFactory struct {
	manager store.Manager
	fs      afero.Fs
}

// NewCredentialServiceFactory creates a new credential service factory
func NewCredentialServiceFactory(manager store.Manager, fs afero.Fs) CredentialServiceFactory {
	return &credentialServiceFactory{
		manager: manager,
		fs:      fs,
	}
}

func (factory *credentialServiceFactory) Create(configFile string) (CredentialService, error) {

	provider := config.NewFsProvider(factory.fs, configFile)
	cfg, err := provider.Get()
	if err != nil {
		return nil, err
	}

	graph, err := config.NewConfigurationGraph(cfg)
	if err != nil {
		return nil, err
	}

	return NewCredentialService(cfg, graph, factory.manager)
}
