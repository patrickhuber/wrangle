package services

import (
	"fmt"
	"strings"

	"github.com/patrickhuber/go-iter"
	"github.com/patrickhuber/wrangle/internal/config"
	"github.com/patrickhuber/wrangle/internal/envdiff"
	"github.com/patrickhuber/wrangle/internal/global"
	"github.com/patrickhuber/wrangle/internal/stores"
	"github.com/patrickhuber/wrangle/internal/template"

	"github.com/patrickhuber/go-shellhook"
	"github.com/patrickhuber/go-xplat/console"
	"github.com/patrickhuber/go-xplat/env"
	"github.com/patrickhuber/go-xplat/filepath"
	"github.com/patrickhuber/go-xplat/os"
)

type ExportRequest struct {
	Shell string
}

type Export interface {
	Execute(r *ExportRequest) error
}

type export struct {
	shells        map[string]shellhook.Shell
	console       console.Console
	configuration Configuration
	registry      stores.Registry
	environment   env.Environment
	os            os.OS
	path          *filepath.Processor
}

func NewExport(
	shells map[string]shellhook.Shell,
	console console.Console,
	configuration Configuration,
	registry stores.Registry,
	environment env.Environment,
	os os.OS,
	path *filepath.Processor) Export {
	return &export{
		shells:        shells,
		console:       console,
		configuration: configuration,
		registry:      registry,
		environment:   environment,
		os:            os,
		path:          path,
	}
}

func (e *export) shouldTakeNoAction(workingDirectory string, localConfig string) (bool, error) {
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

func (e *export) getVariableValues(cfg config.Config) (map[string]string, error) {
	// create variable providers for each store
	variableProviders, err := e.createVariableProviders(cfg)
	if err != nil {
		return nil, err
	}

	// create the template options from the variable providers
	vp := iter.FromSlice(variableProviders)
	optionIter := iter.Select(vp, template.WithProvider)
	options := iter.ToSlice(optionIter)

	vars := map[string]string{}
	for k, v := range cfg.Spec.Environment {

		if !template.HasVariables(v) {
			vars[k] = v
			continue
		}

		// set v as a template and extract any vars
		t := template.New(v, options...)
		value, err := t.Evaluate()
		if err != nil {
			return nil, err
		}
		vars[k] = fmt.Sprintf("%v", value)
	}
	return vars, nil
}

func (e *export) Execute(r *ExportRequest) error {
	shell, ok := e.shells[r.Shell]
	if !ok {
		return fmt.Errorf("invalid shell '%s'", r.Shell)
	}

	wd, err := e.os.WorkingDirectory()
	if err != nil {
		return err
	}

	// is the local config the same as the working directory?
	localConfig, localConfigSet := e.environment.Lookup(global.EnvLocalConfig)
	if localConfigSet {
		noAction, err := e.shouldTakeNoAction(wd, localConfig)
		if err != nil {
			return err
		}
		if noAction {
			return nil
		}
	}

	cfg, err := e.configuration.Get()
	if err != nil {
		return err
	}

	vars, err := e.getVariableValues(cfg)
	if err != nil {
		return err
	}

	vars[global.EnvLocalConfig] = wd

	// convert the current environment to a map
	previous := e.environment.Export()

	// revert the previous state
	d, ok := previous[global.EnvDiff]
	if ok {
		changes, err := envdiff.Decode(d)
		if err != nil {
			return err
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

	// apply difference
	out := e.console.Out()
	for _, change := range changes {
		switch c := change.(type) {
		case envdiff.Add:
			str := shell.Export(c.Key, c.Value)
			_, err = fmt.Fprintln(out, str)
		case envdiff.Remove:
			str := shell.Unset(c.Key)
			_, err = fmt.Fprintln(out, str)
		case envdiff.Update:
			str := shell.Export(c.Key, c.Value)
			_, err = fmt.Fprintln(out, str)
		}
		if err != nil {
			return err
		}
	}

	// save the diff
	diffStr, err := envdiff.Encode(changes)
	if err != nil {
		return err
	}
	exportStr := shell.Export(global.EnvDiff, diffStr)
	_, err = fmt.Fprintln(out, exportStr)
	return err
}

// cleanEnv removes wrangle keys from the map and returns the modified map
func cleanEnv(m map[string]string) map[string]string {
	delete(m, global.EnvDiff)
	delete(m, global.EnvLocalConfig)
	delete(m, global.EnvConfig)
	return m
}

func (e export) isSelfOrSubDirectory(base string, rel string) (bool, error) {
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

func (e export) createVariableProviders(cfg config.Config) ([]template.VariableProvider, error) {
	var variableProviders []template.VariableProvider
	// the registry is responsible for finding the factory to create the store
	for _, store := range cfg.Spec.Stores {
		factory, err := e.registry.Get(store.Type)
		if err != nil {
			return nil, err
		}

		s, err := factory.Create(store.Properties)
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
	result, ok, err := stp.store.Get(k)
	return result, ok, err
}
