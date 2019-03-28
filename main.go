package main

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/spf13/afero"

	"github.com/patrickhuber/wrangle/collections"
	"github.com/patrickhuber/wrangle/config"
	"github.com/patrickhuber/wrangle/crypto"
	"github.com/patrickhuber/wrangle/env"
	"github.com/patrickhuber/wrangle/filesystem"
	"github.com/patrickhuber/wrangle/packages"
	"github.com/patrickhuber/wrangle/processes"
	"github.com/patrickhuber/wrangle/store"
	"github.com/patrickhuber/wrangle/tasks"
	"github.com/patrickhuber/wrangle/ui"

	credhub "github.com/patrickhuber/wrangle/store/credhub"
	store_env "github.com/patrickhuber/wrangle/store/env"
	file "github.com/patrickhuber/wrangle/store/file"
)

func main() {
	// create platform, filesystem, working directory and console
	platform := runtime.GOOS
	fileSystem := filesystem.NewOsFsWrapper(afero.NewOsFs())
	console := ui.NewOSConsole()
	workingDirectory, err := os.Getwd()
	failOnError(err)

	// create env dictionary
	environmentVariables := env.NewDictionary()

	// create config store manager
	configStoreManager, err := createConfigStoreManager(fileSystem, platform, environmentVariables)
	failOnError(err)

	validateConfigStoreManager(configStoreManager)

	// create process factory
	processFactory := processes.NewOsFactory()

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
		workingDirectory,
		configStoreManager,
		fileSystem,
		processFactory,
		console,
		platform,
		environmentVariables,
		loader,
		packagesManager)
	failOnError(err)

	err = app.Run(os.Args)
	failOnError(err)
}

func createConfigStoreManager(fileSystem afero.Fs, platform string, environmentVariables collections.Dictionary) (store.Manager, error) {
	manager := store.NewManager()
	factory, err := crypto.NewPgpFactory(fileSystem, platform)
	if err != nil {
		return nil, err
	}
	manager.Register(credhub.NewCredHubStoreProvider())
	manager.Register(file.NewFileStoreProvider(fileSystem, factory))
	manager.Register(store_env.NewEnvStoreProvider(environmentVariables))
	return manager, nil
}

func validateConfigStoreManager(manager store.Manager) {
	if manager == nil {
		fmt.Printf("unable to create config store manager\n")
		os.Exit(1)
	}
}

func failOnError(err error) {
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
