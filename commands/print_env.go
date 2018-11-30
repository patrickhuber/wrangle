package commands

import (
	"strings"

	"github.com/patrickhuber/wrangle/services"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

// CreatePrintEnvCommand creates a print env command from the cli context
func CreatePrintEnvCommand(
	printService services.PrintService) *cli.Command {

	return &cli.Command{
		Name:  "print-env",
		Usage: "print command environemnt variables",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "name, n",
				Usage: "process named `NAME`",
			},
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
				Format:      format}
			return printService.Print(params)
		},
	}
}
