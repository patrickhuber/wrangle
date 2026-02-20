package commands

import "github.com/urfave/cli/v2"

// Config is the parent command for configuration operations
var Config = &cli.Command{
	Name:               "config",
	Description:        "manage configuration",
	Usage:              "manage configuration",
	Subcommands:        []*cli.Command{Initialize, Export, Interpolate},
	CustomHelpTemplate: CommandHelpTemplate,
}
