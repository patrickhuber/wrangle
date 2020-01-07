package commands

import (
	"github.com/patrickhuber/wrangle/config"
	"github.com/patrickhuber/wrangle/processes"

	"strings"

	"github.com/patrickhuber/wrangle/filesystem"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

// CreateRunCommand creates a run command from the cli context
func CreateRunCommand(
	app *cli.App,
	runService processes.RunService,
	fs filesystem.FileSystem) *cli.Command {

	command := &cli.Command{
		Name:      "run",
		Aliases:   []string{"r"},
		Usage:     "run a command",
		ArgsUsage: "<process name> [arguments]",
		Action: func(context *cli.Context) error {

			processName := context.Args().First()
			if strings.TrimSpace(processName) == "" {
				return errors.New("process name argument is required")
			}

			additionalArguments := context.Args().Tail()
			configFile := context.GlobalString("config")

			configProvider := config.NewFsProvider(fs, configFile)

			cfg, err := configProvider.Get()
			if err != nil {
				return err
			}

			params := processes.NewProcessParams(processName, cfg, additionalArguments...)

			return runService.Run(params)
		},
	}

	setCommandCustomHelpTemplateWithGlobalOptions(app, command)
	return command
}
