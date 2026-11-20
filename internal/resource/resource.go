package resource

import "fmt"

// Resource is a generic way to fetch versions from a given endpoint
type Resource interface {
	Type() string
	Name() string
	Versions(currentVersion string) ([]string, error)
}

type Config struct {
	Name   string
	Type   string
	Source map[string]string
}

type Provider interface {
	Type() string
	Create(config Config) (Resource, error)
}

type Factory interface {
	Create(config Config) (Resource, error)
}

type factory struct {
	providers map[string]Provider
}

func NewFactory(providers ...Provider) Factory {
	providerMap := make(map[string]Provider, len(providers))
	for _, provider := range providers {
		providerMap[provider.Type()] = provider
	}
	return &factory{providers: providerMap}
}

func (f *factory) Create(config Config) (Resource, error) {
	provider, ok := f.providers[config.Type]
	if !ok {
		return nil, fmt.Errorf("resource provider %q not found", config.Type)
	}
	return provider.Create(config)
}
