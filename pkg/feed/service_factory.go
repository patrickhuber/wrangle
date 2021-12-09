package feed

import "github.com/patrickhuber/wrangle/pkg/config"

type ServiceFactory interface {
	Create(f *config.Feed) (Service, error)
}

type serviceFactory struct {
	providers []Provider
}

func NewServiceFactory(providers ...Provider) ServiceFactory {
	return &serviceFactory{
		providers: providers,
	}
}

func (factory *serviceFactory) Create(f *config.Feed) (Service, error) {
	for _, p := range factory.providers {
		if p.Type() == f.Name {
			return p.Create(f)
		}
	}
	return nil, nil
}
