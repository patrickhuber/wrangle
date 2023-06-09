package config

import (
	"github.com/patrickhuber/go-xplat/env"
	"github.com/patrickhuber/go-xplat/filepath"
	"github.com/patrickhuber/go-xplat/os"
	"github.com/patrickhuber/wrangle/pkg/config"
)

func NewTest(os os.OS, environment env.Environment, path filepath.Processor) (*config.Config, error) {
	cfg, err := NewDefault(os, environment, path)
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
