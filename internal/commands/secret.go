package commands

import (
	"fmt"

	"github.com/patrickhuber/go-cross/console"
	"github.com/patrickhuber/go-di"
	"github.com/patrickhuber/wrangle/internal/app"
	"github.com/patrickhuber/wrangle/internal/secret"
	"github.com/urfave/cli/v2"
)

// GetSecret is the get subcommand for secret operations
var GetSecret = &cli.Command{
	Name:               "get",
	Action:             GetSecretAction,
	Description:        "get the specified secret",
	Usage:              "get the specified secret",
	CustomHelpTemplate: CommandHelpTemplate,
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

type GetSecretCommand struct {
	Secret  secret.Service `inject:""`
	Console console.Console `inject:""`
	Options GetSecretOptions
}

type GetSecretOptions struct {
	Key   string
	Store string
}

func GetSecretAction(ctx *cli.Context) error {
	resolver, err := app.GetResolver(ctx)
	if err != nil {
		return err
	}
	cmd := &GetSecretCommand{
		Options: GetSecretOptions{
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

func (cmd *GetSecretCommand) Execute() error {
	err := validate(SetSecretOptionStoreName, cmd.Options.Store)
	if err != nil {
		return err
	}
	err = validate(SetSecretOptionKeyName, cmd.Options.Key)
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

// Secret is the parent command for secret operations
var Secret = &cli.Command{
	Name:               "secret",
	Description:        "manage secrets",
	Usage:              "manage secrets",
	Subcommands:        []*cli.Command{GetSecret, SetSecret},
	CustomHelpTemplate: CommandHelpTemplate,
}
