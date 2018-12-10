package tasks

import (
	"github.com/patrickhuber/wrangle/collections"
)

// MoveTask represents a move task
type MoveTask struct {
	Details MoveTaskDetails `yaml:"move"`
}

// MoveTaskDetails represent a move task parameters
type MoveTaskDetails struct {
	Source      string `yaml:"source"`
	Destination string `yaml:"destination"`
}

func (t *MoveTask) Type() string {
	return "move"
}

func (t *MoveTask) Params() collections.ReadOnlyDictionary {
	dictionary := collections.NewDictionary()
	dictionary.Set("source", t.Details.Source)
	dictionary.Set("destination", t.Details.Destination)
	return dictionary
}

// NewMoveTask returns an instance of a move task
func NewMoveTask(name string, source string, destination string) Task {
	return &MoveTask{
		Details: MoveTaskDetails{
			Destination: destination,
			Source:      source,
		},
	}
}
