package commands

import (
	"github.com/urfave/cli/v2"
)

type ListFeedsCommand struct {
	Options *ListFeedsOptions
}

type ListFeedsOptions struct {
}

func ListFeeds(cli *cli.Context) error {
	return ListFeedsInternal(&ListFeedsCommand{})
}

func ListFeedsInternal(cmd *ListFeedsCommand) error {
	return nil
}
