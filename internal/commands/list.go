package commands

import "github.com/urfave/cli/v2"

// list subcommand
var List = &cli.Command{
	Name:        "list",
	Description: "list available packages, feeds, and variables",
	Usage:       "list available packages, feeds, and variables",
	Subcommands: []*cli.Command{
		ListPackages,
		ListFeeds,
		ListVariables,
	},
	CustomHelpTemplate: CommandHelpTemplate,
	Hidden:             true,
}
