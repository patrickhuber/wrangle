package config

import (
	"fmt"

	"github.com/patrickhuber/go-xplat/env"
	"github.com/patrickhuber/go-xplat/filepath"
	"github.com/patrickhuber/go-xplat/os"
	"github.com/patrickhuber/go-xplat/platform"
	"github.com/patrickhuber/wrangle/pkg/config"
)

func NewDefault(os os.OS, environment env.Environment, path *filepath.Processor) (*config.Config, error) {
	root := "/opt/wrangle"
	plat := os.Platform()

	p := platform.Platform(plat)
	switch {
	case p.IsWindows():
		programData := environment.Get("PROGRAMDATA")
		root = path.Join(programData, "wrangle")
	case p.IsUnix():
		break
	default:
		return nil, fmt.Errorf("%s is unsupported", plat)
	}

	cfg := &config.Config{
		Paths: &config.Paths{
			Root:     root,
			Packages: path.Join(root, "packages"),
			Bin:      path.Join(root, "bin"),
		},
		References: []*config.Reference{
			{
				Name:    "wrangle",
				Version: "latest",
			},
			{
				Name:    "shim",
				Version: "latest",
			},
		},
		Feeds: []*config.Feed{
			{
				Name: "default",
				Type: "git",
				URI:  "https://github.com/patrickhuber/wrangle-packages",
			},
		},
	}
	return cfg, nil
}
