package tasks

import "github.com/mitchellh/mapstructure"

type setProvider struct {
}

// NewSetProvider creates a new set provider
func NewSetProvider() Provider {
	return &setProvider{}
}

func (provider setProvider) TaskType() string {
	return setTaskType
}

func (provider setProvider) Execute(task Task, context TaskContext) error {

	return nil
}

func (provider setProvider) Decode(task interface{}) (Task, error) {
	t := &SetTask{}
	err := mapstructure.Decode(task, t)
	if err != nil {
		return nil, err
	}
	return t, nil
}
