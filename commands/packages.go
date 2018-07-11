package commands

import (
	"fmt"

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
	for _, p := range configuration.Packages {
		fmt.Fprintf(cmd.console.Out(), "%s - %s", p.Name, p.Version)
		fmt.Fprintln(cmd.console.Out())
	}
	return nil
}
