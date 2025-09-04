package interpolate

import "github.com/patrickhuber/wrangle/internal/config"

type Service interface {
	Execute() (config.Config, error)
}

type service struct {
	configuration config.Configuration
}

func NewService(configuration config.Configuration) Service {
	return &service{
		configuration: configuration,
	}
}

func (i *service) Execute() (config.Config, error) {
	return i.configuration.Get()
}
