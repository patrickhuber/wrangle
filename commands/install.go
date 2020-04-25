package commands

import (
	"github.com/patrickhuber/wrangle/global"
	"github.com/patrickhuber/wrangle/packages"
	"github.com/urfave/cli"
)

// CreateInstallCommand creates the install command
func CreateInstallCommand(
	service packages.Service,
	platform string) *cli.Command {
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
			installRequest := &packages.InstallRequest{
				Directories: &packages.InstallRequestDirectories{
					Root:     context.String("root"),
					Bin:      context.String("bin"),
					Packages: context.String("path"),
				},
				Package: &packages.InstallRequestPackage{
					Name:    context.Args().First(),
					Version: context.String("version"),
				},
				Feed: &packages.InstallRequestFeed{
					URL: context.String("url"),
				},
				Platform: platform,
			}

			return service.Install(installRequest)
		},
	}
}
