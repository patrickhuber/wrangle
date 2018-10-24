package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"

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

	cliApp.Commands = []cli.Command{
		*createInitCommand(fileSystem),
		*createRunCommand(manager, fileSystem, processFactory, console, loader),
		*createPrintCommand(manager, fileSystem, console, rendererFactory, loader),
		*createPrintEnvCommand(manager, fileSystem, console, rendererFactory, loader),
		*createPackagesCommand(fileSystem, console, loader),
		*createInstallCommand(fileSystem, console, platform, packagesManager, loader),
		*createEnvCommand(console, envDictionary),
		*createStoresCommand(fileSystem, console, loader),
		*createListProcessesCommand(fileSystem, console, loader),
	}
	return cliApp, nil
}

func createInitCommand(
	fileSystem afero.Fs,
) *cli.Command {

	initCommand := commands.NewInitCommand(fileSystem)
	return &cli.Command{
		Name:  "init",
		Usage: "initialize the wrangle configuration",
		Action: func(context *cli.Context) error {
			return initCommand.Execute(context.GlobalString("config"))
		},
	}
}

func createRunCommand(
	manager store.Manager,
	fileSystem afero.Fs,
	processFactory processes.Factory,
	console ui.Console,
	loader config.Loader) *cli.Command {
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
		Action: func(context *cli.Context) error {
			cfg, err := loadConfiguration(context, loader)
			if err != nil {
				return err
			}

			processName := context.Args().First()
			if strings.TrimSpace(processName) == "" {
				return errors.New("process name argument is required")
			}

			additionalArguments := context.Args().Tail()

			params := commands.NewProcessParams(cfg, processName, additionalArguments...)
			return runCommand.Execute(params)
		},
	}
}

func createPrintCommand(
	manager store.Manager,
	fileSystem afero.Fs,
	console ui.Console,
	rendererFactory renderers.Factory,
	loader config.Loader) *cli.Command {

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
				Name:  "format, f",
				Usage: "Print for with the given format (bash|powershell)",
			},
		},
		Action: func(context *cli.Context) error {
			processName := context.Args().First()
			if strings.TrimSpace(processName) == "" {
				return errors.New("process name argument is required")
			}

			format := context.String("format")

			cfg, err := loadConfiguration(context, loader)
			if err != nil {
				return err
			}
			params := &commands.PrintParams{
				Configuration: cfg,
				ProcessName:   processName,
				Format:        format,
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
	rendererFactory renderers.Factory,
	loader config.Loader) *cli.Command {

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
				Name:  "format, f",
				Usage: "Print for with the given format (bash|powershell)",
			},
		},
		Action: func(context *cli.Context) error {
			processName := context.Args().First()
			if strings.TrimSpace(processName) == "" {
				return errors.New("process name argument is required")
			}
			format := context.String("format")
			cfg, err := loadConfiguration(context, loader)
			if err != nil {
				return err
			}
			params := &commands.PrintParams{
				Configuration: cfg,
				ProcessName:   processName,
				Format:        format}
			return printCommand.Execute(params)
		},
	}
}

func createPackagesCommand(
	fileSystem afero.Fs,
	console ui.Console,
	loader config.Loader) *cli.Command {
	return &cli.Command{
		Name:    "packages",
		Aliases: []string{"k"},
		Usage:   "prints the list of packages and versions in the config file",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "path, p",
				Usage:  "the package install path",
				EnvVar: global.PackagePathKey,
			},
		},
		Action: func(context *cli.Context) error {
			packagesPath := context.String("path")
			packagesCommand := commands.NewPackages(fileSystem, console, packagesPath)
			cfg, err := loadConfiguration(context, loader)
			if err != nil {
				return err
			}
			return packagesCommand.Execute(cfg)
		},
	}
}

func createInstallCommand(
	fileSystem filesystem.FsWrapper,
	console ui.Console,
	platform string,
	manager packages.Manager,
	loader config.Loader) *cli.Command {
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
			cli.StringFlag{
				Name:  "version, v",
				Usage: "the package version",
			},
		},
		Action: func(context *cli.Context) error {
			packagesPath := context.String("path")
			installPackageCommand, err := commands.NewInstall(platform, packagesPath, fileSystem, manager, loader)
			if err != nil {
				return err
			}

			packageName := context.Args().First()
			if strings.TrimSpace(packageName) == "" {
				return errors.New("missing required argument package name")
			}
			packageVersion := context.String("version")
			return installPackageCommand.Execute(packageName, packageVersion)
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

func createListProcessesCommand(
	fileSystem afero.Fs,
	console ui.Console,
	loader config.Loader) *cli.Command {

	listProcessesCommand := commands.NewListProcesses(
		console)

	return &cli.Command{
		Name:  "processes",
		Usage: "prints the list of processes for the given environment in the config file",
		Action: func(context *cli.Context) error {
			cfg, err := loadConfiguration(context, loader)
			if err != nil {
				return err
			}
			return listProcessesCommand.Execute(cfg)
		},
	}
}

func createStoresCommand(
	fileSystem afero.Fs,
	console ui.Console,
	loader config.Loader) *cli.Command {

	listStoresCommand := commands.NewStores(
		console)

	return &cli.Command{
		Name:    "stores",
		Aliases: []string{"s"},
		Usage:   "prints the list of stores in the config file",
		Action: func(context *cli.Context) error {
			cfg, err := loadConfiguration(context, loader)
			if err != nil {
				return err
			}
			return listStoresCommand.Execute(cfg)
		},
	}
}

func loadConfiguration(context *cli.Context, loader config.Loader) (*config.Config, error) {
	configFile := context.GlobalString("config")
	return loader.LoadConfig(configFile)
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
