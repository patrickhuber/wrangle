package templates

import "github.com/patrickhuber/wrangle/collections"

type dictionaryResolver struct {
	dictionary collections.ReadOnlyDictionary
}

// NewDictionaryResolver reuturns a new dictionary resolver
func NewDictionaryResolver(dictionary collections.ReadOnlyDictionary) VariableResolver {
	return &dictionaryResolver{
		dictionary: dictionary,
	}
}

func (resolver *dictionaryResolver) Get(name string) (interface{}, error) {
	return resolver.dictionary.Get(name)
}

func (resolver *dictionaryResolver) Lookup(name string) (interface{}, bool, error) {
	v, ok := resolver.dictionary.Lookup(name)
	return v, ok, nil
}
