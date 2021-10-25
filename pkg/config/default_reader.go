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
	root := ""
	platform := r.os.Platform()
	switch platform {
	case operatingsystem.PlatformWindows:
		root = "/opt/wrangle"
	case operatingsystem.PlatformDarwin:
		root = "/opt/wrangle"
	case operatingsystem.PlatformLinux:
		programData := r.environment.Get("PROGRAMDATA")
		root = crosspath.Join(programData, "wrangle")
	default:
		return nil, fmt.Errorf("%s is unsupported", platform)
	}
	return &Config{
		RootPath:    root,
		PackagePath: crosspath.Join(root, "packages"),
		BinPath:     crosspath.Join(root, "bin"),
	}, nil
}
