package commands

import (
	"fmt"
	"text/tabwriter"

	"github.com/spf13/afero"

	"github.com/patrickhuber/wrangle/config"
	"github.com/patrickhuber/wrangle/ui"
)

type packagesCommand struct {
	console     ui.Console
	fileSystem  afero.Fs
	packagePath string
}

// PackagesCommand lists all packages in the configuration
type PackagesCommand interface {
	Execute(configuration *config.Config) error
}

// NewPackages returns a new packages command object
func NewPackages(fileSystem afero.Fs, console ui.Console, packagePath string) PackagesCommand {
	return &packagesCommand{
		console:     console,
		packagePath: packagePath}
}

func (cmd *packagesCommand) Execute(configuration *config.Config) error {
	w := tabwriter.NewWriter(cmd.console.Out(), 0, 0, 1, ' ', 0)
	fmt.Fprintln(w, "name\tversion")
	fmt.Fprintln(w, "----\t-------")
	packageFolders, err := afero.ReadDir(cmd.fileSystem, cmd.packagePath)
	if err != nil {
		return err
	}
	for _, packageFolder := range packageFolders {
		if !packageFolder.IsDir() {
			continue
		}
		packageVersions, err := afero.ReadDir(cmd.fileSystem, packageFolder.Name())

		if err != nil {
			return err
		}
		for _, packageVersion := range packageVersions {
			fmt.Fprintf(w, "%s\t%s", packageFolder.Name(), packageVersion.Name())
			fmt.Fprintln(w)
		}

	}
	return w.Flush()
}
