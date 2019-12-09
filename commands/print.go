package commands

import "github.com/urfave/cli"

// CreatePrintCommand creates the parent list command
func CreatePrintCommand(
	subcommands ...cli.Command) cli.Command {
	command := cli.Command{
		Name:        "print",
		Aliases:     []string{"p"},
		Subcommands: subcommands,
	}
	return command
}
