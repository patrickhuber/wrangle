package config_test

import (
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/wrangle/config"

	yaml "gopkg.in/yaml.v2"
)

var _ = Describe("Config", func() {
	It("can parse config", func() {

		var data = `
stores:
- name: name
  type: type
  config: config
  params:
    key: value
processes:
- name: go
  args:
  - test
  env:
    KEY: value
imports:
- name: bosh
  version: 6.7  
`
		// vscode likes to be a bad monkey so clean up in case it gets over tabby
		data = strings.Replace(data, "\t", "  ", -1)
		config := config.Config{}
		err := yaml.Unmarshal([]byte(data), &config)
		Expect(err).To(BeNil())

		// config sources)
		Expect(len(config.Stores)).To(Equal(1))
		Expect(len(config.Stores[0].Params)).To(Equal(1))
		Expect(config.Stores[0].Params["key"]).To(Equal("value"))

		// environments
		Expect(len(config.Processes)).To(Equal(1))
		Expect(len(config.Processes[0].Args)).To(Equal(1))
		Expect(len(config.Processes[0].Vars)).To(Equal(1))

		// packages
		Expect(len(config.Imports)).To(Equal(1))

	})

	It("can parse package", func() {
		var data = `
name: test
version: 1.0.0
targets:
- platform: windows
  architecture: amd64
  tasks:
  - download:
      uri: https://test.myfile.com
      out: myfile
`
		pkg := config.Package{}
		err := yaml.UnmarshalStrict([]byte(data), &pkg)
		Expect(err).To(BeNil())
	})
})
