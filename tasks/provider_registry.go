package tasks

import (
	"fmt"
)

type providerRegistry struct {
	providers map[string]Provider
}

// ProviderRegistry defines a registry for task runners
type ProviderRegistry interface {
	Get(taskType string) (Provider, error)
	Register(runner Provider)
	Decode(task interface{}) (Task, error)
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

func (registry *providerRegistry) Decode(task interface{}) (Task, error) {
	m, ok := task.(map[interface{}]interface{})
	if !ok {
		return nil, fmt.Errorf("unable to cast task to string map of interface")
	}
	for providerTaskType, provider := range registry.providers {
		_, ok := m[providerTaskType]
		if !ok {
			continue
		}
		task, err := provider.Decode(task)
		if err != nil {
			return nil, err
		}
		return task, nil
	}
	return nil, fmt.Errorf("no provider is registered for task type")
}
