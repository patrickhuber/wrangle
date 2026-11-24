package commands

import (
	"fmt"

	"github.com/patrickhuber/go-di"
	"github.com/patrickhuber/wrangle/internal/app"
	"github.com/patrickhuber/wrangle/internal/config/edit"
	"github.com/urfave/cli/v2"
)

var PackageAdd = &cli.Command{
	Name:        "add",
	Action:      PackageAddAction,
	Description: "Adds a package to configuration without installing it",
	Usage:       "add a package to configuration",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "version",
			Aliases: []string{"v"},
			Usage:   "version of the package to add (default: latest)",
			Value:   "",
		},
		&cli.StringFlag{
			Name:  "scope",
			Usage: "configuration scope: nearest, local, user, or system (default: nearest)",
		},
	},
}

type PackageAddCommand struct {
	Edit    edit.Service      `inject:""`
	Options PackageAddOptions `options:""`
}

type PackageAddOptions struct {
	Package string `position:"0"`
	Version string `flag:"version"`
	Scope   string `flag:"scope"`
}

func (cmd *PackageAddCommand) Execute() error {
	request := &edit.PackageRequest{
		Package: cmd.Options.Package,
		Version: cmd.Options.Version,
		Scope:   cmd.Options.Scope,
	}
	return cmd.Edit.AddPackage(request)
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
			Scope:   ctx.String("scope"),
		},
	}

	err = di.Inject(resolver, cmd)
	if err != nil {
		return err
	}

	return cmd.Execute()
}
