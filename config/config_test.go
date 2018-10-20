package config

import (
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

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
packages:
- name: bosh
  version: 6.7
  platforms:
  - name: linux
    alias: bosh
    download:
      url: https://s3.amazonaws.com/bosh-cli-artifacts/bosh-cli-((version))-linux-amd64
      out: bosh-cli-((version))-linux-amd64		
  - name: windows
    alias: bosh.exe
    download:
      url: https://s3.amazonaws.com/bosh-cli-artifacts/bosh-cli-((version))-windows-amd64.exe
      out: bosh-cli-((version))-windows-amd64.exe		
  - name: darwin
    alias: bosh
    download:
      url: https://s3.amazonaws.com/bosh-cli-artifacts/bosh-cli-((version))-darwin-amd64
      out: bosh-cli-((version))-darwin-amd64
`
		// vscode likes to be a bad monkey so clean up in case it gets over tabby
		data = strings.Replace(data, "\t", "  ", -1)
		config := Config{}
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
		Expect(len(config.Packages)).To(Equal(1))
		Expect(len(config.Packages[0].Platforms)).To(Equal(3))

	})
})
