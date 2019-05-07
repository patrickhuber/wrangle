package commands

import (
	"github.com/spf13/afero"
	"github.com/patrickhuber/wrangle/services"
	"github.com/patrickhuber/wrangle/config"
	"github.com/urfave/cli"
)

// CreateStoresCommand creates a stores cli command from the cli context
func CreateStoresCommand(
	app *cli.App,
	storesService services.StoresService,
	fs afero.Fs) *cli.Command {
	command := &cli.Command{
		Name:    "stores",
		Aliases: []string{"s"},
		Usage:   "prints the list of stores in the config file",
		Action: func(context *cli.Context) error {			

			configFile := context.GlobalString("config")

			configProvider := config.NewFsProvider(fs, configFile)
			
			cfg, err := configProvider.Get()
			if err != nil{
				return err
			}

			return storesService.List(cfg)
		},
	}		
	setCommandCustomHelpTemplateWithGlobalOptions(app, command)	
	return command
}
