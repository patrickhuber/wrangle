package tasks

import (
	"github.com/patrickhuber/wrangle/filepath"
	"github.com/patrickhuber/wrangle/filesystem"
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/patrickhuber/wrangle/ui"
	"github.com/pkg/errors"	
)

const moveTaskType = "move"

type moveProvider struct {
	fileSystem filesystem.FileSystem
	console    ui.Console
}

// NewMoveProvider creates a new move provider for moving files
func NewMoveProvider(fileSystem filesystem.FileSystem, console ui.Console) Provider {
	return &moveProvider{
		fileSystem: fileSystem,
		console:    console,
	}
}

func (provider *moveProvider) TaskType() string {
	return moveTaskType
}

func (provider *moveProvider) Execute(t Task, context TaskContext) error {			
	source, ok := t.Params()["source"]
	if !ok {
		return errors.New("source parameter is required for move task")
	}
	sourceString, ok := source.(string)
	if!ok{
		return errors.New("move task source parameter is expected to be of type string")
	}
	sourceString = filepath.Join(context.PackageVersionPath(), sourceString)

	destination, ok := t.Params()["destination"]
	if !ok {
		return errors.New("destination parameter is required for move task")
	}
	destinationString, ok := destination.(string)
	if !ok{
		return errors.New("move task destination parameter is expected to be of type string")
	}
	destinationString = filepath.Join(context.PackageVersionPath(), destinationString)

	fmt.Fprintf(provider.console.Out(), "moving '%s' to '%s'", source, destination)
	fmt.Fprintln(provider.console.Out())
	
	return provider.fileSystem.Rename(sourceString, destinationString)
}

func (provider *moveProvider) Decode(task interface{}) (Task, error) {
	var tsk = &MoveTask{}
	err := mapstructure.Decode(task, tsk)
	if err != nil {
		return nil, err
	}
	return tsk, nil
}