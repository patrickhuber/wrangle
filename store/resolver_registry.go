package store

import (
	"github.com/patrickhuber/wrangle/collections"
	"github.com/patrickhuber/wrangle/config"
	"github.com/patrickhuber/wrangle/templates"
)

// ResolverRegistry defines a registry for resolvers
type ResolverRegistry interface {
	GetResolvers(stores []string) ([]templates.VariableResolver, error)
}

type resolverRegistry struct {
	resolvers map[string]templates.VariableResolver	
}

func (reg *resolverRegistry) GetResolvers(stores []string) ([]templates.VariableResolver, error) {
	if stores == nil || len(stores) == 0 {
		return []templates.VariableResolver{}, nil
	}
	resolvers := make([]templates.VariableResolver, 0)
	for _, configuration := range stores {
		resolver := reg.resolvers[configuration]
		resolvers = append(resolvers, resolver)
	}
	return resolvers, nil
}

// NewResolverRegistry creates a new resolver registry
func NewResolverRegistry(cfg *config.Config, graph config.Graph, manager Manager) (ResolverRegistry, error) {
	queue := collections.NewQueue()

	for _, s := range cfg.Stores {
		node := graph.Node(s.Name)
		queue.Enqueue(node)
	}

	reg := &resolverRegistry{
		resolvers: map[string]templates.VariableResolver{},		
	}

	for !queue.Empty() {

		node := queue.Dequeue().(config.Node)

		// dequeue the node, skip processing if it already exists
		if _, ok := reg.resolvers[node.Name()]; ok {
			continue
		}

		// check if any parents ofthe node have not been processed
		var reprocess = false
		for _, parent := range node.Parents() {
			name := parent.Name()
			if _, ok := reg.resolvers[name]; !ok {
				queue.Enqueue(parent)
				reprocess = true
			}
		}

		// the node must be reprocesed
		if reprocess {
			queue.Enqueue(node)
			continue
		}

		// now that all dependencies are loaded, go ahead and reprocess the node
		configSource := graph.Store(node.Name())
		resolver, err := reg.createResolver(configSource, manager)
		if err != nil {
			return nil, err
		}

		// store the node so it isn't processed again
		reg.resolvers[node.Name()] = resolver
	}
	return reg, nil
}

func (reg *resolverRegistry) createResolver(configSource *config.Store, manager Manager) (templates.VariableResolver, error) {
	configSource, err := reg.resolveConfigSourceParameters(configSource)
	if err != nil {
		return nil, err
	}

	// create the config store using the updated configSource
	configStore, err := manager.Create(configSource)
	if err != nil {
		return nil, err
	}

	return NewStoreVariableResolver(configStore), nil
}

func (reg *resolverRegistry) resolveConfigSourceParameters(configSource *config.Store) (*config.Store, error) {
	shouldResolveConfigSourceParameters := configSource.Stores != nil && len(configSource.Stores) > 0
	if !shouldResolveConfigSourceParameters {
		return configSource, nil
	}

	resolvers, err := reg.GetResolvers(configSource.Stores)
	if err != nil {
		return nil, err
	}

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
