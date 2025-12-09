package interpolate

import (
	"encoding/json"
	"fmt"

	"github.com/go-viper/mapstructure/v2"
	"github.com/patrickhuber/wrangle/internal/config"
	"github.com/patrickhuber/wrangle/internal/stores"
	"github.com/patrickhuber/wrangle/internal/template"
)

type Service interface {
	Execute() (config.Config, error)
}

type service struct {
	configuration config.Service
	stores        stores.Service
}

func NewService(configuration config.Service, stores stores.Service) Service {
	return &service{
		configuration: configuration,
		stores:        stores,
	}
}

func (i *service) Execute() (config.Config, error) {
	cfg, err := i.configuration.Get()
	if err != nil {
		return config.Config{}, err
	}

	layers, providers, err := i.createVariableProviders(cfg)
	if err != nil {
		return config.Config{}, err
	}

	if len(layers) == 0 {
		return cfg, nil
	}

	// marshal spec to a map so the template evaluator (map/slice/string aware) can walk all fields
	specMap := map[string]any{}
	{
		bytes, err := json.Marshal(cfg.Spec)
		if err != nil {
			return config.Config{}, err
		}
		if err := json.Unmarshal(bytes, &specMap); err != nil {
			return config.Config{}, err
		}
	}

	var lastResult *template.EvaluationResult

	var unresolved []string
	var accumulated []template.VariableProvider

	for _, layer := range layers {
		for _, name := range layer {
			accumulated = append(accumulated, providers[name])
		}

		var options []template.Option
		for _, vp := range accumulated {
			options = append(options, template.WithProvider(vp))
		}

		result, err := template.New(specMap, options...).Evaluate()
		if err != nil {
			return config.Config{}, err
		}

		lastResult = result
		unresolved = result.Unresolved

		valueMap, ok := result.Value.(map[string]any)
		if !ok {
			return config.Config{}, fmt.Errorf("unexpected evaluation result type %T", result.Value)
		}
		specMap = valueMap

		if len(unresolved) == 0 {
			break
		}
	}

	if len(unresolved) > 0 {
		return config.Config{}, fmt.Errorf("unable to resolve the following variables %v", unresolved)
	}

	var spec config.Spec
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{Result: &spec})
	if err != nil {
		return config.Config{}, err
	}
	if err := decoder.Decode(lastResult.Value); err != nil {
		return config.Config{}, err
	}

	cfg.Spec = spec
	return cfg, nil
}

func (i *service) createVariableProviders(cfg config.Config) ([][]string, map[string]template.VariableProvider, error) {
	layers, err := topologicalLayers(cfg.Spec.Stores)
	if err != nil {
		return nil, nil, err
	}

	providers := map[string]template.VariableProvider{}
	for _, store := range cfg.Spec.Stores {
		s, err := i.stores.Get(store.Name)
		if err != nil {
			return nil, nil, err
		}
		providers[store.Name] = storeToProvider{store: s}
	}

	return layers, providers, nil
}

type storeToProvider struct {
	store stores.Store
}

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

func (stp storeToProvider) Get(key string) (any, bool, error) {
	k, err := stores.ParseKey(key)
	if err != nil {
		return nil, false, err
	}
	return stp.store.Get(k)
}

func topologicalLayers(storesCfg []config.Store) ([][]string, error) {
	graph := map[string][]string{}
	indegree := map[string]int{}

	for _, s := range storesCfg {
		graph[s.Name] = append(graph[s.Name], s.Dependencies...)
		if _, ok := indegree[s.Name]; !ok {
			indegree[s.Name] = 0
		}
		for _, dep := range s.Dependencies {
			indegree[s.Name] = indegree[s.Name] + 1
			if _, ok := indegree[dep]; !ok {
				indegree[dep] = 0
			}
		}
	}

	// detect missing dependencies
	for _, s := range storesCfg {
		for _, dep := range s.Dependencies {
			if _, ok := graph[dep]; !ok {
				return nil, fmt.Errorf("store '%s' depends on unknown store '%s'", s.Name, dep)
			}
		}
	}

	var layers [][]string
	processed := 0
	for {
		var layer []string
		for name, degree := range indegree {
			if degree == 0 {
				layer = append(layer, name)
			}
		}
		if len(layer) == 0 {
			break
		}

		layers = append(layers, layer)
		for _, name := range layer {
			// mark as processed by setting indegree to -1
			indegree[name] = -1
			for storeName, deps := range graph {
				for _, dep := range deps {
					if dep == name {
						indegree[storeName] = indegree[storeName] - 1
					}
				}
			}
			processed++
		}
	}

	if processed != len(indegree) {
		return nil, fmt.Errorf("store dependency cycle detected")
	}

	return layers, nil
}
