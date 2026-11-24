package commands

import (
	"fmt"

	"github.com/patrickhuber/go-di"
	"github.com/patrickhuber/wrangle/internal/app"
	"github.com/patrickhuber/wrangle/internal/config/edit"
	"github.com/patrickhuber/wrangle/internal/install"

	"github.com/urfave/cli/v2"
)

// install subcommand
var PackageInstall = &cli.Command{
	Name:        "install",
	Action:      PackageInstallAction,
	Description: "Installs the specified package",
	Usage:       "install the specified package",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "version",
			Aliases: []string{"v"},
			Usage:   "version of the package to install (default: latest)",
			Value:   "",
		},
		&cli.StringFlag{
			Name:  "scope",
			Usage: "also add to configuration: nearest, local, user, or system",
		},
		&cli.BoolFlag{
			Name:    "force",
			Aliases: []string{"f"},
			Usage:   "force reinstall if already installed",
			Value:   false,
		},
	},
}

type PackageInstallCommand struct {
	Install install.Service       `inject:""`
	Edit    edit.Service          `inject:""`
	Options PackageInstallOptions `options:""`
}

type PackageInstallOptions struct {
	Package string `position:"0"`
	Version string `flag:"version"`
	Force   bool   `flag:"force"`
	Scope   string `flag:"scope"`
}

func (cmd *PackageInstallCommand) Execute() error {
	if cmd.Options.Scope != "" {
		if err := edit.ValidateScope(cmd.Options.Scope); err != nil {
			return err
		}
	}
	request := &install.Request{
		Package: cmd.Options.Package,
		Version: cmd.Options.Version,
		Force:   cmd.Options.Force,
	}
	if err := cmd.Install.Execute(request); err != nil {
		return err
	}
	if cmd.Options.Scope != "" {
		return cmd.Edit.AddPackage(&edit.PackageRequest{
			Package: cmd.Options.Package,
			Version: cmd.Options.Version,
			Scope:   cmd.Options.Scope,
		})
	}
	return nil
}

func PackageInstallAction(ctx *cli.Context) error {

	pkg := ctx.Args().First()
	if len(pkg) == 0 {
		return fmt.Errorf("package name is required")
	}

	resolver, err := app.GetResolver(ctx)
	if err != nil {
		return fmt.Errorf("invalid install configuration. %w", err)
	}

	cmd := &PackageInstallCommand{
		Options: PackageInstallOptions{
			Package: pkg,
			Version: ctx.String("version"),
			Force:   ctx.Bool("force"),
			Scope:   ctx.String("scope"),
		},
	}

	err = di.Inject(resolver, cmd)
	if err != nil {
		return err
	}

	return cmd.Execute()
}
