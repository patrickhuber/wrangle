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
	Hidden:             true,
}

type InstallCommand struct {
	Install install.Service `inject:""`
	Options InstallOptions  `options:""`
}

type InstallOptions struct {
	Package string `position:"0"`
}

func (cmd *InstallCommand) Execute() error {

	request := &install.Request{
		Package: cmd.Options.Package,
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
		},
	}

	err = di.Inject(resolver, cmd)
	if err != nil {
		return err
	}

	return cmd.Execute()
}
