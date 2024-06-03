package main

import (
	"io"
	"log"

	"github.com/patrickhuber/go-di"
	"github.com/patrickhuber/go-xplat/console"
	"github.com/patrickhuber/go-xplat/env"
	"github.com/patrickhuber/go-xplat/filepath"
	"github.com/patrickhuber/go-xplat/os"
	"github.com/patrickhuber/go-xplat/platform"
	"github.com/patrickhuber/wrangle/internal/app"
	"github.com/patrickhuber/wrangle/internal/commands"
	"github.com/patrickhuber/wrangle/internal/enums"
	"github.com/patrickhuber/wrangle/internal/global"
	setup "github.com/patrickhuber/wrangle/internal/host"
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
		Metadata: map[string]any{
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
			if globalConfigFile == "" {
				return nil
			}

			resolver, err := app.GetResolver(ctx)
			if err != nil {
				return err
			}

			environment, err := di.Resolve[env.Environment](resolver)
			if err != nil {
				return err
			}

			_, ok := environment.Lookup(global.EnvConfig)
			if !ok {
				environment.Set(global.EnvConfig, globalConfigFile)
			}

			return nil
		},
	}

	// register
	app.Commands = []*cli.Command{
		commands.Bootstrap,
		commands.List,
		commands.Get,
		commands.Set,
		commands.Initialize,
		commands.Export,
		commands.Hook,
	}

	// this is a hack to get global options printed in the commands
	// see https://github.com/urfave/cli/issues/734#issuecomment-597344796
	globalOptionsTemplate := `
{{if .VisibleFlags}}GLOBAL OPTIONS:
   {{range $index, $option := .VisibleFlags}}{{if $index}}
   {{end}}{{$option}}{{end}}
{{end}}
`
	origHelpPrinterCustom := cli.HelpPrinterCustom
	defer func() {
		cli.HelpPrinterCustom = origHelpPrinterCustom
	}()
	cli.HelpPrinterCustom = func(out io.Writer, tmpl string, data any, customFuncs map[string]any) {

		// inject the application name
		appName := func() string {
			return app.Name
		}
		if customFuncs == nil {
			customFuncs = map[string]any{}
		}
		customFuncs["appname"] = appName

		// map the data to a map?
		origHelpPrinterCustom(out, tmpl, data, customFuncs)

		// run on the app context if this is a command
		if data != app {
			origHelpPrinterCustom(app.Writer, globalOptionsTemplate, app, nil)
		}
	}

	err = app.Run(console.Args())
	handle(err)
}

func handle(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
