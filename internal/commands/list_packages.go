package commands

import (
	"github.com/patrickhuber/wrangle/internal/types"
	"github.com/patrickhuber/wrangle/pkg/console"
	"github.com/patrickhuber/wrangle/pkg/di"
	"github.com/patrickhuber/wrangle/pkg/feed"
	"github.com/patrickhuber/wrangle/pkg/global"
	"github.com/patrickhuber/wrangle/pkg/ilog"
	"github.com/patrickhuber/wrangle/pkg/packages"
	"github.com/patrickhuber/wrangle/pkg/structio"
	"github.com/urfave/cli/v2"
)

type ListPackagesCommand struct {
	FeedService feed.Service
	Logger      ilog.Logger
	Console     console.Console
	Options     *ListPackagesOptions
}

type ListPackagesOptions struct {
	Output string
}

func ListPackages(ctx *cli.Context) error {
	resolver := ctx.App.Metadata[global.MetadataDependencyInjection].(di.Resolver)

	feedService, err := resolver.Resolve(types.FeedService)
	if err != nil {
		return err
	}

	logger, err := resolver.Resolve(types.Logger)
	if err != nil {
		return err
	}

	con, err := resolver.Resolve(types.Console)
	if err != nil {
		return err
	}

	return ListPackagesInternal(
		&ListPackagesCommand{
			FeedService: feedService.(feed.Service),
			Logger:      logger.(ilog.Logger),
			Console:     con.(console.Console),
			Options: &ListPackagesOptions{
				Output: ctx.String("output"),
			},
		})
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
