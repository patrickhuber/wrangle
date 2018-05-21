package config

import (
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
- name: name
  environments:
  - name: lab
    config: name
    process: go
    args:
    - version
`
	config := Config{}
	err := yaml.Unmarshal([]byte(data), &config)
	require.Nil(err)
	require.Equal(1, len(config.ConfigSources))
	require.Equal(1, len(config.Processes))
	require.Equal(1, len(config.Processes[0].Environments))
}
