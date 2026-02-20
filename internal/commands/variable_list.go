package commands

import (
	"fmt"

	"github.com/patrickhuber/go-cross/console"
	"github.com/patrickhuber/go-di"
	"github.com/patrickhuber/wrangle/internal/app"
	"github.com/patrickhuber/wrangle/internal/config"
	"github.com/patrickhuber/wrangle/internal/structio"
	"github.com/urfave/cli/v2"
)

var VariableList = &cli.Command{
	Name: "list",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name: "output",
		},
	},
	Action:      VariableListAction,
	Description: "list available variables",
	Usage:       "list available variables",
}

type VariableListCommand struct {
	Configuration config.Service       `inject:""`
	Console       console.Console      `inject:""`
	Options       *VariableListOptions `options:""`
}

type VariableListOptions struct {
	Output string `flag:"output"`
}

func VariableListAction(ctx *cli.Context) error {
	resolver, err := app.GetResolver(ctx)
	if err != nil {
		return fmt.Errorf("invalid list variable command. %w", err)
	}
	cmd := &VariableListCommand{
		Options: &VariableListOptions{
			Output: ctx.String("output"),
		},
	}
	err = di.Inject(resolver, cmd)
	if err != nil {
		return err
	}
	return cmd.Execute()
}

func (cmd *VariableListCommand) Execute() error {
	cfg, err := cmd.Configuration.Get()
	if err != nil {
		return err
	}
	w := cmd.Console.Out()
	output := ""
	if cmd.Options != nil {
		output = cmd.Options.Output
	}
	writer := structio.NewWriter(w, output)
	return writer.Write(cfg.Spec.Variables)
}
