package tasks

import (
	"github.com/patrickhuber/wrangle/collections"
)

// LinkTask represents a symlink
type LinkTask struct {
	Details LinkTaskDetails `yaml:"link" mapstructure:"link"`
}

// LinkTaskDetails contain the details for the link task
type LinkTaskDetails struct {
	Source      string `yaml:"source"`
	Destination string `yaml:"destination"`
}

func (task *LinkTask) Type() string {
	return "link"
}

func (task *LinkTask) Params() collections.ReadOnlyDictionary {
	dictionary := collections.NewDictionary()
	dictionary.Set("source", task.Details.Source)
	dictionary.Set("destination", task.Details.Destination)
	return dictionary
}

func NewLinkTask(source string, destination string) Task {
	return &LinkTask{
		Details: LinkTaskDetails{
			Source:      source,
			Destination: destination,
		},
	}
}
