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
	test        bool
}

func NewDefaultReaderWithTestMode(os operatingsystem.OS, environment env.Environment) Reader {
	return &defaultReader{
		os:          os,
		environment: environment,
		test:        true,
	}
}
func NewDefaultReader(os operatingsystem.OS, environment env.Environment) Reader {
	return &defaultReader{
		os:          os,
		environment: environment,
		test:        false,
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

	cfg := &Config{
		Paths: &Paths{
			Root:     root,
			Packages: crosspath.Join(root, "packages"),
			Bin:      crosspath.Join(root, "bin"),
		},
	}
	if r.test {
		cfg.Feeds = []*Feed{
			{
				Name: "default",
				Type: "memory",
			},
		}
	} else {
		cfg.Feeds = []*Feed{
			{
				Name: "default",
				Type: "git",
				URI:  "git://github.com/patrickhuber/wrangle-packages.git",
			},
		}
	}
	return cfg, nil
}
