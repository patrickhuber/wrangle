package commands

import (
	"github.com/patrickhuber/wrangle/credentials"
	"github.com/patrickhuber/wrangle/ui"
	"github.com/urfave/cli"
)

// CreateGetSecretCommand creates an get command from the cli context
func CreateGetSecretCommand(app *cli.App, credentialServiceFactory credentials.ServiceFactory, console ui.Console) *cli.Command {
	return &cli.Command{
		Name:  "secret",
		Usage: "gets a the secret from the store",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "store, s",
				Usage: "the store containing the credential",
			},
			cli.StringFlag{
				Name:  "path, p",
				Usage: "the path or key to the credential",
			},
		},
		Action: func(context *cli.Context) error {

			configFile := context.GlobalString("config")
			storeName := context.String("store")
			path := context.String("path")

			credentialService, err := credentialServiceFactory.Create(configFile)
			if err != nil {
				return err
			}

			data, err := credentialService.Get(storeName, path)
			if err != nil {
				return err
			}

			j, err := data.Json()
			if err != nil {
				return err
			}

			_, err = console.Out().Write(j)
			return err
		},
	}
}
