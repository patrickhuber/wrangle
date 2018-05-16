package config

import (
	"testing"

	yaml "gopkg.in/yaml.v2"
)

func TestCanParseConfig(t *testing.T) {
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
	if err := yaml.Unmarshal([]byte(data), &config); err != nil {
		t.Error(err)
	}

	expectedConfigSourceCount := 1
	actualConfigSourceCount := len(config.ConfigSources)
	if expectedConfigSourceCount != actualConfigSourceCount {
		t.Errorf("Expected %d config sources, found %d", expectedConfigSourceCount, actualConfigSourceCount)
	}

	expectedProcessCount := 1
	actualProcessCount := len(config.Processes)
	if expectedProcessCount != actualProcessCount {
		t.Errorf("Expected %d process sources, found %d", expectedProcessCount, actualProcessCount)
	}

	expectedEnvironmentCount := 1
	actualEnvironmentCount := len(config.Processes[0].Environments)
	if expectedEnvironmentCount != actualEnvironmentCount {
		t.Errorf("Expected %d process[0] environments, found %d", expectedEnvironmentCount, actualEnvironmentCount)
	}
}
