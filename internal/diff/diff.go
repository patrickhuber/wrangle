package diff

import (
	"fmt"

	"github.com/patrickhuber/wrangle/internal/config"
	"github.com/patrickhuber/wrangle/internal/envdiff"
	"github.com/patrickhuber/wrangle/internal/global"
	"github.com/patrickhuber/wrangle/internal/stores"
	"github.com/patrickhuber/wrangle/internal/template"

	"github.com/patrickhuber/go-cross/env"
	"github.com/patrickhuber/go-cross/filepath"
	"github.com/patrickhuber/go-cross/os"
)

type Service interface {
	Execute() ([]envdiff.Change, error)
}

type diff struct {
	configuration config.Configuration
	store         stores.Service
	os            os.OS
	path          filepath.Provider
	environment   env.Environment
}

func NewService(
	configuration config.Configuration,
	store stores.Service,
	os os.OS,
	environment env.Environment,
	path filepath.Provider) Service {
	return &diff{
		configuration: configuration,
		store:         store,
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
	cfg, err := e.configuration.Get()
	if err != nil {
		return nil, err
	}

	vars, err := e.getVariableValues(cfg)
	if err != nil {
		return nil, err
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

func (e *diff) getVariableValues(cfg config.Config) (map[string]string, error) {
	// create variable providers for each store
	variableProviders, err := e.createVariableProviders(cfg)
	if err != nil {
		return nil, err
	}

	// create the template options from the variable providers
	var options []template.Option
	for _, vp := range variableProviders {
		options = append(options, template.WithProvider(vp))
	}

	vars := map[string]string{}
	var unresolved []string
	for k, v := range cfg.Spec.Environment {

		if !template.HasVariables(v) {
			vars[k] = v
			continue
		}

		// set v as a template and extract any vars
		t := template.New(v, options...)
		result, err := t.Evaluate()
		if err != nil {
			return nil, err
		}
		if len(result.Unresolved) > 0 {
			unresolved = append(unresolved, result.Unresolved...)
			continue
		}
		vars[k] = fmt.Sprintf("%v", result.Value)
	}
	if len(unresolved) > 0 {
		return nil, fmt.Errorf("unable to resolve the following variables %v", unresolved)
	}
	return vars, nil
}

// cleanEnv removes wrangle keys from the map and returns the modified map
func cleanEnv(m map[string]string) map[string]string {
	delete(m, global.EnvDiff)
	delete(m, global.EnvLocalConfig)
	delete(m, global.EnvSystemConfig)
	return m
}

func (e diff) createVariableProviders(cfg config.Config) ([]template.VariableProvider, error) {
	var variableProviders []template.VariableProvider

	// the registry is responsible for finding the factory to create the store
	for _, store := range cfg.Spec.Stores {
		s, err := e.store.Get(store.Name)
		if err != nil {
			return nil, err
		}
		variableProviders = append(variableProviders, storeToProvider{store: s})
	}
	return variableProviders, nil
}

type storeToProvider struct {
	store stores.Store
}

// List implements template.VariableProvider.
func (stp storeToProvider) List() ([]string, error) {
	result, err := stp.store.List()
	if err != nil {
		return nil, err
	}
	var names []string
	for _, l := range result {
		names = append(names, l.Data.Name)
	}
	return names, nil
}

// Get implements template.VariableProvider.
func (stp storeToProvider) Get(key string) (any, bool, error) {
	k, err := stores.ParseKey(key)
	if err != nil {
		return nil, false, err
	}
	return stp.store.Get(k)
}
