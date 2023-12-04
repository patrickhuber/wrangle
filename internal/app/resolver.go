package app

import (
	"fmt"

	"github.com/patrickhuber/go-di"
	"github.com/patrickhuber/wrangle/internal/global"
	"github.com/urfave/cli/v2"
)

func GetResolver(ctx *cli.Context) (di.Resolver, error) {
	if ctx == nil || ctx.App == nil || ctx.App.Metadata == nil {
		return nil, fmt.Errorf("Application Context, Application or Metadata is null")
	}
	return ctx.App.Metadata[global.MetadataDependencyInjection].(di.Resolver), nil
}
