package commands

import (
	"github.com/patrickhuber/wrangle/feed"
	"github.com/patrickhuber/wrangle/global"
	"github.com/patrickhuber/wrangle/ui"
	"github.com/urfave/cli"
)

// CreateListPackagesCommand creates a packages command from the cli context
func CreateListPackagesCommand(
	console ui.Console,
	feedServiceFactory feed.ServiceFactory) cli.Command {
	return cli.Command{
		Name:    "packages",
		Aliases: []string{"k"},
		Usage:   "prints the list of packages and versions in the feed directory. If the feed directory isn't specified, uses the feed URL instead.",
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
				Value:  global.PackageFeedURL,
			},
		},
		Action: func(context *cli.Context) error {
			packagesPath := context.String("path")
			feedURL := context.String("url")

			feedService, err := feedServiceFactory.Get(packagesPath, feedURL)
			if err != nil {
				return err
			}

			listResponse, err := feedService.List(&feed.ListRequest{})
			if err != nil {
				return err
			}

			writer := feed.NewTableWriter(console.Out())
			return writer.Write(listResponse.Packages)
		},
	}
}
