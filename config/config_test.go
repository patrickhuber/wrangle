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
processes:
- name: go
  environments:
    - name: lab
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
	require.Equal(1, len(config.Processes))
	require.Equal(1, len(config.Processes[0].Environments))
	require.Equal(1, len(config.Processes[0].Environments[0].Args))
	require.Equal(1, len(config.Processes[0].Environments[0].Vars))
}
