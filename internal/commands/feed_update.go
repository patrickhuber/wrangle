package commands

import (
	"fmt"
	"os"

	"github.com/patrickhuber/go-di"
	"github.com/patrickhuber/wrangle/internal/app"
	"github.com/patrickhuber/wrangle/internal/feed"
	"github.com/urfave/cli/v2"
)

var FeedUpdate = &cli.Command{
	Name:        "update",
	Action:      FeedUpdateAction,
	Description: "update feed packages in a local packages directory",
	Usage:       "update feed packages in a local packages directory",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:        "directory",
			Aliases:     []string{"d"},
			Usage:       "packages repository directory (defaults to the current directory)",
			DefaultText: "current directory",
		},
	},
}

type FeedUpdateCommand struct {
	UpdatePackages feed.UpdatePackages `inject:""`
}

func (cmd *FeedUpdateCommand) Execute(directory string) error {
	_, err := cmd.UpdatePackages.Execute(&feed.UpdatePackagesRequest{
		Directory: directory,
	})
	return err
}

func FeedUpdateAction(ctx *cli.Context) error {
	directory := ctx.String("directory")
	if directory == "" {
		var err error
		directory, err = os.Getwd()
		if err != nil {
			return fmt.Errorf("get current directory: %w", err)
		}
	}

	resolver, err := app.GetResolver(ctx)
	if err != nil {
		return fmt.Errorf("resolve feed update dependencies: %w", err)
	}
	cmd := &FeedUpdateCommand{}
	err = di.Inject(resolver, cmd)
	if err != nil {
		return err
	}
	return cmd.Execute(directory)
}
