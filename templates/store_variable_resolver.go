package templates

import "github.com/patrickhuber/cli-mgr/store"

type storeVariableResolver struct {
	store store.Store
}

// NewStoreVariableResolver returns new store variable resolver
func NewStoreVariableResolver(store store.Store) VariableResolver {
	return &storeVariableResolver{store: store}
}

func (resolver *storeVariableResolver) Get(name string) (interface{}, error) {
	data, err := resolver.store.GetByName(name)
	if err != nil {
		return nil, err
	}
	return data.Value(), nil
}
