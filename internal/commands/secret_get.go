package commands

import (
	"fmt"

	"github.com/patrickhuber/go-cross/console"
	"github.com/patrickhuber/go-di"
	"github.com/patrickhuber/wrangle/internal/app"
	"github.com/patrickhuber/wrangle/internal/secret"
	"github.com/urfave/cli/v2"
)

// SecretGet is the get subcommand for secret operations
var SecretGet = &cli.Command{
	Name:        "get",
	Action:      SecretGetAction,
	Description: "get the specified secret",
	Usage:       "get the specified secret",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "key",
			Aliases: []string{"k"},
		},
		&cli.StringFlag{
			Name:    "store",
			Aliases: []string{"s"},
		},
	},
}

type SecretGetCommand struct {
	Secret  secret.Service  `inject:""`
	Console console.Console `inject:""`
	Options SecretGetOptions
}

type SecretGetOptions struct {
	Key   string
	Store string
}

func SecretGetAction(ctx *cli.Context) error {
	resolver, err := app.GetResolver(ctx)
	if err != nil {
		return err
	}
	cmd := &SecretGetCommand{
		Options: SecretGetOptions{
			Key:   ctx.String("key"),
			Store: ctx.String("store"),
		},
	}
	err = di.Inject(resolver, cmd)
	if err != nil {
		return err
	}
	return cmd.Execute()
}

func (cmd *SecretGetCommand) Execute() error {
	err := validate(SecretSetOptionStoreName, cmd.Options.Store)
	if err != nil {
		return err
	}
	err = validate(SecretSetOptionKeyName, cmd.Options.Key)
	if err != nil {
		return err
	}
	value, err := cmd.Secret.Get(cmd.Options.Store, cmd.Options.Key)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(cmd.Console.Out(), value)
	return err
}
