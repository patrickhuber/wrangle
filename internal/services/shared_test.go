package services_test

import (
	"fmt"

	"github.com/patrickhuber/wrangle/internal/services"
	"github.com/patrickhuber/wrangle/internal/types"
	"github.com/patrickhuber/wrangle/pkg/config"
	"github.com/patrickhuber/wrangle/pkg/di"
	"github.com/patrickhuber/wrangle/pkg/filesystem"
	"github.com/patrickhuber/wrangle/pkg/operatingsystem"
)

func ResolveFileSystem(container di.Container) (filesystem.FileSystem, error) {
	obj, err := container.Resolve(types.FileSystem)
	if err != nil {
		return nil, err
	}
	fs, ok := obj.(filesystem.FileSystem)
	if !ok {
		return nil, fmt.Errorf("Unable to cast filesystem")
	}
	return fs, nil
}

func ResolveOperatingSystem(container di.Container) (operatingsystem.OS, error) {
	obj, err := container.Resolve(types.OS)
	if err != nil {
		return nil, err
	}
	opsys, ok := obj.(operatingsystem.OS)
	if !ok {
		return nil, fmt.Errorf("Unable to cast filesystem")
	}
	return opsys, nil
}

func ResolveConfigReader(container di.Container) (config.Reader, error) {
	obj, err := container.Resolve(types.ConfigReader)
	if err != nil {
		return nil, err
	}
	reader, ok := obj.(config.Reader)
	if !ok {
		return nil, fmt.Errorf("Unable to cast config reader")
	}
	return reader, nil
}

func ResolveInstallService(container di.Container) (services.Install, error) {
	obj, err := container.Resolve(types.InstallService)
	if err != nil {
		return nil, err
	}
	reader, ok := obj.(services.Install)
	if !ok {
		return nil, fmt.Errorf("Unable to cast install service")
	}
	return reader, nil
}

func ResolveBootstrapService(container di.Container) (services.Bootstrap, error) {
	obj, err := container.Resolve(types.BootstrapService)
	if err != nil {
		return nil, err
	}
	bootstrap, ok := obj.(services.Bootstrap)
	if !ok {
		return nil, fmt.Errorf("Unable to cast bootstrap service")
	}
	return bootstrap, nil
}