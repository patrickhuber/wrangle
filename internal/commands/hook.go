package commands

import (
	"fmt"

	"github.com/patrickhuber/go-di"
	"github.com/patrickhuber/go-xplat/console"
	"github.com/patrickhuber/wrangle/internal/app"
	"github.com/patrickhuber/wrangle/internal/services"
	"github.com/urfave/cli/v2"
)

var Hook = &cli.Command{
	Name:   "hook",
	Action: HookAction,
	Flags:  []cli.Flag{},
}

type HookCommand struct {
	Hook    services.Hook   `inject:""`
	Console console.Console `inject:""`
	Options HookOptions
}

type HookOptions struct {
	Shell string
}

func HookAction(ctx *cli.Context) error {
	resolver, err := app.GetResolver(ctx)
	if err != nil {
		return fmt.Errorf("invalid initialize command configuration. %w", err)
	}
	if ctx.Args().Len() < 3 {
		return fmt.Errorf("expected <shell> argument")
	}
	cmd := &HookCommand{
		Options: HookOptions{
			Shell: ctx.Args().Get(2),
		},
	}
	err = di.Inject(resolver, cmd)
	if err != nil {
		return err
	}
	return cmd.Execute()
}

func (cmd *HookCommand) Execute() error {
	executable, err := cmd.Console.Executable()
	if err != nil {
		return err
	}
	req := &services.HookRequest{
		Executable: executable,
		Shell:      cmd.Options.Shell,
	}
	return cmd.Hook.Execute(req)
}
