package main

import (
	"io"
	"log"

	"github.com/patrickhuber/go-cross/console"
	"github.com/patrickhuber/go-cross/env"
	"github.com/patrickhuber/go-cross/filepath"
	"github.com/patrickhuber/go-cross/os"
	"github.com/patrickhuber/go-di"
	"github.com/patrickhuber/wrangle/internal/app"
	"github.com/patrickhuber/wrangle/internal/commands"
	"github.com/patrickhuber/wrangle/internal/config"
	"github.com/patrickhuber/wrangle/internal/enums"
	"github.com/patrickhuber/wrangle/internal/global"
	"github.com/patrickhuber/wrangle/internal/host"
	"github.com/urfave/cli/v2"
)

// set with -ldflags
var version = ""

func main() {
	h := host.New()

	container := h.Container()

	o, err := di.Resolve[os.OS](container)
	handle(err)

	path, err := di.Resolve[filepath.Provider](container)
	handle(err)

	console, err := di.Resolve[console.Console](container)
	handle(err)

	environment, err := di.Resolve[env.Environment](container)
	handle(err)

	appName, err := config.GetAppName("wrangle", o.Platform())
	handle(err)

	root, err := config.GetRoot(environment, o.Platform())
	handle(err)

	home, err := o.Home()
	handle(err)

	app := &cli.App{
		Metadata: map[string]any{
			global.MetadataDependencyInjection: container,
		},
		Name:        appName,
		Version:     version,
		Usage:       "A DevOps Environment Management CLI",
		Description: "A DevOps Environment Management CLI",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    global.FlagBin,
				Aliases: []string{"b"},
				EnvVars: []string{global.EnvBin},
				Value:   config.GetDefaultBinPath(path, root),
			},
			&cli.StringFlag{
				Name:    global.FlagRoot,
				Aliases: []string{"r"},
				EnvVars: []string{global.EnvRoot},
				Value:   root,
			},
			&cli.StringFlag{
				Name:    global.FlagPackages,
				Aliases: []string{"p"},
				EnvVars: []string{global.EnvPackages},
				Value:   config.GetDefaultPackagesPath(path, root),
			},
			&cli.StringFlag{
				Name:    global.FlagSystemConfig,
				Aliases: []string{"g"},
				EnvVars: []string{global.EnvSystemConfig},
				Value:   config.GetDefaultSystemConfigPath(path, root),
			},
			&cli.StringFlag{
				Name:    global.FlagUserConfig,
				Aliases: []string{"u"},
				EnvVars: []string{global.EnvUserConfig},
				Value:   config.GetDefaultUserConfigPath(path, home),
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

			container, err := app.GetContainer(ctx)
			if err != nil {
				return err
			}

			// register the cli context
			di.RegisterInstance(container, config.CliContext(ctx))
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
		commands.Interpolate,
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
