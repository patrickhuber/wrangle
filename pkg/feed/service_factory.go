package feed

import "github.com/patrickhuber/wrangle/pkg/config"

type ServiceFactory interface {
	Create(f *config.Feed) Service
}

type serviceFactory struct {
	services []Service
}

func NewServiceFactory(services []Service) ServiceFactory {
	return &serviceFactory{
		services: services,
	}
}

func (factory *serviceFactory) Create(f *config.Feed) Service {
	for _, s := range factory.services {
		if s.Name() == f.Name {
			return s
		}
	}
	return nil
}
