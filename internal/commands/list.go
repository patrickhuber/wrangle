package commands

import "github.com/urfave/cli/v2"

// list subcommand
var List = &cli.Command{
	Name: "list",
	Subcommands: []*cli.Command{
		ListPackages,
		ListFeeds,
		ListVariables,
	},
}
