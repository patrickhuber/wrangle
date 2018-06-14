package store

import "github.com/patrickhuber/cli-mgr/templates"

type storeVariableResolver struct {
	store Store
}

// NewStoreVariableResolver returns new store variable resolver
func NewStoreVariableResolver(store Store) templates.VariableResolver {
	return &storeVariableResolver{store: store}
}

func (resolver *storeVariableResolver) Get(name string) (interface{}, error) {
	data, err := resolver.store.GetByName(name)
	if err != nil {
		return nil, err
	}
	return data.Value(), nil
}
