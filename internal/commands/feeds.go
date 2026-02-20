package commands

import "github.com/urfave/cli/v2"

// Feeds is the parent command for feed listing operations
var Feeds = &cli.Command{
	Name:               "feeds",
	Description:        "manage feeds",
	Usage:              "manage feeds",
	Subcommands:        []*cli.Command{ListFeeds},
	CustomHelpTemplate: CommandHelpTemplate,
}
