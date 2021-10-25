package main

import (
	"log"
	"os"
	"strings"

	"github.com/patrickhuber/wrangle/internal/commands"
	"github.com/patrickhuber/wrangle/pkg/config"
	"github.com/patrickhuber/wrangle/pkg/console"
	"github.com/patrickhuber/wrangle/pkg/crosspath"
	"github.com/patrickhuber/wrangle/pkg/env"
	"github.com/patrickhuber/wrangle/pkg/feed"
	"github.com/patrickhuber/wrangle/pkg/filesystem"
	"github.com/patrickhuber/wrangle/pkg/ilog"
	"github.com/patrickhuber/wrangle/pkg/models"
	"github.com/patrickhuber/wrangle/pkg/operatingsystem"
	"github.com/spf13/afero"
	"github.com/urfave/cli/v2"
)

func main() {

	logger := ilog.Default()

	o := operatingsystem.New()
	environment := env.New()
	fs := filesystem.FromAferoFS(afero.NewOsFs())

	globalConfigProvider := config.NewDefaultReader(o, environment)
	globalConfig, err := globalConfigProvider.Get()
	handle(err)

	feedService := feed.NewMemoryService()
	feedManager := feed.NewManager(globalConfig)

	appName := "wrangle"
	if strings.EqualFold(o.Platform(), operatingsystem.PlatformWindows) {
		appName = appName + ".exe"
	}

	console := console.NewOS()

	app := &cli.App{
		Metadata: map[string]interface{}{
			"feedService": feedService,
			"feedManager": feedManager,
			"logger":      logger,
			"fileSystem":  fs,
			"os":          o,
			"console":     console,
		},
		Name: appName,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "bin",
				Aliases: []string{"b"},
				Value:   globalConfig.BinPath,
				EnvVars: []string{config.EnvironmentVariableBinPathKey},
			},
			&cli.StringFlag{
				Name:    "root",
				Aliases: []string{"r"},
				Value:   globalConfig.RootPath,
				EnvVars: []string{config.EnvironmentVariableRootPathKey},
			},
			&cli.StringFlag{
				Name:    "packages",
				Aliases: []string{"p"},
				Value:   globalConfig.PackagePath,
				EnvVars: []string{config.EnvironmentVariablePackagesPathKey},
			},
			&cli.StringFlag{
				Name:    "global",
				Aliases: []string{"g"},
				Value:   crosspath.Join(o.Home(), ".wrangle", "config.yml"),
			},
			&cli.GenericFlag{
				Name:        "output",
				Aliases:     []string{"o"},
				Value:       models.NewFormatEnum(),
				Required:    false,
				DefaultText: "table",
			},
		},
	}

	// list subcommand
	list := &cli.Command{
		Name: "list",
		Subcommands: []*cli.Command{
			{
				Name:   "packages",
				Action: commands.ListPackages,
			},
			{
				Name: "processes",
			},
			{
				Name: "stores",
			},
			{
				Name:   "feeds",
				Action: commands.NewListFeeds(feedManager).Execute,
			},
		},
	}

	// get subcommand
	get := &cli.Command{
		Name: "get",
	}

	// install subcommand
	install := &cli.Command{
		Name: "install",
	}

	// bootstrap subcommand
	bootstrap := &cli.Command{
		Name:   "bootstrap",
		Action: commands.Bootstrap,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "force",
				Aliases: []string{"f"},
				Value:   false,
			},
		},
	}

	// register
	app.Commands = []*cli.Command{
		bootstrap,
		list,
		get,
		install,
	}
	err = app.Run(os.Args)
	handle(err)
}

func handle(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
