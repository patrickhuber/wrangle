package commands

import "github.com/urfave/cli/v2"

// list subcommand
var Set = &cli.Command{
	Name:        "set",
	Description: "set the specified resource",
	Usage:       "set the specified resource",
	Subcommands: []*cli.Command{
		SetSecret,
	},
	CustomHelpTemplate: CommandHelpTemplate,
}
