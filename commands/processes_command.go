package commands

import (
	"fmt"

	"github.com/patrickhuber/cli-mgr/config"
	"github.com/patrickhuber/cli-mgr/ui"
	"github.com/spf13/afero"
)

// ProcessesCommand is an interface for a ProcessesCommand
type ProcessesCommand interface {
	ExecuteCommand(configFile string) error
}

type processesCommand struct {
	fileSystem afero.Fs
	console    ui.Console
}

// NewProcessesCommand returns a new process command for the given filesystem and console
func NewProcessesCommand(fileSystem afero.Fs, console ui.Console) ProcessesCommand {
	return &processesCommand{
		fileSystem: fileSystem,
		console:    console,
	}
}

func (cmd *processesCommand) ExecuteCommand(configFile string) error {
	loader := config.NewLoader(cmd.fileSystem)
	cfg, err := loader.Load(configFile)
	if err != nil {
		return err
	}
	for _, process := range cfg.Processes {
		fmt.Fprintln(cmd.console.Out(), process.Name)
	}
	return nil
}
