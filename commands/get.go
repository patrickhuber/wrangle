package commands

import "github.com/urfave/cli"

func CreateGetCommand(
	subcommands ...cli.Command) cli.Command {
	command := cli.Command{
		Name:        "get",
		Subcommands: subcommands,
	}
	return command
}
