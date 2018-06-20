package main

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/spf13/afero"

	"github.com/patrickhuber/cli-mgr/commands"
	"github.com/patrickhuber/cli-mgr/config"
	"github.com/patrickhuber/cli-mgr/option"
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

	processFactory := processes.NewOsProcessFactory()
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

	err = app.cliApplication.Run(os.Args)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func createApplication(
	manager store.Manager,
	fileSystem afero.Fs,
	processFactory processes.ProcessFactory,
	console ui.Console,
	platform string) (*application, error) {

	defaultConfigPath, err := config.GetConfigPath(&option.Options{})
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
		*createEnvCommand(manager, fileSystem, platform, console),
	}

	app := &application{cliApplication: cliApp}

	cliApp.Before = func(context *cli.Context) error {
		configFile := context.GlobalString("config")
		if configFile == "" {
			configFile = defaultConfigPath
		}
		configLoader := config.NewConfigLoader(fileSystem)
		configuration, err := configLoader.Load(configFile)
		if err != nil {
			return err
		}
		app.configuration = configuration
		return nil
	}

	return app, nil
}

func createRunCommand(
	manager store.Manager,
	fileSystem afero.Fs,
	processFactory processes.ProcessFactory) *cli.Command {
	runCommand := commands.NewRunCommand(
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
			configFile := context.GlobalString("config")
			processName := context.String("name")
			environmentName := context.String("environment")
			params := commands.NewRunCommandParams(configFile, processName, environmentName)
			return runCommand.ExecuteCommand(params)
		},
	}
}

func createEnvCommand(
	manager store.Manager,
	fileSystem afero.Fs,
	platform string,
	console ui.Console) *cli.Command {

	envCommand := commands.NewEnvCommand(
		manager,
		fileSystem,
		platform,
		console)

	return &cli.Command{
		Name:    "env",
		Aliases: []string{"e"},
		Usage:   "print command environemnt variables",
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
			configFile := context.GlobalString("config")
			processName := context.String("name")
			environmentName := context.String("environment")
			params := commands.NewRunCommandParams(configFile, processName, environmentName)
			return envCommand.ExecuteCommand(params)
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
