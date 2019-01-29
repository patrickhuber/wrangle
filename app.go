package main

import (
	"github.com/patrickhuber/wrangle/feed"
	"github.com/patrickhuber/wrangle/collections"
	"github.com/patrickhuber/wrangle/commands"
	"github.com/patrickhuber/wrangle/config"
	"github.com/patrickhuber/wrangle/env"
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
	loader config.Loader,
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

	initService := services.NewInitService(fileSystem)
	runService := services.NewRunService(manager, fileSystem, processFactory, console, loader)
	printService := services.NewPrintService(manager, fileSystem, console, rendererFactory, loader)
	packagesServiceFactory := services.NewPackageServiceFactory(console)
	feedServiceFactory := feed.NewFeedServiceFactory(fileSystem)
	installService, err := services.NewInstallService(platform, fileSystem, packagesManager, loader)
	envService := services.NewEnvService(console, envDictionary)
	storesService := services.NewStoresService(console, loader)
	processesService := services.NewProcessesService(console, loader)
	credentialServiceFactory := services.NewCredentialServiceFactory(manager, loader)

	if err != nil {
		return nil, err
	}

	cliApp.Commands = []cli.Command{
		*commands.CreateInitCommand(cliApp, initService),
		*commands.CreateRunCommand(cliApp, runService),
		*commands.CreatePrintCommand(cliApp, printService),
		*commands.CreatePrintEnvCommand(cliApp, printService),
		*commands.CreatePackagesCommand(packagesServiceFactory, feedServiceFactory),
		*commands.CreateInstallCommand(installService),
		*commands.CreateEnvCommand(envService),
		*commands.CreateStoresCommand(cliApp, storesService),
		*commands.CreateListProcessesCommand(cliApp, processesService),
		*commands.CreateMoveCommand(cliApp, credentialServiceFactory),
		*commands.CreateCopyCommand(cliApp, credentialServiceFactory),
	}
	
	
	return cliApp, nil
}
