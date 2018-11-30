package main

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/patrickhuber/wrangle/tasks"

	"github.com/patrickhuber/wrangle/collections"
	"github.com/patrickhuber/wrangle/packages"

	"github.com/spf13/afero"

	"github.com/patrickhuber/wrangle/commands"
	"github.com/patrickhuber/wrangle/config"
	"github.com/patrickhuber/wrangle/crypto"
	"github.com/patrickhuber/wrangle/env"
	"github.com/patrickhuber/wrangle/filesystem"
	"github.com/patrickhuber/wrangle/global"
	"github.com/patrickhuber/wrangle/processes"
	"github.com/patrickhuber/wrangle/renderers"
	"github.com/patrickhuber/wrangle/store"
	"github.com/patrickhuber/wrangle/ui"
	"github.com/patrickhuber/wrangle/services"

	credhub "github.com/patrickhuber/wrangle/store/credhub"
	store_env "github.com/patrickhuber/wrangle/store/env"
	file "github.com/patrickhuber/wrangle/store/file"

	"github.com/urfave/cli"
)

type application struct {
	cliApplication *cli.App
	configuration  *config.Config
}

func main() {
	// create platform, filesystem and console
	platform := runtime.GOOS
	fileSystem := filesystem.NewOsFsWrapper(afero.NewOsFs())
	console := ui.NewOSConsole()

	// create config store m anager
	configStoreManager, err := createConfigStoreManager(fileSystem, platform)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	validateConfigStoreManager(configStoreManager)

	// create process factory
	processFactory := processes.NewOsFactory()

	// create env dictionary
	envDictionary := env.NewDictionary()

	// create config loader
	loader := config.NewLoader(fileSystem)

	// create task providers
	taskProviders := tasks.NewProviderRegistry()
	taskProviders.Register(tasks.NewDownloadProvider(fileSystem, console))
	taskProviders.Register(tasks.NewExtractProvider(fileSystem, console))
	taskProviders.Register(tasks.NewLinkProvider(fileSystem, console))
	taskProviders.Register(tasks.NewMoveProvider(fileSystem, console))

	// create package manager
	packagesManager := packages.NewManager(fileSystem, taskProviders)

	// creates the app
	// see https://github.com/urfave/cli#customization-1 for template
	app, err := createApplication(
		configStoreManager,
		fileSystem,
		processFactory,
		console,
		platform,
		envDictionary,
		loader,
		packagesManager)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

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
	cliApp.Version = "0.9.0"

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
	packagesService := services.NewPackagesService(fileSystem, console)
	installService, err := services.NewInstallService(platform, fileSystem, packagesManager, loader)
	envService := services.NewEnvService(console, envDictionary)
	storesService := services.NewStoresService(console, loader)
	processesService := services.NewProcessesService( console, loader)

	if err != nil{
		return nil, err
	}

	cliApp.Commands = []cli.Command{
		*commands.CreateInitCommand(initService),
		*commands.CreateRunCommand(runService),
		*commands.CreatePrintCommand(printService),
		*commands.CreatePrintEnvCommand(printService),
		*commands.CreatePackagesCommand(packagesService),
		*commands.CreateInstallCommand(installService),
		*commands.CreateEnvCommand(envService),
		*commands.CreateStoresCommand(storesService),
		*commands.CreateListProcessesCommand(processesService),
	}
	return cliApp, nil
}




func createConfigStoreManager(fileSystem afero.Fs, platform string) (store.Manager, error) {
	manager := store.NewManager()
	factory, err := crypto.NewPgpFactory(fileSystem, platform)
	if err != nil {
		return nil, err
	}
	manager.Register(credhub.NewCredHubStoreProvider())
	manager.Register(file.NewFileStoreProvider(fileSystem, factory))
	manager.Register(store_env.NewEnvStoreProvider())
	return manager, nil
}

func validateConfigStoreManager(manager store.Manager) {
	if manager == nil {
		fmt.Printf("unable to create config store manager\n")
		os.Exit(1)
	}
}
