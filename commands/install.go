package commands

import (
	"errors"
	"strings"

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
				Name:  "version, v",
				Usage: "the package version",
			},
		},
		Action: func(context *cli.Context) error {
			pacakgeRoot := context.String("path")
			packageName := context.Args().First()
			if strings.TrimSpace(packageName) == "" {
				return errors.New("missing required argument package name")
			}
			packageVersion := context.String("version")
			return installService.Install(pacakgeRoot, packageName, packageVersion)
		},
	}
}
