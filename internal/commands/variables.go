package commands

import "github.com/urfave/cli/v2"

// Variables is the parent command for variable listing operations
var Variables = &cli.Command{
	Name:               "variables",
	Description:        "manage variables",
	Usage:              "manage variables",
	Subcommands:        []*cli.Command{ListVariables},
	CustomHelpTemplate: CommandHelpTemplate,
}
