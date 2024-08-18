package commands

import "github.com/urfave/cli/v2"

// get subcommand
var Get = &cli.Command{
	Name:               "get",
	Hidden:             true,
	CustomHelpTemplate: CommandHelpTemplate,
}
