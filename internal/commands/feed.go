package commands

import "github.com/urfave/cli/v2"

// Feed is the parent command for feed listing operations
var Feed = &cli.Command{
	Name:               "feed",
	Description:        "manage feeds",
	Usage:              "manage feeds",
	Subcommands:        []*cli.Command{FeedList},
	CustomHelpTemplate: CommandHelpTemplate,
}
