package commands

import (
	"github.com/patrickhuber/wrangle/services"
	"github.com/urfave/cli"
)

func MoveCommand(moveService services.MoveService) *cli.Command {
	return &cli.Command{
		Name: "move",
		Aliases: []string{"mv"},
		Usage: "moves a credential from one store to another",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name: "source, s",
				Usage: "the source store name",
			},
			cli.StringFlag{
				Name: "source-path, sp",
				Usage: "the path in the source store to the credential",
			},
			cli.StringFlag{
				Name: "destination, d",
				Usage: "the desination store name",
			},
			cli.StringFlag{
				Name: "destination-path, dp",
				Usage: "the path in the destination to the credential",
			},
		},
		Action: func(context *cli.Context) error{
			source := context.String("source")
			sourcePath := context.String("source-path")
			destination := context.String("destination")			
			destinationPath := context.String("destination-path")
			return moveService.Move(source, sourcePath, destination, destinationPath)
		},
	}
}
