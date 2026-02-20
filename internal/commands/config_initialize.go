package commands

import (
	"fmt"

	"github.com/patrickhuber/go-di"
	"github.com/patrickhuber/wrangle/internal/app"
	"github.com/patrickhuber/wrangle/internal/initialize"
	"github.com/urfave/cli/v2"
)

// initialize subcommand
var ConfigInitialize = &cli.Command{
	Name:    "initialize",
	Aliases: []string{"init"},
	Action:  ConfigInitializeAction,
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "force",
			Aliases: []string{"f"},
			Value:   false,
		},
	},
	Description:        "Initialize local configuration in the current directory",
	Usage:              "initialize local configuration in the current directory",
	CustomHelpTemplate: CommandHelpTemplate,
}

type ConfigInitializeCommand struct {
	Initialize initialize.Service `inject:""`
	Options    ConfigInitializeOptions
}

type ConfigInitializeOptions struct {
	Force bool `flag:"force"`
}

func ConfigInitializeAction(ctx *cli.Context) error {
	resolver, err := app.GetResolver(ctx)
	if err != nil {
		return fmt.Errorf("invalid initialize command configuration. %w", err)
	}

	cmd := &ConfigInitializeCommand{
		Options: ConfigInitializeOptions{
			Force: ctx.Bool("force"),
		},
	}
	err = di.Inject(resolver, cmd)
	if err != nil {
		return err
	}
	return cmd.Execute()
}

func (cmd *ConfigInitializeCommand) Execute() error {
	req := &initialize.Request{
		Force: cmd.Options.Force,
	}
	return cmd.Initialize.Execute(req)
}
