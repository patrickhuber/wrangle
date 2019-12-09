package main

import (
	"github.com/patrickhuber/wrangle/collections"
	"github.com/patrickhuber/wrangle/commands"
	"github.com/patrickhuber/wrangle/config"
	"github.com/patrickhuber/wrangle/credentials"
	"github.com/patrickhuber/wrangle/env"
	"github.com/patrickhuber/wrangle/feed"
	"github.com/patrickhuber/wrangle/filesystem"
	"github.com/patrickhuber/wrangle/global"
	"github.com/patrickhuber/wrangle/initialize"
	"github.com/patrickhuber/wrangle/packages"
	"github.com/patrickhuber/wrangle/processes"
	"github.com/patrickhuber/wrangle/renderers"
	"github.com/patrickhuber/wrangle/renderers/items"
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
	fileSystem filesystem.FileSystem,
	processFactory processes.Factory,
	console ui.Console,
	platform string,
	envDictionary collections.Dictionary,
	packagesManager packages.Manager) (*cli.App, error) {

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

	initService := initialize.NewService(fileSystem)
	runService := services.NewRunService(manager, fileSystem, processFactory, console)
	printService := services.NewPrintService(manager, fileSystem, console, rendererFactory)
	packagesServiceFactory := packages.NewServiceFactory(console)
	feedServiceFactory := feed.NewFeedServiceFactory(fileSystem)
	installService, err := packages.NewInstallService(platform, fileSystem, packagesManager)
	envService := env.NewService(console, envDictionary)
	storesService := store.NewService(console)
	processesService := processes.NewService(console)
	credentialServiceFactory := credentials.NewServiceFactory(manager, fileSystem)
	renderFactory := items.NewFactory()

	if err != nil {
		return nil, err
	}

	cliApp.Commands = []cli.Command{
		commands.CreateListCommand(
			commands.CreateListPackagesCommand(packagesServiceFactory, feedServiceFactory),
			commands.CreateListProcessesCommand(cliApp, processesService, fileSystem),
			commands.CreateListStoresCommand(cliApp, storesService, fileSystem),
			commands.CreateListSecretsCommand(cliApp, console, credentialServiceFactory, renderFactory),
		),
		commands.CreateGetCommand(),
		commands.CreatePrintCommand(
			commands.CreatePrintProcessCommand(cliApp, printService, fileSystem),
			commands.CreatePrintEnvCommand(cliApp, printService, fileSystem),
		),
		*commands.CreateInitCommand(cliApp, initService),
		*commands.CreateRunCommand(cliApp, runService, fileSystem),
		*commands.CreateInstallCommand(installService),
		*commands.CreateEnvCommand(envService),
		*commands.CreateMoveCommand(cliApp, credentialServiceFactory),
		*commands.CreateCopyCommand(cliApp, credentialServiceFactory),
	}

	return cliApp, nil
}
