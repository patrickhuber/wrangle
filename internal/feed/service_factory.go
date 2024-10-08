package feed

import (
	"fmt"

	"github.com/patrickhuber/wrangle/internal/config"
)

type ServiceFactory interface {
	Create(f config.Feed) (Service, error)
}

type serviceFactory struct {
	providers []Provider
}

func NewServiceFactory(providers ...Provider) ServiceFactory {
	return &serviceFactory{
		providers: providers,
	}
}

func (factory *serviceFactory) Create(f config.Feed) (Service, error) {
	for _, p := range factory.providers {
		if p.Type() == f.Type {
			return p.Create(f)
		}
	}
	return nil, fmt.Errorf("there is no provider registered to create feed.Service of type '%s'", f.Type)
}
