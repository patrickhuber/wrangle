package commands

import (
	"fmt"

	"github.com/patrickhuber/go-di"
	"github.com/patrickhuber/wrangle/internal/app"
	"github.com/patrickhuber/wrangle/internal/restore"
	"github.com/urfave/cli/v2"
)

var PackageRestore = &cli.Command{
	Name:        "restore",
	Action:      PackageRestoreAction,
	Description: "Restores all packages defined in the configuration",
	Usage:       "restore packages from configuration",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "force",
			Aliases: []string{"f"},
			Usage:   "force reinstall if already installed",
			Value:   false,
		},
	},
}

type PackageRestoreCommand struct {
	Restore restore.Service       `inject:""`
	Options PackageRestoreOptions `options:""`
}

type PackageRestoreOptions struct {
	Force bool `flag:"force"`
}

func (cmd *PackageRestoreCommand) Execute() error {
	request := &restore.Request{
		Force: cmd.Options.Force,
	}
	return cmd.Restore.Execute(request)
}

func PackageRestoreAction(ctx *cli.Context) error {
	resolver, err := app.GetResolver(ctx)
	if err != nil {
		return fmt.Errorf("invalid restore configuration. %w", err)
	}

	cmd := &PackageRestoreCommand{
		Options: PackageRestoreOptions{
			Force: ctx.Bool("force"),
		},
	}

	err = di.Inject(resolver, cmd)
	if err != nil {
		return err
	}

	return cmd.Execute()
}
