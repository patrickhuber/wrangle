package main

import (
	"log"
	"os"
	"strings"

	"github.com/patrickhuber/wrangle/internal/commands"
	"github.com/patrickhuber/wrangle/internal/services"
	"github.com/patrickhuber/wrangle/internal/types"
	"github.com/patrickhuber/wrangle/pkg/console"
	"github.com/patrickhuber/wrangle/pkg/crosspath"
	"github.com/patrickhuber/wrangle/pkg/di"
	"github.com/patrickhuber/wrangle/pkg/env"
	"github.com/patrickhuber/wrangle/pkg/feed"
	"github.com/patrickhuber/wrangle/pkg/feed/git"
	"github.com/patrickhuber/wrangle/pkg/filesystem"
	"github.com/patrickhuber/wrangle/pkg/global"
	"github.com/patrickhuber/wrangle/pkg/ilog"
	"github.com/patrickhuber/wrangle/pkg/models"
	"github.com/patrickhuber/wrangle/pkg/operatingsystem"
	"github.com/spf13/afero"
	"github.com/urfave/cli/v2"
)

// set with -ldflags
var version = ""

func main() {

	container := di.NewContainer()
	container.RegisterConstructor(ilog.Default)
	container.RegisterConstructor(operatingsystem.New)
	container.RegisterConstructor(env.New)
	container.RegisterConstructor(afero.NewOsFs)
	container.RegisterConstructor(filesystem.FromAferoFS)
	container.RegisterConstructor(console.NewOS)
	container.RegisterDynamic(types.FeedServiceFactory, func(r di.Resolver) interface{} {
		return feed.NewServiceFactory(git.NewProvider())
	})
	container.RegisterConstructor(services.NewInstall)

	o := container.Resolve(types.OS).(operatingsystem.OS)
	appName := "wrangle"
	if strings.EqualFold(o.Platform(), operatingsystem.PlatformWindows) {
		appName = appName + ".exe"
	}

	app := &cli.App{
		Metadata: map[string]interface{}{
			global.MetadataDependencyInjection: container,
		},
		Name:        appName,
		Version:     version,
		Description: "A DevOps Environment Management CLI",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    global.FlagBin,
				Aliases: []string{"b"},
				EnvVars: []string{global.EnvBin},
			},
			&cli.StringFlag{
				Name:    global.FlagRoot,
				Aliases: []string{"r"},
				EnvVars: []string{global.EnvRoot},
			},
			&cli.StringFlag{
				Name:    global.FlagPackages,
				Aliases: []string{"p"},
				EnvVars: []string{global.EnvPackages},
			},
			&cli.StringFlag{
				Name:    global.FlagConfig,
				Aliases: []string{"g"},
				Value:   crosspath.Join(o.Home(), ".wrangle", "config.yml"),
				EnvVars: []string{global.EnvConfig},
			},
			&cli.GenericFlag{
				Name:        global.FlagOutput,
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
				Action: commands.ListFeeds,
			},
		},
	}

	// get subcommand
	get := &cli.Command{
		Name: "get",
	}

	// install subcommand
	install := &cli.Command{
		Name:   "install",
		Action: commands.Install,
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
	err := app.Run(os.Args)
	handle(err)
}

func handle(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
