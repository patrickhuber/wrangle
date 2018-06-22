package store

import (
	"fmt"

	"github.com/patrickhuber/cli-mgr/config"
	"github.com/patrickhuber/cli-mgr/templates"
)

type pipeline struct {
	manager       Manager
	configuration *config.Config
}

// Pipeline resolves an environment configuration against the pipeline configuration
type Pipeline interface {
	Run(processName string, environmentName string) (*config.Environment, error)
}

// NewPipeline creates a new pipeline with the given manager and configuration
func NewPipeline(manager Manager, configuration *config.Config) Pipeline {
	return &pipeline{
		manager:       manager,
		configuration: configuration,
	}
}

func (p *pipeline) Run(processName string, environmentName string) (*config.Environment, error) {
	for _, process := range p.configuration.Processes {
		if process.Name == processName {
			for j := range process.Environments {
				environment := &process.Environments[j]
				if environment.Name == environmentName {
					return p.run(environment)
				}
			}
		}
	}
	return nil, fmt.Errorf("Unable to find environment '%s' for process '%s'", processName, environmentName)
}

func (p *pipeline) run(environment *config.Environment) (*config.Environment, error) {
	if environment.Config == "" {
		return environment, nil
	}
	resolvers, err := p.createResolvers(environment.Config)
	if err != nil {
		return nil, err
	}
	document := map[string]interface{}{
		"args": environment.Args,
		"vars": environment.Vars,
	}
	template := templates.NewTemplate(document)
	resolved, err := template.Evaluate(resolvers...)
	if err != nil {
		return nil, err
	}
	resolvedMap := resolved.(map[string]interface{})
	args := resolvedMap["args"]
	vars := resolvedMap["vars"]

	environment.Args, err = normalizeToStringSlice(args)
	if err != nil {
		return nil, err
	}

	environment.Vars, err = normalizeToMapStringOfString(vars)
	if err != nil {
		return nil, err
	}
	return environment, nil
}

func (p *pipeline) createResolvers(configName string) ([]templates.VariableResolver, error) {
	// create an ordered list of config sources based on the dependency order
	orderedConfigSources := p.createOrderedListOfConfigSources(configName)

	// traverse the list of config sources in reverse order creating resolvers
	// and using those resolvers on any configuration params
	resolvers := make([]templates.VariableResolver, 0)
	for i := len(orderedConfigSources) - 1; i >= 0; i-- {
		resolver, err := p.createResolver(orderedConfigSources[i], resolvers)
		if err != nil {
			return nil, err
		}
		// prepend the resolver to the list of resolvers because we are traveling backwards
		// through the config sources
		resolvers = append([]templates.VariableResolver{resolver}, resolvers...)
	}
	return resolvers, nil
}

func (p *pipeline) createOrderedListOfConfigSources(configName string) []*config.ConfigSource {
	configSourceMap := p.createMapOfConfigSources()

	current := configName
	orderedConfigSources := make([]*config.ConfigSource, 0)
	for current != "" {
		configSource := configSourceMap[current]
		orderedConfigSources = append(orderedConfigSources, configSource)
		current = configSource.Config
	}
	return orderedConfigSources
}

func (p *pipeline) createMapOfConfigSources() map[string]*config.ConfigSource {
	configSourceMap := make(map[string]*config.ConfigSource)
	for i, configSource := range p.configuration.ConfigSources {
		configSourceMap[configSource.Name] = &p.configuration.ConfigSources[i]
	}
	return configSourceMap
}

func (p *pipeline) resolveConfigSourceParameters(configSource *config.ConfigSource, resolvers []templates.VariableResolver) (*config.ConfigSource, error) {
	shouldResolveConfigSourceParameters := configSource.Config != ""
	if shouldResolveConfigSourceParameters {
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
	}
	return configSource, nil
}

func normalizeToMapStringOfString(document interface{}) (map[string]string, error) {
	switch t := document.(type) {
	case (map[string]interface{}):
		newMap := make(map[string]string)
		for k, v := range t {
			newMap[k] = v.(string)
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
		for _, item := range t {
			slice = append(slice, item.(string))
		}
		return slice, nil
	}
	return nil, fmt.Errorf("Unable to normalize type '%T'", document)
}

func (p *pipeline) createResolver(
	configSource *config.ConfigSource,
	existingResolvers []templates.VariableResolver) (templates.VariableResolver, error) {
	configSource, err := p.resolveConfigSourceParameters(configSource, existingResolvers)
	if err != nil {
		return nil, err
	}

	// create the config store using the updated configSource
	configStore, err := p.manager.Create(configSource)
	if err != nil {
		return nil, err
	}

	return NewStoreVariableResolver(configStore), nil
}
