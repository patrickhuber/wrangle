package commands

import (
	"fmt"

	"github.com/patrickhuber/wrangle/credentials"
	"github.com/patrickhuber/wrangle/renderers/items"
	"github.com/patrickhuber/wrangle/ui"
	"github.com/urfave/cli"
)

// CreateListSecretsCommand command that lists secrets for the given store
func CreateListSecretsCommand(app *cli.App,
	console ui.Console,
	credentialServiceFactory credentials.ServiceFactory,
	renderFactory items.Factory) cli.Command {
	command := cli.Command{
		Name: "secrets",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "store, s",
				Usage: "the store where the credentials will be fetched",
			},
			cli.StringFlag{
				Name:  "path, p",
				Usage: "the path or key to fetch the credentials",
			},
			cli.StringFlag{
				Name:  "format, f",
				Usage: "the format to return the list. possible values: json, yaml, table, tree",
			},
		},
		Action: func(context cli.Context) error {
			// if path is not specified, just list all the credentials
			path := context.String("path")

			storeName := context.String("store")
			if storeName == "" {
				return fmt.Errorf("missing required flag 'store'")
			}

			configPath := context.GlobalString("config")
			credentialService, err := credentialServiceFactory.Create(configPath)
			if err != nil {
				return err
			}

			credentials, err := credentialService.List(storeName, path)
			if err != nil {
				return err
			}

			format := context.String("format")
			renderer, err := renderFactory.Create(format)
			if err != nil {
				return err
			}
			return renderer.RenderItems(credentials, console.Out())
		},
	}
	return command
}
