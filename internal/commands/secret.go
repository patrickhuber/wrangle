package commands

import "github.com/urfave/cli/v2"

// Secret is the parent command for secret operations
var Secret = &cli.Command{
	Name:        "secret",
	Description: "manage secrets",
	Usage:       "manage secrets",
	Subcommands: []*cli.Command{SecretGet, SecretSet},
}
