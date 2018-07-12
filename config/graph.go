package config

import "fmt"

type node struct {
	name     string
	children map[string]*node
	parents  map[string]*node
}

type graph struct {
	nodes  map[string]*node
	stores map[string]*Store
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
	Store(name string) *Store
}

// NewConfigurationGraph creates a graph from the configuration
func NewConfigurationGraph(configuration *Config) (Graph, error) {

	nodes := make(map[string]*node)
	stores := make(map[string]*Store)

	// create the nodes
	for i := range configuration.Stores {

		// create the node and assign the config source
		store := configuration.Stores[i]
		n := newNode(store.Name)

		// create maps for easy lookup when resolving references
		stores[store.Name] = &store
		nodes[store.Name] = n
	}

	// loop over the config sources now that we have a full
	// node list to use when building links
	for i := range configuration.Stores {
		store := configuration.Stores[i]
		n := nodes[store.Name]
		for _, parentName := range store.Stores {
			parent := nodes[parentName]
			parent.children[n.name] = n
			n.parents[parentName] = parent
		}
	}

	g := &graph{
		nodes:  nodes,
		stores: stores,
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

func (g *graph) Store(name string) *Store {
	return g.stores[name]
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
