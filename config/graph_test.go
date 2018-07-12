package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfigurationGraphCanLoadLinear(t *testing.T) {
	r := require.New(t)
	data := `
stores:
- name: head
  type: file
  stores:
  params:
    path: /test/head.yml
- name: middle
  type: file
  stores:
  - head
  params:
    path: /test/middle.yml
- name: tail
  type: file
  stores:
  - middle
  params:
    path: /test/tail.yml
`
	cfg, err := SerializeString(data)
	r.Nil(err)
	graph, err := NewConfigurationGraph(cfg)
	r.Nil(err)
	r.NotNil(graph)

	head := graph.Node("head")
	r.NotNil(head)

	middle := graph.Node("middle")
	r.NotNil(middle)

	tail := graph.Node("tail")
	r.NotNil(tail)

	r.Equal(1, len(head.Children()))
	r.Equal(0, len(head.Parents()))
	r.Equal(1, len(middle.Children()))
	r.Equal(1, len(middle.Parents()))
	r.Equal(0, len(tail.Children()))
	r.Equal(1, len(tail.Parents()))

	r.NotNil(graph.Store("tail"))
	r.NotNil(graph.Store("head"))
	r.NotNil(graph.Store("middle"))
}

func TestConfigurationGraphCanLoadTree(t *testing.T) {
	r := require.New(t)
	data := `
stores:
- name: root
  type: file
  stores:
  params:
    path: /test/root.yml
- name: left-child
  type: file
  stores:
  - root
  params:
    path: /test/left-child.yml
- name: right-child
  type: file
  stores:
  - root
  params:
    path: /test/right-child.yml
`

	cfg, err := SerializeString(data)
	r.Nil(err)
	graph, err := NewConfigurationGraph(cfg)
	r.Nil(err)
	r.NotNil(graph)

	root := graph.Node("root")
	r.NotNil(root)

	leftChild := graph.Node("left-child")
	r.NotNil(leftChild)

	rightChild := graph.Node("right-child")
	r.NotNil(rightChild)

	r.Equal(2, len(root.Children()))
	r.Equal(0, len(root.Parents()))
	r.Equal(0, len(leftChild.Children()))
	r.Equal(1, len(leftChild.Parents()))
	r.Equal(0, len(rightChild.Children()))
	r.Equal(1, len(rightChild.Parents()))

	r.NotNil(graph.Store("root"))
	r.NotNil(graph.Store("left-child"))
	r.NotNil(graph.Store("right-child"))
}

func TestConfigurationGraphCanLoadGraph(t *testing.T) {
	r := require.New(t)
	data := `
stores:
- name: root
  type: file
  stores:
  params:
    path: /test/root.yml
- name: left-child
  type: file
  stores:
  - root
  params:
    path: /test/left-child.yml
- name: right-child
  type: file
  stores:
  - root
  - left-child
  params:
    path: /test/right-child.yml
`
	cfg, err := SerializeString(data)
	r.Nil(err)

	graph, err := NewConfigurationGraph(cfg)
	r.Nil(err)
	r.NotNil(graph)

	root := graph.Node("root")
	r.NotNil(root)

	leftChild := graph.Node("left-child")
	r.NotNil(leftChild)

	rightChild := graph.Node("right-child")
	r.NotNil(rightChild)

	r.Equal(2, len(root.Children()))
	r.Equal(0, len(root.Parents()))
	r.Equal(1, len(leftChild.Children()))
	r.Equal(1, len(leftChild.Parents()))
	r.Equal(0, len(rightChild.Children()))
	r.Equal(2, len(rightChild.Parents()))

	r.NotNil(graph.Store("root"))
	r.NotNil(graph.Store("left-child"))
	r.NotNil(graph.Store("right-child"))
}

func TestConfigurationGraphFailsToCreateCycles(t *testing.T) {
	r := require.New(t)
	data := `
stores:
- name: one
  type: file
  stores:
  - three
  params:
    path: /test/one.yml
- name: two
  type: file
  stores:
  - one
  params:
    path: /test/two.yml
- name: three
  type: file
  stores:
  - two
  params:
    path: /test/three.yml
`
	cfg, err := SerializeString(data)
	r.Nil(err)

	_, err = NewConfigurationGraph(cfg)
	r.NotNil(err)
}
