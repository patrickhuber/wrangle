package commands

import (
	"fmt"

	"github.com/patrickhuber/go-di"
	"github.com/patrickhuber/wrangle/internal/app"
	"github.com/patrickhuber/wrangle/internal/feed"
	"github.com/urfave/cli/v2"
)

var FeedUpdate = &cli.Command{
	Name:        "update",
	Action:      FeedUpdateAction,
	Description: "update feed packages to their latest versions based on resource configurations",
	Usage:       "update feed packages to their latest versions",
}

type FeedUpdateCommand struct {
	UpdatePackages feed.UpdatePackages `inject:""`
}

func (cmd *FeedUpdateCommand) Execute() error {
	_, err := cmd.UpdatePackages.Execute(&feed.UpdatePackagesRequest{})
	return err
}

func FeedUpdateAction(ctx *cli.Context) error {
	resolver, err := app.GetResolver(ctx)
	if err != nil {
		return fmt.Errorf("invalid feed update configuration. %w", err)
	}
	cmd := &FeedUpdateCommand{}
	err = di.Inject(resolver, cmd)
	if err != nil {
		return err
	}
	return cmd.Execute()
}
