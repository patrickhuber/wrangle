package commands

import (
	"fmt"

	"github.com/patrickhuber/go-di"
	"github.com/patrickhuber/wrangle/internal/app"
	"github.com/patrickhuber/wrangle/internal/install"

	"github.com/urfave/cli/v2"
)

// install subcommand
var Install = &cli.Command{
	Name:               "install",
	Action:             InstallAction,
	CustomHelpTemplate: CommandHelpTemplate,
	Description:        "Installs the specified package",
	Usage:              "install the specified package",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "version",
			Aliases: []string{"v"},
			Usage:   "version of the package to install (default: latest)",
			Value:   "",
		},
		&cli.BoolFlag{
			Name:    "force",
			Aliases: []string{"f"},
			Usage:   "force reinstall if already installed",
			Value:   false,
		},
	},
}

type InstallCommand struct {
	Install install.Service `inject:""`
	Options InstallOptions  `options:""`
}

type InstallOptions struct {
	Package string `position:"0"`
	Version string `flag:"version"`
	Force   bool   `flag:"force"`
}

func (cmd *InstallCommand) Execute() error {

	request := &install.Request{
		Package: cmd.Options.Package,
		Version: cmd.Options.Version,
		Force:   cmd.Options.Force,
	}
	return cmd.Install.Execute(request)
}

func InstallAction(ctx *cli.Context) error {

	pkg := ctx.Args().First()
	if len(pkg) == 0 {
		return fmt.Errorf("package name is required")
	}

	resolver, err := app.GetResolver(ctx)
	if err != nil {
		return fmt.Errorf("invalid install configuration. %w", err)
	}

	cmd := &InstallCommand{
		Options: InstallOptions{
			Package: pkg,
			Version: ctx.String("version"),
			Force:   ctx.Bool("force"),
		},
	}

	err = di.Inject(resolver, cmd)
	if err != nil {
		return err
	}

	return cmd.Execute()
}
