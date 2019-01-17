package commands

import (
	"github.com/patrickhuber/wrangle/feed"
	"github.com/patrickhuber/wrangle/global"
	"github.com/patrickhuber/wrangle/services"
	"github.com/urfave/cli"
)

// CreatePackagesCommand creates a packages command from the cli context
func CreatePackagesCommand(
	packageServiceFactory services.PackageServiceFactory,
	feedServiceFactory feed.FeedServiceFactory) *cli.Command {
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
			cli.StringFlag{
				Name:   "url, u",
				Usage:  "the feed url",
				EnvVar: global.PackageFeedURLKey,
			},
		},
		Action: func(context *cli.Context) error {
			packagesPath := context.String("path")
			feedURL := context.String("url")
			feedService := feedServiceFactory.Get(packagesPath, feedURL)
			packagesService := packageServiceFactory.Get(feedService)
			return packagesService.List()
		},
	}
}
