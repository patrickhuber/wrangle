package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/patrickhuber/wrangle/collections"

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
	platform := runtime.GOOS
	fileSystem := filesystem.NewOsFsWrapper(afero.NewOsFs())
	configStoreManager, err := createConfigStoreManager(fileSystem, platform)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	validateConfigStoreManager(configStoreManager)

	processFactory := processes.NewOsFactory()
	console := ui.NewOSConsole()
	envDictionary := env.NewDictionary()

	// creates the app
	// see https://github.com/urfave/cli#customization-1 for template
	app, err := createApplication(
		configStoreManager,
		fileSystem,
		processFactory,
		console,
		platform,
		envDictionary)
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
	envDictionary collections.Dictionary) (*cli.App, error) {

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

	cliApp.Commands = []cli.Command{
		*createRunCommand(manager, fileSystem, processFactory, console),
		*createPrintCommand(manager, fileSystem, console, rendererFactory),
		*createPrintEnvCommand(manager, fileSystem, console, rendererFactory),
		*createEnvironmentsCommand(fileSystem, console),
		*createPackagesCommand(fileSystem, console),
		*createInstallCommand(fileSystem, platform),
		*createEnvCommand(console, envDictionary),
		*createStoresCommand(fileSystem, console),
		*createListProcessesCommand(fileSystem, console),
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
		Name:      "run",
		Aliases:   []string{"r"},
		Usage:     "run a command",
		ArgsUsage: "<process name> [arguments]",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "environment, e",
				Usage:  "Use environment named `ENVIRONMENT`",
				EnvVar: "WRANGLE_ENVIRONMENT",
			},
		},
		Action: func(context *cli.Context) error {
			cfg, err := createConfiguration(context, fileSystem)
			if err != nil {
				return err
			}

			processName := context.Args().First()
			if strings.TrimSpace(processName) == "" {
				return errors.New("process name argument is required")
			}

			environmentName := context.String("environment")

			additionalArguments := context.Args().Tail()

			params := commands.NewProcessParams(cfg, environmentName, processName, additionalArguments...)
			return runCommand.Execute(params)
		},
	}
}

func createPrintCommand(
	manager store.Manager,
	fileSystem afero.Fs,
	console ui.Console,
	rendererFactory renderers.Factory) *cli.Command {

	printCommand := commands.NewPrint(
		manager,
		fileSystem,
		console,
		rendererFactory)

	return &cli.Command{
		Name:      "print",
		Aliases:   []string{"p"},
		Usage:     "prints the process as it would be executed",
		ArgsUsage: "<process name> [arguments]",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "environment, e",
				Usage: "Use environment named `ENVIRONMENT`",
			},
			cli.StringFlag{
				Name:  "format, f",
				Usage: "Print for with the given format (bash|powershell)",
			},
		},
		Action: func(context *cli.Context) error {
			processName := context.Args().First()
			if strings.TrimSpace(processName) == "" {
				return errors.New("process name argument is required")
			}

			environmentName := context.String("environment")
			format := context.String("format")

			cfg, err := createConfiguration(context, fileSystem)
			if err != nil {
				return err
			}
			params := &commands.PrintParams{
				Configuration:   cfg,
				EnvironmentName: environmentName,
				ProcessName:     processName,
				Format:          format,
				Include: commands.PrintParamsInclude{
					ProcessAndArgs: true,
				},
			}

			return printCommand.Execute(params)
		},
	}
}

func createPrintEnvCommand(
	manager store.Manager,
	fileSystem afero.Fs,
	console ui.Console,
	rendererFactory renderers.Factory) *cli.Command {

	printCommand := commands.NewPrint(
		manager,
		fileSystem,
		console,
		rendererFactory)

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
			cli.StringFlag{
				Name:  "format, f",
				Usage: "Print for with the given format (bash|powershell)",
			},
		},
		Action: func(context *cli.Context) error {
			processName := context.Args().First()
			if strings.TrimSpace(processName) == "" {
				return errors.New("process name argument is required")
			}
			environmentName := context.String("environment")
			format := context.String("format")
			cfg, err := createConfiguration(context, fileSystem)
			if err != nil {
				return err
			}
			params := &commands.PrintParams{
				Configuration:   cfg,
				EnvironmentName: environmentName,
				ProcessName:     processName,
				Format:          format}
			return printCommand.Execute(params)
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
		Name:      "install",
		Aliases:   []string{"i"},
		Usage:     "installs the package with the given `NAME` for the current platform",
		ArgsUsage: "<package name>",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "path, p",
				Usage:  "the package install path",
				EnvVar: global.PackagePathKey,
			},
		},
		Action: func(context *cli.Context) error {
			cfg, err := createConfiguration(context, fileSystem)
			packageName := context.Args().First()
			if strings.TrimSpace(packageName) == "" {
				return errors.New("missing required argument package name")
			}
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

func createEnvCommand(console ui.Console, dictionary collections.Dictionary) *cli.Command {
	return &cli.Command{
		Name:  "env",
		Usage: "prints values of all associated environment variables",
		Action: func(context *cli.Context) error {
			return commands.NewEnv(console, dictionary).Execute()
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
			return environmentsCommand.Execute(cfg)
		},
	}
}

func createListProcessesCommand(
	fileSystem afero.Fs,
	console ui.Console) *cli.Command {

	listProcessesCommand := commands.NewListProcesses(
		console)

	return &cli.Command{
		Name:  "processes",
		Usage: "prints the list of processes for the given environment in the config file",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "environment, e",
				Usage: "the environment name",
			},
		},
		Action: func(context *cli.Context) error {
			if !context.IsSet("environment") {
				return fmt.Errorf("environment flag is required")
			}
			environmentName := context.String("environment")
			cfg, err := createConfiguration(context, fileSystem)
			if err != nil {
				return err
			}
			return listProcessesCommand.Execute(cfg, environmentName)
		},
	}
}

func createStoresCommand(
	fileSystem afero.Fs,
	console ui.Console) *cli.Command {

	listStoresCommand := commands.NewStores(
		console)

	return &cli.Command{
		Name:    "stores",
		Aliases: []string{"s"},
		Usage:   "prints the list of stores in the config file",
		Action: func(context *cli.Context) error {
			cfg, err := createConfiguration(context, fileSystem)
			if err != nil {
				return err
			}
			return listStoresCommand.Execute(cfg)
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
