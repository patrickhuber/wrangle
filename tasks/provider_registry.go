package tasks

import "fmt"

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
	taskRunner, ok := registry.providers[taskType]
	if !ok {
		return nil, fmt.Errorf("unable to locate task factory '%s'", taskType)
	}
	return taskRunner, nil
}
