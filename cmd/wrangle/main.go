package main

import (
	"log"
	"os"
	"strings"

	"github.com/patrickhuber/go-di"
	"github.com/patrickhuber/wrangle/internal/app"
	"github.com/patrickhuber/wrangle/internal/commands"
	"github.com/patrickhuber/wrangle/internal/setup"
	"github.com/patrickhuber/wrangle/pkg/config"
	"github.com/patrickhuber/wrangle/pkg/crosspath"
	"github.com/patrickhuber/wrangle/pkg/enums"
	"github.com/patrickhuber/wrangle/pkg/global"
	"github.com/patrickhuber/wrangle/pkg/operatingsystem"
	"github.com/urfave/cli/v2"
)

// set with -ldflags
var version = ""

func main() {
	s := setup.New()
	container := s.Container()
	o, err := di.Resolve[operatingsystem.OS](container)
	handle(err)

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
	err = app.Run(os.Args)
	handle(err)
}

func handle(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
