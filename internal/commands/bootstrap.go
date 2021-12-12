package commands

import (
	"fmt"

	"github.com/patrickhuber/wrangle/internal/services"
	"github.com/patrickhuber/wrangle/internal/types"
	"github.com/patrickhuber/wrangle/pkg/config"
	"github.com/patrickhuber/wrangle/pkg/di"
	"github.com/patrickhuber/wrangle/pkg/filesystem"
	"github.com/patrickhuber/wrangle/pkg/global"
	"github.com/urfave/cli/v2"
)

func Bootstrap(ctx *cli.Context) error {
	if ctx == nil || ctx.App == nil || ctx.App.Metadata == nil {
		return fmt.Errorf("invalid bootstrap command configuration. Application Context, Application or Metadata is null")
	}
	resolver := ctx.App.Metadata[global.MetadataDependencyInjection].(di.Resolver)

	o, err := resolver.Resolve(types.InstallService)
	if err != nil {
		return err
	}
	i := o.(services.Install)

	o, err = resolver.Resolve(types.FileSystem)
	if err != nil {
		return err
	}
	fs := o.(filesystem.FileSystem)

	o, err = resolver.Resolve(types.ConfigReader)
	if err != nil {
		return err
	}
	defaultReader := o.(config.Reader)

	cfg, err := defaultReader.Get()
	if err != nil {
		return err
	}
	cfg.Feeds = []*config.Feed{
		{
			Name: "default",
			Type: "git",
			URI:  "git://github.com/patrickhuber/wrangle-packages.git",
		},
	}
	bootstrap := services.NewBootstrap(i, fs, cfg)
	req := &services.BootstrapRequest{
		Force:            ctx.Bool("force"),
		GlobalConfigFile: ctx.String(global.FlagConfig),
	}
	return bootstrap.Execute(req)
}
