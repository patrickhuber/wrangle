package commands

import (
	"fmt"

	"github.com/patrickhuber/wrangle/credentials"
	"github.com/patrickhuber/wrangle/store"

	"github.com/urfave/cli"
)

func CreateMoveCommand(
	app *cli.App,
	credentialServiceFactory credentials.ServiceFactory) *cli.Command {
	command := &cli.Command{
		Name:      "move",
		Aliases:   []string{"mv"},
		Usage:     "moves a credential from one store to another",
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

			configFile := context.GlobalString("config")
			credentialService, err := credentialServiceFactory.Create(configFile)
			if err != nil {
				return err
			}
			return credentialService.Move(
				sourceStoreAndPath.Store, sourceStoreAndPath.Path,
				destinationStoreAndPath.Store, destinationStoreAndPath.Path)
		},
	}

	setCommandCustomHelpTemplateWithGlobalOptions(app, command)
	return command
}
