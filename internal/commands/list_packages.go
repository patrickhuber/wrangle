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

var ListPackages = &cli.Command{
	Name: "packages",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name: "output",
		},
	},
	Action: ListPackagesAction,
}

type ListPackagesCommand struct {
	Service feed.ListPackages    `inject:""`
	Console console.Console      `inject:""`
	Options *ListPackagesOptions `options:""`
}

type ListPackagesOptions struct {
	Output string `flag:"output"`
}

func ListPackagesAction(ctx *cli.Context) error {
	resolver, err := app.GetResolver(ctx)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	listPackagesCommand := &ListPackagesCommand{
		Options: &ListPackagesOptions{
			Output: ctx.String("output"),
		},
	}
	err = di.Inject(resolver, listPackagesCommand)
	if err != nil {
		return err
	}
	return listPackagesCommand.Execute()
}

func (cmd *ListPackagesCommand) Execute() error {
	request := &feed.ListPackagesRequest{}
	response, err := cmd.Service.Execute(request)
	if err != nil {
		return err
	}
	w := cmd.Console.Out()
	writer := structio.NewWriter(w, cmd.Options.Output)
	return writer.Write(response.Items)
}
