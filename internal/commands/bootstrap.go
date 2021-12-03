package commands

import (
	"github.com/patrickhuber/wrangle/internal/services"
	"github.com/patrickhuber/wrangle/pkg/config"
	"github.com/patrickhuber/wrangle/pkg/di"
	"github.com/patrickhuber/wrangle/pkg/filesystem"
	"github.com/patrickhuber/wrangle/pkg/global"
	"github.com/urfave/cli/v2"
)

func Bootstrap(ctx *cli.Context) error {
	resolver := ctx.App.Metadata[global.MetadataDependencyInjection].(di.Resolver)
	i := resolver.Resolve(global.InstallService).(services.Install)
	fs := resolver.Resolve(global.FileSystem).(filesystem.FileSystem)
	defaultReader := resolver.Resolve(global.DefaultConfigReader).(config.Reader)
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
