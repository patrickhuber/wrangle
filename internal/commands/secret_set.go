package commands

import (
	"fmt"
	"strings"

	"github.com/patrickhuber/go-di"
	"github.com/patrickhuber/wrangle/internal/app"
	"github.com/patrickhuber/wrangle/internal/secret"
	"github.com/urfave/cli/v2"
)

var SecretSet = &cli.Command{
	Name:        "set",
	Action:      SecretSetAction,
	Description: "set the specified secret",
	Usage:       "set the specified secret",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "key",
			Aliases: []string{"k"},
		},
		&cli.StringFlag{
			Name:    "value",
			Aliases: []string{"v"},
		},
		&cli.StringFlag{
			Name:    "store",
			Aliases: []string{"s"},
		},
	},
}

type SecretSetCommand struct {
	Secret  secret.Service `inject:""`
	Options SecretSetOptions
}

type SecretSetOptions struct {
	Key   string
	Value string
	Store string
}

const SecretSetOptionKeyName = "key"
const SecretSetOptionValueName = "value"
const SecretSetOptionStoreName = "store"

func SecretSetAction(ctx *cli.Context) error {
	resolver, err := app.GetResolver(ctx)
	if err != nil {
		return err
	}
	cmd := &SecretSetCommand{
		Options: SecretSetOptions{
			Key:   ctx.String(SecretSetOptionKeyName),
			Value: ctx.String(SecretSetOptionValueName),
			Store: ctx.String(SecretSetOptionStoreName),
		},
	}
	err = di.Inject(resolver, cmd)
	if err != nil {
		return err
	}

	return cmd.Execute()
}

func (cmd *SecretSetCommand) Execute() error {
	err := validate(SecretSetOptionStoreName, cmd.Options.Store)
	if err != nil {
		return err
	}
	err = validate(SecretSetOptionKeyName, cmd.Options.Key)
	if err != nil {
		return err
	}
	err = validate(SecretSetOptionValueName, cmd.Options.Value)
	if err != nil {
		return err
	}
	return cmd.Secret.Set(
		cmd.Options.Store,
		cmd.Options.Key,
		cmd.Options.Value)
}

func validate(key string, value string) error {
	if !isEmpty(value) {
		return nil
	}
	return fmt.Errorf("invalid value for flag %s", key)
}

func isEmpty(s string) bool {
	s = strings.TrimSpace(s)
	return strings.EqualFold(s, "")
}
