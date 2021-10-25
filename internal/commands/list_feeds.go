package commands

import (
	"github.com/patrickhuber/wrangle/pkg/feed"
	"github.com/urfave/cli/v2"
)

type listFeeds struct {
}

func NewListFeeds(feedManager feed.Manager) Command {
	return &listFeeds{}
}

func (c *listFeeds) Execute(ctx *cli.Context) error {
	return nil
}
