package commands

import (
	"fmt"

	"github.com/patrickhuber/go-di"
	"github.com/patrickhuber/wrangle/internal/app"
	"github.com/patrickhuber/wrangle/pkg/console"
	"github.com/patrickhuber/wrangle/pkg/feed"
	"github.com/patrickhuber/wrangle/pkg/ilog"
	"github.com/patrickhuber/wrangle/pkg/packages"
	"github.com/patrickhuber/wrangle/pkg/structio"
	"github.com/urfave/cli/v2"
)

var ListPackages = &cli.Command{
	Name: "packages",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name: "output",
		},
	},
}

type ListPackagesCommand struct {
	FeedService feed.Service         `inject:""`
	Logger      ilog.Logger          `inject:""`
	Console     console.Console      `inject:""`
	Options     *ListPackagesOptions `options:""`
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
	di.Inject(resolver, listPackagesCommand)
	return listPackagesCommand.Execute()
}

func (cmd *ListPackagesCommand) Execute() error {
	request := &feed.ListRequest{}
	response, err := cmd.FeedService.List(request)
	if err != nil {
		return err
	}
	w := cmd.Console.Out()
	writer := structio.NewWriter(w, cmd.Options.Output)

	packages := []*packages.Package{}
	for _, i := range response.Items {
		packages = append(packages, i.Package)
	}
	return writer.Write(packages)
}
