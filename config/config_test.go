package config

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	yaml "gopkg.in/yaml.v2"
)

func TestCanParseConfig(t *testing.T) {
	require := require.New(t)

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
`
	// vscode likes to be a bad monkey so clean up in case it gets over tabby
	data = strings.Replace(data, "\t", "  ", -1)
	config := Config{}
	err := yaml.Unmarshal([]byte(data), &config)
	require.Nil(err)
	require.Equal(1, len(config.ConfigSources))
	require.Equal(1, len(config.ConfigSources[0].Params))
	require.Equal("value", config.ConfigSources[0].Params["key"])
	require.Equal(1, len(config.Environments))
	require.Equal(1, len(config.Environments[0].Processes))
	require.Equal(1, len(config.Environments[0].Processes[0].Args))
	require.Equal(1, len(config.Environments[0].Processes[0].Vars))
}
