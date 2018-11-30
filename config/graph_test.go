package config_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/wrangle/config"
)

var _ = Describe("graph", func() {
	It("can load linear graph", func() {
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
		cfg, err := config.DeserializeConfigString(data)
		Expect(err).To(BeNil())

		graph, err := config.NewConfigurationGraph(cfg)
		Expect(err).To(BeNil())
		Expect(graph).ToNot(BeNil())

		head := graph.Node("head")
		Expect(head).ToNot(BeNil())

		middle := graph.Node("middle")
		Expect(middle).ToNot(BeNil())

		tail := graph.Node("tail")
		Expect(tail).ToNot(BeNil())

		Expect(len(head.Children())).To(Equal(1))
		Expect(len(head.Parents())).To(Equal(0))
		Expect(len(middle.Children())).To(Equal(1))
		Expect(len(middle.Parents())).To(Equal(1))
		Expect(len(tail.Children())).To(Equal(0))
		Expect(len(tail.Parents())).To(Equal(1))

		Expect(graph.Store("tail")).ToNot(BeNil())
		Expect(graph.Store("head")).ToNot(BeNil())
		Expect(graph.Store("middle")).ToNot(BeNil())
	})
	It("can load tree", func() {
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

		cfg, err := config.DeserializeConfigString(data)
		Expect(err).To(BeNil())

		graph, err := config.NewConfigurationGraph(cfg)
		Expect(err).To(BeNil())
		Expect(graph).ToNot(BeNil())

		root := graph.Node("root")
		Expect(root).ToNot(BeNil())

		leftChild := graph.Node("left-child")
		Expect(leftChild).ToNot(BeNil())

		rightChild := graph.Node("right-child")
		Expect(rightChild).ToNot(BeNil())

		Expect(len(root.Children())).To(Equal(2))
		Expect(len(root.Parents())).To(Equal(0))
		Expect(len(leftChild.Children())).To(Equal(0))
		Expect(len(leftChild.Parents())).To(Equal(1))
		Expect(len(rightChild.Children())).To(Equal(0))
		Expect(len(rightChild.Parents())).To(Equal(1))

		Expect(graph.Store("root")).ToNot(BeNil())
		Expect(graph.Store("left-child")).ToNot(BeNil())
		Expect(graph.Store("right-child")).ToNot(BeNil())
	})

	It("can load graph", func() {

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
		cfg, err := config.DeserializeConfigString(data)
		Expect(err).To(BeNil())

		graph, err := config.NewConfigurationGraph(cfg)
		Expect(err).To(BeNil())
		Expect(graph).ToNot(BeNil())

		root := graph.Node("root")
		Expect(root).ToNot(BeNil())

		leftChild := graph.Node("left-child")
		Expect(leftChild).ToNot(BeNil())

		rightChild := graph.Node("right-child")
		Expect(rightChild).ToNot(BeNil())

		Expect(len(root.Children())).To(Equal(2))
		Expect(len(root.Parents())).To(Equal(0))
		Expect(len(leftChild.Children())).To(Equal(1))
		Expect(len(leftChild.Parents())).To(Equal(1))
		Expect(len(rightChild.Children())).To(Equal(0))
		Expect(len(rightChild.Parents())).To(Equal(2))

		Expect(graph.Store("root")).ToNot(BeNil())
		Expect(graph.Store("left-child")).ToNot(BeNil())
		Expect(graph.Store("right-child")).ToNot(BeNil())

	})

	It("fails to create cycle", func() {
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
		cfg, err := config.DeserializeConfigString(data)
		Expect(err).To(BeNil())

		_, err = config.NewConfigurationGraph(cfg)
		Expect(err).ToNot(BeNil())
	})
})
