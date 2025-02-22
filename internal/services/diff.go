package services

import (
	"fmt"
	"strings"

	"github.com/patrickhuber/wrangle/internal/config"
	"github.com/patrickhuber/wrangle/internal/envdiff"
	"github.com/patrickhuber/wrangle/internal/global"
	"github.com/patrickhuber/wrangle/internal/stores"
	"github.com/patrickhuber/wrangle/internal/template"

	"github.com/patrickhuber/go-cross/env"
	"github.com/patrickhuber/go-cross/filepath"
	"github.com/patrickhuber/go-cross/os"
)

type Diff interface {
	Execute() ([]envdiff.Change, error)
}

type diff struct {
	configuration Configuration
	store         Store
	environment   env.Environment
	os            os.OS
	path          filepath.Provider
}

func NewDiff(
	configuration Configuration,
	store Store,
	environment env.Environment,
	os os.OS,
	path filepath.Provider) Diff {
	return &diff{
		configuration: configuration,
		store:         store,
		environment:   environment,
		os:            os,
		path:          path,
	}
}

func (e *diff) Execute() ([]envdiff.Change, error) {
	wd, err := e.os.WorkingDirectory()
	if err != nil {
		return nil, err
	}

	// is the local config the same as the working directory?
	localConfig, localConfigSet := e.environment.Lookup(global.EnvLocalConfig)
	if localConfigSet {
		noAction, err := e.shouldTakeNoAction(wd, localConfig)
		if err != nil {
			return nil, err
		}
		if noAction {
			return nil, nil
		}
	}

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

func (e *diff) shouldTakeNoAction(workingDirectory string, localConfig string) (bool, error) {
	if localConfig == workingDirectory {
		// do nothing
		return true, nil
	}
	// is the working directory a sub directory of the last configuration file?
	filePaths, err := e.configuration.LocalConfigurationFiles()
	if err != nil {
		return false, err
	}
	if len(filePaths) > 0 {
		lastFilePath := filePaths[len(filePaths)-1]
		inSubOfLoadedConfig, err := e.isSelfOrSubDirectory(lastFilePath, workingDirectory)
		if err != nil {
			return false, err
		}
		if inSubOfLoadedConfig {
			return true, nil
		}
	}
	return false, nil
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

func (e diff) isSelfOrSubDirectory(base string, rel string) (bool, error) {
	if base == rel {
		return true, nil
	}

	// are we in a sub directory?
	rel, err := e.path.Rel(base, rel)
	if err != nil {
		return false, err
	}

	// this is not a sub directory becuase it contains ".."
	return !strings.Contains(rel, ".."), nil
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
