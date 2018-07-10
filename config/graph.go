package config

import "fmt"

type node struct {
	name     string
	children map[string]*node
	parents  map[string]*node
}

type graph struct {
	nodes   map[string]*node
	sources map[string]*ConfigSource
}

// Node represents a node in the graph
type Node interface {
	Children() []Node
	Child(name string) Node
	Parents() []Node
	Parent(name string) Node
	Name() string
}

// Graph represents a configuration graph
type Graph interface {
	Nodes() []Node
	Node(name string) Node
	Source(name string) *ConfigSource
}

// NewConfigurationGraph creates a graph from the configuration
func NewConfigurationGraph(configuration *Config) (Graph, error) {

	nodes := make(map[string]*node)
	sources := make(map[string]*ConfigSource)

	// create the nodes
	for i := range configuration.ConfigSources {

		// create the node and assign the config source
		source := configuration.ConfigSources[i]
		n := newNode(source.Name)

		// create maps for easy lookup when resolving references
		sources[source.Name] = &source
		nodes[source.Name] = n
	}

	// loop over the config sources now that we have a full
	// node list to use when building links
	for i := range configuration.ConfigSources {
		source := configuration.ConfigSources[i]
		n := nodes[source.Name]
		for _, parentName := range source.Configurations {
			parent := nodes[parentName]
			parent.children[n.name] = n
			n.parents[parentName] = parent
		}
	}

	g := &graph{
		nodes:   nodes,
		sources: sources,
	}

	if g.isCyclic() {
		return nil, fmt.Errorf("the configuration graph contains cycles. Review the configuration sources and elimintate all reference loops")
	}

	return g, nil
}

func (g *graph) isCyclic() bool {
	visited := make(map[string]bool)
	stack := make(map[string]bool)

	for name := range g.nodes {
		if g.isCyclicRecursive(name, visited, stack) {
			return true
		}
	}
	return false
}

func (g *graph) isCyclicRecursive(name string, visited map[string]bool, stack map[string]bool) bool {
	if visited, ok := stack[name]; visited && ok {
		return true
	}
	if visited, ok := visited[name]; visited && ok {
		return false
	}
	visited[name] = true
	stack[name] = true

	n := g.nodes[name]
	for key := range n.children {
		if g.isCyclicRecursive(key, visited, stack) {
			return true
		}
	}
	stack[name] = false
	return false
}

func (g *graph) Nodes() []Node {
	var nodes = make([]Node, 0)
	for _, value := range g.nodes {
		nodes = append(nodes, value)
	}
	return nodes
}

func (g *graph) Node(name string) Node {
	return g.nodes[name]
}

func (g *graph) Source(name string) *ConfigSource {
	return g.sources[name]
}

func newNode(name string) *node {
	return &node{
		name:     name,
		parents:  make(map[string]*node, 0),
		children: make(map[string]*node, 0),
	}
}

func (n *node) Children() []Node {
	var nodes = make([]Node, 0)
	for _, value := range n.children {
		nodes = append(nodes, value)
	}
	return nodes
}

func (n *node) Child(name string) Node {
	return n.children[name]
}

func (n *node) Parents() []Node {
	var nodes = make([]Node, 0)
	for _, value := range n.parents {
		nodes = append(nodes, value)
	}
	return nodes
}

func (n *node) Parent(name string) Node {
	return n.parents[name]
}

func (n *node) Name() string {
	return n.name
}
