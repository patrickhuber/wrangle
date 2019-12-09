package commands

import (
	"github.com/patrickhuber/wrangle/initialize"
	"github.com/urfave/cli"
)

func CreateInitCommand(
	app *cli.App,
	initService initialize.Service,
) *cli.Command {
	command := &cli.Command{
		Name:  "init",
		Usage: "initialize the wrangle configuration",
		Action: func(context *cli.Context) error {
			return initService.Init(context.GlobalString("config"))
		},
	}

	setCommandCustomHelpTemplateWithGlobalOptions(app, command)
	return command
}
