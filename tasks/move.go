package tasks

// MoveTask represents a move task
type MoveTask struct {
	Details MoveTaskDetails `yaml:"move" mapstructure:"move"`
}

// MoveTaskDetails represent a move task parameters
type MoveTaskDetails struct {
	Source      string `yaml:"source"`
	Destination string `yaml:"destination"`
}

func (t *MoveTask) Type() string {
	return "move"
}

func (t *MoveTask) Params() map[string]interface{} {
	dictionary := make(map[string]interface{})
	dictionary["source"] = t.Details.Source
	dictionary["destination"] = t.Details.Destination
	return dictionary
}

// NewMoveTask returns an instance of a move task
func NewMoveTask(source string, destination string) Task {
	return &MoveTask{
		Details: MoveTaskDetails{
			Destination: destination,
			Source:      source,
		},
	}
}
