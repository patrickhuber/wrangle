package tasks

import (
	"fmt"

	"github.com/patrickhuber/wrangle/ui"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	yaml "gopkg.in/yaml.v2"
)

const moveTaskType = "move"

type moveProvider struct {
	fileSystem afero.Fs
	console    ui.Console
}

// NewMoveProvider creates a new move provider for moving files
func NewMoveProvider(fileSystem afero.Fs, console ui.Console) Provider {
	return &moveProvider{
		fileSystem: fileSystem,
		console:    console,
	}
}

func (provider *moveProvider) TaskType() string {
	return moveTaskType
}

func (provider *moveProvider) Execute(t Task) error {
	source, ok := t.Params().Lookup("source")
	if !ok {
		return errors.New("source parameter is required for move task")
	}

	destination, ok := t.Params().Lookup("destination")
	if !ok {
		return errors.New("destination parameter is required for move task")
	}

	fmt.Fprintf(provider.console.Out(), "moving '%s' to '%s'", source, destination)
	fmt.Fprintln(provider.console.Out())

	return provider.fileSystem.Rename(source, destination)
}

func (provider *moveProvider) Unmarshal(data string) (Task, error) {
	var task = &MoveTask{}
	err := yaml.Unmarshal([]byte(data), task)
	if err != nil {
		return nil, err
	}
	return task, nil
}
