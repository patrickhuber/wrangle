package services

import "github.com/patrickhuber/wrangle/internal/config"

type Interpolate interface {
	Execute() (config.Config, error)
}

type interpolate struct {
	configuration Configuration
}

func NewInterpolate(configuration Configuration) Interpolate {
	return &interpolate{
		configuration: configuration,
	}
}

func (i *interpolate) Execute() (config.Config, error) {
	return i.configuration.Get()
}
