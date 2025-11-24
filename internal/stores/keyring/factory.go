package keyring

import (
	"github.com/99designs/keyring"
	"github.com/patrickhuber/wrangle/internal/stores"
)

const name string = "keyring"
const AllowedBackends string = "allowed_backends"
const ServiceProperty string = "service"
const FileDirectory string = "file.directory"
const FilePassword string = "file.password"
const PassDirectory string = "pass.directory"
const PassPrefix string = "pass.prefix"
const PassCmd string = "pass.command"

type Factory struct {
}

func NewFactory() stores.Factory {
	return Factory{}
}

func (f Factory) Name() string {
	return name
}

func (f Factory) Create(properties map[string]any) (stores.Store, error) {
	config := &keyring.Config{}
	service, err := stores.GetRequiredProperty[string](properties, ServiceProperty)
	if err != nil {
		return nil, err
	}
	config.ServiceName = service

	fileDirectory, ok, err := stores.GetOptionalProperty[string](properties, FileDirectory)
	if err != nil {
		return nil, err
	}
	if ok {
		config.FileDir = fileDirectory
	}
	filePassword, ok, err := stores.GetOptionalProperty[string](properties, FilePassword)
	if err != nil {
		return nil, err
	}
	if ok {
		config.FilePasswordFunc = func(s string) (string, error) { return filePassword, nil }
	}
	passDirectory, ok, err := stores.GetOptionalProperty[string](properties, PassDirectory)
	if err != nil {
		return nil, err
	}
	if ok {
		config.PassDir = passDirectory
	}
	passPrefix, ok, err := stores.GetOptionalProperty[string](properties, PassPrefix)
	if err != nil {
		return nil, err
	}
	if ok {
		config.PassPrefix = passPrefix
	}
	passCmd, ok, err := stores.GetOptionalProperty[string](properties, PassCmd)
	if err != nil {
		return nil, err
	}
	if ok {
		config.PassCmd = passCmd
	}
	allowedBackends, ok, err := stores.GetOptionalProperty[[]string](properties, AllowedBackends)
	if err != nil {
		return nil, err
	}
	config.AllowedBackends = nil
	if ok {
		for _, backend := range allowedBackends {
			config.AllowedBackends = append(config.AllowedBackends, keyring.BackendType(backend))
		}
	}
	return NewVault(*config), nil
}
