package tasks

const setTaskType = "set"

// SetTask represents a symlink
type SetTask struct {
	Details SetTaskDetails `yaml:"set" mapstructure:"set"`
}

// SetTaskDetails contain the details for the link task
type SetTaskDetails struct {
	Variables map[string]string `yaml:"variables"`
}

// NewSetTask creates a new set task
func NewSetTask(variables map[string]string) Task {
	return &SetTask{
		Details: SetTaskDetails{
			Variables: variables,
		},
	}
}

func (task *SetTask) Params() map[string]interface{} {
	var dictionary = make(map[string]interface{})
	dictionary["variables"] = task.Details.Variables
	return dictionary
}

func (task *SetTask) Type() string {
	return setTaskType
}
