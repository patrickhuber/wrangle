package commands

import (
	"strings"

	"github.com/patrickhuber/wrangle/services"

	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

// CreatePrintCommand Creates a Print Command from the cli context
func CreatePrintCommand(
	app *cli.App,
	printService services.PrintService) *cli.Command {

	command := &cli.Command{
		Name:      "print",
		Aliases:   []string{"p"},
		Usage:     "prints the process as it would be executed",
		ArgsUsage: "<process name> [arguments]",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "format, f",
				Usage: "Print for with the given format (bash|powershell)",
			},
		},
		Action: func(context *cli.Context) error {
			processName := context.Args().First()
			if strings.TrimSpace(processName) == "" {
				return errors.New("process name argument is required")
			}

			format := context.String("format")

			configFile := context.GlobalString("config")

			params := &services.PrintParams{
				ConfigFile:  configFile,
				ProcessName: processName,
				Format:      format,
				Include: services.PrintParamsInclude{
					ProcessAndArgs: true,
				},
			}

			return printService.Print(params)
		},
	}
	
	setCommandCustomHelpTemplateWithGlobalOptions(app, command)	
	return command
}
