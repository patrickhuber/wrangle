package commands

import (
	"fmt"

	"github.com/patrickhuber/go-di"
	"github.com/patrickhuber/wrangle/internal/app"
	"github.com/patrickhuber/wrangle/internal/services"
	"github.com/urfave/cli/v2"
)

func Bootstrap(ctx *cli.Context) error {
	if ctx == nil || ctx.App == nil || ctx.App.Metadata == nil {
		return fmt.Errorf("invalid bootstrap command configuration. Application Context, Application or Metadata is null")
	}
	resolver := app.GetResolver(ctx)
	bootstrap, err := di.Resolve[services.Bootstrap](resolver)
	if err != nil {
		return err
	}
	req := &services.BootstrapRequest{
		Force: ctx.Bool("force"),
	}
	return bootstrap.Execute(req)
}
