package commands

import (
	"fmt"

	"github.com/patrickhuber/wrangle/renderers/items"
	"github.com/patrickhuber/wrangle/services"
	"github.com/patrickhuber/wrangle/ui"
	"github.com/urfave/cli"
)

func CreateListCommand(app *cli.App,
	console ui.Console,
	credentialService services.CredentialService,
	renderFactory items.Factory) *cli.Command {
	command := &cli.Command{
		Name:    "list",
		Aliases: []string{"ls"},
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

			format := context.String("format")

			credentials, err := credentialService.List(storeName, path)
			if err != nil {
				return err
			}

			renderer, err := renderFactory.Create(format)
			if err != nil {
				return err
			}
			return renderer.RenderItems(credentials, console.Out())
		},
	}
	return command
}
