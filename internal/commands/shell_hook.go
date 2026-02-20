package commands

import (
	"fmt"

	"github.com/patrickhuber/go-cross/console"
	"github.com/patrickhuber/go-di"
	"github.com/patrickhuber/wrangle/internal/app"
	"github.com/patrickhuber/wrangle/internal/hook"
	"github.com/urfave/cli/v2"
)

var ShellHook = &cli.Command{
	Name:        "hook",
	Action:      ShellHookAction,
	Description: "Generates the shell hook script for the specified shell",
	Usage:       "generate the shell hook script for the specified shell",
	Flags:       []cli.Flag{},
	CustomHelpTemplate: CommandHelpTemplate + `
ARGS:
   shell	(bash|powershell)
`,
	ArgsUsage: "<shell>",
}

type ShellHookCommand struct {
	Hook    hook.Service    `inject:""`
	Console console.Console `inject:""`
	Options ShellHookOptions
}

type ShellHookOptions struct {
	Shell string
}

func ShellHookAction(ctx *cli.Context) error {
	resolver, err := app.GetResolver(ctx)
	if err != nil {
		return fmt.Errorf("invalid initialize command configuration. %w", err)
	}
	if ctx.Args().Len() < 1 {
		return fmt.Errorf("expected <shell> argument")
	}
	cmd := &ShellHookCommand{
		Options: ShellHookOptions{
			Shell: ctx.Args().Get(0),
		},
	}
	err = di.Inject(resolver, cmd)
	if err != nil {
		return err
	}
	return cmd.Execute()
}

func (cmd *ShellHookCommand) Execute() error {
	executable, err := cmd.Console.Executable()
	if err != nil {
		return err
	}
	req := &hook.Request{
		Executable: executable,
		Shell:      cmd.Options.Shell,
	}
	return cmd.Hook.Execute(req)
}
