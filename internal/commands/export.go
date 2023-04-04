package commands

import (
	"fmt"

	"github.com/patrickhuber/go-di"
	"github.com/patrickhuber/wrangle/internal/app"
	"github.com/patrickhuber/wrangle/internal/services"
	"github.com/urfave/cli/v2"
)

var Export = &cli.Command{
	Name:    "Export",
	Aliases: []string{"e"},
	Action:  ExportAction,
	Flags:   []cli.Flag{},
}

type ExportCommand struct {
	Export  services.Export `inject:""`
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
	if ctx.Args().Len() < 3 {
		return fmt.Errorf("expected <shell> argument")
	}
	cmd := &ExportCommand{
		Options: ExportOptions{
			Shell: ctx.Args().Get(2),
		},
	}
	err = di.Inject(resolver, cmd)
	if err != nil {
		return err
	}
	return cmd.Execute()
}

func (cmd *ExportCommand) Execute() error {
	req := &services.ExportRequest{
		Shell: cmd.Options.Shell,
	}
	return cmd.Export.Execute(req)
}
