package commands

import (
	"github.com/patrickhuber/wrangle/services"
	"github.com/urfave/cli"
)

// CreateStoresCommand creates a stores cli command from the cli context
func CreateStoresCommand(
	app *cli.App,
	storesService services.StoresService) *cli.Command {
	command := &cli.Command{
		Name:    "stores",
		Aliases: []string{"s"},
		Usage:   "prints the list of stores in the config file",
		Action: func(context *cli.Context) error {			
			configFile := context.GlobalString("config")
			return storesService.List(configFile)
		},
	}		
	setCommandCustomHelpTemplateWithGlobalOptions(app, command)	
	return command
}
