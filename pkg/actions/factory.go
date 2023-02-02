package actions

import "fmt"

type factory struct {
	providers map[string]Provider
}

// Factory defines a task factory for creating tasks
type Factory interface {
	Create(name string) (Provider, error)
}

func NewFactory(providers ...Provider) Factory {
	providerMap := map[string]Provider{}
	for _, p := range providers {
		providerMap[p.Type()] = p
	}
	return &factory{
		providers: providerMap,
	}
}

func (f *factory) Create(name string) (Provider, error) {
	provider, ok := f.providers[name]
	if !ok {
		return nil, fmt.Errorf("provider %s not found", name)
	}
	return provider, nil
}
