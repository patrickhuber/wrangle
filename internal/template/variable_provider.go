package template

// VariableProvider provides a variable value and lists available variables
type VariableProvider interface {
	// Get returns the value for the variable from this provider
	Get(key string) (any, bool, error)
	// List lists the variables contained in this provider
	List() ([]string, error)
}

type MapProvider map[string]any

func (mp MapProvider) Get(key string) (any, bool, error) {
	if mp == nil {
		return nil, false, nil
	}
	v, ok := mp[key]
	return v, ok, nil
}

func (mp MapProvider) List() ([]string, error) {
	var keys []string
	for k := range mp {
		keys = append(keys, k)
	}
	return keys, nil
}
