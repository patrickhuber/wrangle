package commands

import (
	"fmt"

	"github.com/patrickhuber/wrangle/pkg/config"
	"github.com/patrickhuber/wrangle/pkg/di"
	"github.com/patrickhuber/wrangle/pkg/env"
	"github.com/patrickhuber/wrangle/pkg/filesystem"
	"github.com/patrickhuber/wrangle/pkg/global"
	"github.com/patrickhuber/wrangle/pkg/operatingsystem"
	"github.com/urfave/cli/v2"
)

type BootstrapCommand struct {
	FileSystem      filesystem.FileSystem
	OperatingSystem operatingsystem.OS
	Environment     env.Environment
	Config          *config.Config
	Options         *BootstrapCommandOptions
}

type BootstrapCommandOptions struct {
	ApplicationName  string
	GlobalConfigFile string
	Force            bool
}

func Bootstrap(ctx *cli.Context) error {
	resolver := ctx.App.Metadata[global.MetadataDependencyInjection].(di.Resolver)
	cmd := CreateBootstrapCommand(resolver, ctx)
	return BootstrapInternal(cmd)
}

func CreateBootstrapCommand(resolver di.Resolver, ctx *cli.Context) *BootstrapCommand {
	cmd := &BootstrapCommand{
		FileSystem:      resolver.Resolve(global.FileSystem).(filesystem.FileSystem),
		OperatingSystem: resolver.Resolve(global.OperatingSystem).(operatingsystem.OS),
		Environment:     resolver.Resolve(global.Environment).(env.Environment),
		Options: &BootstrapCommandOptions{
			ApplicationName:  ctx.App.Name,
			Force:            ctx.Bool("force"),
			GlobalConfigFile: ctx.String(global.FlagConfig),
		},
	}

	return cmd
}

func BootstrapInternal(opt *BootstrapCommand) error {
	err := ValidateBootstrapOptions(opt)
	if err != nil {
		return err
	}

	err = BootstrapInternalCreateGlobalConfig(opt)
	if err != nil {
		return err
	}
	return BootstrapInternalInstallPackages(opt)
}

func ValidateBootstrapOptions(opt *BootstrapCommand) error {
	if opt == nil {
		return fmt.Errorf("BootstrapOptions must not be nil")
	}
	if opt.FileSystem == nil {
		return fmt.Errorf("BootstrapOptions.FileSystem must not be nil")
	}
	if opt.OperatingSystem == nil {
		return fmt.Errorf("BootstrapOptions.OperatingSystem must not be nil")
	}
	return nil
}

func BootstrapInternalCreateGlobalConfig(opt *BootstrapCommand) error {

	// if the file exists or force is disabled, return
	if !opt.Options.Force {
		return nil
	}

	// create the config provider from the global path option
	configProvider := config.NewFileProvider(opt.FileSystem, opt.Options.GlobalConfigFile)

	// set the global config
	return configProvider.Set(opt.Config)
}

func BootstrapInternalInstallPackages(opt *BootstrapCommand) error {
	packageList := []string{"wrangle", "shim"}
	for _, p := range packageList {
		err := InstallInternal(
			&InstallCommand{
				Options: &InstallOptions{
					Package: p,
				},
			})
		if err != nil {
			return err
		}
	}
	return nil
}

func BootstrapInternalAddBinPathToProfile(opt *BootstrapCommand) error {
	return nil
}
