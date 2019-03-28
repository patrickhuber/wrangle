package env

import (
	"fmt"

	"github.com/patrickhuber/wrangle/collections"

	"github.com/patrickhuber/wrangle/config"
	"github.com/patrickhuber/wrangle/store"
)

type envStoreProvider struct {
	variables collections.Dictionary
}

func NewEnvStoreProvider(variables collections.Dictionary) store.Provider {
	return &envStoreProvider{
		variables: variables,
	}
}

func (p *envStoreProvider) Name() string {
	return "env"
}

func (p *envStoreProvider) Create(source *config.Store) (store.Store, error) {
	name := source.Name
	lookup := source.Params
	envStore := NewEnvStore(name, lookup, p.variables)
	if envStore.Type() != source.StoreType {
		return nil, fmt.Errorf(
			"provider '%s' can not create stores of type '%s'",
			p.Name(),
			source.StoreType)
	}
	return envStore, nil
}
