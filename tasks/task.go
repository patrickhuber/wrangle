package tasks

import "github.com/patrickhuber/wrangle/collections"

// Task defines a unit of operation for a pacakge
type Task interface {
	Name() string
	Type() string
	Params() collections.ReadOnlyDictionary
}

type task struct {
	name     string
	taskType string
	params   collections.Dictionary
}

func (t *task) Name() string {
	return t.name
}

func (t *task) Type() string {
	return t.taskType
}

func (t *task) Params() collections.ReadOnlyDictionary {
	return t.params
}

// NewTask creates a new task from the given parameters
func NewTask(name string, taskType string, params map[string]string) Task {
	dictionary := collections.NewDictionary()
	for k, v := range params {
		dictionary.Set(k, v)
	}
	return &task{
		name:     name,
		taskType: taskType,
		params:   dictionary,
	}
}
