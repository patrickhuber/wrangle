package tasks

import (
	"fmt"

	"github.com/patrickhuber/wrangle/filesystem"
	"github.com/patrickhuber/wrangle/ui"
)

const linkTaskType = "link"

type linkProvider struct {
	fileSystem filesystem.FsWrapper
	console    ui.Console
}

// NewLinkProvider creates a new provider for creating process linkes (symlinks)
func NewLinkProvider(fileSystem filesystem.FsWrapper, console ui.Console) Provider {
	return &linkProvider{
		fileSystem: fileSystem,
		console:    console,
	}
}

// NewLinkTask returns an instance of a link task
func (provider *linkProvider) TaskType() string {
	return linkTaskType
}

func (provider *linkProvider) Execute(task Task) error {

	source, ok := task.Params().Lookup("source")
	if !ok {
		return fmt.Errorf("source parameter is required for link tasks")
	}

	destination, ok := task.Params().Lookup("destination")
	if !ok {
		return fmt.Errorf("destination parameter is required for link tasks")
	}

	return provider.fileSystem.Symlink(source, destination)
}
