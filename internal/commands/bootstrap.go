package commands

import (
	"fmt"

	"github.com/patrickhuber/wrangle/internal/services"
	"github.com/patrickhuber/wrangle/internal/types"
	"github.com/patrickhuber/wrangle/pkg/di"
	"github.com/patrickhuber/wrangle/pkg/global"
	"github.com/urfave/cli/v2"
)

func Bootstrap(ctx *cli.Context) error {
	if ctx == nil || ctx.App == nil || ctx.App.Metadata == nil {
		return fmt.Errorf("invalid bootstrap command configuration. Application Context, Application or Metadata is null")
	}
	resolver := ctx.App.Metadata[global.MetadataDependencyInjection].(di.Resolver)

	obj, err := resolver.Resolve(types.BootstrapService)
	if err != nil {
		return err
	}
	bootstrap := obj.(services.Bootstrap)
	req := &services.BootstrapRequest{
		Force:            ctx.Bool("force"),
		GlobalConfigFile: ctx.String(global.FlagConfig),
	}
	return bootstrap.Execute(req)
}
