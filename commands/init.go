package commands

import (
	"github.com/patrickhuber/wrangle/services"
	"github.com/urfave/cli"
)

func CreateInitCommand(
	initService services.InitService,
) *cli.Command {
	return &cli.Command{
		Name:  "init",
		Usage: "initialize the wrangle configuration",
		Action: func(context *cli.Context) error {
			return initService.Init(context.GlobalString("config"))
		},
	}
}
