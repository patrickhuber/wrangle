package commands

import (
	"fmt"
	"strings"

	"github.com/patrickhuber/go-di"
	"github.com/patrickhuber/wrangle/internal/app"
	"github.com/patrickhuber/wrangle/internal/services"
	"github.com/urfave/cli/v2"
)

var SetSerect = &cli.Command{
	Name:   "secret",
	Action: SetSecretAction,
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

type SetSecretCommand struct {
	Secret  services.Secret `inject:""`
	Options SetSecretOptions
}

type SetSecretOptions struct {
	Key   string
	Value string
	Store string
}

const SetSecretOptionKeyName = "key"
const SetSecretOptionValueName = "value"
const SetSecretOptionStoreName = "store"

func SetSecretAction(ctx *cli.Context) error {
	resolver, err := app.GetResolver(ctx)
	if err != nil {
		return err
	}
	cmd := &SetSecretCommand{
		Options: SetSecretOptions{
			Key:   ctx.String(SetSecretOptionKeyName),
			Value: ctx.String(SetSecretOptionValueName),
			Store: ctx.String(SetSecretOptionStoreName),
		},
	}
	err = di.Inject(resolver, cmd)
	if err != nil {
		return err
	}

	return cmd.Execute()
}

func (cmd *SetSecretCommand) Execute() error {
	err := validate(SetSecretOptionStoreName, cmd.Options.Store)
	if err != nil {
		return err
	}
	err = validate(SetSecretOptionKeyName, cmd.Options.Key)
	if err != nil {
		return err
	}
	err = validate(SetSecretOptionValueName, cmd.Options.Value)
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
