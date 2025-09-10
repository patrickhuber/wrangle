package commands

import (
	"fmt"

	"github.com/patrickhuber/go-di"
	"github.com/patrickhuber/wrangle/internal/app"
	"github.com/patrickhuber/wrangle/internal/diff"
	"github.com/patrickhuber/wrangle/internal/export"
	"github.com/urfave/cli/v2"
)

var Export = &cli.Command{
	Name:        "export",
	Action:      ExportAction,
	Description: "Exports the current environment variables to in the format of the specified shell",
	Usage:       "export aggregated environment variables to the format of the specified shell",
	CustomHelpTemplate: CommandHelpTemplate + `
ARGS:
   shell	(bash|powershell)
`,
	ArgsUsage: "<shell>",
}

type ExportCommand struct {
	Export  export.Service `inject:""`
	Diff    diff.Service   `inject:""`
	Options ExportOptions
}

type ExportOptions struct {
	Shell string
}

func ExportAction(ctx *cli.Context) error {

	resolver, err := app.GetResolver(ctx)
	if err != nil {
		return fmt.Errorf("invalid initialize command configuration. %w", err)
	}
	if ctx.Args().Len() < 1 {
		return fmt.Errorf("expected <shell> argument")
	}
	cmd := &ExportCommand{
		Options: ExportOptions{
			Shell: ctx.Args().Get(0),
		},
	}
	err = di.Inject(resolver, cmd)
	if err != nil {
		return err
	}
	return cmd.Execute()
}

func (cmd *ExportCommand) Execute() error {
	changes, err := cmd.Diff.Execute()
	if err != nil {
		return err
	}
	return cmd.Export.Execute(cmd.Options.Shell, changes)
}
