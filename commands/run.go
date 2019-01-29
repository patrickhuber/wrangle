package commands

import (	
	"github.com/patrickhuber/wrangle/services"
	"strings"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

// CreateRunCommand creates a run command from the cli context
func CreateRunCommand(
	app *cli.App,
	runService services.RunService) *cli.Command {		

	command := &cli.Command{
		Name:      "run",
		Aliases:   []string{"r"},
		Usage:     "run a command",
		ArgsUsage: "<process name> [arguments]",
		Action: func(context *cli.Context) error {
			configFile := context.GlobalString("config")

			processName := context.Args().First()
			if strings.TrimSpace(processName) == "" {
				return errors.New("process name argument is required")
			}

			additionalArguments := context.Args().Tail()

			params := services.NewProcessParams(processName, additionalArguments...)			

			return runService.Run(configFile, params)
		},
	}
	
	setCommandCustomHelpTemplateWithGlobalOptions(app, command)	
	return command
}