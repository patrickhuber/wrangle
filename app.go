package main

import (
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
	manager store.Manager,
	fileSystem filesystem.FsWrapper,
	processFactory processes.Factory,
	console ui.Console,
	platform string,
	envDictionary collections.Dictionary,
	loader config.Loader,
	packagesManager packages.Manager) (*cli.App, error) {

	rendererFactory := renderers.NewFactory(env.NewDictionary())

	defaultConfigPath, err := config.GetDefaultConfigPath()
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
	feedServiceFactory := services.NewFeedServiceFactory(fileSystem)
	installService, err := services.NewInstallService(platform, fileSystem, packagesManager, loader)
	envService := services.NewEnvService(console, envDictionary)
	storesService := services.NewStoresService(console, loader)
	processesService := services.NewProcessesService(console, loader)

	if err != nil {
		return nil, err
	}

	cliApp.Commands = []cli.Command{
		*commands.CreateInitCommand(initService),
		*commands.CreateRunCommand(runService),
		*commands.CreatePrintCommand(printService),
		*commands.CreatePrintEnvCommand(printService),
		*commands.CreatePackagesCommand(packagesServiceFactory, feedServiceFactory),
		*commands.CreateInstallCommand(installService),
		*commands.CreateEnvCommand(envService),
		*commands.CreateStoresCommand(storesService),
		*commands.CreateListProcessesCommand(processesService),
	}
	return cliApp, nil
}
