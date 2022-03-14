package commands

import (
	"fmt"

	"github.com/patrickhuber/wrangle/internal/app"
	"github.com/patrickhuber/wrangle/internal/services"
	"github.com/patrickhuber/wrangle/internal/types"

	"github.com/urfave/cli/v2"
)

func Install(ctx *cli.Context) error {

	pkg := ctx.Args().First()
	if len(pkg) == 0 {
		return fmt.Errorf("package name is required")
	}

	resolver := app.GetResolver(ctx)
	o, err := resolver.Resolve(types.InstallService)
	if err != nil {
		return err
	}
	service := o.(services.Install)
	request := &services.InstallRequest{
		Package: pkg,
	}
	return service.Execute(request)
}
