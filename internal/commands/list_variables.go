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

var ListVariables = &cli.Command{
	Name: "variables",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name: "output",
		},
	},
	Action:      ListVariablesAction,
	Description: "list available variables",
	Usage:       "list available variables",
}

type ListVariablesCommand struct {
	Configuration config.Service        `inject:""`
	Console       console.Console       `inject:""`
	Options       *ListVariablesOptions `options:""`
}

type ListVariablesOptions struct {
	Output string `flag:"output"`
}

func ListVariablesAction(ctx *cli.Context) error {
	resolver, err := app.GetResolver(ctx)
	if err != nil {
		return fmt.Errorf("invalid list variable command. %w", err)
	}
	listVariablesCommand := &ListVariablesCommand{
		Options: &ListVariablesOptions{
			Output: ctx.String("output"),
		},
	}
	err = di.Inject(resolver, listVariablesCommand)
	if err != nil {
		return err
	}
	return listVariablesCommand.Execute()
}

func (cmd *ListVariablesCommand) Execute() error {
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
