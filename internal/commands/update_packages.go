package commands

import (
	"fmt"

	"github.com/patrickhuber/go-di"
	"github.com/patrickhuber/wrangle/internal/app"
	"github.com/patrickhuber/wrangle/internal/feed"
	"github.com/urfave/cli/v2"
)

var UpdatePackages = &cli.Command{
	Name:        "update-packages",
	Action:      UpdatePackagesAction,
	Description: "updates feed packages from github releases and renders package manifests",
	Usage:       "update package manifests in a feed directory",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "feed-directory",
			Usage: "relative path to the feed directory from the current working directory",
			Value: "feed",
		},
	},
	CustomHelpTemplate: CommandHelpTemplate,
}

type UpdatePackagesCommand struct {
	Service feed.UpdatePackages `inject:""`
	Options UpdatePackagesOptions
}

type UpdatePackagesOptions struct {
	FeedDirectory string
}

func UpdatePackagesAction(ctx *cli.Context) error {
	resolver, err := app.GetResolver(ctx)
	if err != nil {
		return fmt.Errorf("invalid update-packages command configuration. %w", err)
	}

	cmd := &UpdatePackagesCommand{
		Options: UpdatePackagesOptions{
			FeedDirectory: ctx.String("feed-directory"),
		},
	}

	err = di.Inject(resolver, cmd)
	if err != nil {
		return err
	}

	return cmd.Execute()
}

func (c *UpdatePackagesCommand) Execute() error {
	_, err := c.Service.Execute(&feed.UpdatePackagesRequest{
		FeedDirectory: c.Options.FeedDirectory,
	})
	return err
}
