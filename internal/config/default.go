package config

import (
	"fmt"

	"github.com/patrickhuber/wrangle/pkg/config"
	"github.com/patrickhuber/wrangle/pkg/crosspath"
	"github.com/patrickhuber/wrangle/pkg/env"
	"github.com/patrickhuber/wrangle/pkg/operatingsystem"
)

func NewDefault(os operatingsystem.OS, environment env.Environment) (*config.Config, error) {
	root := "/opt/wrangle"
	platform := os.Platform()

	switch platform {
	case operatingsystem.PlatformWindows:
		programData := environment.Get("PROGRAMDATA")
		root = crosspath.Join(programData, "wrangle")
	case operatingsystem.PlatformDarwin:
	case operatingsystem.PlatformLinux:
		break
	default:
		return nil, fmt.Errorf("%s is unsupported", platform)
	}

	cfg := &config.Config{
		Paths: &config.Paths{
			Root:     root,
			Packages: crosspath.Join(root, "packages"),
			Bin:      crosspath.Join(root, "bin"),
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
