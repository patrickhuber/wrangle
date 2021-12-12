package config

import (
	"fmt"

	"github.com/patrickhuber/wrangle/pkg/crosspath"
	"github.com/patrickhuber/wrangle/pkg/env"
	"github.com/patrickhuber/wrangle/pkg/operatingsystem"
)

type defaultReader struct {
	os          operatingsystem.OS
	environment env.Environment
}

func NewDefaultReader(os operatingsystem.OS, environment env.Environment) Reader {
	return &defaultReader{
		os:          os,
		environment: environment,
	}
}

func (r *defaultReader) Get() (*Config, error) {

	root := "/opt/wrangle"
	platform := r.os.Platform()

	switch platform {
	case operatingsystem.PlatformWindows:
		programData := r.environment.Get("PROGRAMDATA")
		root = crosspath.Join(programData, "wrangle")
	case operatingsystem.PlatformDarwin:
	case operatingsystem.PlatformLinux:
		break
	default:
		return nil, fmt.Errorf("%s is unsupported", platform)
	}

	return &Config{
		Paths: &Paths{
			Root:    root,
			Packages: crosspath.Join(root, "packages"),
			Bin:     crosspath.Join(root, "bin"),
		},
	}, nil
}
