package commands

import (
	"fmt"

	"github.com/patrickhuber/wrangle/services"
	"github.com/patrickhuber/wrangle/store"
	"github.com/urfave/cli"
)

func CreateCopyCommand(
	app *cli.App,
	credentialServiceFactory services.CredentialServiceFactory) *cli.Command {
	command := &cli.Command{
		Name:      "copy",
		Aliases:   []string{"cp"},
		Usage:     "copies a credential from one store to another",
		ArgsUsage: "<source>:<key> <destination>:<key>",
		Action: func(context *cli.Context) error {
			source := context.Args().Get(0)
			if len(source) == 0 {
				return fmt.Errorf("missing source")
			}

			destination := context.Args().Get(1)
			if len(destination) == 0 {
				return fmt.Errorf("missing destination")
			}

			sourceStoreAndPath, err := store.ParsePath(source)
			if err != nil {
				return err
			}

			destinationStoreAndPath, err := store.ParsePath(destination)
			if err != nil {
				return err
			}

			// need to pass this to the credential service
			configPath := context.GlobalString("config")

			credentialService, err := credentialServiceFactory.Create(configPath)
			if err != nil {
				return err
			}
			return credentialService.Copy(
				sourceStoreAndPath.Store, sourceStoreAndPath.Path,
				destinationStoreAndPath.Store, destinationStoreAndPath.Path)
		},
	}

	setCommandCustomHelpTemplateWithGlobalOptions(app, command)
	return command
}
