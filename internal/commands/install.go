package commands

import (
	"fmt"

	"github.com/patrickhuber/go-di"
	"github.com/patrickhuber/wrangle/internal/app"
	"github.com/patrickhuber/wrangle/internal/services"

	"github.com/urfave/cli/v2"
)

func Install(ctx *cli.Context) error {

	pkg := ctx.Args().First()
	if len(pkg) == 0 {
		return fmt.Errorf("package name is required")
	}

	resolver := app.GetResolver(ctx)
	service, err := di.Resolve[services.Install](resolver)
	if err != nil {
		return err
	}
	request := &services.InstallRequest{
		Package: pkg,
	}
	return service.Execute(request)
}
