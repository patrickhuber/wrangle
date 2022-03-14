package commands

import (
	"github.com/patrickhuber/go-di"
	"github.com/patrickhuber/wrangle/internal/app"
	"github.com/patrickhuber/wrangle/pkg/console"
	"github.com/patrickhuber/wrangle/pkg/ilog"
	"github.com/urfave/cli/v2"
)

type ListVariablesCommand struct {
	Logger  ilog.Logger     `inject:""`
	Console console.Console `inject:""`
	Options ListVariablesOptions
}

type ListVariablesOptions struct {
	Output string
}

func ListVariables(ctx *cli.Context) error {
	resolver := app.GetResolver(ctx)
	listVariablesCommand := &ListVariablesCommand{}
	err := di.Inject(resolver, listVariablesCommand)
	if err != nil {
		return err
	}
	return ListVariablesInternal(listVariablesCommand)
}

func ListVariablesInternal(cmd *ListVariablesCommand) error {
	return nil
}
