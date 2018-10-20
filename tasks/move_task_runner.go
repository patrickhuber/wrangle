package tasks

import (
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

type moveTaskRunner struct {
	fileSystem afero.Fs
}

// NewMoveTaskRunner creates a new task runner for moving files
func NewMoveTaskRunner(fileSystem afero.Fs) TaskRunner {
	return &moveTaskRunner{
		fileSystem: fileSystem,
	}
}

func (runner *moveTaskRunner) Execute(t Task) error {
	source, ok := t.Params().Lookup("source")
	if !ok {
		return errors.New("source parameter is required for move task")
	}

	destination, ok := t.Params().Lookup("destination")
	if !ok {
		return errors.New("destination parameter is required for move task")
	}

	return runner.fileSystem.Rename(source, destination)
}
