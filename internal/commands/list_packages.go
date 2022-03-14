package commands

import (
	"github.com/patrickhuber/go-di"
	"github.com/patrickhuber/wrangle/internal/app"
	"github.com/patrickhuber/wrangle/pkg/console"
	"github.com/patrickhuber/wrangle/pkg/feed"
	"github.com/patrickhuber/wrangle/pkg/ilog"
	"github.com/patrickhuber/wrangle/pkg/packages"
	"github.com/patrickhuber/wrangle/pkg/structio"
	"github.com/urfave/cli/v2"
)

type ListPackagesCommand struct {
	FeedService feed.Service    `inject:""`
	Logger      ilog.Logger     `inject:""`
	Console     console.Console `inject:""`
	Options     *ListPackagesOptions
}

type ListPackagesOptions struct {
	Output string
}

func ListPackages(ctx *cli.Context) error {
	resolver := app.GetResolver(ctx)
	listPackagesCommand := &ListPackagesCommand{
		Options: &ListPackagesOptions{
			Output: ctx.String("output"),
		},
	}
	di.Inject(resolver, listPackagesCommand)
	return ListPackagesInternal(listPackagesCommand)
}

func ListPackagesInternal(cmd *ListPackagesCommand) error {
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
