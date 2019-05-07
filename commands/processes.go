package commands

import (
	"github.com/spf13/afero"
	"github.com/patrickhuber/wrangle/config"
	"github.com/patrickhuber/wrangle/services"
	"github.com/urfave/cli"
)

// CreateListProcessesCommand creates cli command for listing processes from the cli context
func CreateListProcessesCommand(
	app *cli.App,
	processesService services.ProcessesService,
	fs afero.Fs) *cli.Command {

	command := &cli.Command{
		Name:  "processes",
		Usage: "prints the list of processes for the given environment in the config file",
		Action: func(context *cli.Context) error {
			configFile := context.GlobalString("config")			

			configProvider := config.NewFsProvider(fs, configFile)
			
			cfg, err := configProvider.Get()
			if err != nil{
				return err
			}
			return processesService.List(cfg)
		},
	}
	
	setCommandCustomHelpTemplateWithGlobalOptions(app, command)	
	return command
}
