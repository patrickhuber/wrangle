package tasks

import (
	"fmt"

	"github.com/mitchellh/mapstructure"

	"github.com/patrickhuber/wrangle/filepath"
	"github.com/patrickhuber/wrangle/filesystem"
	"github.com/patrickhuber/wrangle/ui"
)

const linkTaskType = "link"

type linkProvider struct {
	fileSystem filesystem.FileSystem
	console    ui.Console
}

// NewLinkProvider creates a new provider for creating process linkes (symlinks)
func NewLinkProvider(fileSystem filesystem.FileSystem, console ui.Console) Provider {
	return &linkProvider{
		fileSystem: fileSystem,
		console:    console,
	}
}

// NewLinkTask returns an instance of a link task
func (provider *linkProvider) TaskType() string {
	return linkTaskType
}

func (provider *linkProvider) Execute(task Task, context TaskContext) error {

	sourceInterface, ok := task.Params()["source"]
	if !ok {
		return fmt.Errorf("source parameter is required for link tasks")
	}
	source, ok := sourceInterface.(string)
	if !ok{
		return fmt.Errorf("source parameter is expected to be of type string")
	}
	source = filepath.Join(context.PackageVersionPath(), source)

	aliasInterface, ok := task.Params()["alias"]
	if !ok {
		return fmt.Errorf("alias parameter is required for link tasks")
	}	
	alias, ok := aliasInterface.(string)
	if !ok{
		return fmt.Errorf("alias parameter is expected to be of type string")
	}

	alias = filepath.Join(context.Bin(), alias)

	return provider.fileSystem.Symlink(source, alias)
}

func (provider *linkProvider) Decode(task interface{}) (Task, error) {
	var tsk = &LinkTask{}
	err := mapstructure.Decode(task, tsk)
	if err != nil {
		return nil, err
	}
	return tsk, nil
}
