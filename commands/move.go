package commands

import (
	"github.com/patrickhuber/wrangle/services"
	"github.com/urfave/cli"
)

func CreateMoveCommand(
	app *cli.App,
	credentialServiceFactory services.CredentialServiceFactory) *cli.Command {
	command := &cli.Command{
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
			// need to pass this to the service
			config := context.GlobalString("config")			

			credentialService, err := credentialServiceFactory.Create(config)
			if err != nil{
				return err
			}
			return credentialService.Move(source, sourcePath, destination, destinationPath)
		},
	}

	setCommandCustomHelpTemplateWithGlobalOptions(app, command)	
	return command
}
