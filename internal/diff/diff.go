package diff

import (
	"github.com/patrickhuber/wrangle/internal/config"
	"github.com/patrickhuber/wrangle/internal/envdiff"
	"github.com/patrickhuber/wrangle/internal/global"
	"github.com/patrickhuber/wrangle/internal/interpolate"

	"github.com/patrickhuber/go-cross/env"
	"github.com/patrickhuber/go-cross/filepath"
	"github.com/patrickhuber/go-cross/os"
)

type Service interface {
	Execute() ([]envdiff.Change, error)
}

type diff struct {
	configuration config.Service
	interpolate   interpolate.Service
	os            os.OS
	path          filepath.Provider
	environment   env.Environment
}

func NewService(
	configuration config.Service,
	interpolate interpolate.Service,
	os os.OS,
	environment env.Environment,
	path filepath.Provider) Service {
	return &diff{
		configuration: configuration,
		interpolate:   interpolate,
		os:            os,
		path:          path,
		environment:   environment,
	}
}

func (e *diff) Execute() ([]envdiff.Change, error) {
	wd, err := e.os.WorkingDirectory()
	if err != nil {
		return nil, err
	}

	// configuration get uses the default configuration provider to load configurations
	// this also looks at the environment and working directory to determin if the config should change
	cfg, err := e.interpolate.Execute()
	if err != nil {
		return nil, err
	}

	vars := map[string]string{}
	for k, v := range cfg.Spec.Environment {
		vars[k] = v
	}

	vars[global.EnvLocalConfig] = wd

	// convert the current environment to a map
	previous := e.environment.Export()

	// revert the previous state
	d, ok := previous[global.EnvDiff]
	if ok {
		changes, err := envdiff.Decode(d)
		if err != nil {
			return nil, err
		}
		envdiff.Revert(previous, changes)
	}

	// apply the vars to the current
	for k, v := range vars {
		previous[k] = v
	}

	// compute the difference
	current := cleanEnv(e.environment.Export())

	changes := envdiff.Diff(current, previous)

	// save the diff
	diffStr, err := envdiff.Encode(changes)
	if err != nil {
		return nil, err
	}

	changes = append(changes, envdiff.Add{
		Key:   global.EnvDiff,
		Value: diffStr,
	})

	return changes, err
}

// cleanEnv removes wrangle keys from the map and returns the modified map
func cleanEnv(m map[string]string) map[string]string {
	delete(m, global.EnvDiff)
	delete(m, global.EnvLocalConfig)
	delete(m, global.EnvSystemConfig)
	return m
}
