package app

import (
	"fmt"

	"github.com/patrickhuber/go-di"
	"github.com/patrickhuber/wrangle/internal/global"
	"github.com/urfave/cli/v2"
)

func GetContainer(ctx *cli.Context) (di.Container, error) {
	if ctx == nil || ctx.App == nil || ctx.App.Metadata == nil {
		return nil, fmt.Errorf("application context, application or metadata is null")
	}
	return ctx.App.Metadata[global.MetadataDependencyInjection].(di.Container), nil
}
