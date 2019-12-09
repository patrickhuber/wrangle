package commands

import (
	"github.com/patrickhuber/wrangle/config"
	"github.com/patrickhuber/wrangle/filesystem"
	"github.com/patrickhuber/wrangle/store"
	"github.com/urfave/cli"
)

// CreateListStoresCommand creates a stores cli command from the cli context
func CreateListStoresCommand(
	app *cli.App,
	storesService store.Service,
	fs filesystem.FileSystem) cli.Command {
	command := cli.Command{
		Name:    "stores",
		Aliases: []string{"s"},
		Usage:   "prints the list of stores in the config file",
		Action: func(context *cli.Context) error {

			configFile := context.GlobalString("config")

			configProvider := config.NewFsProvider(fs, configFile)

			cfg, err := configProvider.Get()
			if err != nil {
				return err
			}

			return storesService.List(cfg)
		},
	}
	setCommandCustomHelpTemplateWithGlobalOptions(app, &command)
	return command
}
