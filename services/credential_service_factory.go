package services

import (
	"github.com/patrickhuber/wrangle/config"
	"github.com/patrickhuber/wrangle/store"
	"github.com/patrickhuber/wrangle/templates"
	"github.com/spf13/afero"
)

// CredentialServiceFactory defines a factory for creating credential services
type CredentialServiceFactory interface {
	Create(configFile string) (CredentialService, error)
}

type credentialServiceFactory struct {
	manager         store.Manager
	fs              afero.Fs
	templateFactory templates.Factory
}

// NewCredentialServiceFactory creates a new credential service factory
func NewCredentialServiceFactory(manager store.Manager, fs afero.Fs, templateFactory templates.Factory) CredentialServiceFactory {
	return &credentialServiceFactory{
		manager:         manager,
		fs:              fs,
		templateFactory: templateFactory,
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

	return NewCredentialService(cfg, graph, factory.manager, factory.templateFactory)
}
