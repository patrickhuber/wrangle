package commands

import (
	"fmt"
	"strings"

	"github.com/patrickhuber/go-di"
	"github.com/patrickhuber/wrangle/internal/app"
	"github.com/patrickhuber/wrangle/internal/config"
	"github.com/patrickhuber/wrangle/internal/config/edit"
	"github.com/urfave/cli/v2"
)

var StoreAdd = &cli.Command{
	Name:        "add",
	Action:      StoreAddAction,
	Description: "adds a store to configuration",
	Usage:       "add a store to configuration",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "type",
			Usage:    "store type",
			Required: true,
		},
		&cli.StringSliceFlag{
			Name:  "property",
			Usage: "store property in key=value form (repeatable)",
		},
		&cli.StringSliceFlag{
			Name:  "include",
			Usage: "environment variable names to include for env stores",
		},
		&cli.StringSliceFlag{
			Name:  "exclude",
			Usage: "environment variable names to exclude for env stores",
		},
		&cli.StringSliceFlag{
			Name:  "dependency",
			Usage: "name of a store this store depends on (repeatable)",
		},
		&cli.StringFlag{
			Name:  "scope",
			Usage: "configuration scope: nearest, local, user, or system (default: nearest)",
		},
	},
}

type StoreAddCommand struct {
	Edit    edit.Service    `inject:""`
	Options StoreAddOptions `options:""`
}

func (cmd *StoreAddCommand) Execute() error {
	properties, err := parseStoreProperties(cmd.Options.Properties)
	if err != nil {
		return err
	}
	if len(cmd.Options.Include) > 0 {
		if _, exists := properties["include"]; exists {
			return fmt.Errorf("store property \"include\" was specified more than once")
		}
		properties["include"] = cmd.Options.Include
	}
	if len(cmd.Options.Exclude) > 0 {
		if _, exists := properties["exclude"]; exists {
			return fmt.Errorf("store property \"exclude\" was specified more than once")
		}
		properties["exclude"] = cmd.Options.Exclude
	}

	return cmd.Edit.AddStore(&edit.StoreRequest{
		Store: config.Store{
			Name:         cmd.Options.Name,
			Type:         cmd.Options.Type,
			Properties:   properties,
			Dependencies: cmd.Options.Dependencies,
		},
		Scope: cmd.Options.Scope,
	})
}

func StoreAddAction(ctx *cli.Context) error {
	name := ctx.Args().First()
	if name == "" {
		return fmt.Errorf("store name is required")
	}

	resolver, err := app.GetResolver(ctx)
	if err != nil {
		return fmt.Errorf("invalid store add configuration. %w", err)
	}

	cmd := &StoreAddCommand{
		Options: StoreAddOptions{
			Name:         name,
			Type:         ctx.String("type"),
			Properties:   ctx.StringSlice("property"),
			Include:      ctx.StringSlice("include"),
			Exclude:      ctx.StringSlice("exclude"),
			Dependencies: ctx.StringSlice("dependency"),
			Scope:        ctx.String("scope"),
		},
	}
	if err := di.Inject(resolver, cmd); err != nil {
		return err
	}

	return cmd.Execute()
}

func parseStoreProperties(values []string) (map[string]any, error) {
	properties := make(map[string]any, len(values))
	for _, value := range values {
		key, propertyValue, ok := strings.Cut(value, "=")
		key = strings.TrimSpace(key)
		if !ok || key == "" {
			return nil, fmt.Errorf("invalid store property %q (expected key=value)", value)
		}
		if _, exists := properties[key]; exists {
			return nil, fmt.Errorf("store property %q was specified more than once", key)
		}
		properties[key] = propertyValue
	}
	return properties, nil
}
