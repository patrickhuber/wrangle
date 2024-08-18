package commands

import (
	"fmt"

	"github.com/patrickhuber/go-di"
	"github.com/patrickhuber/go-xplat/console"
	"github.com/patrickhuber/wrangle/internal/app"
	"github.com/urfave/cli/v2"
)

// get subcommand
var Version = &cli.Command{
	Name:               "version",
	Action:             VersionAction,
	CustomHelpTemplate: CommandHelpTemplate,
}

func VersionAction(ctx *cli.Context) error {
	resolver, err := app.GetResolver(ctx)
	if err != nil {
		return err
	}
	c, err := di.Resolve[console.Console](resolver)
	if err != nil {
		return err
	}
	fmt.Fprintln(c.Out(), ctx.App.Version)
	return nil
}
