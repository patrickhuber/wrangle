package store

import "github.com/patrickhuber/wrangle/templates"

type storeVariableResolver struct {
	store Store
}

// NewStoreVariableResolver returns a new variable resolver that resolves variables from a store
func NewStoreVariableResolver(store Store) templates.VariableResolver {
	return &storeVariableResolver{store: store}
}

func (resolver *storeVariableResolver) Get(name string) (interface{}, error) {
	data, err := resolver.store.Get(name)
	if err != nil {
		return nil, err
	}
	return data.Value(), nil
}

func (resolver *storeVariableResolver) Lookup(name string) (interface{}, bool, error){
	data, ok, err := resolver.store.Lookup(name)
	if err != nil {
		return nil, false, err
	}
	if !ok{
		return nil, false, nil
	}
	return data.Value(), true, nil
}