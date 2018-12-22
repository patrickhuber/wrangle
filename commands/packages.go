package commands

import (
	"github.com/patrickhuber/wrangle/global"
	"github.com/patrickhuber/wrangle/services"
	"github.com/urfave/cli"
)

// CreatePackagesCommand creates a packages command from the cli context
func CreatePackagesCommand(
	packagesService services.PackagesService) *cli.Command {
	return &cli.Command{
		Name:    "packages",
		Aliases: []string{"k"},
		Usage:   "prints the list of packages and versions in the feed directory",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "path, p",
				Usage:  "the package install path",
				EnvVar: global.PackagePathKey,
			},
		},
		Action: func(context *cli.Context) error {
			packagesPath := context.String("path")
			return packagesService.List(packagesPath)
		},
	}
}
