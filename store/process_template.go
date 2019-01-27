package store

import (
	"fmt"

	"github.com/patrickhuber/wrangle/templates"

	"github.com/patrickhuber/wrangle/config"
)

type processTemplate struct {
	cfg      *config.Config
	registry ResolverRegistry
}

// ProcessTemplate defines a template for processes
type ProcessTemplate interface {
	Evaluate(processName string) (*config.Process, error)
}

// NewProcessTemplate  creates a new process template with the given config and manager
func NewProcessTemplate(cfg *config.Config, manager Manager) (ProcessTemplate, error) {

	g, err := config.NewConfigurationGraph(cfg)
	if err != nil {
		return nil, err
	}

	registry, err := NewResolverRegistry(cfg, g, manager)
	if err != nil {
		return nil, err
	}

	return &processTemplate{
		registry: registry,
		cfg:      cfg,
	}, nil
}

func (t *processTemplate) Evaluate(processName string) (*config.Process, error) {
	for _, process := range t.cfg.Processes {
		if process.Name == processName {
			return t.evaluate(&process)
		}
	}
	return nil, fmt.Errorf("unable to find process '%s'", processName)
}

func (t *processTemplate) evaluate(process *config.Process) (*config.Process, error) {

	const argsKey = "args"
	const envKey = "env"

	if process.Stores == nil {
		return process, nil
	}

	document := map[string]interface{}{
		argsKey: process.Args,
		envKey:  process.Vars,
	}

	resolvers, err := t.registry.GetResolvers(process.Stores)
	if err != nil {
		return nil, err
	}

	template := templates.NewTemplate(document)
	resolved, err := template.Evaluate(resolvers...)
	if err != nil {
		return nil, err
	}

	resolvedMap := resolved.(map[string]interface{})
	args := resolvedMap[argsKey]
	env := resolvedMap[envKey]

	process.Args, err = normalizeToStringSlice(args)
	if err != nil {
		return nil, err
	}

	process.Vars, err = normalizeToMapStringOfString(env)
	if err != nil {
		return nil, err
	}
	return process, nil
}

func normalizeToMapStringOfString(document interface{}) (map[string]string, error) {
	switch t := document.(type) {
	case (map[string]interface{}):
		newMap := make(map[string]string)
		for k, v := range t {
			newMap[k] = fmt.Sprintf("%v", v)
		}
		return newMap, nil
	case (map[string]string):
		return t, nil
	}

	return nil, fmt.Errorf("Unable to normalize type '%T'", document)
}

func normalizeToStringSlice(document interface{}) ([]string, error) {
	switch t := document.(type) {
	case ([]string):
		return t, nil
	case ([]interface{}):
		slice := make([]string, len(t))
		for i, item := range t {
			slice[i] = fmt.Sprintf("%v", item)
		}
		return slice, nil
	}
	return nil, fmt.Errorf("Unable to normalize type '%T'", document)
}
