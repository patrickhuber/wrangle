package commands

import "github.com/urfave/cli/v2"

var Store = &cli.Command{
	Name:        "store",
	Description: "manage configuration stores",
	Usage:       "manage configuration stores",
	Subcommands: []*cli.Command{StoreAdd},
}
