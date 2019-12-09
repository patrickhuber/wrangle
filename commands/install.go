package commands

import (
	"github.com/patrickhuber/wrangle/global"
	"github.com/patrickhuber/wrangle/packages"
	"github.com/urfave/cli"
)

// CreateInstallCommand creates the install command
func CreateInstallCommand(
	installService packages.InstallService) *cli.Command {
	return &cli.Command{
		Name:      "install",
		Aliases:   []string{"i"},
		Usage:     "installs the package with the given `NAME` for the current platform",
		ArgsUsage: "<package name> [options]",
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
				Name:   "url, u",
				Usage:  "the feed url",
				EnvVar: global.PackageFeedURLKey,
				Value:  global.PackageFeedURL,
			},
			cli.StringFlag{
				Name:  "version, v",
				Usage: "the package version",
			},
		},
		Action: func(context *cli.Context) error {
			installServiceRequest := &packages.InstallServiceRequest{
				Directories: &packages.InstallServiceRequestDirectories{
					Root:     context.String("root"),
					Bin:      context.String("bin"),
					Packages: context.String("path"),
				},
				Package: &packages.InstallServiceRequestPackage{
					Name:    context.Args().First(),
					Version: context.String("version"),
				},
				Feed: &packages.InstallServiceRequestFeed{
					URL: context.String("url"),
				},
			}

			return installService.Install(installServiceRequest)
		},
	}
}
