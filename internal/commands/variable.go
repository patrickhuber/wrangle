package commands

import "github.com/urfave/cli/v2"

// Variable is the parent command for variable operations
var Variable = &cli.Command{
	Name:               "variable",
	Description:        "manage variables",
	Usage:              "manage variables",
	Subcommands:        []*cli.Command{VariableList},
	CustomHelpTemplate: CommandHelpTemplate,
}
