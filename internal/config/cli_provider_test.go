package config_test

import (
	"fmt"
	"testing"

	"github.com/patrickhuber/go-config"
	"github.com/patrickhuber/go-dataptr"
	wrangle_config "github.com/patrickhuber/wrangle/internal/config"
	"github.com/patrickhuber/wrangle/internal/global"
	"github.com/urfave/cli/v2"
)

func TestCliProvider(t *testing.T) {
	// create a fake cli context and use it to create a cli provdier
	app := cli.NewApp()
	app.Flags = append(app.Flags, &cli.StringFlag{
		Name: global.FlagUserConfig,
	})
	app.Action = func(ctx *cli.Context) error {
		provider := wrangle_config.NewCliProvider(ctx)
		value, err := provider.Get(&config.GetContext{})
		if err != nil {
			return err
		}
		if value == nil {
			return fmt.Errorf("expected value to not be nil")
		}
		localConfig, err := dataptr.GetAs[string]("/spec/env/"+global.EnvUserConfig, value)
		if err != nil {
			return err
		}
		if localConfig == "" {
			return fmt.Errorf("expected local config to have a value")
		}
		return nil
	}
	err := app.Run([]string{
		"myapp", "--" + global.FlagUserConfig, "test",
	})
	if err != nil {
		t.Fatal(err)
	}
}
