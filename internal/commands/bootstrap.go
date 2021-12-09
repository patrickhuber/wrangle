package commands

import (
	"github.com/patrickhuber/wrangle/internal/services"
	"github.com/patrickhuber/wrangle/internal/types"
	"github.com/patrickhuber/wrangle/pkg/config"
	"github.com/patrickhuber/wrangle/pkg/di"
	"github.com/patrickhuber/wrangle/pkg/filesystem"
	"github.com/patrickhuber/wrangle/pkg/global"
	"github.com/urfave/cli/v2"
)

func Bootstrap(ctx *cli.Context) error {
	resolver := ctx.App.Metadata[global.MetadataDependencyInjection].(di.Resolver)

	i := resolver.Resolve(types.InstallService).(services.Install)
	fs := resolver.Resolve(types.FileSystem).(filesystem.FileSystem)
	defaultReader := resolver.Resolve(types.ConfigReader).(config.Reader)

	cfg, err := defaultReader.Get()
	if err != nil {
		return err
	}
	bootstrap := services.NewBootstrap(i, fs, cfg)
	req := &services.BootstrapRequest{
		Force:            ctx.Bool("force"),
		GlobalConfigFile: ctx.String(global.FlagConfig),
	}
	return bootstrap.Execute(req)
}
