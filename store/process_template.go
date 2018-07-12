package store

import (
	"fmt"

	"github.com/patrickhuber/wrangle/collections"
	"github.com/patrickhuber/wrangle/config"
	"github.com/patrickhuber/wrangle/templates"
)

type processTemplate struct {
	graph     config.Graph
	manager   Manager
	cfg       *config.Config
	resolvers map[string]templates.VariableResolver
}

type ProcessTemplate interface {
	Evaluate(environmentName string, processName string) (*config.Process, error)
}

func NewProcessTemplate(cfg *config.Config, manager Manager) (ProcessTemplate, error) {
	g, err := config.NewConfigurationGraph(cfg)
	if err != nil {
		return nil, err
	}
	return &processTemplate{
		graph:     g,
		manager:   manager,
		cfg:       cfg,
		resolvers: make(map[string]templates.VariableResolver),
	}, nil
}

func (t *processTemplate) Evaluate(environmentName string, processName string) (*config.Process, error) {
	for _, environment := range t.cfg.Environments {
		if environment.Name == environmentName {
			for _, process := range environment.Processes {
				if process.Name == processName {
					return t.evaluate(&process)
				}
			}
		}
	}
	return nil, fmt.Errorf("unable to find process '%s' in environment '%s'", processName, environmentName)
}

func (t *processTemplate) evaluate(process *config.Process) (*config.Process, error) {
	// create the resolvers needed to evaluate just this process template
	err := t.populateResolvers(process.Stores)
	if err != nil {
		return nil, err
	}
	resolvers := t.getResolvers(process.Stores)

	const argsKey = "args"
	const envKey = "env"

	document := map[string]interface{}{
		argsKey: process.Args,
		envKey:  process.Vars,
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

func (t *processTemplate) populateResolvers(configurations []string) error {
	queue := collections.NewQueue()
	for _, configuration := range configurations {
		node := t.graph.Node(configuration)
		queue.Enqueue(node)
	}
	for !queue.Empty() {

		// dequeue the node, skip processing if its resolver exists
		node := queue.Dequeue().(config.Node)
		if _, ok := t.resolvers[node.Name()]; ok {
			continue
		}

		// check if any parents of the node have not been processed
		var reprocess = false
		for _, parent := range node.Parents() {
			name := parent.Name()
			if _, ok := t.resolvers[name]; !ok {
				queue.Enqueue(parent)
				reprocess = true
			}
		}

		// the node must be reprocessed
		if reprocess {
			queue.Enqueue(node)
			continue
		}

		configSource := t.graph.Store(node.Name())
		resolver, err := t.createResolver(configSource)
		if err != nil {
			return err
		}
		t.resolvers[node.Name()] = resolver
	}
	return nil
}

func (t *processTemplate) createResolver(
	configSource *config.Store) (templates.VariableResolver, error) {
	configSource, err := t.resolveConfigSourceParameters(configSource)
	if err != nil {
		return nil, err
	}

	// create the config store using the updated configSource
	configStore, err := t.manager.Create(configSource)
	if err != nil {
		return nil, err
	}

	return NewStoreVariableResolver(configStore), nil
}

func (t *processTemplate) resolveConfigSourceParameters(configSource *config.Store) (*config.Store, error) {
	shouldResolveConfigSourceParameters := configSource.Stores != nil && len(configSource.Stores) > 0
	if !shouldResolveConfigSourceParameters {
		return configSource, nil
	}

	resolvers := t.getResolvers(configSource.Stores)
	// create a template and use the template to resolve the params with the current in-order resolvers
	template := templates.NewTemplate(configSource.Params)
	document, err := template.Evaluate(resolvers...)
	if err != nil {
		return nil, err
	}
	params, err := normalizeToMapStringOfString(document)
	if err != nil {
		return nil, err
	}
	configSource.Params = params

	return configSource, nil
}

func (t *processTemplate) getResolvers(configurations []string) []templates.VariableResolver {
	resolvers := make([]templates.VariableResolver, 0)
	for _, configuration := range configurations {
		resolver := t.resolvers[configuration]
		resolvers = append(resolvers, resolver)
	}
	return resolvers
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
