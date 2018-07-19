package commands

import (
	"fmt"
	"text/tabwriter"

	"github.com/patrickhuber/wrangle/config"
	"github.com/patrickhuber/wrangle/ui"
	"github.com/spf13/afero"
)

// Environments is an interface for a ProcessesCommand
type Environments interface {
	Execute(configuration *config.Config) error
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

func (cmd *environments) Execute(configuration *config.Config) error {
	w := tabwriter.NewWriter(cmd.console.Out(), 0, 0, 1, ' ', 0)
	fmt.Fprintln(w, "name")
	fmt.Fprintln(w, "----")
	for _, item := range configuration.Environments {
		fmt.Fprintf(w, "%s", item.Name)
		fmt.Fprintln(w)
	}
	return w.Flush()
}
