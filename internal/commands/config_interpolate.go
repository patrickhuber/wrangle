package commands

import (
	"fmt"

	"github.com/patrickhuber/go-cross/console"
	"github.com/patrickhuber/go-di"
	"github.com/patrickhuber/wrangle/internal/app"
	"github.com/patrickhuber/wrangle/internal/interpolate"
	"github.com/patrickhuber/wrangle/internal/structio"
	"github.com/urfave/cli/v2"
)

var ConfigInterpolate = &cli.Command{
	Name:        "interpolate",
	Aliases:     []string{"int"},
	Action:      ConfigInterpolateAction,
	Description: "Interpolate generates the aggregated configuration from all configurations",
	Usage:       "generate the aggregated configuration from all configurations",
}

type ConfigInterpolateCommand struct {
	Interpolate interpolate.Service `inject:""`
	Console     console.Console     `inject:""`
	Options     ConfigInterpolateOptions
}

type ConfigInterpolateOptions struct{}

func ConfigInterpolateAction(ctx *cli.Context) error {
	resolver, err := app.GetResolver(ctx)
	if err != nil {
		return fmt.Errorf("invalid list variable command. %w", err)
	}
	cmd := &ConfigInterpolateCommand{}
	err = di.Inject(resolver, cmd)
	if err != nil {
		return err
	}
	return cmd.Execute()
}

func (cmd *ConfigInterpolateCommand) Execute() error {
	cfg, err := cmd.Interpolate.Execute()
	if err != nil {
		return err
	}
	// TODO: switch on output type
	writer := structio.NewYamlWriter(cmd.Console.Out())
	return writer.Write(cfg)
}
