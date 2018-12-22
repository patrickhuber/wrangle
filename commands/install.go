package commands

import (
	"github.com/patrickhuber/wrangle/global"
	"github.com/patrickhuber/wrangle/services"
	"github.com/urfave/cli"
)

// CreateInstallCommand creates the install command
func CreateInstallCommand(
	installService services.InstallService) *cli.Command {
	return &cli.Command{
		Name:      "install",
		Aliases:   []string{"i"},
		Usage:     "installs the package with the given `NAME` for the current platform",
		ArgsUsage: "<package name>",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "path, p",
				Usage:  "the package install path",
				EnvVar: global.PackagePathKey,
			},
			cli.StringFlag{
				Name:   "bin, b",
				Usage:  "the packages bin directory",
				EnvVar: global.BinPathKey,
			},
			cli.StringFlag{
				Name:   "root, r",
				Usage:  "the wrangle root directory",
				EnvVar: global.RootPathKey,
			},
			cli.StringFlag{
				Name:  "version, v",
				Usage: "the package version",
			},
		},
		Action: func(context *cli.Context) error {
			packagesRoot := context.String("path")
			root := context.String("root")
			bin := context.String("bin")
			packageName := context.Args().First()
			packageVersion := context.String("version")

			return installService.Install(root, bin, packagesRoot, packageName, packageVersion)
		},
	}
}
