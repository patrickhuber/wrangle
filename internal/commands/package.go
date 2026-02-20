package commands

import "github.com/urfave/cli/v2"

// Package is the parent command for all package operations
var Package = &cli.Command{
	Name:        "package",
	Description: "manage packages",
	Usage:       "manage packages",
	Subcommands: []*cli.Command{PackageList, PackageInstall, PackageRestore, PackageAdd},
}
