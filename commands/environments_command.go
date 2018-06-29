package commands

import (
	"fmt"

	"github.com/patrickhuber/cli-mgr/config"
	"github.com/patrickhuber/cli-mgr/ui"
	"github.com/spf13/afero"
)

// EnvironmentsCommand is an interface for a ProcessesCommand
type EnvironmentsCommand interface {
	ExecuteCommand(configuration *config.Config) error
}

type environmentsCommand struct {
	fileSystem afero.Fs
	console    ui.Console
}

// NewEnvironmentsCommand returns a new process command for the given filesystem and console
func NewEnvironmentsCommand(fileSystem afero.Fs, console ui.Console) EnvironmentsCommand {
	return &environmentsCommand{
		fileSystem: fileSystem,
		console:    console,
	}
}

func (cmd *environmentsCommand) ExecuteCommand(configuration *config.Config) error {
	for _, environment := range configuration.Environments {
		fmt.Fprintln(cmd.console.Out(), environment.Name)
	}
	return nil
}
