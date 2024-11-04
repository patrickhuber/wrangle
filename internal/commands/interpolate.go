package commands

import (
	"fmt"

	"github.com/patrickhuber/go-cross/console"
	"github.com/patrickhuber/go-di"
	"github.com/patrickhuber/wrangle/internal/app"
	"github.com/patrickhuber/wrangle/internal/services"
	"github.com/patrickhuber/wrangle/internal/structio"
	"github.com/urfave/cli/v2"
)

var Interpolate = &cli.Command{
	Name:    "interpolate",
	Aliases: []string{"int"},
	Action:  InterpolateAction,
}

type InterpolateCommand struct {
	Interpolate services.Interpolate `inject:""`
	Console     console.Console      `inject:""`
	Options     InterpolateOptions
}

type InterpolateOptions struct{}

func InterpolateAction(ctx *cli.Context) error {
	resolver, err := app.GetResolver(ctx)
	if err != nil {
		return fmt.Errorf("invalid list variable command. %w", err)
	}
	interpolateCommand := &InterpolateCommand{}
	err = di.Inject(resolver, interpolateCommand)
	if err != nil {
		return err
	}
	return interpolateCommand.Execute()
}

func (cmd *InterpolateCommand) Execute() error {
	cfg, err := cmd.Interpolate.Execute()
	if err != nil {
		return err
	}
	// TODO: switch on output type
	writer := structio.NewYamlWriter(cmd.Console.Out())
	return writer.Write(cfg)
}
