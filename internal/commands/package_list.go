package commands

import (
	"fmt"

	"github.com/patrickhuber/go-cross/console"
	"github.com/patrickhuber/go-di"
	"github.com/patrickhuber/wrangle/internal/app"
	"github.com/patrickhuber/wrangle/internal/feed"
	"github.com/patrickhuber/wrangle/internal/structio"
	"github.com/urfave/cli/v2"
)

var PackageList = &cli.Command{
	Name: "list",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name: "output",
		},
	},
	Action:      PackageListAction,
	Description: "list available packages",
	Usage:       "list available packages",
}

type PackageListCommand struct {
	Service feed.ListPackages   `inject:""`
	Console console.Console     `inject:""`
	Options *PackageListOptions `options:""`
}

type PackageListOptions struct {
	Output string `flag:"output"`
}

func PackageListAction(ctx *cli.Context) error {
	resolver, err := app.GetResolver(ctx)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	cmd := &PackageListCommand{
		Options: &PackageListOptions{
			Output: ctx.String("output"),
		},
	}
	err = di.Inject(resolver, cmd)
	if err != nil {
		return err
	}
	return cmd.Execute()
}

func (cmd *PackageListCommand) Execute() error {
	request := &feed.ListPackagesRequest{}
	response, err := cmd.Service.Execute(request)
	if err != nil {
		return err
	}
	w := cmd.Console.Out()
	writer := structio.NewWriter(w, cmd.Options.Output)
	return writer.Write(response.Items)
}
