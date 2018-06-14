package templates

type chainVariableResolver struct {
	delegate VariableResolver
	main     VariableResolver
}

// NewChainVariableResolver creates a chain resolver using the existing resolver as a delegate
func NewChainVariableResolver(main VariableResolver, delegate VariableResolver) VariableResolver {
	return &chainVariableResolver{main: main, delegate: delegate}
}

func (resolver *chainVariableResolver) Get(name string) (interface{}, error) {
	data, err := resolver.main.Get(name)
	if err != nil {
		return nil, err
	}
	template := NewTemplate(data)
	document, err := template.Evaluate(resolver.delegate)
	if err != nil {
		return nil, err
	}
	return document, nil
}
