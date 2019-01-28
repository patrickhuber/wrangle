package commands

import (
	"github.com/patrickhuber/wrangle/services"
	"github.com/urfave/cli"
)

func CreateCopyCommand(credentialServiceFactory services.CredentialServiceFactory) *cli.Command {
	return &cli.Command{
		Name: "copy",
		Aliases: []string{"cp"},
		Usage: "copies a credential from one store to another",
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
			// need to pass this to the credential service
			config := context.GlobalString("config")

			credentialService, err := credentialServiceFactory.Create(config)
			if err != nil{
				return err
			}
			return credentialService.Copy(source, sourcePath, destination, destinationPath)
		},
	}
}
