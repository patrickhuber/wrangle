package env

import (
	"fmt"

	"github.com/patrickhuber/wrangle/config"
	"github.com/patrickhuber/wrangle/store"
)

type envStoreProvider struct {
}

func NewEnvStoreProvider() store.Provider {
	return &envStoreProvider{}
}

func (p *envStoreProvider) GetName() string {
	return "env"
}

func (p *envStoreProvider) Create(source *config.Store) (store.Store, error) {
	name := source.Name
	lookup := source.Params
	envStore := NewEnvStore(name, lookup)
	if envStore.Type() != source.StoreType {
		return nil, fmt.Errorf(
			"provider '%s' can not create stores of type '%s'",
			p.GetName(),
			source.StoreType)
	}
	return envStore, nil
}
