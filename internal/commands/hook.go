package commands

import (
	"fmt"
	"os"

	"github.com/patrickhuber/go-di"
	"github.com/patrickhuber/wrangle/internal/app"
	"github.com/patrickhuber/wrangle/internal/services"
	"github.com/urfave/cli/v2"
)

var Hook = &cli.Command{
	Name:    "hook",
	Aliases: []string{"h"},
	Action:  HookAction,
	Flags:   []cli.Flag{},
}

type HookCommand struct {
	Hook    services.Hook `inject:""`
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
	executable, err := os.Executable()
	if err != nil {
		return err
	}
	req := &services.HookRequest{
		Executable: executable,
		Shell:      cmd.Options.Shell,
	}
	return cmd.Hook.Execute(req)
}
