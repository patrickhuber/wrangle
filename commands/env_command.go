package commands

import "github.com/urfave/cli"

type EnvCommand struct {
}

func NewEnvCommand() *EnvCommand {
	return &EnvCommand{}
}

func (cmd *EnvCommand) ExecuteCommand(c *cli.Context) error {
	return nil
}
