package commands

import (
	"github.com/urfave/cli"
)

// CreateListCommand creates the parent list command
func CreateListCommand(
	subcommands ...cli.Command) cli.Command {
	command := cli.Command{
		Name:        "list",
		Aliases:     []string{"ls"},
		Subcommands: subcommands,
	}
	return command
}
