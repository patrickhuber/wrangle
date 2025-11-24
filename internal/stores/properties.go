package stores

import "fmt"

func GetRequiredProperty[T any](properties map[string]any, key string) (T, error) {
	value, ok := properties[key]
	if !ok {
		var zero T
		return zero, fmt.Errorf("missing required property '%s'", key)
	}
	typedValue, ok := value.(T)
	if !ok {
		var zero T
		return zero, fmt.Errorf("property '%s' is not of expected type", key)
	}
	return typedValue, nil
}

func GetOptionalProperty[T any](properties map[string]any, key string) (T, bool, error) {
	value, ok := properties[key]
	if !ok {
		var zero T
		return zero, false, nil
	}
	typedValue, ok := value.(T)
	if !ok {
		var zero T
		return zero, false, fmt.Errorf("property '%s' is not of expected type", key)
	}
	return typedValue, true, nil
}
