package environment

import (
	"fmt"
	"path"

	"github.com/patrickhuber/go-cross/env"
	"github.com/patrickhuber/wrangle/internal/stores"
)

const (
	includeProperty = "include"
	excludeProperty = "exclude"
)

type Factory struct {
	environment env.Environment
}

func NewFactory(environment env.Environment) stores.Factory {
	return Factory{environment: environment}
}

func (f Factory) Name() string {
	return "env"
}

func (f Factory) Create(properties map[string]any) (stores.Store, error) {
	include, includeSet, err := stringListProperty(properties, includeProperty)
	if err != nil {
		return nil, err
	}
	exclude, _, err := stringListProperty(properties, excludeProperty)
	if err != nil {
		return nil, err
	}

	if err := validatePatterns(includeProperty, include); err != nil {
		return nil, err
	}
	if err := validatePatterns(excludeProperty, exclude); err != nil {
		return nil, err
	}

	return &Store{
		environment: f.environment,
		include:     include,
		includeSet:  includeSet,
		exclude:     exclude,
	}, nil
}

func stringListProperty(properties map[string]any, name string) ([]string, bool, error) {
	value, exists := properties[name]
	if !exists {
		return nil, false, nil
	}

	switch list := value.(type) {
	case []string:
		return list, true, nil
	case []any:
		result := make([]string, len(list))
		for i, item := range list {
			pattern, ok := item.(string)
			if !ok {
				return nil, false, fmt.Errorf("store property '%s' item %d is not a string", name, i)
			}
			result[i] = pattern
		}
		return result, true, nil
	default:
		return nil, false, fmt.Errorf("store property '%s' must be a list of strings", name)
	}
}

func validatePatterns(name string, patterns []string) error {
	for _, pattern := range patterns {
		if _, err := path.Match(pattern, ""); err != nil {
			return fmt.Errorf("invalid glob in store property '%s': %w", name, err)
		}
	}
	return nil
}
