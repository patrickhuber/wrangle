package commands

import "github.com/urfave/cli/v2"

// Command defines a cli command execution context
type Command interface {
	Execute(ctx *cli.Context) error
}
