package main

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/spf13/afero"

	"github.com/patrickhuber/wrangle/commands"
	"github.com/patrickhuber/wrangle/config"
	"github.com/patrickhuber/wrangle/filesystem"
	"github.com/patrickhuber/wrangle/global"
	"github.com/patrickhuber/wrangle/processes"
	"github.com/patrickhuber/wrangle/store"
	"github.com/patrickhuber/wrangle/ui"

	credhub "github.com/patrickhuber/wrangle/store/credhub"
	"github.com/patrickhuber/wrangle/store/env"
	file "github.com/patrickhuber/wrangle/store/file"

	"github.com/urfave/cli"
)

type application struct {
	cliApplication *cli.App
	configuration  *config.Config
}

func main() {
	fileSystem := filesystem.NewOsFsWrapper(afero.NewOsFs())
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
	fileSystem filesystem.FsWrapper,
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
	cliApp.Version = "0.4.3"

	cliApp.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "config, c",
			Usage:  "Load configuration from `FILE`",
			EnvVar: global.ConfigFileKey,
			Value:  defaultConfigPath,
		},
	}

	cliApp.Commands = []cli.Command{
		*createRunCommand(manager, fileSystem, processFactory, console),
		*createPrintCommand(manager, fileSystem, platform, console),
		*createPrintEnvCommand(manager, fileSystem, platform, console),
		*createEnvironmentsCommand(fileSystem, console),
		*createPackagesCommand(fileSystem, console),
		*createInstallCommand(fileSystem, platform),
		*createEnvCommand(console),
	}
	return cliApp, nil
}

func createRunCommand(
	manager store.Manager,
	fileSystem afero.Fs,
	processFactory processes.Factory,
	console ui.Console) *cli.Command {
	runCommand := commands.NewRun(
		manager,
		fileSystem,
		processFactory,
		console)

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
			cfg, err := createConfiguration(context, fileSystem)
			if err != nil {
				return err
			}
			processName := context.String("name")
			environmentName := context.String("environment")
			params := commands.NewProcessParams(cfg, environmentName, processName)
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
		Usage:   "prints the process as it would be executed",
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
			cfg, err := createConfiguration(context, fileSystem)
			if err != nil {
				return err
			}
			params := commands.NewProcessParams(cfg, environmentName, processName)
			return printCommand.Execute(params)
		},
	}
}

func createPrintEnvCommand(
	manager store.Manager,
	fileSystem afero.Fs,
	platform string,
	console ui.Console) *cli.Command {

	printEnvCommand := commands.NewPrintEnv(
		manager,
		fileSystem,
		platform,
		console)

	return &cli.Command{
		Name:  "print-env",
		Usage: "print command environemnt variables",
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
			cfg, err := createConfiguration(context, fileSystem)
			if err != nil {
				return err
			}
			params := commands.NewProcessParams(cfg, environmentName, processName)
			return printEnvCommand.Execute(params)
		},
	}
}

func createPackagesCommand(
	fileSystem afero.Fs,
	console ui.Console) *cli.Command {
	packagesCommand := commands.NewPackages(console)
	return &cli.Command{
		Name:    "packages",
		Aliases: []string{"k"},
		Usage:   "prints the list of packages and versions in the config file",
		Action: func(context *cli.Context) error {
			cfg, err := createConfiguration(context, fileSystem)
			if err != nil {
				return err
			}
			return packagesCommand.Execute(cfg)
		},
	}
}

func createInstallCommand(
	fileSystem filesystem.FsWrapper,
	platform string) *cli.Command {
	return &cli.Command{
		Name:    "install",
		Aliases: []string{"i"},
		Usage:   "installs the package with the given `NAME` for the current platform",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "name, n",
				Usage: "package named `NAME`",
			},
			cli.StringFlag{
				Name:   "path, p",
				Usage:  "the package install path",
				EnvVar: global.PackagePathKey,
			},
		},
		Action: func(context *cli.Context) error {
			cfg, err := createConfiguration(context, fileSystem)
			packageName := context.String("name")
			packageInstallPath := context.String("path")
			if err != nil {
				return err
			}
			installPackageCommand, err := commands.NewInstall(platform, packageInstallPath, fileSystem, ui.NewOSConsole())
			if err != nil {
				return err
			}
			return installPackageCommand.Execute(cfg, packageName)
		},
	}
}

func createEnvCommand(console ui.Console) *cli.Command {
	return &cli.Command{
		Name:  "env",
		Usage: "prints values of all associated environment variables",
		Action: func(context *cli.Context) error {
			return commands.NewEnv(console).Execute()
		},
	}
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
			cfg, err := createConfiguration(context, fileSystem)
			if err != nil {
				return err
			}
			return environmentsCommand.ExecuteCommand(cfg)
		},
	}
}

func createConfiguration(context *cli.Context, fileSystem afero.Fs) (*config.Config, error) {
	configFile := context.GlobalString("config")
	var err error
	if configFile == "" {
		configFile, err = config.GetDefaultConfigPath()
		if err != nil {
			return nil, err
		}
	}
	configLoader := config.NewLoader(fileSystem)
	return configLoader.Load(configFile)
}

func createConfigStoreManager(fileSystem afero.Fs) store.Manager {
	manager := store.NewManager()
	manager.Register(credhub.NewCredHubStoreProvider())
	manager.Register(file.NewFileStoreProvider(fileSystem))
	manager.Register(env.NewEnvStoreProvider())
	return manager
}

func validateConfigStoreManager(manager store.Manager) {
	if manager == nil {
		fmt.Printf("unable to create config store manager\n")
		os.Exit(1)
	}
}
