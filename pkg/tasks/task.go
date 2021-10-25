package tasks

import "fmt"

// Task defines a unit of work for a package
type Task struct {
	Type       string
	Parameters map[string]interface{}
	Outputs    map[string]interface{}
}

// GetStringParameter casts the given parameter as a string and returns the value
func (t *Task) GetStringParameter(name string) (string, error) {
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
