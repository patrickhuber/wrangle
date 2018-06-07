package commands

import "github.com/urfave/cli"

type Command interface {
	ExecuteCommand(c *cli.Context) error
}
