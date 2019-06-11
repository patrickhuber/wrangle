package main

import (
	"github.com/patrickhuber/wrangle/templates"
	"github.com/patrickhuber/wrangle/collections"
	"github.com/patrickhuber/wrangle/commands"
	"github.com/patrickhuber/wrangle/config"
	"github.com/patrickhuber/wrangle/env"
	"github.com/patrickhuber/wrangle/feed"
	"github.com/patrickhuber/wrangle/filesystem"
	"github.com/patrickhuber/wrangle/global"
	"github.com/patrickhuber/wrangle/packages"
	"github.com/patrickhuber/wrangle/processes"
	"github.com/patrickhuber/wrangle/renderers"
	"github.com/patrickhuber/wrangle/services"
	"github.com/patrickhuber/wrangle/store"
	"github.com/patrickhuber/wrangle/ui"
	"github.com/urfave/cli"
)

// set with -ldflags
var version = ""

func createApplication(
	workingDirectory string,
	manager store.Manager,
	fileSystem filesystem.FsWrapper,
	processFactory processes.Factory,
	console ui.Console,
	platform string,
	envDictionary collections.Dictionary,
	packagesManager packages.Manager,
	templateFactory templates.Factory) (*cli.App, error) {

	rendererFactory := renderers.NewFactory(env.NewDictionary())

	defaultConfigPath, err := config.GetDefaultConfigPath(workingDirectory)
	if err != nil {
		return nil, err
	}

	cliApp := cli.NewApp()
	cliApp.Usage = "a cli management tool"
	cliApp.Writer = console.Out()
	cliApp.ErrWriter = console.Error()
	cliApp.Version = version

	cliApp.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "config, c",
			Usage:  "Load configuration from `FILE`",
			EnvVar: global.ConfigFileKey,
			Value:  defaultConfigPath,
		},
	}

	initService := services.NewInitService(fileSystem)
	runService := services.NewRunService(manager, fileSystem, processFactory, console, templateFactory)
	printService := services.NewPrintService(manager, fileSystem, console, rendererFactory, templateFactory)
	packagesServiceFactory := services.NewPackageServiceFactory(console)
	feedServiceFactory := feed.NewFeedServiceFactory(fileSystem)
	installService, err := services.NewInstallService(platform, fileSystem, packagesManager)
	envService := services.NewEnvService(console, envDictionary)
	storesService := services.NewStoresService(console)
	processesService := services.NewProcessesService(console)
	credentialServiceFactory := services.NewCredentialServiceFactory(manager, fileSystem, templateFactory)

	if err != nil {
		return nil, err
	}

	cliApp.Commands = []cli.Command{
		*commands.CreateInitCommand(cliApp, initService),
		*commands.CreateRunCommand(cliApp, runService, fileSystem),
		*commands.CreatePrintCommand(cliApp, printService, fileSystem),
		*commands.CreatePrintEnvCommand(cliApp, printService, fileSystem),
		*commands.CreatePackagesCommand(packagesServiceFactory, feedServiceFactory),
		*commands.CreateInstallCommand(installService),
		*commands.CreateEnvCommand(envService),
		*commands.CreateStoresCommand(cliApp, storesService, fileSystem),
		*commands.CreateListProcessesCommand(cliApp, processesService, fileSystem),
		*commands.CreateMoveCommand(cliApp, credentialServiceFactory),
		*commands.CreateCopyCommand(cliApp, credentialServiceFactory),
	}

	return cliApp, nil
}
