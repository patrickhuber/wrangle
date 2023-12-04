package services

import (
	"github.com/patrickhuber/wrangle/internal/config"
)

type Configuration struct {
	Local  config.LocalProvider
	Global config.GlobalProvider
}

func (c Configuration) Get() (config.Config, error) {
	globalConfig, err := c.Global.Get()
	if err != nil {
		return config.Config{}, err
	}
	localConfigurations, err := c.Local.Get()
	if err != nil {
		return config.Config{}, err
	}
	if len(localConfigurations) == 0 {
		return globalConfig, nil
	}
	// TODO merge configs
	return config.Config{}, nil
}
