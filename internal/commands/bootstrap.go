package commands

import (
	"fmt"

	"github.com/patrickhuber/go-di"
	"github.com/patrickhuber/wrangle/internal/app"
	"github.com/patrickhuber/wrangle/internal/services"
	"github.com/urfave/cli/v2"
)

// bootstrap subcommand
var Bootstrap = &cli.Command{
	Name:        "bootstrap",
	Action:      BootstrapAction,
	Description: "bootstrap creates the global configuration file and installs the base packages",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "force",
			Aliases: []string{"f"},
			Value:   false,
		},
	},
}

type BootstrapCommand struct {
	Bootstrap services.Bootstrap `inject:""`
	Options   BootstarapOptions  `options:""`
}

type BootstarapOptions struct {
	Force bool `flag:"force"`
}

func (cmd *BootstrapCommand) Execute() error {
	req := &services.BootstrapRequest{
		Force: cmd.Options.Force,
	}
	return cmd.Bootstrap.Execute(req)
}

func BootstrapAction(ctx *cli.Context) error {
	resolver, err := app.GetResolver(ctx)
	if err != nil {
		return fmt.Errorf("invalid bootstrap command configuration. %w", err)
	}
	cmd := &BootstrapCommand{
		Options: BootstarapOptions{
			Force: ctx.Bool("force"),
		},
	}
	err = di.Inject(resolver, cmd)
	if err != nil {
		return err
	}
	return cmd.Execute()
}
