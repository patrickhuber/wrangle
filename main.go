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

func main() {
	configStoreManager := createConfigStoreManager()
	validateConfigStoreManager(configStoreManager)

	fileSystem := afero.NewOsFs()
	processFactory := processes.NewOsProcessFactory()
	console := ui.NewOSConsole()

	// creates the app
	// see https://github.com/urfave/cli#customization-1 for template
	app, err := createCliApplication(
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

func createCliApplication(
	manager store.Manager,
	fileSystem afero.Fs,
	processFactory processes.ProcessFactory,
	console ui.Console,
	platform string,
) (*cli.App, error) {

	defaultConfigPath, err := config.GetConfigPath(&option.Options{})
	if err != nil {
		return nil, err
	}

	app := cli.NewApp()
	app.Usage = "a cli management tool"
	app.Writer = console.Out()
	app.ErrWriter = console.Error()

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "config, c",
			Usage:  "Load configuration from `FILE`",
			EnvVar: "CLI_MGR_CONFIG",
			Value:  defaultConfigPath,
		},
	}

	app.Commands = []cli.Command{
		*createRunCommand(manager, fileSystem, processFactory),
		*createEnvCommand(manager, fileSystem, platform, console),
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

func createConfigStoreManager() store.Manager {
	manager := store.NewManager()
	manager.Register(&credhub.CredHubConfigStoreProvider{})
	manager.Register(&file.FileConfigStoreProvider{})
	return manager
}

func validateConfigStoreManager(manager store.Manager) {
	if manager == nil {
		fmt.Printf("unable to create config store manager\n")
		os.Exit(1)
	}
}
