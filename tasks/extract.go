package tasks

import (
	"github.com/patrickhuber/wrangle/collections"
)

// ExtractTask represents an extract task
type ExtractTask struct {
	Details ExtractTaskDetails `yaml:"extract" mapstructure:"extract"`
}

// ExtractTaskDetails represent extract parameters
type ExtractTaskDetails struct {
	Archive     string `yaml:"archive"`
	Destination string `yaml:"destination"`
}

func (t *ExtractTask) Type() string {
	return "extract"
}

func (t *ExtractTask) Params() collections.ReadOnlyDictionary {
	dictionary := collections.NewDictionary()
	dictionary.Set("archive", t.Details.Archive)
	dictionary.Set("destination", t.Details.Destination)
	return dictionary
}

// NewExtractTask returns a new instance of a extract task
func NewExtractTask(archive string, destination string) Task {
	return &ExtractTask{
		Details: ExtractTaskDetails{
			Archive:     archive,
			Destination: destination,
		},
	}
}
