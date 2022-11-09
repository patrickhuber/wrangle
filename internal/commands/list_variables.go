package commands

import (
	"fmt"

	"github.com/patrickhuber/go-di"
	"github.com/patrickhuber/go-log"
	"github.com/patrickhuber/wrangle/internal/app"
	"github.com/patrickhuber/wrangle/pkg/console"
	"github.com/urfave/cli/v2"
)

var ListVariables = &cli.Command{
	Name:   "variables",
	Action: ListVariablesAction,
}

type ListVariablesCommand struct {
	Logger  log.Logger           `inject:""`
	Console console.Console      `inject:""`
	Options ListVariablesOptions `options:""`
}

type ListVariablesOptions struct {
	Output string `flag:"output"`
}

func ListVariablesAction(ctx *cli.Context) error {
	resolver, err := app.GetResolver(ctx)
	if err != nil {
		return fmt.Errorf("invalid list variable command. %w", err)
	}
	listVariablesCommand := &ListVariablesCommand{}
	err = di.Inject(resolver, listVariablesCommand)
	if err != nil {
		return err
	}
	return listVariablesCommand.Execute()
}

func (cmd *ListVariablesCommand) Execute() error {
	return nil
}
