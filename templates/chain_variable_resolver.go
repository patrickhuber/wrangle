package templates

import "github.com/patrickhuber/cli-mgr/store"

type chainVariableResolver struct {
	delegate VariableResolver
	store    store.Store
}

// NewChainVariableResolver creates a chain resolver using the existing resolver as a delegate
func NewChainVariableResolver(store store.Store, delegate VariableResolver) VariableResolver {
	return &chainVariableResolver{store: store, delegate: delegate}
}

func (resolver *chainVariableResolver) Get(name string) (interface{}, error) {
	data, err := resolver.store.GetByName(name)
	if err != nil {
		return nil, err
	}
	template := NewTemplate(data.Value())
	document, err := template.Evaluate(resolver.delegate)
	if err != nil {
		return nil, err
	}
	return document, nil
}
