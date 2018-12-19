package tasks

import (
	"fmt"

	yaml "gopkg.in/yaml.v2"
)

type providerRegistry struct {
	providers map[string]Provider
}

// ProviderRegistry defines a registry for task runners
type ProviderRegistry interface {
	Get(taskType string) (Provider, error)
	Register(runner Provider)
}

// NewProviderRegistry creates a new provider registry instance
func NewProviderRegistry() ProviderRegistry {
	return &providerRegistry{
		providers: make(map[string]Provider),
	}
}

func (registry *providerRegistry) Register(provider Provider) {
	registry.providers[provider.TaskType()] = provider
}

func (registry *providerRegistry) Get(taskType string) (Provider, error) {
	provider, ok := registry.providers[taskType]
	if !ok {
		return nil, fmt.Errorf("unable to locate task factory '%s'", taskType)
	}
	return provider, nil
}

func (registry *providerRegistry) Parse(data string) (Task, error) {
	m := make(map[interface{}]interface{})
	err := yaml.Unmarshal([]byte(data), m)
	if err != nil {
		return nil, err
	}
	for key := range m {
		stringKey := key.(string)
		provider, err := registry.Get(stringKey)
		if err != nil {
			return nil, err
		}
		return provider.Unmarshal(data)
	}
	return nil, fmt.Errorf("unable to parse task")
}
