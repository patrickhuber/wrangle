package app

import (
	"github.com/patrickhuber/go-di"
	"github.com/patrickhuber/wrangle/pkg/global"
	"github.com/urfave/cli/v2"
)

func GetResolver(ctx *cli.Context) di.Resolver {
	return ctx.App.Metadata[global.MetadataDependencyInjection].(di.Resolver)
}
