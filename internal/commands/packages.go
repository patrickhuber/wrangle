package commands

import "github.com/urfave/cli/v2"

// Packages is the parent command for package listing operations
var Packages = &cli.Command{
	Name:               "packages",
	Description:        "manage packages",
	Usage:              "manage packages",
	Subcommands:        []*cli.Command{ListPackages},
	CustomHelpTemplate: CommandHelpTemplate,
}
