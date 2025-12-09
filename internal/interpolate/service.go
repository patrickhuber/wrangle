package interpolate

import "github.com/patrickhuber/wrangle/internal/config"

type Service interface {
	Execute() (config.Config, error)
}

type service struct {
	configuration config.Service
}

func NewService(configuration config.Service) Service {
	return &service{
		configuration: configuration,
	}
}

func (i *service) Execute() (config.Config, error) {
	return i.configuration.Get()
}
