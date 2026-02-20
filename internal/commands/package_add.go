package commands

import (
	"fmt"

	"github.com/patrickhuber/go-di"
	"github.com/patrickhuber/wrangle/internal/add"
	"github.com/patrickhuber/wrangle/internal/app"
	"github.com/urfave/cli/v2"
)

var PackageAdd = &cli.Command{
	Name:               "add",
	Action:             PackageAddAction,
	CustomHelpTemplate: CommandHelpTemplate,
	Description:        "Adds a package to the local configuration file without installing it",
	Usage:              "add a package to the local configuration",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "version",
			Aliases: []string{"v"},
			Usage:   "version of the package to add (default: latest)",
			Value:   "",
		},
	},
}

type PackageAddCommand struct {
	Add     add.Service       `inject:""`
	Options PackageAddOptions `options:""`
}

type PackageAddOptions struct {
	Package string `position:"0"`
	Version string `flag:"version"`
}

func (cmd *PackageAddCommand) Execute() error {
	request := &add.Request{
		Package: cmd.Options.Package,
		Version: cmd.Options.Version,
	}
	return cmd.Add.Execute(request)
}

func PackageAddAction(ctx *cli.Context) error {
	pkg := ctx.Args().First()
	if len(pkg) == 0 {
		return fmt.Errorf("package name is required")
	}

	resolver, err := app.GetResolver(ctx)
	if err != nil {
		return fmt.Errorf("invalid add configuration. %w", err)
	}

	cmd := &PackageAddCommand{
		Options: PackageAddOptions{
			Package: pkg,
			Version: ctx.String("version"),
		},
	}

	err = di.Inject(resolver, cmd)
	if err != nil {
		return err
	}

	return cmd.Execute()
}
