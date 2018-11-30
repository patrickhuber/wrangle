package commands

import (
	"github.com/patrickhuber/wrangle/services"
	"github.com/urfave/cli"
)

// CreateEnvCommand creates an env command from the cli context
func CreateEnvCommand(envService services.EnvService) *cli.Command {
	return &cli.Command{
		Name:  "env",
		Usage: "prints values of all associated environment variables",
		Action: func(context *cli.Context) error {
			return envService.Execute()
		},
	}
}
