package commands_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v2"

	"github.com/patrickhuber/go-cross/platform"
	"github.com/patrickhuber/go-di"
	"github.com/patrickhuber/wrangle/internal/commands"
	"github.com/patrickhuber/wrangle/internal/config"
	"github.com/patrickhuber/wrangle/internal/global"

	"github.com/patrickhuber/wrangle/internal/host"
)

func TestBootstrap(t *testing.T) {

	s := host.NewTest(platform.Linux, nil, nil)
	container := s.Container()

	app := cli.NewApp()
	app.Name = "wrangle"
	app.Metadata = map[string]interface{}{
		global.MetadataDependencyInjection: container,
	}
	app.Commands = []*cli.Command{commands.Bootstrap}
	app.Before = func(ctx *cli.Context) error {
		cliContext := config.CliContext(ctx)
		di.RegisterInstance(container, cliContext)
		return nil
	}
	err := app.Run([]string{"wrangle", "bootstrap"})
	require.NoError(t, err)
}
