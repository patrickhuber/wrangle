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

	metadataDependencyInjection, ok := ctx.App.Metadata[global.MetadataDependencyInjection]
	if !ok {
		return nil, fmt.Errorf("dependency injection metadata is missing")
	}

	container, ok := metadataDependencyInjection.(di.Container)
	if !ok {
		return nil, fmt.Errorf("dependency injection metadata is not a container")
	}
	return container, nil
}
