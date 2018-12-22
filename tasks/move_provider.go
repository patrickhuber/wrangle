package tasks

import (
	"github.com/patrickhuber/wrangle/filepath"
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/patrickhuber/wrangle/ui"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
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

func (provider *moveProvider) Execute(t Task, context TaskContext) error {
	source, ok := t.Params().Lookup("source")
	if !ok {
		return errors.New("source parameter is required for move task")
	}
	source = filepath.Join(context.PackageVersionPath(), source)

	destination, ok := t.Params().Lookup("destination")
	if !ok {
		return errors.New("destination parameter is required for move task")
	}
	destination = filepath.Join(context.PackageVersionPath(), destination)

	fmt.Fprintf(provider.console.Out(), "moving '%s' to '%s'", source, destination)
	fmt.Fprintln(provider.console.Out())

	return provider.fileSystem.Rename(source, destination)
}

func (provider *moveProvider) Decode(task interface{}) (Task, error) {
	var tsk = &MoveTask{}
	err := mapstructure.Decode(task, tsk)
	if err != nil {
		return nil, err
	}
	return tsk, nil
}
