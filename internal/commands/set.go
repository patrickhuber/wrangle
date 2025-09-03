package commands

import "github.com/urfave/cli/v2"

// list subcommand
var Set = &cli.Command{
	Name: "set",
	Subcommands: []*cli.Command{
		SetSecret,
	},
	CustomHelpTemplate: CommandHelpTemplate,
}
