package commands

import (
	"github.com/urfave/cli/v2"
)

var ListFeeds = &cli.Command{
	Name:        "feeds",
	Action:      ListFeedsAction,
	Description: "list available feeds",
	Usage:       "list available feeds",
}

type ListFeedsCommand struct {
	Options *ListFeedsOptions
}

type ListFeedsOptions struct {
}

func ListFeedsAction(cli *cli.Context) error {
	cmd := &ListFeedsCommand{}
	return (cmd).Execute()
}

func (cmd *ListFeedsCommand) Execute() error {
	return nil
}
