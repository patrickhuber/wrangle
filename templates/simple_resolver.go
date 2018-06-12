package templates

import "fmt"

type simpleResolver struct {
	innerMap map[string]interface{}
}

func newSimpleResolver(values ...interface{}) (VariableResolver, error) {
	innerMap := make(map[string]interface{})
	if len(values)%2 == 1 {
		return nil, fmt.Errorf("values must be a list of key value pairs. ex: key1, value1, key2, value2")
	}
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, fmt.Errorf("key '%v' is not a string", key)
		}
		value := values[i+1]
		innerMap[key] = value
	}
	return &simpleResolver{innerMap: innerMap}, nil
}

func (resolver *simpleResolver) Get(key string) (interface{}, error) {
	value, ok := resolver.innerMap[key]
	if !ok {
		return nil, fmt.Errorf("unable to find key '%s'", key)
	}
	return value, nil
}
