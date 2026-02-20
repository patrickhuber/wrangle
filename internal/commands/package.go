package commands

import "github.com/urfave/cli/v2"

// Package is the parent command for single package operations
var Package = &cli.Command{
	Name:               "package",
	Description:        "manage a package",
	Usage:              "manage a package",
	Subcommands:        []*cli.Command{Install},
	CustomHelpTemplate: CommandHelpTemplate,
}
