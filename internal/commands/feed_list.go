package commands

import (
	"github.com/urfave/cli/v2"
)

var FeedList = &cli.Command{
	Name:        "list",
	Action:      FeedListAction,
	Description: "list available feeds",
	Usage:       "list available feeds",
}

type FeedListCommand struct {
	Options *FeedListOptions
}

type FeedListOptions struct {
}

func FeedListAction(cli *cli.Context) error {
	cmd := &FeedListCommand{}
	return (cmd).Execute()
}

func (cmd *FeedListCommand) Execute() error {
	return nil
}
