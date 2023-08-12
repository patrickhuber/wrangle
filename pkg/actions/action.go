package actions

import "fmt"

// Action defines a unit of work for a package
type Action struct {
	Type       string
	Parameters map[string]any
	Outputs    map[string]any
}

// GetStringParameter casts the given parameter as a string and returns the value
func (t *Action) GetStringParameter(name string) (string, error) {
	inter, ok := t.Parameters[name]
	if !ok {
		return "", fmt.Errorf("parameter %s was not found in the task", name)
	}
	value, ok := inter.(string)
	if !ok {
		return "", fmt.Errorf("parameter %s is not a string", name)
	}
	return value, nil
}

func (t *Action) GetOptionalStringParameter(name string) (string, bool, error) {
	inter, ok := t.Parameters[name]
	if !ok {
		return "", false, nil
	}
	value, ok := inter.(string)
	if !ok {
		return "", false, fmt.Errorf("parameter %s is not a string", name)
	}
	return value, true, nil
}
