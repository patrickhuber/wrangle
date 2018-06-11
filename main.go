package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/afero"

	"github.com/patrickhuber/cli-mgr/commands"
	"github.com/patrickhuber/cli-mgr/config"
	"github.com/patrickhuber/cli-mgr/option"
	"github.com/patrickhuber/cli-mgr/store"

	credhub "github.com/patrickhuber/cli-mgr/store/credhub"
	file "github.com/patrickhuber/cli-mgr/store/file"

	"github.com/urfave/cli"
)

func main() {
	configStoreManager := createConfigStoreManager()
	validateConfigStoreManager(configStoreManager)

	fileSystem := afero.NewOsFs()

	configPath, err := config.GetConfigPath(&option.Options{})
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	// creates the app
	// see https://github.com/urfave/cli#customization-1 for template
	app := cli.NewApp()
	app.Usage = "a cli management tool"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "config, c",
			Usage:  "Load configuration from `FILE`",
			EnvVar: "CLI_MGR_CONFIG",
			Value:  configPath,
		},
	}
	runCommand := commands.NewRunCommand(
		configStoreManager,
		fileSystem,
		commands.NewOsProcessFactory())

	app.Commands = []cli.Command{
		{
			Name:    "run",
			Aliases: []string{"r"},
			Usage:   "run a command",
			Action:  runCommand.ExecuteCommand,
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
		},
		{
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
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
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
