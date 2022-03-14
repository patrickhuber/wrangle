package config

import (
	"github.com/patrickhuber/wrangle/pkg/config"
	"github.com/patrickhuber/wrangle/pkg/env"
	"github.com/patrickhuber/wrangle/pkg/operatingsystem"
)

func NewTest(os operatingsystem.OS, environment env.Environment) (*config.Config, error) {
	cfg, err := NewDefault(os, environment)
	if err != nil {
		return nil, err
	}
	cfg.Feeds = []*config.Feed{
		{
			Name: "default",
			Type: "memory",
		},
	}
	return cfg, nil
}
