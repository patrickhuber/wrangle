package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/spf13/afero"

	"github.com/patrickhuber/cli-mgr/commands"
	"github.com/patrickhuber/cli-mgr/config"
	"github.com/patrickhuber/cli-mgr/processes"
	"github.com/patrickhuber/cli-mgr/store"
	"github.com/patrickhuber/cli-mgr/ui"

	credhub "github.com/patrickhuber/cli-mgr/store/credhub"
	file "github.com/patrickhuber/cli-mgr/store/file"

	"github.com/urfave/cli"
)

type application struct {
	cliApplication *cli.App
	configuration  *config.Config
}

func main() {
	fileSystem := afero.NewOsFs()
	configStoreManager := createConfigStoreManager(fileSystem)
	validateConfigStoreManager(configStoreManager)

	processFactory := processes.NewOsFactory()
	console := ui.NewOSConsole()

	// creates the app
	// see https://github.com/urfave/cli#customization-1 for template
	app, err := createApplication(
		configStoreManager,
		fileSystem,
		processFactory,
		console,
		runtime.GOOS)
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
	fileSystem afero.Fs,
	processFactory processes.Factory,
	console ui.Console,
	platform string) (*cli.App, error) {

	defaultConfigPath, err := config.GetDefaultConfigPath()
	if err != nil {
		return nil, err
	}

	cliApp := cli.NewApp()
	cliApp.Usage = "a cli management tool"
	cliApp.Writer = console.Out()
	cliApp.ErrWriter = console.Error()

	cliApp.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "config, c",
			Usage:  "Load configuration from `FILE`",
			EnvVar: "CLI_MGR_CONFIG",
			Value:  defaultConfigPath,
		},
	}

	cliApp.Commands = []cli.Command{
		*createRunCommand(manager, fileSystem, processFactory),
		*createPrintCommand(manager, fileSystem, platform, console),
		*createEnvironmentsCommand(fileSystem, console),
		*createPackagesCommand(console),
	}

	cliApp.Before = func(context *cli.Context) error {

		configFile := context.GlobalString("config")
		if configFile == "" {
			configFile = defaultConfigPath
		}
		configLoader := config.NewLoader(fileSystem)
		configuration, err := configLoader.Load(configFile)
		if err != nil {
			return err
		}

		app := context.App
		if app.Metadata == nil {
			app.Metadata = make(map[string]interface{})
		}

		app.Metadata["configuration"] = configuration

		return nil
	}

	return cliApp, nil
}

func createRunCommand(
	manager store.Manager,
	fileSystem afero.Fs,
	processFactory processes.Factory) *cli.Command {
	runCommand := commands.NewRun(
		manager,
		fileSystem,
		processFactory)

	return &cli.Command{
		Name:    "run",
		Aliases: []string{"r"},
		Usage:   "run a command",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "name, n",
				Usage: "Execute command named `NAME`",
			},
			cli.StringFlag{
				Name:  "environment, e",
				Usage: "Use environment named `ENVIRONMENT`",
			},
		},
		Action: func(context *cli.Context) error {
			cfg, err := getConfigurationFromCliContext(context)
			if err != nil {
				return err
			}
			processName := context.String("name")
			environmentName := context.String("environment")
			params := commands.NewRunCommandParams(cfg, processName, environmentName)
			return runCommand.Execute(params)
		},
	}
}

func createPrintCommand(
	manager store.Manager,
	fileSystem afero.Fs,
	platform string,
	console ui.Console) *cli.Command {

	printCommand := commands.NewPrint(
		manager,
		fileSystem,
		platform,
		console)

	return &cli.Command{
		Name:    "print",
		Aliases: []string{"p"},
		Usage:   "print command environemnt variables",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "name, n",
				Usage: "process named `NAME`",
			},
			cli.StringFlag{
				Name:  "environment, e",
				Usage: "Use environment named `ENVIRONMENT`",
			},
		},
		Action: func(context *cli.Context) error {
			processName := context.String("name")
			environmentName := context.String("environment")
			cfg, err := getConfigurationFromCliContext(context)
			if err != nil {
				return err
			}
			params := commands.NewRunCommandParams(cfg, processName, environmentName)
			return printCommand.Execute(params)
		},
	}
}

func createPackagesCommand(
	console ui.Console) *cli.Command {
	packagesCommand := commands.NewPackages(console)
	return &cli.Command{
		Name:    "packages",
		Aliases: []string{"k"},
		Usage:   "prints the list of packages and versions in the config file",
		Action: func(context *cli.Context) error {
			cfg, err := getConfigurationFromCliContext(context)
			if err != nil {
				return err
			}
			return packagesCommand.Execute(cfg)
		},
	}
}

func getConfigurationFromCliContext(context *cli.Context) (*config.Config, error) {
	configuration, ok := context.App.Metadata["configuration"]
	if !ok {
		return nil, errors.New("unable to load configuration from configuration metadata")
	}
	cfg, ok := configuration.(*config.Config)
	if !ok {
		return nil, errors.New("configuration loaded from metadata is not the expected type of *config.Config")
	}
	return cfg, nil
}

func createEnvironmentsCommand(
	fileSystem afero.Fs,
	console ui.Console) *cli.Command {

	environmentsCommand := commands.NewEnvironments(
		fileSystem,
		console)

	return &cli.Command{
		Name:    "environments",
		Aliases: []string{"e"},
		Usage:   "prints the list of environments in the config file",
		Action: func(context *cli.Context) error {
			cfg, err := getConfigurationFromCliContext(context)
			if err != nil {
				return err
			}
			return environmentsCommand.ExecuteCommand(cfg)
		},
	}
}

func createConfigStoreManager(fileSystem afero.Fs) store.Manager {
	manager := store.NewManager()
	manager.Register(credhub.NewCredHubStoreProvider())
	manager.Register(file.NewFileStoreProvider(fileSystem))
	return manager
}

func validateConfigStoreManager(manager store.Manager) {
	if manager == nil {
		fmt.Printf("unable to create config store manager\n")
		os.Exit(1)
	}
}
