package commands

import "github.com/urfave/cli/v2"

// Shell is the parent command for shell operations
var Shell = &cli.Command{
	Name:               "shell",
	Description:        "manage shell integration",
	Usage:              "manage shell integration",
	Subcommands:        []*cli.Command{ShellHook},
	CustomHelpTemplate: CommandHelpTemplate,
}
