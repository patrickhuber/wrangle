package commands

import (
	"github.com/patrickhuber/go-di"
	"github.com/urfave/cli/v2"
)

type Factory interface {
	Resolver() di.Resolver
}

type factory struct {
	resolver di.Resolver
}

func (f *factory) Resolver() di.Resolver {
	return f.resolver
}

func NewFactory(ctx *cli.Context) Factory {
	return &factory{}
}

func Create[T any](factory Factory) (*T, error) {
	var instance T
	err := di.Inject(factory.Resolver(), instance)
	if err != nil {
		return nil, err
	}
	return &instance, nil
}
