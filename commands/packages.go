package commands

import (
	"fmt"
	"text/tabwriter"

	"github.com/patrickhuber/wrangle/config"
	"github.com/patrickhuber/wrangle/ui"
)

type packagesCommand struct {
	console ui.Console
}

// PackagesCommand lists all packages in the configuration
type PackagesCommand interface {
	Execute(configuration *config.Config) error
}

// NewPackages returns a new packages command object
func NewPackages(console ui.Console) PackagesCommand {
	return &packagesCommand{console: console}
}

func (cmd *packagesCommand) Execute(configuration *config.Config) error {
	w := tabwriter.NewWriter(cmd.console.Out(), 0, 0, 1, ' ', 0)
	fmt.Fprintln(w, "name\tversion")
	fmt.Fprintln(w, "----\t-------")
	for _, item := range configuration.Packages {
		fmt.Fprintf(w, "%s\t%s", item.Name, item.Version)
		fmt.Fprintln(w)
	}
	return w.Flush()
}
