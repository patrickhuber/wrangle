package commands

import (
	"fmt"

	"github.com/patrickhuber/cli-mgr/config"
	"github.com/patrickhuber/cli-mgr/ui"
	"github.com/spf13/afero"
)

// Environments is an interface for a ProcessesCommand
type Environments interface {
	ExecuteCommand(configuration *config.Config) error
}

type environments struct {
	fileSystem afero.Fs
	console    ui.Console
}

// NewEnvironments returns a new process command for the given filesystem and console
func NewEnvironments(fileSystem afero.Fs, console ui.Console) Environments {
	return &environments{
		fileSystem: fileSystem,
		console:    console,
	}
}

func (cmd *environments) ExecuteCommand(configuration *config.Config) error {
	for _, environment := range configuration.Environments {
		fmt.Fprintln(cmd.console.Out(), environment.Name)
	}
	return nil
}
