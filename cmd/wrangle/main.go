package main

import (
	"log"

	"github.com/patrickhuber/go-di"
	"github.com/patrickhuber/go-xplat/console"
	"github.com/patrickhuber/go-xplat/filepath"
	"github.com/patrickhuber/go-xplat/os"
	"github.com/patrickhuber/go-xplat/platform"
	"github.com/patrickhuber/wrangle/internal/app"
	"github.com/patrickhuber/wrangle/internal/commands"
	setup "github.com/patrickhuber/wrangle/internal/host"
	"github.com/patrickhuber/wrangle/pkg/config"
	"github.com/patrickhuber/wrangle/pkg/enums"
	"github.com/patrickhuber/wrangle/pkg/global"
	"github.com/urfave/cli/v2"
)

// set with -ldflags
var version = ""

func main() {
	s := setup.New()
	container := s.Container()

	o, err := di.Resolve[os.OS](container)
	handle(err)

	path, err := di.Resolve[*filepath.Processor](container)
	handle(err)

	console, err := di.Resolve[console.Console](container)
	handle(err)

	appName := "wrangle"
	plat := platform.Platform(o.Platform())
	if plat.IsWindows() {
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
				Value:   path.Join(o.Home(), ".wrangle", "config.yml"),
				EnvVars: []string{global.EnvConfig},
			},
			&cli.GenericFlag{
				Name:        global.FlagOutput,
				Aliases:     []string{"o"},
				Value:       enums.NewFormatEnum(),
				Required:    false,
				DefaultText: "table",
			},
		},
		Before: func(ctx *cli.Context) error {
			globalConfigFile := ctx.String(global.FlagConfig)
			resolver, err := app.GetResolver(ctx)
			if err != nil {
				return err
			}
			properties, err := di.Resolve[config.Properties](resolver)
			if err != nil {
				return err
			}
			properties.Set(config.GlobalConfigFilePathProperty, globalConfigFile)
			return nil
		},
	}

	// register
	app.Commands = []*cli.Command{
		commands.Bootstrap,
		commands.List,
		commands.Get,
		commands.List,
		commands.Initialize,
	}
	err = app.Run(console.Args())
	handle(err)
}

func handle(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
