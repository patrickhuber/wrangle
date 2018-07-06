package config

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	yaml "gopkg.in/yaml.v2"
)

func TestCanParseConfig(t *testing.T) {
	r := require.New(t)

	var data = `
config-sources:
- name: name
  type: type
  config: config
  params:
    key: value
environments:
- name: lab
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
	r.Nil(err)

	// config sources)
	r.Equal(1, len(config.ConfigSources))
	r.Equal(1, len(config.ConfigSources[0].Params))
	r.Equal("value", config.ConfigSources[0].Params["key"])

	// environments
	r.Equal(1, len(config.Environments))
	r.Equal(1, len(config.Environments[0].Processes))
	r.Equal(1, len(config.Environments[0].Processes[0].Args))
	r.Equal(1, len(config.Environments[0].Processes[0].Vars))

	// packages
	r.Equal(1, len(config.Packages))
	r.Equal(3, len(config.Packages[0].Platforms))
}
