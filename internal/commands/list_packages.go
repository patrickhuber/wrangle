package commands

import (
	"github.com/patrickhuber/wrangle/pkg/console"
	"github.com/patrickhuber/wrangle/pkg/feed"
	"github.com/patrickhuber/wrangle/pkg/ilog"
	"github.com/patrickhuber/wrangle/pkg/packages"
	"github.com/patrickhuber/wrangle/pkg/structio"
	"github.com/urfave/cli/v2"
)

type ListPackagesOptions struct {
	FeedService feed.Service
	Logger      ilog.Logger
	Output      string
	Console     console.Console
}

func ListPackages(ctx *cli.Context) error {
	return ListPackagesInternal(&ListPackagesOptions{
		Output:      ctx.String("output"),
		FeedService: ctx.App.Metadata["feedService"].(feed.Service),
		Logger:      ctx.App.Metadata["logger"].(ilog.Logger),
		Console:     ctx.App.Metadata["console"].(console.Console),
	})
}

func ListPackagesInternal(ctx *ListPackagesOptions) error {
	request := &feed.ListRequest{}
	response, err := ctx.FeedService.List(request)
	if err != nil {
		return err
	}
	w := ctx.Console.Out()
	writer := structio.NewWriter(w, ctx.Output)

	packages := []*packages.Package{}
	for _, i := range response.Items {
		packages = append(packages, i.Package)
	}
	return writer.Write(packages)
}
