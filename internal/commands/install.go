package commands

import (
	"fmt"

	"github.com/patrickhuber/wrangle/internal/services"
	"github.com/patrickhuber/wrangle/internal/types"
	"github.com/patrickhuber/wrangle/pkg/di"
	"github.com/patrickhuber/wrangle/pkg/global"

	"github.com/urfave/cli/v2"
)

func Install(ctx *cli.Context) error {

	pkg := ctx.Args().First()
	if len(pkg) == 0 {
		return fmt.Errorf("package name is required")
	}

	resolver := ctx.App.Metadata[global.MetadataDependencyInjection].(di.Resolver)
	service := resolver.Resolve(types.InstallService).(services.Install)
	request := &services.InstallRequest{
		GlobalConfigFile: ctx.String(global.FlagConfig),
		Package:          pkg,
	}
	return service.Execute(request)
}
