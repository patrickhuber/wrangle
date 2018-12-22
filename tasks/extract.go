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
	Archive string `yaml:"archive"`
}

func (t *ExtractTask) Type() string {
	return "extract"
}

func (t *ExtractTask) Params() collections.ReadOnlyDictionary {
	dictionary := collections.NewDictionary()
	dictionary.Set("archive", t.Details.Archive)
	return dictionary
}

// NewExtractTask returns a new instance of a extract task
func NewExtractTask(archive string) Task {
	return &ExtractTask{
		Details: ExtractTaskDetails{
			Archive: archive,
		},
	}
}
