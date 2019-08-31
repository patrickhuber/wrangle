package tasks

import (
	"errors"
	"fmt"
)

// Task defines a unit of operation for a pacakge
type Task interface {
	Type() string
	Params() map[string]interface{}
}

type task struct {
	taskType string
	params   map[string]interface{}
	outputs map[string]interface{}
}

func (t *task) Type() string {
	return t.taskType
}

func(t *task) Params() map[string]interface{}{
	return t.params
}

// NewTask creates a new task from the given parameters
func NewTask(taskType string, params map[string]interface{}) Task {
	dictionary := make(map[string]interface{})
	for k, v := range params {
		dictionary[k] = v
	}
	return &task{
		taskType: taskType,
		params:   dictionary,
		outputs: make(map[string]interface{}),
	}
}

func (t *task) ValidateAndCastAsString(task, field string) (string, error){
	interfaceValue, ok := t.Params()[field]
	if !ok {
		message := fmt.Sprintf("%s parameter is required for %s task", field, task)
		return "", errors.New(message)
	}
	stringValue, ok := interfaceValue.(string)
	if!ok{
		message := fmt.Sprintf("%s task %s parameter is expected to be of type string", task, field)
		return "", errors.New(message)
	}
	return stringValue, nil
}