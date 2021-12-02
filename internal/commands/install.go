package commands

import (
	"fmt"

	"github.com/patrickhuber/wrangle/pkg/config"
	"github.com/patrickhuber/wrangle/pkg/di"
	"github.com/patrickhuber/wrangle/pkg/feed"
	"github.com/patrickhuber/wrangle/pkg/filesystem"
	"github.com/patrickhuber/wrangle/pkg/global"

	"github.com/urfave/cli/v2"
)

type InstallCommand struct {
	Options        *InstallOptions
	FileSystem     filesystem.FileSystem
	ServiceFactory feed.ServiceFactory
}

type InstallOptions struct {
	Package          string
	GlobalConfigFile string
}

func Install(ctx *cli.Context) error {

	pkg := ctx.Args().First()
	if len(pkg) == 0 {
		return fmt.Errorf("package name is required")
	}

	resolver := ctx.App.Metadata[global.MetadataDependencyInjection].(di.Resolver)

	return InstallInternal(
		&InstallCommand{
			FileSystem: resolver.Resolve(global.FileSystem).(filesystem.FileSystem),
			Options: &InstallOptions{
				Package: pkg,
			},
		})
}

func InstallInternal(cmd *InstallCommand) error {
	configProvider := config.NewFileProvider(cmd.FileSystem, cmd.Options.GlobalConfigFile)
	cfg, err := configProvider.Get()
	if err != nil {
		return err
	}

	for _, f := range cfg.Feeds {
		svc := cmd.ServiceFactory.Create(f)
		svc.List(&feed.ListRequest{
			
		})
	}

	return nil
}
